/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefact

import (
	"io"
	"os"
	"sync"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ArtefactReceiver interface {
	ReceiveArtefactToData() (*[]byte, error)
	ReceiveArtefactToFile(filename string) error
	ReceiveArtefactToWriter(writer io.Writer) error
	UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error
}

type ArtefactReceiverImpl struct {
	stream     grpc.BidiStreamingServer[generated.ArtefactChunk, generated.ArtefactMessage]
	streamLock sync.Mutex
}

func NewArtefactReceiver(stream generated.ArtefactUpdateApi_PushArtefactServer) *ArtefactReceiverImpl {
	artefactReceiver := &ArtefactReceiverImpl{stream: stream}
	return artefactReceiver
}

func (ar *ArtefactReceiverImpl) receiveArtefactChunk() (*generated.ArtefactChunk, error) {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()
	return ar.stream.Recv()
}

func (ar *ArtefactReceiverImpl) ReceiveArtefactMetaData() (ArtefactMetaData, error) {
	chunk, err := ar.receiveArtefactChunk()
	if err != nil {
		return nil, err
	}

	return NewArtefactMetaDataFromInternal(chunk.GetMetadata())
}

func (ar *ArtefactReceiverImpl) ReceiveArtefactToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return ar.ReceiveArtefactToWriter(file)
}

func (ar *ArtefactReceiverImpl) ReceiveArtefactToWriter(writer io.Writer) error {
	err := ar.issueRequest(generated.ArtefactOperationRequestType_AORT_ARTEFACT_TRANSMISSION)
	if err != nil {
		return err
	}

	for {
		chunk, err := ar.receiveArtefactChunk()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		data := chunk.GetFileContent()
		lenData := len(data)
		if lenData > 0 {
			start := 0
			for {
				remainingData := data[start:]
				n, err := writer.Write(remainingData)
				if err != nil {
					return err
				}

				start += n
				if start == lenData {
					break
				}
			}
		}
	}

	return nil
}

func (ar *ArtefactReceiverImpl) ReceiveArtefactToData() (*[]byte, error) {
	err := ar.issueRequest(generated.ArtefactOperationRequestType_AORT_ARTEFACT_TRANSMISSION)
	if err != nil {
		return nil, err
	}

	var data []byte
	for {
		chunk, err := ar.receiveArtefactChunk()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		newData := chunk.GetFileContent()
		data = append(data, newData...)
	}

	return &data, nil
}

func (ar *ArtefactReceiverImpl) transmitRequest(request *generated.ArtefactOperationRequest) error {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()

	message := &generated.ArtefactMessage{Message: &generated.ArtefactMessage_Request{Request: request}}
	return ar.stream.Send(message)
}

func (ar *ArtefactReceiverImpl) issueRequest(requestType generated.ArtefactOperationRequestType) error {
	request := &generated.ArtefactOperationRequest{Type: requestType}

	return ar.transmitRequest(request)
}

func (ar *ArtefactReceiverImpl) transmitStatusUpdate(status *generated.ArtefactOperationStatus) error {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()

	message := &generated.ArtefactMessage{Message: &generated.ArtefactMessage_Status{Status: status}}
	return ar.stream.Send(message)
}

func (ar *ArtefactReceiverImpl) UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error {
	log.Info().Str("Phase", phase.String()).Str("State", state.String()).Str("Message", message).Uint8("Progress", progress).Msg("Status Update")

	statusMessage := &generated.ArtefactOperationStatus{
		Phase:    phase,
		State:    state,
		Message:  message,
		Progress: uint32(progress),
	}

	return ar.transmitStatusUpdate(statusMessage)
}
