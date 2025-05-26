/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package client

import (
	"fmt"
	"io"
	"os"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"google.golang.org/grpc"
)

type ArtefactReceiver struct {
	stream           grpc.ServerStreamingClient[generated.ArtefactChunk]
	artefactFilename string
	artefactMetaData *generated.ArtefactMetaData
	handler          ArtefactMessageHandler
}

func NewArtefactReceiver(stream grpc.ServerStreamingClient[generated.ArtefactChunk], artefactFilename string, artefactMetaData *generated.ArtefactMetaData, handler ArtefactMessageHandler) *ArtefactReceiver {
	artefactReceiver := &ArtefactReceiver{
		stream:           stream,
		artefactFilename: artefactFilename,
		artefactMetaData: artefactMetaData,
		handler:          handler,
	}
	return artefactReceiver
}

func (ar *ArtefactReceiver) HandleInteraction() error {
	artefactFileOut, err := os.Create(ar.artefactFilename)
	if err != nil {
		err = fmt.Errorf("failed to create artefact file: %w", err)
		ar.handler.HandleError(err)
		return err
	}
	defer artefactFileOut.Close()

	for {
		chunk, err := ar.stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			err = fmt.Errorf("failed to receive artefact chunk: %w", err)
			ar.handler.HandleError(err)
			return err
		}

		if chunk == nil {
			continue
		}

		data := chunk.GetFileContent()
		if data != nil {
			lenData := len(data)
			if lenData > 0 {
				start := 0
				for {
					remainingData := data[start:]
					n, err := artefactFileOut.Write(remainingData)
					if err != nil {
						err = fmt.Errorf("failed to write artefact file: %w", err)
						ar.handler.HandleError(err)
						return err
					}

					start += n
					if start == lenData {
						break
					}
				}
			}
		}

		statusUpdate := chunk.GetStatus()
		if statusUpdate != nil {
			ar.handler.HandleStatusUpdate(statusUpdate)
		}
	}

	return nil
}
