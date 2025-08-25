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

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"google.golang.org/grpc"
)

type StatusReceiver struct {
	stream  grpc.ServerStreamingClient[generated.ArtefactOperationStatus]
	handler ArtefactMessageHandler
}

func NewStatusReceiver(stream generated.ArtefactUpdateApi_CancelUpdateClient, handler ArtefactMessageHandler) *StatusReceiver {
	statusReceiver := &StatusReceiver{
		stream:  stream,
		handler: handler,
	}
	return statusReceiver
}

func (sr *StatusReceiver) HandleInteraction() error {
	for {
		statusUpdate, err := sr.stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			err = fmt.Errorf("failed to receive status update: %w", err)
			sr.handler.HandleError(err)
			return err
		}

		if statusUpdate != nil {
			sr.handler.HandleStatusUpdate(statusUpdate)
		}
	}

	return nil
}
