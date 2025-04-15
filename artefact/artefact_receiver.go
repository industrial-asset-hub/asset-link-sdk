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

type ArtefactReceiver struct {
	stream     generated.ArtefactUpdateApi_PushArtefactServer
	streamLock sync.Mutex
}

func NewArtefactReceiver(stream generated.ArtefactUpdateApi_PushArtefactServer) *ArtefactReceiver {
	artefactReceiver := &ArtefactReceiver{stream: stream}
	return artefactReceiver
}

func (ar *ArtefactReceiver) ReceiveArtefactChunk() (*generated.ArtefactChunk, error) {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()
	return ar.stream.Recv()
}

func (ar *ArtefactReceiver) ReceiveArtefactMetaData() (*ArtefactMetaData, error) {
	chunk, err := ar.ReceiveArtefactChunk()
	if err != nil {
		return nil, err
	}

	internalMetaData := chunk.GetMetadata()

	metaData := NewArtefactMetaData(internalMetaData.DeviceIdentifier, internalMetaData.ArtefactIdentifier.Type)

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
	for {
		chunk, err := ar.ReceiveArtefactChunk()
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
	var data []byte

	for {
		chunk, err := ar.ReceiveArtefactChunk()
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

func (ar *ArtefactReceiver) TransmitStatusUpdate(status *generated.ArtefactUpdateStatus) error {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()
	return ar.stream.Send(status)
}

func (ar *ArtefactReceiver) UpdateStatus(phase generated.ArtefactUpdateState, status generated.TransferStatus, message string, progress uint8) error {
	statusMessage := &generated.ArtefactUpdateStatus{
		Status: &generated.Status{
			Status:  status,
			Message: message,
		},
		State:    phase,
		Progress: int32(progress),
	}

	return ar.TransmitStatusUpdate(statusMessage)
}
