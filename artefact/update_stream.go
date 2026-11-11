// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package artefact

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// ProgressReporter cannot send final results. The SDK owns stream completion.
type ProgressReporter interface {
	Context() context.Context
	ReportProgress(step string, percent uint32) error
	ReportIndeterminate(step string) error
}

type updateTransmitter struct {
	stream    grpc.ServerStreamingServer[generated.UpdateEvent]
	mu        sync.Mutex
	finished  bool
	sendError error
}

func (t *updateTransmitter) Context() context.Context { return t.stream.Context() }
func (t *updateTransmitter) ReportProgress(step string, percent uint32) error {
	return t.progress(&generated.UpdateProgress{Step: step, Percent: &percent})
}
func (t *updateTransmitter) ReportIndeterminate(step string) error {
	return t.progress(&generated.UpdateProgress{Step: step})
}
func (t *updateTransmitter) progress(progress *generated.UpdateProgress) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.finished {
		return fmt.Errorf("update result already sent")
	}
	if t.sendError != nil {
		return t.sendError
	}
	if err := ValidateProgress(progress); err != nil {
		t.sendError = err
		return err
	}
	t.sendError = t.stream.Send(&generated.UpdateEvent{Event: &generated.UpdateEvent_Progress{Progress: progress}})
	return t.sendError
}

// RunUpdate is also available to custom servers and conformance tests. Ordinary
// handler errors, nil/invalid results and panics default to UNKNOWN, never FAILED.
func RunUpdate(metadata *ArtefactMetaData, stream grpc.ServerStreamingServer[generated.UpdateEvent], handler func(*ArtefactMetaData, ProgressReporter) (*Result, error)) error {
	t := &updateTransmitter{stream: stream}
	result, err := invokeUpdate(handler, metadata, t)
	if err != nil {
		log.Warn().Err(err).Str("job_id", metadata.GetJobId()).Msg("Asset Link update handler returned an error")
	}
	if validationError := ValidateResult(result); validationError != nil {
		log.Warn().Err(validationError).Str("job_id", metadata.GetJobId()).Msg("Invalid Asset Link result; reporting unknown outcome")
		result = Unknown("The operation outcome could not be confirmed. Check the Asset Link logs.")
	} else if err != nil && result.GetOutcome() != generated.UpdateOutcome_UPDATE_OUTCOME_FAILED {
		result = Unknown("The operation outcome could not be confirmed. Check the Asset Link logs.")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.finished = true
	if t.sendError != nil {
		return t.sendError
	}
	return stream.Send(&generated.UpdateEvent{Event: &generated.UpdateEvent_Result{Result: result}})
}

func invokeUpdate(handler func(*ArtefactMetaData, ProgressReporter) (*Result, error), metadata *ArtefactMetaData, reporter ProgressReporter) (result *Result, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = nil
			err = fmt.Errorf("update handler panic: %v", recovered)
		}
	}()
	return handler(metadata, reporter)
}

type UpdateReceiver interface {
	Recv() (*generated.UpdateEvent, error)
}

// ReadResult validates the entire stream, including normal EOF after exactly
// one result. A result followed by a transport error is not authoritative.
func ReadResult(stream UpdateReceiver, observe func(*generated.UpdateProgress)) (*Result, error) {
	var result *Result
	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			if result == nil {
				return nil, fmt.Errorf("update stream ended without a final result")
			}
			return result, nil
		}
		if err != nil {
			return nil, fmt.Errorf("update stream interrupted: %w", err)
		}
		if event == nil || event.Event == nil {
			return nil, fmt.Errorf("empty update event")
		}
		if result != nil {
			return nil, fmt.Errorf("update event received after final result")
		}
		switch payload := event.Event.(type) {
		case *generated.UpdateEvent_Progress:
			if err := ValidateProgress(payload.Progress); err != nil {
				return nil, err
			}
			if observe != nil {
				observe(payload.Progress)
			}
		case *generated.UpdateEvent_Result:
			if err := ValidateResult(payload.Result); err != nil {
				return nil, err
			}
			result = payload.Result
		default:
			return nil, fmt.Errorf("unknown update event")
		}
	}
}
