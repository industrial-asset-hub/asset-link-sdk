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
)

type UpdateActivateReceiver struct {
	stream     generated.ArtefactUpdateApi_PrepareUpdateServer
	streamLock sync.Mutex
}

func NewUpdateActivateReceiver(stream generated.ArtefactUpdateApi_PrepareUpdateServer) *UpdateActivateReceiver {
	updateActivateReceiver := &UpdateActivateReceiver{stream: stream}
	return updateActivateReceiver
}

func (uar *UpdateActivateReceiver) ReceiveUpdateChunk() (*generated.ArtefactChunk, error) {
	uar.streamLock.Lock()
	defer uar.streamLock.Unlock()
	return uar.stream.Recv()
}

func (uar *UpdateActivateReceiver) ReceiveUpdateMetaData() (*UpdateMetaData, error) {
	chunk, err := uar.ReceiveUpdateChunk()
	if err != nil {
		return nil, err
	}

	internalMetaData := chunk.GetMetadata()

	metaData := NewUpdateMetaData(internalMetaData.JobIdentifier.JobId, internalMetaData.DeviceIdentifier.Blob, internalMetaData.ArtefactIdentifier.Type)

	return metaData, nil
}

func (uar *UpdateActivateReceiver) ReceiveUpdateToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return uar.ReceiveUpdateToWriter(file)
}

func (uar *UpdateActivateReceiver) ReceiveUpdateToWriter(writer io.Writer) error {
	for {
		chunk, err := uar.ReceiveUpdateChunk()
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

func (uar *UpdateActivateReceiver) ReceiveUpdateToData() (*[]byte, error) {
	var data []byte

	for {
		chunk, err := uar.ReceiveUpdateChunk()
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

func (uar *UpdateActivateReceiver) TransmitRequest(request *generated.ArtefactOperationRequest) error {
	uar.streamLock.Lock()
	defer uar.streamLock.Unlock()

	message := &generated.ArtefactUpdateMessage{Message: &generated.ArtefactUpdateMessage_Request{Request: request}}
	return uar.stream.Send(message)
}

func (uar *UpdateActivateReceiver) IssueRequest(requestType generated.ArtefactOperationRequestType) error {
	request := &generated.ArtefactOperationRequest{Type: requestType}

	return uar.TransmitRequest(request)
}

func (uar *UpdateActivateReceiver) TransmitStatusUpdate(status *generated.ArtefactOperationStatus) error {
	uar.streamLock.Lock()
	defer uar.streamLock.Unlock()

	message := &generated.ArtefactUpdateMessage{Message: &generated.ArtefactUpdateMessage_Status{Status: status}}
	return uar.stream.Send(message)
}

func (uar *UpdateActivateReceiver) UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error {
	statusMessage := &generated.ArtefactOperationStatus{
		Phase:    phase,
		State:    state,
		Message:  message,
		Progress: uint32(progress),
	}

	return uar.TransmitStatusUpdate(statusMessage)
}
