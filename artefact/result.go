// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package artefact

import (
	"fmt"
	"unicode/utf8"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
)

type Result = generated.UpdateResult
type FailureReason = generated.FailureReason

const (
	InvalidRequest      = generated.FailureReason_FAILURE_REASON_INVALID_REQUEST
	NotSupported        = generated.FailureReason_FAILURE_REASON_NOT_SUPPORTED
	TargetNotFound      = generated.FailureReason_FAILURE_REASON_TARGET_NOT_FOUND
	ArtifactUnavailable = generated.FailureReason_FAILURE_REASON_ARTIFACT_UNAVAILABLE
	ArtifactInvalid     = generated.FailureReason_FAILURE_REASON_ARTIFACT_INVALID
	DeviceRejected      = generated.FailureReason_FAILURE_REASON_DEVICE_REJECTED
	OperationFailed     = generated.FailureReason_FAILURE_REASON_OPERATION_FAILED
	OutcomeNotConfirmed = generated.FailureReason_FAILURE_REASON_OUTCOME_NOT_CONFIRMED
)

func Succeeded() *Result {
	return &Result{Outcome: generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED}
}

// Failed must only be used when the requested change never started, or failure
// and the absence of ongoing device work have been confirmed. It promises no rollback.
func Failed(reason FailureReason, publicMessage string) *Result {
	return &Result{Outcome: generated.UpdateOutcome_UPDATE_OUTCOME_FAILED, Failure: &generated.UpdateFailure{Reason: reason, PublicMessage: publicMessage}}
}

func Unknown(publicMessage string) *Result {
	return &Result{Outcome: generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN, Failure: &generated.UpdateFailure{Reason: generated.FailureReason_FAILURE_REASON_OUTCOME_NOT_CONFIRMED, PublicMessage: publicMessage}}
}

func ValidateResult(result *Result) error {
	if result == nil {
		return fmt.Errorf("missing update result")
	}
	switch result.Outcome {
	case generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED:
		if result.Failure != nil {
			return fmt.Errorf("successful update result contains failure information")
		}
	case generated.UpdateOutcome_UPDATE_OUTCOME_FAILED, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN:
		failure := result.Failure
		if failure == nil || failure.Reason <= generated.FailureReason_FAILURE_REASON_UNSPECIFIED || failure.Reason > generated.FailureReason_FAILURE_REASON_OUTCOME_NOT_CONFIRMED {
			return fmt.Errorf("unsuccessful update result requires a known failure reason")
		}
		if len(failure.PublicMessage) > 512 || !utf8.ValidString(failure.PublicMessage) {
			return fmt.Errorf("invalid public failure message")
		}
		if result.Outcome == generated.UpdateOutcome_UPDATE_OUTCOME_FAILED && failure.Reason == OutcomeNotConfirmed {
			return fmt.Errorf("unconfirmed outcome cannot be a confirmed failure")
		}
	default:
		return fmt.Errorf("unknown or unspecified update outcome")
	}
	return nil
}

func ValidateProgress(progress *generated.UpdateProgress) error {
	if progress == nil || len(progress.Step) > 64 || !utf8.ValidString(progress.Step) {
		return fmt.Errorf("invalid update progress step")
	}
	if progress.Percent != nil && *progress.Percent > 100 {
		return fmt.Errorf("update progress must be between 0 and 100")
	}
	return nil
}
