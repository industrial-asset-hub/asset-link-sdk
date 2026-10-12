// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package artefactupdate

import (
	"context"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artefact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	common "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type updateHandler struct {
	handle func(*artefact.ArtefactMetaData, artefact.StatusTransmitter) error
}

func (h *updateHandler) HandlePrepareUpdate(m *artefact.ArtefactMetaData, s artefact.StatusTransmitter) error {
	return h.handle(m, s)
}
func (h *updateHandler) HandleActivateUpdate(m *artefact.ArtefactMetaData, s artefact.StatusTransmitter) error {
	return h.handle(m, s)
}
func (h *updateHandler) HandleCancelUpdate(m *artefact.ArtefactMetaData, s artefact.StatusTransmitter) error {
	return h.handle(m, s)
}

func updateClient(t *testing.T, server *ArtefactUpdateServerEntity) (generated.ArtefactUpdateApiClient, context.Context) {
	t.Helper()
	listener := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	generated.RegisterArtefactUpdateApiServer(s, server)
	go func() { _ = s.Serve(listener) }()
	t.Cleanup(func() { s.Stop(); _ = listener.Close() })
	c, err := grpc.NewClient("passthrough:///update", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return listener.DialContext(ctx) }))
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return generated.NewArtefactUpdateApiClient(c), ctx
}

func wireMetadata() *generated.ArtefactMetaData {
	return &generated.ArtefactMetaData{
		JobId: "job-1", AssetId: "asset-1", ViewId: "view-1", ArtefactType: generated.ArtefactType_AT_FIRMWARE,
		ConnectionParameterSet: &common.ConnectionParameterSet{SchemaUri: "urn:driver", ParameterJson: `{"key":"device-1"}`, Credentials: []*common.ConnectionCredential{{SchemaUri: "urn:driver", Credentials: "device-secret"}}},
		OciReference:           &generated.OciArtefactReference{RegistryUrl: "http://zot:5000", Repository: "firmware", ManifestDigest: "sha256:" + strings.Repeat("a", 64)},
	}
}

func TestSingleRequestStreamsStatuses(t *testing.T) {
	handler := &updateHandler{handle: func(m *artefact.ArtefactMetaData, s artefact.StatusTransmitter) error {
		if m.GetJobId() != "job-1" || m.GetViewId() != "view-1" || m.GetAssetId() != "asset-1" || m.GetConnectionParameterSet().GetParameterJson() != `{"key":"device-1"}` || m.GetConnectionParameterSet().Credentials[0].Credentials != "device-secret" || m.GetOciReference().GetRepository() != "firmware" {
			return status.Error(codes.Internal, "incorrect metadata")
		}
		if err := s.Context().Err(); err != nil {
			return err
		}
		if err := s.UpdateStatus(generated.ArtefactOperationPhase_AOP_PREPARE, generated.ArtefactOperationState_AOS_OK, "preparing", 50); err != nil {
			return err
		}
		return s.UpdateStatus(generated.ArtefactOperationPhase_AOP_PREPARE, generated.ArtefactOperationState_AOS_OK, "prepared", 100)
	}}
	client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: handler})
	for _, open := range []func(context.Context, *generated.ArtefactMetaData, ...grpc.CallOption) (grpc.ServerStreamingClient[generated.ArtefactOperationStatus], error){client.PrepareUpdate, client.ActivateUpdate, client.CancelUpdate} {
		stream, err := open(ctx, wireMetadata())
		require.NoError(t, err)
		first, err := stream.Recv()
		require.NoError(t, err)
		require.EqualValues(t, 50, first.Progress)
		last, err := stream.Recv()
		require.NoError(t, err)
		require.EqualValues(t, 100, last.Progress)
		_, err = stream.Recv()
		require.ErrorIs(t, err, io.EOF)
	}
}

func TestInvalidMetadataNeverCallsHandler(t *testing.T) {
	for _, mutate := range []func(*generated.ArtefactMetaData){
		func(m *generated.ArtefactMetaData) { m.JobId = "" },
		func(m *generated.ArtefactMetaData) { m.ConnectionParameterSet = nil },
		func(m *generated.ArtefactMetaData) { m.ConnectionParameterSet.ParameterJson = "invalid" },
		func(m *generated.ArtefactMetaData) { m.OciReference = nil },
		func(m *generated.ArtefactMetaData) { m.OciReference.ManifestDigest = "latest" },
		func(m *generated.ArtefactMetaData) { m.ArtefactType = 999 },
	} {
		m := wireMetadata()
		mutate(m)
		client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: &updateHandler{handle: func(*artefact.ArtefactMetaData, artefact.StatusTransmitter) error {
			t.Error("invalid request reached handler")
			return nil
		}}})
		for _, open := range []func(context.Context, *generated.ArtefactMetaData, ...grpc.CallOption) (grpc.ServerStreamingClient[generated.ArtefactOperationStatus], error){client.PrepareUpdate, client.ActivateUpdate, client.CancelUpdate} {
			stream, err := open(ctx, m)
			require.NoError(t, err)
			_, err = stream.Recv()
			require.Equal(t, codes.InvalidArgument, status.Code(err))
		}
	}
}

func TestHandlerErrorAfterProgress(t *testing.T) {
	client, ctx := updateClient(t, &ArtefactUpdateServerEntity{Update: &updateHandler{handle: func(_ *artefact.ArtefactMetaData, s artefact.StatusTransmitter) error {
		if err := s.UpdateStatus(generated.ArtefactOperationPhase_AOP_PREPARE, generated.ArtefactOperationState_AOS_OK, "checking", 50); err != nil {
			return err
		}
		return status.Error(codes.FailedPrecondition, "device changed")
	}}})
	stream, err := client.PrepareUpdate(ctx, wireMetadata())
	require.NoError(t, err)
	_, err = stream.Recv()
	require.NoError(t, err)
	_, err = stream.Recv()
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestUnimplementedUpdate(t *testing.T) {
	client, ctx := updateClient(t, &ArtefactUpdateServerEntity{})
	for _, open := range []func(context.Context, *generated.ArtefactMetaData, ...grpc.CallOption) (grpc.ServerStreamingClient[generated.ArtefactOperationStatus], error){client.PrepareUpdate, client.ActivateUpdate, client.CancelUpdate} {
		stream, err := open(ctx, wireMetadata())
		require.NoError(t, err)
		_, err = stream.Recv()
		require.Equal(t, codes.Unimplemented, status.Code(err))
	}
}
