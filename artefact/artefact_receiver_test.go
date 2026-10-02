/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artefact

import (
	"bytes"
	"context"
	"io"
	"testing"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type receiverStream struct {
	grpc.ServerStream
	chunks   []*generated.ArtefactChunk
	messages []*generated.ArtefactMessage
	err      error
}

func (s *receiverStream) Context() context.Context { return context.Background() }
func (s *receiverStream) Send(m *generated.ArtefactMessage) error {
	s.messages = append(s.messages, m)
	return s.err
}
func (s *receiverStream) Recv() (*generated.ArtefactChunk, error) {
	if len(s.chunks) == 0 {
		return nil, io.EOF
	}
	chunk := s.chunks[0]
	s.chunks = s.chunks[1:]
	return chunk, nil
}

type shortWriter struct{}

func (shortWriter) Write([]byte) (int, error) { return 0, nil }

func TestReceiverRejectsShortWrite(t *testing.T) {
	stream := &receiverStream{chunks: []*generated.ArtefactChunk{{Data: &generated.ArtefactChunk_FileContent{FileContent: []byte("firmware")}}}}
	err := NewArtefactReceiver(stream).ReceiveArtefactToWriter(shortWriter{})
	require.ErrorIs(t, err, io.ErrShortWrite)
}

func TestReceiverRejectsUnexpectedChunks(t *testing.T) {
	for _, chunk := range []*generated.ArtefactChunk{
		{},
		{Data: &generated.ArtefactChunk_Metadata{Metadata: &generated.ArtefactMetaData{}}},
		{Data: &generated.ArtefactChunk_Status{Status: &generated.ArtefactOperationStatus{}}},
	} {
		stream := &receiverStream{chunks: []*generated.ArtefactChunk{chunk}}
		err := NewArtefactReceiver(stream).ReceiveArtefactToWriter(&bytes.Buffer{})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	}
}

func TestReceiverPropagatesSendError(t *testing.T) {
	stream := &receiverStream{err: status.Error(codes.Canceled, "cancelled")}
	err := NewArtefactReceiver(stream).ReceiveArtefactToWriter(&bytes.Buffer{})
	require.Equal(t, codes.Canceled, status.Code(err))
}

func TestReceiverValidatesProgress(t *testing.T) {
	stream := &receiverStream{}
	receiver := NewArtefactReceiver(stream)
	err := receiver.UpdateStatus(generated.ArtefactOperationPhase_AOP_PREPARE, generated.ArtefactOperationState_AOS_OK, "invalid", 101)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Empty(t, stream.messages)
}
