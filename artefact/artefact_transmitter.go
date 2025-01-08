/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefact

import (
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

func (ar *ArtefactReceiver) TransmitArtefact() (*Artefact, error) {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()
	// TODO: implement
	return nil, nil
}
