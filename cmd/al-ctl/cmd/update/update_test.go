// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package update

import (
	"bytes"
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const requestJSON = `{
	"jobId": "demo",
	"connectionParameterSet": {"schemaUri": "urn:simulator", "parameterJson": "{\"device\":\"A0\"}"},
	"ociReference": {"registryUrl": "http://registry:5000", "repository": "firmware",
	"manifestDigest": "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
}`

func TestUpdateCommands(t *testing.T) {
	for _, tc := range []struct {
		name, operation, method string
		result                  *artifact.Result
		streamError             error
		wantError               string
	}{
		{name: "prepare", operation: "prepare", method: generated.ArtifactUpdateApi_PrepareUpdate_FullMethodName, result: artifact.Succeeded()},
		{name: "activate", operation: "activate", method: generated.ArtifactUpdateApi_ActivateUpdate_FullMethodName, result: artifact.Succeeded()},
		{name: "cancel", operation: "cancel", method: generated.ArtifactUpdateApi_CancelUpdate_FullMethodName, result: artifact.Succeeded()},
		{name: "confirmed rejection", operation: "activate", method: generated.ArtifactUpdateApi_ActivateUpdate_FullMethodName, result: artifact.Failed(artifact.DeviceRejected, "Incompatible firmware"), wantError: `FAILED (DEVICE_REJECTED): "Incompatible firmware"`},
		{name: "unknown device outcome", operation: "activate", method: generated.ArtifactUpdateApi_ActivateUpdate_FullMethodName, result: artifact.Unknown(""), wantError: "UNKNOWN (OUTCOME_NOT_CONFIRMED)"},
		{name: "result followed by transport error", operation: "activate", method: generated.ArtifactUpdateApi_ActivateUpdate_FullMethodName, result: artifact.Succeeded(), streamError: status.Error(codes.Unavailable, "connection lost"), wantError: "UNKNOWN: could not confirm update outcome"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requestFile := writeRequest(t, requestJSON)
			wantRequest, err := readRequest(requestFile)
			require.NoError(t, err)
			calls := startUpdateServer(t, tc.result, tc.streamError)
			out, progress, err := runCommand(t, tc.operation, "--request", requestFile)
			require.Equal(t, "Progress: \"connecting\"\nProgress: \"download\" 100%\n", progress)
			call := <-calls
			require.Equal(t, tc.method, call.method)
			require.True(t, proto.Equal(wantRequest, call.request))
			if tc.wantError != "" {
				require.ErrorContains(t, err, tc.wantError)
				require.Empty(t, out, "must not print success")
				return
			}
			require.NoError(t, err)
			require.Equal(t, "SUCCEEDED\n", out)
		})
	}
}

func TestUpdateCommandRejectsInvalidRequest(t *testing.T) {
	for _, tc := range []struct {
		name, input, wantError string
	}{
		{"bad JSON", `{"jobId": "demo", "connectionParameterSet":"device-secret"}`, "not valid ArtifactMetaData JSON"},
		{"invalid digest", strings.Replace(requestJSON, strings.Repeat("a", 64), "latest", 1), "sha256 manifest digest"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, _, err := runCommand(t, "prepare", "--request", writeRequest(t, tc.input))
			require.ErrorContains(t, err, tc.wantError)
			require.NotContains(t, err.Error(), "device-secret")
			require.Empty(t, out)
		})
	}
	_, _, err := runCommand(t, "prepare", "--request", filepath.Join(t.TempDir(), "missing.json"))
	require.ErrorContains(t, err, "read update request")
	_, _, err = runCommand(t, "prepare")
	require.ErrorContains(t, err, `required flag(s) "request" not set`)
}

func TestUpdateCommandHonorsTimeout(t *testing.T) {
	oldTimeout := shared.TimeoutSeconds
	shared.TimeoutSeconds = 0.1
	t.Cleanup(func() { shared.TimeoutSeconds = oldTimeout })
	startUpdateServer(t, nil, nil)
	out, _, err := runCommand(t, "activate", "--request", writeRequest(t, requestJSON))
	require.ErrorContains(t, err, "UNKNOWN: could not confirm update outcome")
	require.ErrorContains(t, err, "DeadlineExceeded")
	require.Empty(t, out)
}

func runCommand(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := NewCommand()
	cmd.SilenceErrors = true
	cmd.SetArgs(args)
	var out, stderr bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	err := cmd.ExecuteContext(ctx)
	return out.String(), stderr.String(), err
}

func writeRequest(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "update.json")
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
	return path
}

type updateCall struct {
	method  string
	request *artifact.ArtifactMetaData
}

func startUpdateServer(t *testing.T, result *artifact.Result, streamError error) <-chan updateCall {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	oldEndpoint := shared.AssetLinkEndpoint
	shared.AssetLinkEndpoint = listener.Addr().String()
	calls := make(chan updateCall, 1)
	server := grpc.NewServer(grpc.StreamInterceptor(func(_ any, stream grpc.ServerStream, info *grpc.StreamServerInfo, _ grpc.StreamHandler) error {
		return serveUpdate(stream, info.FullMethod, calls, result, streamError)
	}))
	generated.RegisterArtifactUpdateApiServer(server, &generated.UnimplementedArtifactUpdateApiServer{})
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
		shared.AssetLinkEndpoint = oldEndpoint
	})
	return calls
}

func serveUpdate(stream grpc.ServerStream, method string, calls chan<- updateCall, result *artifact.Result, streamError error) error {
	request := new(artifact.ArtifactMetaData)
	if err := stream.RecvMsg(request); err != nil {
		return err
	}
	calls <- updateCall{method, request}
	if result == nil {
		<-stream.Context().Done()
		return stream.Context().Err()
	}
	for _, progress := range []*generated.UpdateProgress{{Step: "connecting"}, {Step: "download", Percent: proto.Uint32(100)}} {
		if err := stream.SendMsg(&generated.UpdateEvent{Event: &generated.UpdateEvent_Progress{Progress: progress}}); err != nil {
			return err
		}
	}
	if err := stream.SendMsg(&generated.UpdateEvent{Event: &generated.UpdateEvent_Result{Result: result}}); err != nil {
		return err
	}
	return streamError
}
