/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artefactupdate

import (
	"context"
	"encoding/base64"
	"io"
	"net"
	"testing"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artefact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

type updateHandler struct {
	prepare  func(artefact.ArtefactMetaData, artefact.ArtefactReceiver) error
	activate func(artefact.ArtefactMetaData, artefact.ArtefactReceiver) error
	cancel   func(artefact.ArtefactMetaData, artefact.StatusTransmitter) error
}

func (h *updateHandler) HandlePrepareUpdate(m artefact.ArtefactMetaData, r artefact.ArtefactReceiver) error {
	return h.prepare(m, r)
}
func (h *updateHandler) HandleActivateUpdate(m artefact.ArtefactMetaData, r artefact.ArtefactReceiver) error {
	return h.activate(m, r)
}
func (h *updateHandler) HandleCancelUpdate(m artefact.ArtefactMetaData, s artefact.StatusTransmitter) error {
	return h.cancel(m, s)
}

type updateRequirementsHandler struct {
	*updateHandler
	requirements func(context.Context, *generated.GetUpdateTargetRequirementsRequest) (*generated.GetUpdateTargetRequirementsResponse, error)
}

func (h *updateRequirementsHandler) GetUpdateTargetRequirements(ctx context.Context, request *generated.GetUpdateTargetRequirementsRequest) (*generated.GetUpdateTargetRequirementsResponse, error) {
	return h.requirements(ctx, request)
}

func updateClient(t *testing.T, server *ArtefactUpdateServerEntity) (generated.ArtefactUpdateApiClient, context.Context) {
	t.Helper()
	listener := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	generated.RegisterArtefactUpdateApiServer(grpcServer, server)
	go func() { _ = grpcServer.Serve(listener) }()
	t.Cleanup(func() { grpcServer.Stop(); _ = listener.Close() })
	connection, err := grpc.NewClient("passthrough:///update", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return listener.DialContext(ctx) }))
	require.NoError(t, err)
	t.Cleanup(func() { _ = connection.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return generated.NewArtefactUpdateApiClient(connection), ctx
}

func wireMetadata() *generated.ArtefactMetaData {
	return &generated.ArtefactMetaData{
		JobIdentifier:      &generated.JobIdentifier{JobId: "job-1"},
		ArtefactIdentifier: &generated.ArtefactIdentifier{Type: generated.ArtefactType_AT_FIRMWARE},
		DeviceIdentifier:   &generated.DeviceIdentifier{Blob: []byte(base64.StdEncoding.EncodeToString([]byte("device-1")))},
		DeviceCredentials:  &generated.DeviceCredentials{Username: "user", Password: "secret"},
	}
}

func TestUpdateStreams(t *testing.T) {
	for _, operation := range []string{"prepare", "activate"} {
		t.Run(operation, func(t *testing.T) {
			result := make(chan []byte, 1)
			handler := func(metadata artefact.ArtefactMetaData, receiver artefact.ArtefactReceiver) error {
				if metadata.GetJobId() != "job-1" || string(metadata.GetDeviceIdentifierBlob()) != "device-1" || metadata.GetDeviceCredentials().Password != "secret" {
					return status.Error(codes.Internal, "incorrect metadata")
				}
				// Progress must be sendable even while Recv is waiting for file content.
				downloaded := make(chan error, 1)
				go func() {
					data, err := receiver.ReceiveArtefactToData()
					if err == nil {
						result <- *data
					}
					downloaded <- err
				}()
				if err := receiver.UpdateStatus(generated.ArtefactOperationPhase_AOP_DOWNLOAD, generated.ArtefactOperationState_AOS_OK, "downloading", 0); err != nil {
					return err
				}
				if err := <-downloaded; err != nil {
					return err
				}
				return receiver.UpdateStatus(generated.ArtefactOperationPhase_AOP_INSTALLATION, generated.ArtefactOperationState_AOS_OK, "done", 100)
			}
			client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: &updateHandler{prepare: handler, activate: handler}})
			var stream grpc.BidiStreamingClient[generated.ArtefactChunk, generated.ArtefactMessage]
			var err error
			if operation == "prepare" {
				stream, err = client.PrepareUpdate(ctx)
			} else {
				stream, err = client.ActivateUpdate(ctx)
			}
			require.NoError(t, err)
			require.NoError(t, stream.Send(&generated.ArtefactChunk{Data: &generated.ArtefactChunk_Metadata{Metadata: wireMetadata()}}))
			var sawRequest, sawProgress bool
			for range 2 {
				message, err := stream.Recv()
				require.NoError(t, err)
				sawRequest = sawRequest || message.GetRequest() != nil
				sawProgress = sawProgress || message.GetStatus() != nil
			}
			require.True(t, sawRequest)
			require.True(t, sawProgress)
			for _, content := range []string{"firm", "ware"} {
				require.NoError(t, stream.Send(&generated.ArtefactChunk{Data: &generated.ArtefactChunk_FileContent{FileContent: []byte(content)}}))
			}
			require.NoError(t, stream.CloseSend())
			message, err := stream.Recv()
			require.NoError(t, err)
			require.Equal(t, uint32(100), message.GetStatus().Progress)
			_, err = stream.Recv()
			require.ErrorIs(t, err, io.EOF)
			require.Equal(t, []byte("firmware"), <-result)
		})
	}
}

