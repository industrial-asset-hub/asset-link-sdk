// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

// Package contracttest exercises fixture-backed handlers through the SDK wrapper.
// Tests must independently establish the expected device outcome.
package contracttest

import (
	"context"
	"io"
	"testing"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	"google.golang.org/grpc"
)

type recorder struct {
	grpc.ServerStreamingServer[generated.UpdateEvent]
	ctx    context.Context
	events []*generated.UpdateEvent
	index  int
}

func (r *recorder) Context() context.Context { return r.ctx }
func (r *recorder) Send(event *generated.UpdateEvent) error {
	if err := r.ctx.Err(); err != nil {
		return err
	}
	r.events = append(r.events, event)
	return nil
}
func (r *recorder) Recv() (*generated.UpdateEvent, error) {
	if err := r.ctx.Err(); err != nil {
		return nil, err
	}
	if r.index == len(r.events) {
		return nil, io.EOF
	}
	event := r.events[r.index]
	r.index++
	return event, nil
}

// Check validates one scenario and requires the expected outcome. Supply a
// bounded context and independent success, failure and interruption fixtures.
// Cancellation fails the check. Handlers run synchronously and must observe
// the reporter context; Check cannot interrupt a handler that ignores it.
func Check(t testing.TB, ctx context.Context, metadata *artifact.ArtifactMetaData, handler func(*artifact.ArtifactMetaData, artifact.ProgressReporter) (*artifact.Result, error), expected generated.UpdateOutcome) *artifact.Result {
	t.Helper()
	if err := ctx.Err(); err != nil {
		t.Fatalf("contract-test context ended: %v", err)
	}
	if err := artifact.ValidateMetaData(metadata); err != nil {
		t.Fatalf("invalid contract-test request: %v", err)
	}
	stream := &recorder{ctx: ctx}
	if err := artifact.RunUpdate(metadata, stream, handler); err != nil {
		t.Fatalf("update stream failed: %v", err)
	}
	result, err := artifact.ReadResult(stream, nil)
	if err != nil {
		t.Fatalf("invalid update protocol: %v", err)
	}
	if result.GetOutcome() != expected {
		t.Fatalf("update outcome %s; expected %s", result.GetOutcome(), expected)
	}
	return result
}
