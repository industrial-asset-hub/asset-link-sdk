// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT
package contracttest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	common "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/stretchr/testify/require"
)

type fatalProbe struct{ testing.TB }

func (fatalProbe) Fatalf(format string, args ...any) { panic(fmt.Sprintf(format, args...)) }

func contractMetadata() *artifact.ArtifactMetaData {
	return &artifact.ArtifactMetaData{JobId: "job", ConnectionParameterSet: &common.ConnectionParameterSet{SchemaUri: "urn:fixture", ParameterJson: `{}`}, OciReference: &generated.OciArtifactReference{RegistryUrl: "http://registry:5000", Repository: "firmware", ManifestDigest: "sha256:" + strings.Repeat("a", 64)}}
}

func TestCheckFixtures(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	metadata := contractMetadata()
	handler := func(_ *artifact.ArtifactMetaData, p artifact.ProgressReporter) (*artifact.Result, error) {
		return artifact.Succeeded(), p.ReportProgress("download", 100)
	}
	result := Check(t, ctx, metadata, handler, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED)
	require.Equal(t, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED, result.Outcome)
	probe := fatalProbe{TB: t}
	require.Panics(t, func() { Check(probe, ctx, nil, handler, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED) })
	require.Panics(t, func() { Check(probe, ctx, metadata, handler, generated.UpdateOutcome_UPDATE_OUTCOME_FAILED) })
	require.Panics(t, func() {
		Check(probe, ctx, metadata, func(_ *artifact.ArtifactMetaData, p artifact.ProgressReporter) (*artifact.Result, error) {
			return nil, p.ReportProgress("download", 101)
		}, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN)
	})
}

func TestRecorderRespectsContext(t *testing.T) {
	for _, expired := range []bool{false, true} {
		name := "canceled"
		if expired {
			name = "deadline exceeded"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if expired {
				var deadlineCancel context.CancelFunc
				ctx, deadlineCancel = context.WithDeadline(ctx, time.Now().Add(-time.Second))
				defer deadlineCancel()
			} else {
				cancel()
			}
			event := &generated.UpdateEvent{}
			stream := &recorder{ctx: ctx, events: []*generated.UpdateEvent{event}}
			require.ErrorIs(t, stream.Send(event), ctx.Err())
			require.Len(t, stream.events, 1)
			got, err := stream.Recv()
			require.ErrorIs(t, err, ctx.Err())
			require.Nil(t, got)
		})
	}
}

func TestCheckRejectsCanceledContextBeforeHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	handler := func(*artifact.ArtifactMetaData, artifact.ProgressReporter) (*artifact.Result, error) {
		t.Error("handler invoked with canceled context")
		return artifact.Succeeded(), nil
	}
	require.PanicsWithValue(t, "contract-test context ended: context canceled", func() {
		Check(fatalProbe{TB: t}, ctx, contractMetadata(), handler, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED)
	})
}

func TestCheckRejectsCancellationDuringHandler(t *testing.T) {
	for _, reportProgress := range []bool{false, true} {
		name := "final result"
		if reportProgress {
			name = "progress"
		}
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			handler := func(_ *artifact.ArtifactMetaData, reporter artifact.ProgressReporter) (*artifact.Result, error) {
				cancel()
				if reportProgress {
					require.ErrorIs(t, reporter.ReportProgress("download", 100), context.Canceled)
				}
				return artifact.Succeeded(), nil
			}
			require.PanicsWithValue(t, "update stream failed: context canceled", func() {
				Check(fatalProbe{TB: t}, ctx, contractMetadata(), handler, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED)
			})
		})
	}
}
