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

func (ar *ArtefactReceiver) ReceiveArtefact() (*Artefact, error) {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()
	// TODO: implement
	return nil, nil
}

func (ar *ArtefactReceiver) UpdateStatus(status *generated.ArtefactUpdateStatus) error {
	ar.streamLock.Lock()
	defer ar.streamLock.Unlock()
	return ar.stream.Send(status)
}
