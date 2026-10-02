/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artefact

import (
	"bytes"
	"context"
	"io"
	"os"
	"sync"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ArtefactReceiver requests file content from the client and reports progress.
// File content ends when the client closes its send side of the stream.
type ArtefactReceiver interface {
	StatusTransmitter
	ReceiveArtefactToData() (*[]byte, error)
	ReceiveArtefactToFile(filename string) error
	ReceiveArtefactToWriter(writer io.Writer) error
}

type ArtefactReceiverImpl struct {
	stream      grpc.BidiStreamingServer[generated.ArtefactChunk, generated.ArtefactMessage]
	receiveLock sync.Mutex
	sendLock    sync.Mutex
}

func NewArtefactReceiver(stream grpc.BidiStreamingServer[generated.ArtefactChunk, generated.ArtefactMessage]) *ArtefactReceiverImpl {
	return &ArtefactReceiverImpl{stream: stream}
}

func (ar *ArtefactReceiverImpl) Context() context.Context { return ar.stream.Context() }

func (ar *ArtefactReceiverImpl) ReceiveArtefactMetaData() (ArtefactMetaData, error) {
	ar.receiveLock.Lock()
	defer ar.receiveLock.Unlock()
	chunk, err := ar.stream.Recv()
	if err == io.EOF {
		return nil, status.Error(codes.InvalidArgument, "metadata must be the first update chunk")
	}
	if err != nil {
		return nil, err
	}
	return NewArtefactMetaDataFromInternal(chunk.GetMetadata())
}

func (ar *ArtefactReceiverImpl) ReceiveArtefactToFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()
	return ar.ReceiveArtefactToWriter(file)
}

func (ar *ArtefactReceiverImpl) ReceiveArtefactToWriter(writer io.Writer) error {
	ar.receiveLock.Lock()
	defer ar.receiveLock.Unlock()
	if err := ar.send(&generated.ArtefactMessage{Message: &generated.ArtefactMessage_Request{
		Request: &generated.ArtefactOperationRequest{Type: generated.ArtefactOperationRequestType_AORT_ARTEFACT_TRANSMISSION},
	}}); err != nil {
		return err
	}
	for {
		chunk, err := ar.stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		content, ok := chunk.GetData().(*generated.ArtefactChunk_FileContent)
		if !ok {
			return status.Error(codes.InvalidArgument, "expected file content after update metadata")
		}
		if _, err := io.Copy(writer, bytes.NewReader(content.FileContent)); err != nil {
			return err
		}
	}
}

func (ar *ArtefactReceiverImpl) ReceiveArtefactToData() (*[]byte, error) {
	var buffer bytes.Buffer
	if err := ar.ReceiveArtefactToWriter(&buffer); err != nil {
		return nil, err
	}
	data := buffer.Bytes()
	return &data, nil
}

func (ar *ArtefactReceiverImpl) send(message *generated.ArtefactMessage) error {
	ar.sendLock.Lock()
	defer ar.sendLock.Unlock()
	return ar.stream.Send(message)
}

func (ar *ArtefactReceiverImpl) UpdateStatus(phase generated.ArtefactOperationPhase, state generated.ArtefactOperationState, message string, progress uint8) error {
	update, err := operationStatus(phase, state, message, progress)
	if err != nil {
		return err
	}
	return ar.send(&generated.ArtefactMessage{Message: &generated.ArtefactMessage_Status{Status: update}})
}