func TestInvalidUpdateMetadata(t *testing.T) {
	for _, test := range []struct {
		name     string
		metadata *generated.ArtefactMetaData
	}{
		{"missing identifiers", &generated.ArtefactMetaData{}},
		{"invalid encoding", &generated.ArtefactMetaData{JobIdentifier: &generated.JobIdentifier{}, ArtefactIdentifier: &generated.ArtefactIdentifier{}, DeviceIdentifier: &generated.DeviceIdentifier{Blob: []byte("!")}}},
		{"unknown type", &generated.ArtefactMetaData{JobIdentifier: &generated.JobIdentifier{}, ArtefactIdentifier: &generated.ArtefactIdentifier{Type: 999}, DeviceIdentifier: &generated.DeviceIdentifier{}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: &updateHandler{}})
			stream, err := client.PrepareUpdate(ctx)
			require.NoError(t, err)
			require.NoError(t, stream.Send(&generated.ArtefactChunk{Data: &generated.ArtefactChunk_Metadata{Metadata: test.metadata}}))
			_, err = stream.Recv()
			require.Equal(t, codes.InvalidArgument, status.Code(err))
			cancel, err := client.CancelUpdate(ctx, test.metadata)
			require.NoError(t, err)
			_, err = cancel.Recv()
			require.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	}
}

func TestMissingFirstMetadata(t *testing.T) {
	for _, empty := range []bool{true, false} {
		client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: &updateHandler{}})
		stream, err := client.ActivateUpdate(ctx)
		require.NoError(t, err)
		if !empty {
			require.NoError(t, stream.Send(&generated.ArtefactChunk{Data: &generated.ArtefactChunk_FileContent{FileContent: []byte("firmware")}}))
		}
		require.NoError(t, stream.CloseSend())
		_, err = stream.Recv()
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	}
}

func TestCancelUpdateStatusAndHandlerError(t *testing.T) {
	client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: &updateHandler{
		cancel: func(m artefact.ArtefactMetaData, s artefact.StatusTransmitter) error {
			if m.GetJobId() != "job-1" {
				return status.Error(codes.Internal, "incorrect job")
			}
			if err := s.UpdateStatus(generated.ArtefactOperationPhase_AOP_CANCELLATION, generated.ArtefactOperationState_AOS_OK, "cancelling", 50); err != nil {
				return err
			}
			return status.Error(codes.FailedPrecondition, "already activated")
		},
	}})
	stream, err := client.CancelUpdate(ctx, wireMetadata())
	require.NoError(t, err)
	message, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, generated.ArtefactOperationPhase_AOP_CANCELLATION, message.Phase)
	_, err = stream.Recv()
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestUnimplementedUpdate(t *testing.T) {
	client, ctx := updateClient(t, &ArtefactUpdateServerEntity{})
	prepare, err := client.PrepareUpdate(ctx)
	require.NoError(t, err)
	_, err = prepare.Recv()
	require.Equal(t, codes.Unimplemented, status.Code(err))
	activate, err := client.ActivateUpdate(ctx)
	require.NoError(t, err)
	_, err = activate.Recv()
	require.Equal(t, codes.Unimplemented, status.Code(err))
	cancel, err := client.CancelUpdate(ctx, wireMetadata())
	require.NoError(t, err)
	_, err = cancel.Recv()
	require.Equal(t, codes.Unimplemented, status.Code(err))
}

