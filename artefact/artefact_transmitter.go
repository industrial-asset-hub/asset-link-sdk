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

type ArtefactTransmitter struct {
	stream     generated.ArtefactUpdateApi_PullArtefactServer
	streamLock sync.Mutex
}

func NewArtefactTransmitter(stream generated.ArtefactUpdateApi_PullArtefactServer) *ArtefactTransmitter {
	artefactTransmitter := &ArtefactTransmitter{stream: stream}
	return artefactTransmitter
}

func (at *ArtefactTransmitter) TransmitArtefactChunk(artefactChunk *generated.ArtefactChunk) error {
	at.streamLock.Lock()
	defer at.streamLock.Unlock()
	return at.stream.Send(artefactChunk)
}

func (at *ArtefactTransmitter) TransmitArtefactMetaData(artefactMetaData *generated.ArtefactMetaData) error {
	return at.TransmitArtefactChunk(&generated.ArtefactChunk{Data: &generated.ArtefactChunk_Metadata{Metadata: artefactMetaData}})
}

func (at *ArtefactTransmitter) TransmitArtefactFromFile(filename string, maxChunkSize uint64) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return at.TransmitArtefactFromReader(file, maxChunkSize)
}

func (at *ArtefactTransmitter) TransmitArtefactFromReader(reader io.Reader, maxChunkSize uint64) error {
	var data = make([]byte, maxChunkSize)
	for {
		n, err := reader.Read(data)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if n > 0 {
			actualData := data[:n]
			chunk := generated.ArtefactChunk{Data: &generated.ArtefactChunk_FileContent{FileContent: actualData}}
			err = at.TransmitArtefactChunk(&chunk)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (at *ArtefactTransmitter) TransmitArtefactFromData(data *[]byte, maxChunkSize uint64) error {
	chunkSize := int(maxChunkSize)
	lenData := len(*data)
	for start := 0; start < lenData; start += chunkSize {
		end := start + chunkSize

		if end > lenData {
			end = lenData
		}

		chunkData := (*data)[start:end]

		chunk := generated.ArtefactChunk{Data: &generated.ArtefactChunk_FileContent{FileContent: chunkData}}
		err := at.TransmitArtefactChunk(&chunk)
		if err != nil {
			return err
		}
	}

	return nil
}
