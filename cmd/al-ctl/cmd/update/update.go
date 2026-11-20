// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package update

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type updateRPC func(generated.ArtifactUpdateApiClient, context.Context, *artifact.ArtifactMetaData, ...grpc.CallOption) (grpc.ServerStreamingClient[generated.UpdateEvent], error)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "update", Short: "Run artifact updates on an Asset Link"}
	cmd.AddCommand(
		newOperation("prepare", "Prepare an OCI artifact without installing it", generated.ArtifactUpdateApiClient.PrepareUpdate),
		newOperation("activate", "Activate a prepared update", generated.ArtifactUpdateApiClient.ActivateUpdate),
		newOperation("cancel", "Discard a prepared update", generated.ArtifactUpdateApiClient.CancelUpdate),
	)
	return cmd
}

func newOperation(name, short string, rpc updateRPC) *cobra.Command {
	var requestFile string
	cmd := &cobra.Command{
		Use:          name,
		Short:        short,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runOperation(cmd, requestFile, rpc)
		},
	}
	cmd.Flags().StringVarP(&requestFile, "request", "f", "", "path to the ArtifactMetaData JSON request")
	_ = cmd.MarkFlagRequired("request")
	return cmd
}

func runOperation(cmd *cobra.Command, requestFile string, rpc updateRPC) error {
	request, err := readRequest(requestFile)
	if err != nil {
		return err
	}
	ctx := cmd.Context()
	if shared.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(shared.TimeoutSeconds*float64(time.Second)))
		defer cancel()
	}
	result, err := invokeUpdate(ctx, shared.AssetLinkEndpoint, request, rpc, printProgress(cmd.ErrOrStderr()))
	if err != nil {
		return err
	}
	return printResult(cmd.OutOrStdout(), result)
}

func readRequest(path string) (*artifact.ArtifactMetaData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read update request: %w", err)
	}
	request := new(artifact.ArtifactMetaData)
	if err := protojson.Unmarshal(data, request); err != nil {
		// Unmarshal errors can quote input, which may contain device credentials.
		return nil, errors.New("request file is not valid ArtifactMetaData JSON")
	}
	if err := artifact.ValidateMetaData(request); err != nil {
		return nil, fmt.Errorf("invalid update request: %w", err)
	}
	return request, nil
}

func invokeUpdate(ctx context.Context, endpoint string, request *artifact.ArtifactMetaData, rpc updateRPC, observe func(*generated.UpdateProgress)) (*artifact.Result, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDisableRetry())
	if err != nil {
		return nil, fmt.Errorf("create Asset Link client: %w", err)
	}
	defer conn.Close()

	stream, err := rpc(generated.NewArtifactUpdateApiClient(conn), ctx, request)
	if err != nil {
		return nil, fmt.Errorf("UNKNOWN: could not confirm update outcome: %w", err)
	}
	result, err := artifact.ReadResult(stream, observe)
	if err != nil {
		return nil, fmt.Errorf("UNKNOWN: could not confirm update outcome: %w", err)
	}
	return result, nil
}

func printProgress(w io.Writer) func(*generated.UpdateProgress) {
	return func(progress *generated.UpdateProgress) {
		if progress.Percent == nil {
			fmt.Fprintf(w, "Progress: %q\n", progress.GetStep())
			return
		}
		fmt.Fprintf(w, "Progress: %q %d%%\n", progress.GetStep(), progress.GetPercent())
	}
}

func printResult(w io.Writer, result *artifact.Result) error {
	outcome := strings.TrimPrefix(result.GetOutcome().String(), "UPDATE_OUTCOME_")
	if failure := result.GetFailure(); failure != nil {
		reason := strings.TrimPrefix(failure.GetReason().String(), "FAILURE_REASON_")
		if message := failure.GetPublicMessage(); message != "" {
			return fmt.Errorf("%s (%s): %q", outcome, reason, message)
		}
		return fmt.Errorf("%s (%s)", outcome, reason)
	}
	_, err := fmt.Fprintln(w, outcome)
	return err
}