func TestGetUpdateTargetRequirements(t *testing.T) {
	expected := &generated.GetUpdateTargetRequirementsResponse{
		IdentifierFormatUri: "urn:siemens:iah:sat-driver:update-target:v1",
		Fields: []*generated.UpdateTargetField{
			{Name: "article_number", AssetJsonPointer: "/product_instance_information/manufacturer_product/product_id"},
			{Name: "serial_number", AssetJsonPointer: "/product_instance_information/serial_number"},
		},
	}
	client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: &updateRequirementsHandler{
		updateHandler: &updateHandler{},
		requirements: func(_ context.Context, request *generated.GetUpdateTargetRequirementsRequest) (*generated.GetUpdateTargetRequirementsResponse, error) {
			require.Equal(t, generated.ArtefactType_AT_FIRMWARE, request.ArtefactIdentifier.Type)
			return expected, nil
		},
	}})
	response, err := client.GetUpdateTargetRequirements(ctx, &generated.GetUpdateTargetRequirementsRequest{
		ArtefactIdentifier: &generated.ArtefactIdentifier{Type: generated.ArtefactType_AT_FIRMWARE},
	})
	require.NoError(t, err)
	require.True(t, proto.Equal(expected, response))
}

func TestGetUpdateTargetRequirementsUnavailable(t *testing.T) {
	for _, update := range []struct {
		name    string
		handler *ArtefactUpdateServerEntity
	}{
		{"no update", &ArtefactUpdateServerEntity{}},
		{"legacy update", &ArtefactUpdateServerEntity{Update: &updateHandler{}}},
	} {
		t.Run(update.name, func(t *testing.T) {
			client, ctx := updateClient(t, update.handler)
			_, err := client.GetUpdateTargetRequirements(ctx, &generated.GetUpdateTargetRequirementsRequest{
				ArtefactIdentifier: &generated.ArtefactIdentifier{Type: generated.ArtefactType_AT_FIRMWARE},
			})
			require.Equal(t, codes.Unimplemented, status.Code(err))
		})
	}
}

func TestGetUpdateTargetRequirementsRejectsInvalidType(t *testing.T) {
	client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: &updateRequirementsHandler{
		updateHandler: &updateHandler{},
		requirements: func(context.Context, *generated.GetUpdateTargetRequirementsRequest) (*generated.GetUpdateTargetRequirementsResponse, error) {
			return nil, status.Error(codes.Internal, "unexpected call")
		},
	}})
	for _, request := range []*generated.GetUpdateTargetRequirementsRequest{
		{},
		{ArtefactIdentifier: &generated.ArtefactIdentifier{Type: 999}},
	} {
		_, err := client.GetUpdateTargetRequirements(ctx, request)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	}
}

func TestGetUpdateTargetRequirementsPropagatesHandlerError(t *testing.T) {
	client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: &updateRequirementsHandler{
		updateHandler: &updateHandler{},
		requirements: func(context.Context, *generated.GetUpdateTargetRequirementsRequest) (*generated.GetUpdateTargetRequirementsResponse, error) {
			return nil, status.Error(codes.FailedPrecondition, "no identity available")
		},
	}})
	_, err := client.GetUpdateTargetRequirements(ctx, &generated.GetUpdateTargetRequirementsRequest{
		ArtefactIdentifier: &generated.ArtefactIdentifier{Type: generated.ArtefactType_AT_FIRMWARE},
	})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
}
