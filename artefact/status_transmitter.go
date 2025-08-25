/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefact

import (
	"sync"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type StatusTransmitter interface {
	UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error
}

type StatusTransmitterImpl struct {
	stream     grpc.ServerStreamingServer[generated.ArtefactOperationStatus]
	streamLock sync.Mutex
}

func NewStatusTransmitter(stream generated.ArtefactUpdateApi_CancelUpdateServer) *StatusTransmitterImpl {
	statusTransmitter := &StatusTransmitterImpl{stream: stream}
	return statusTransmitter
}

func (st *StatusTransmitterImpl) transmitStatusUpdate(status *generated.ArtefactOperationStatus) error {
	st.streamLock.Lock()
	defer st.streamLock.Unlock()

	return st.stream.Send(status)
}

func (st *StatusTransmitterImpl) UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error {
	log.Info().Str("Phase", phase.String()).Str("State", state.String()).Str("Message", message).Uint8("Progress", progress).Msg("Status Update")

	statusMessage := &generated.ArtefactOperationStatus{
		Phase:    phase,
		State:    state,
		Message:  message,
		Progress: uint32(progress),
	}

	return st.transmitStatusUpdate(statusMessage)
}
