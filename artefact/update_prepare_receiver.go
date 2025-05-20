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

type UpdatePrepareReceiver struct {
	stream     generated.ArtefactUpdateApi_PrepareUpdateServer
	streamLock sync.Mutex
}

func NewUpdatePrepareReceiver(stream generated.ArtefactUpdateApi_PrepareUpdateServer) *UpdatePrepareReceiver {
	updatePrepareReceiver := &UpdatePrepareReceiver{stream: stream}
	return updatePrepareReceiver
}

func (upr *UpdatePrepareReceiver) ReceiveUpdateChunk() (*generated.ArtefactChunk, error) {
	upr.streamLock.Lock()
	defer upr.streamLock.Unlock()
	return upr.stream.Recv()
}

func (upr *UpdatePrepareReceiver) ReceiveUpdateMetaData() (*UpdateMetaData, error) {
	chunk, err := upr.ReceiveUpdateChunk()
	if err != nil {
		return nil, err
	}

	internalMetaData := chunk.GetMetadata()

	metaData := NewUpdateMetaData(internalMetaData.JobIdentifier.JobId, internalMetaData.DeviceIdentifier.Blob, internalMetaData.ArtefactIdentifier.Type)

	return metaData, nil
}

func (upr *UpdatePrepareReceiver) ReceiveUpdateToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return upr.ReceiveUpdateToWriter(file)
}

func (upr *UpdatePrepareReceiver) ReceiveUpdateToWriter(writer io.Writer) error {
	for {
		chunk, err := upr.ReceiveUpdateChunk()
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

func (upr *UpdatePrepareReceiver) ReceiveUpdateToData() (*[]byte, error) {
	var data []byte

	for {
		chunk, err := upr.ReceiveUpdateChunk()
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

func (upr *UpdatePrepareReceiver) TransmitRequest(request *generated.ArtefactOperationRequest) error {
	upr.streamLock.Lock()
	defer upr.streamLock.Unlock()

	message := &generated.ArtefactUpdateMessage{Message: &generated.ArtefactUpdateMessage_Request{Request: request}}
	return upr.stream.Send(message)
}

func (upr *UpdatePrepareReceiver) IssueRequest(requestType generated.ArtefactOperationRequestType) error {
	request := &generated.ArtefactOperationRequest{Type: requestType}

	return upr.TransmitRequest(request)
}

func (upr *UpdatePrepareReceiver) TransmitStatusUpdate(status *generated.ArtefactOperationStatus) error {
	upr.streamLock.Lock()
	defer upr.streamLock.Unlock()

	message := &generated.ArtefactUpdateMessage{Message: &generated.ArtefactUpdateMessage_Status{Status: status}}
	return upr.stream.Send(message)
}

func (upr *UpdatePrepareReceiver) UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error {
	statusMessage := &generated.ArtefactOperationStatus{
		Phase:    phase,
		State:    state,
		Message:  message,
		Progress: uint32(progress),
	}

	return upr.TransmitStatusUpdate(statusMessage)
}
