// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

// Package contracttest exercises fixture-backed handlers through the SDK wrapper.
// Tests must independently establish the expected device outcome.
package contracttest

import (
	"context"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artefact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"google.golang.org/grpc"
	"io"
	"testing"
)

type recorder struct {
	grpc.ServerStreamingServer[generated.UpdateEvent]
	ctx    context.Context
	events []*generated.UpdateEvent
	index  int
}

func (r *recorder) Context() context.Context { return r.ctx }
func (r *recorder) Send(event *generated.UpdateEvent) error {
	r.events = append(r.events, event)
	return nil
}
func (r *recorder) Recv() (*generated.UpdateEvent, error) {
	if r.index == len(r.events) {
		return nil, io.EOF
	}
	event := r.events[r.index]
	r.index++
	return event, nil
}

// Check validates one scenario and requires the expected outcome. Supply a
// bounded context and independent success, failure and interruption fixtures.
func Check(t testing.TB, ctx context.Context, metadata *artefact.ArtefactMetaData, handler func(*artefact.ArtefactMetaData, artefact.ProgressReporter) (*artefact.Result, error), expected generated.UpdateOutcome) *artefact.Result {
	t.Helper()
	if err := artefact.ValidateMetaData(metadata); err != nil {
		t.Fatalf("invalid contract-test request: %v", err)
	}
	stream := &recorder{ctx: ctx}
	if err := artefact.RunUpdate(metadata, stream, handler); err != nil {
		t.Fatalf("update stream failed: %v", err)
	}
	result, err := artefact.ReadResult(stream, nil)
	if err != nil {
		t.Fatalf("invalid update protocol: %v", err)
	}
	if result.GetOutcome() != expected {
		t.Fatalf("update outcome %s; expected %s", result.GetOutcome(), expected)
	}
	return result
}
