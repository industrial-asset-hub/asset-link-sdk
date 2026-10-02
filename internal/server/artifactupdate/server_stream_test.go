// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package artifactupdate

import (
	"context"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	common "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type updateHandler struct {
	handle func(*artifact.ArtifactMetaData, artifact.ProgressReporter) (*artifact.Result, error)
}

func (h *updateHandler) HandlePrepareUpdate(m *artifact.ArtifactMetaData, s artifact.ProgressReporter) (*artifact.Result, error) {
	return h.handle(m, s)
}
func (h *updateHandler) HandleActivateUpdate(m *artifact.ArtifactMetaData, s artifact.ProgressReporter) (*artifact.Result, error) {
	return h.handle(m, s)
}
func (h *updateHandler) HandleCancelUpdate(m *artifact.ArtifactMetaData, s artifact.ProgressReporter) (*artifact.Result, error) {
	return h.handle(m, s)
}

func updateClient(t *testing.T, server *ArtifactUpdateServerEntity) (generated.ArtifactUpdateApiClient, context.Context) {
	t.Helper()
	listener := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	generated.RegisterArtifactUpdateApiServer(s, server)
	go func() { _ = s.Serve(listener) }()
	t.Cleanup(func() { s.Stop(); _ = listener.Close() })
	c, err := grpc.NewClient("passthrough:///update", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return listener.DialContext(ctx) }))
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return generated.NewArtifactUpdateApiClient(c), ctx
}

func wireMetadata() *generated.ArtifactMetaData {
	return &generated.ArtifactMetaData{
		JobId:                  "job-1",
		ConnectionParameterSet: &common.ConnectionParameterSet{SchemaUri: "urn:driver", ParameterJson: `{"key":"device-1"}`, Credentials: []*common.ConnectionCredential{{SchemaUri: "urn:driver", Credentials: "device-secret"}}},
		OciReference:           &generated.OciArtifactReference{RegistryUrl: "http://oci-registry:5000", Repository: "firmware", ManifestDigest: "sha256:" + strings.Repeat("a", 64)},
	}
}

func TestSingleRequestStreamsStatuses(t *testing.T) {
	handler := &updateHandler{handle: func(m *artifact.ArtifactMetaData, s artifact.ProgressReporter) (*artifact.Result, error) {
		if m.GetJobId() != "job-1" || m.GetConnectionParameterSet().GetParameterJson() != `{"key":"device-1"}` || m.GetConnectionParameterSet().Credentials[0].Credentials != "device-secret" || m.GetOciReference().GetRepository() != "firmware" {
			return nil, status.Error(codes.Internal, "incorrect metadata")
		}
		if err := s.Context().Err(); err != nil {
			return nil, err
		}
		if err := s.ReportProgress("preparing", 50); err != nil {
			return nil, err
		}
		return artifact.Succeeded(), nil
	}}
	client, ctx := updateClient(t, &ArtifactUpdateServerEntity{Update: handler})
	for _, open := range []func(context.Context, *generated.ArtifactMetaData, ...grpc.CallOption) (grpc.ServerStreamingClient[generated.UpdateEvent], error){client.PrepareUpdate, client.ActivateUpdate, client.CancelUpdate} {
		stream, err := open(ctx, wireMetadata())
		require.NoError(t, err)
		first, err := stream.Recv()
		require.NoError(t, err)
		require.EqualValues(t, 50, first.GetProgress().GetPercent())
		last, err := stream.Recv()
		require.NoError(t, err)
		require.Equal(t, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED, last.GetResult().GetOutcome())
		_, err = stream.Recv()
		require.ErrorIs(t, err, io.EOF)
	}
}

func TestInvalidMetadataNeverCallsHandler(t *testing.T) {
	for _, mutate := range []func(*generated.ArtifactMetaData){
		func(m *generated.ArtifactMetaData) { m.JobId = "" },
		func(m *generated.ArtifactMetaData) { m.ConnectionParameterSet = nil },
		func(m *generated.ArtifactMetaData) { m.ConnectionParameterSet.ParameterJson = "invalid" },
		func(m *generated.ArtifactMetaData) { m.OciReference = nil },
		func(m *generated.ArtifactMetaData) { m.OciReference.ManifestDigest = "latest" },
	} {
		m := wireMetadata()
		mutate(m)
		client, ctx := updateClient(t, &ArtifactUpdateServerEntity{Update: &updateHandler{handle: func(*artifact.ArtifactMetaData, artifact.ProgressReporter) (*artifact.Result, error) {
			t.Error("invalid request reached handler")
			return nil, nil
		}}})
		for _, open := range []func(context.Context, *generated.ArtifactMetaData, ...grpc.CallOption) (grpc.ServerStreamingClient[generated.UpdateEvent], error){client.PrepareUpdate, client.ActivateUpdate, client.CancelUpdate} {
			stream, err := open(ctx, m)
			require.NoError(t, err)
			_, err = stream.Recv()
			require.Equal(t, codes.InvalidArgument, status.Code(err))
		}
	}
}

func TestHandlerErrorAfterProgress(t *testing.T) {
	client, ctx := updateClient(t, &ArtifactUpdateServerEntity{Update: &updateHandler{handle: func(_ *artifact.ArtifactMetaData, s artifact.ProgressReporter) (*artifact.Result, error) {
		if err := s.ReportProgress("checking", 50); err != nil {
			return nil, err
		}
		return nil, status.Error(codes.FailedPrecondition, "device changed")
	}}})
	stream, err := client.PrepareUpdate(ctx, wireMetadata())
	require.NoError(t, err)
	_, err = stream.Recv()
	require.NoError(t, err)
	result, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN, result.GetResult().GetOutcome())
	_, err = stream.Recv()
	require.ErrorIs(t, err, io.EOF)
}

func TestUnimplementedUpdate(t *testing.T) {
	client, ctx := updateClient(t, &ArtifactUpdateServerEntity{})
	for _, open := range []func(context.Context, *generated.ArtifactMetaData, ...grpc.CallOption) (grpc.ServerStreamingClient[generated.UpdateEvent], error){client.PrepareUpdate, client.ActivateUpdate, client.CancelUpdate} {
		stream, err := open(ctx, wireMetadata())
		require.NoError(t, err)
		_, err = stream.Recv()
		require.Equal(t, codes.Unimplemented, status.Code(err))
	}
}
