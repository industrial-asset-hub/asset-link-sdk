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

type ArtefactTransmitter struct {
	stream           grpc.BidiStreamingClient[generated.ArtefactChunk, generated.ArtefactMessage]
	artefactFilename string
	artefactMetaData *generated.ArtefactMetaData
	handler          ArtefactMessageHandler
}

type TransferState uint8

const (
	NoTransfer       TransferState = 0
	TransferOngoing  TransferState = 1
	TransferComplete TransferState = 2
)

func NewArtefactTransmitter(stream grpc.BidiStreamingClient[generated.ArtefactChunk, generated.ArtefactMessage], artefactFilename string, artefactMetaData *generated.ArtefactMetaData, handler ArtefactMessageHandler) *ArtefactTransmitter {
	artefactTransmitter := &ArtefactTransmitter{
		stream:           stream,
		artefactFilename: artefactFilename,
		artefactMetaData: artefactMetaData,
		handler:          handler,
	}
	return artefactTransmitter
}

func (at *ArtefactTransmitter) HandleInteraction() error {
	chunk := generated.ArtefactChunk{Data: &generated.ArtefactChunk_Metadata{Metadata: at.artefactMetaData}}
	err := at.stream.Send(&chunk)
	if err != nil {
		err = fmt.Errorf("failed to send meta data: %w", err)
		at.handler.HandleError(err)
		return err
	}

	artefactFileIn, err := os.Open(at.artefactFilename)
	if err != nil {
		err = fmt.Errorf("failed to open artefact file: %w", err)
		at.handler.HandleError(err)
		return err
	}
	defer artefactFileIn.Close()

	transferState := NoTransfer
	maxChunkSize := 1024
	var data = make([]byte, maxChunkSize)
	for {
		if transferState == TransferOngoing {
			n, err := artefactFileIn.Read(data)
			if err == io.EOF {
				transferState = TransferComplete
				_ = at.stream.CloseSend()
			} else if err != nil {
				err = fmt.Errorf("failed to read from artefact file: %w", err)
				at.handler.HandleError(err)
				return err
			}

			if n > 0 {
				actualData := data[:n]
				chunk = generated.ArtefactChunk{Data: &generated.ArtefactChunk_FileContent{FileContent: actualData}}
				err = at.stream.Send(&chunk)
				if err != nil {
					err = fmt.Errorf("failed to send artefact chunk: %w", err)
					at.handler.HandleError(err)
					return err
				}
			}
			continue
		}

		message, err := at.stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			err = fmt.Errorf("failed to receive message: %w", err)
			at.handler.HandleError(err)
			return err
		}
		if message != nil {
			statusUpdate := message.GetStatus()
			if statusUpdate == nil {
				at.handler.HandleStatusUpdate(statusUpdate)
			}

			request := message.GetRequest()
			if request != nil {
				switch request.GetType() {
				case generated.ArtefactOperationRequestType_AORT_ARTEFACT_TRANSMISSION:
					if transferState == NoTransfer {
						transferState = TransferOngoing
					} else {
						err = fmt.Errorf("received request to start artefact transmission, but transfer is already ongoing/complete")
						at.handler.HandleError(err)
						return err
					}
				default:
				}
			}
		}
	}

	return nil
}
