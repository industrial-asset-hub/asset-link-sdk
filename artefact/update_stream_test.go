// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package artefact

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type eventStream struct {
	grpc.ServerStreamingServer[generated.UpdateEvent]
	events []*generated.UpdateEvent
	err    error
	index  int
}

func (*eventStream) Context() context.Context { return context.Background() }
func (s *eventStream) Send(event *generated.UpdateEvent) error {
	s.events = append(s.events, event)
	return s.err
}
func (s *eventStream) Recv() (*generated.UpdateEvent, error) {
	if s.index < len(s.events) {
		event := s.events[s.index]
		s.index++
		return event, nil
	}
	if s.err != nil {
		return nil, s.err
	}
	return nil, io.EOF
}
func resultEvent(result *Result) *generated.UpdateEvent {
	return &generated.UpdateEvent{Event: &generated.UpdateEvent_Result{Result: result}}
}
func progressEvent(percent *uint32) *generated.UpdateEvent {
	return &generated.UpdateEvent{Event: &generated.UpdateEvent_Progress{Progress: &generated.UpdateProgress{Step: "download", Percent: percent}}}
}

func TestRunUpdateConservativeDefaults(t *testing.T) {
	for name, tc := range map[string]struct {
		result  *Result
		err     error
		outcome generated.UpdateOutcome
	}{
		"success":                       {Succeeded(), nil, generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED},
		"confirmed failure":             {Failed(generated.FailureReason_FAILURE_REASON_ARTIFACT_UNAVAILABLE, "Missing firmware"), errors.New("internal diagnostic"), generated.UpdateOutcome_UPDATE_OUTCOME_FAILED},
		"unknown":                       {Unknown("Check device"), nil, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN},
		"unknown plus error":            {Unknown("Check device"), errors.New("private connection diagnostic"), generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN},
		"unknown with confirmed reason": {&Result{Outcome: generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN, Failure: &generated.UpdateFailure{Reason: ArtifactUnavailable}}, nil, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN},
		"nil result":                    {nil, nil, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN},
		"ordinary error":                {nil, errors.New("not found"), generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN},
		"success plus error":            {Succeeded(), errors.New("connection lost"), generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN},
		"unspecified":                   {&Result{}, nil, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN},
		"malformed failure":             {Failed(0, "bad"), nil, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN},
	} {
		t.Run(name, func(t *testing.T) {
			stream := &eventStream{}
			err := RunUpdate(&ArtefactMetaData{JobId: "job"}, stream, func(_ *ArtefactMetaData, p ProgressReporter) (*Result, error) {
				require.NoError(t, p.ReportIndeterminate("connecting"))
				require.NoError(t, p.ReportProgress("download", 100))
				return tc.result, tc.err
			})
			require.NoError(t, err)
			require.Len(t, stream.events, 3)
			result, err := ReadResult(stream, nil)
			require.NoError(t, err)
			require.Equal(t, tc.outcome, result.Outcome)
			if result.Outcome == generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN {
				require.Equal(t, OutcomeNotConfirmed, result.GetFailure().GetReason())
			}
			if name == "unknown plus error" {
				require.Equal(t, "Check device", result.GetFailure().GetPublicMessage())
			}
		})
	}
}

func TestRunUpdateHandlesPanicAndRetainedReporter(t *testing.T) {
	stream := &eventStream{}
	var reporter ProgressReporter
	require.NoError(t, RunUpdate(nil, stream, func(_ *ArtefactMetaData, p ProgressReporter) (*Result, error) { reporter = p; panic("driver bug") }))
	result, err := ReadResult(stream, nil)
	require.NoError(t, err)
	require.Equal(t, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN, result.Outcome)
	require.Error(t, reporter.ReportProgress("late", 50))
	require.Len(t, stream.events, 1)
}

func TestRunUpdateDoesNotHideSendOrProgressErrors(t *testing.T) {
	for _, invalid := range []bool{false, true} {
		stream := &eventStream{}
		if !invalid {
			stream.err = errors.New("send failed")
		}
		err := RunUpdate(nil, stream, func(_ *ArtefactMetaData, p ProgressReporter) (*Result, error) {
			percent := uint32(50)
			if invalid {
				percent = 101
			}
			_ = p.ReportProgress("download", percent)
			return Succeeded(), nil
		})
		require.Error(t, err)
		if !invalid {
			require.ErrorIs(t, err, stream.err)
		}
	}
}

func TestReadResultRejectsInvalidStreams(t *testing.T) {
	invalid := uint32(101)
	for name, events := range map[string][]*generated.UpdateEvent{
		"empty stream": nil, "nil event": {nil}, "empty event": {{}},
		"progress only": {progressEvent(nil)}, "invalid percentage": {progressEvent(&invalid)},
		"nil progress": {{Event: &generated.UpdateEvent_Progress{}}},
		"nil result":   {resultEvent(nil)}, "unspecified result": {resultEvent(&Result{})},
		"unknown outcome":               {resultEvent(&Result{Outcome: 999})},
		"failure without reason":        {resultEvent(Failed(0, ""))},
		"failure unknown reason":        {resultEvent(Failed(999, ""))},
		"failure unconfirmed outcome":   {resultEvent(Failed(OutcomeNotConfirmed, ""))},
		"unknown with confirmed reason": {resultEvent(&Result{Outcome: generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN, Failure: &generated.UpdateFailure{Reason: ArtifactUnavailable}})},
		"success with failure":          {resultEvent(&Result{Outcome: generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED, Failure: &generated.UpdateFailure{}})},
		"long public message":           {resultEvent(Unknown(strings.Repeat("x", 513)))},
		"long step":                     {{Event: &generated.UpdateEvent_Progress{Progress: &generated.UpdateProgress{Step: strings.Repeat("x", 65)}}}},
		"duplicate result":              {resultEvent(Succeeded()), resultEvent(Succeeded())},
		"progress after result":         {resultEvent(Succeeded()), progressEvent(nil)},
	} {
		t.Run(name, func(t *testing.T) {
			result, err := ReadResult(&eventStream{events: events}, nil)
			require.Error(t, err)
			require.Nil(t, result)
		})
	}
	stream := &eventStream{events: []*generated.UpdateEvent{resultEvent(Succeeded())}, err: errors.New("connection lost")}
	result, err := ReadResult(stream, nil)
	require.ErrorIs(t, err, stream.err)
	require.Nil(t, result)
}
