// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT
package artefact

import (
	"context"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

type statusStream struct {
	grpc.ServerStream
	messages []*generated.ArtefactOperationStatus
	err      error
	ctx      context.Context
}

func (s *statusStream) Context() context.Context { return s.ctx }
func (s *statusStream) Send(m *generated.ArtefactOperationStatus) error {
	s.messages = append(s.messages, m)
	return s.err
}
func TestStatusTransmitter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := &statusStream{ctx: ctx}
	transmitter := NewStatusTransmitter(stream)
	require.Equal(t, ctx, transmitter.Context())
	require.NoError(t, transmitter.UpdateStatus(generated.ArtefactOperationPhase_AOP_PREPARE, generated.ArtefactOperationState_AOS_OK, "prepared", 100))
	require.Len(t, stream.messages, 1)
	err := transmitter.UpdateStatus(generated.ArtefactOperationPhase_AOP_PREPARE, generated.ArtefactOperationState_AOS_OK, "invalid", 101)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Len(t, stream.messages, 1)
	stream.err = status.Error(codes.Canceled, "client gone")
	err = transmitter.UpdateStatus(generated.ArtefactOperationPhase_AOP_PREPARE, generated.ArtefactOperationState_AOS_OK, "prepared", 100)
	require.Equal(t, codes.Canceled, status.Code(err))
	cancel()
	require.ErrorIs(t, transmitter.Context().Err(), context.Canceled)
}
