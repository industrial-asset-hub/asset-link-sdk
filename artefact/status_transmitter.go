/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artefact

import (
	"context"
	"sync"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatusTransmitter reports update progress. Context allows handlers to stop
// work when the RPC is cancelled or its deadline expires.
type StatusTransmitter interface {
	Context() context.Context
	UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error
}

type StatusTransmitterImpl struct {
	stream   grpc.ServerStreamingServer[generated.ArtefactOperationStatus]
	sendLock sync.Mutex
}

func NewStatusTransmitter(stream grpc.ServerStreamingServer[generated.ArtefactOperationStatus]) *StatusTransmitterImpl {
	return &StatusTransmitterImpl{stream: stream}
}

func (st *StatusTransmitterImpl) Context() context.Context { return st.stream.Context() }

func (st *StatusTransmitterImpl) UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error {
	update, err := operationStatus(phase, state, message, progress)
	if err != nil {
		return err
	}
	st.sendLock.Lock()
	defer st.sendLock.Unlock()
	return st.stream.Send(update)
}

func operationStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) (*generated.ArtefactOperationStatus, error) {
	if progress > 100 {
		return nil, status.Error(codes.InvalidArgument, "update progress must be between 0 and 100")
	}
	return &generated.ArtefactOperationStatus{Phase: phase, State: state, Message: message, Progress: uint32(progress)}, nil
}
