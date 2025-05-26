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
	"google.golang.org/grpc"
)

type ArtefactReceiver struct {
	stream     grpc.BidiStreamingServer[generated.ArtefactChunk, generated.ArtefactMessage]
	streamLock sync.Mutex
}

func NewArtefactReceiver(stream generated.ArtefactUpdateApi_PushArtefactServer) *ArtefactReceiver {
	artefactReceiver := &ArtefactReceiver{stream: stream}
	return artefactReceiver
}

func (ar *ArtefactReceiver) receiveArtefactChunk() (*generated.ArtefactChunk, error) {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()
	return ar.stream.Recv()
}

func (ar *ArtefactReceiver) ReceiveArtefactMetaData() (*ArtefactMetaData, error) {
	chunk, err := ar.receiveArtefactChunk()
	if err != nil {
		return nil, err
	}

	internalMetaData := chunk.GetMetadata()

	metaData := NewArtefactMetaData(internalMetaData.JobIdentifier.JobId, internalMetaData.DeviceIdentifier.Blob, internalMetaData.ArtefactIdentifier.Type)

	return metaData, nil
}

func (ar *ArtefactReceiver) ReceiveArtefactToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return ar.ReceiveArtefactToWriter(file)
}

func (ar *ArtefactReceiver) ReceiveArtefactToWriter(writer io.Writer) error {
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

func (ar *ArtefactReceiver) ReceiveArtefactToData() (*[]byte, error) {
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

func (ar *ArtefactReceiver) transmitRequest(request *generated.ArtefactOperationRequest) error {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()

	message := &generated.ArtefactMessage{Message: &generated.ArtefactMessage_Request{Request: request}}
	return ar.stream.Send(message)
}

func (ar *ArtefactReceiver) issueRequest(requestType generated.ArtefactOperationRequestType) error {
	request := &generated.ArtefactOperationRequest{Type: requestType}

	return ar.transmitRequest(request)
}

func (ar *ArtefactReceiver) transmitStatusUpdate(status *generated.ArtefactOperationStatus) error {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()

	message := &generated.ArtefactMessage{Message: &generated.ArtefactMessage_Status{Status: status}}
	return ar.stream.Send(message)
}

func (ar *ArtefactReceiver) UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error {
	statusMessage := &generated.ArtefactOperationStatus{
		Phase:    phase,
		State:    state,
		Message:  message,
		Progress: uint32(progress),
	}

	return ar.transmitStatusUpdate(statusMessage)
}
