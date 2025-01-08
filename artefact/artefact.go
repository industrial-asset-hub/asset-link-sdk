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

type Artefact struct {
	artefactChunks []*generated.ArtefactChunk
	artefactLock   sync.Mutex
}

func NewArtefact() *Artefact {
	var artefactChunks []*generated.ArtefactChunk
	artefactIdentifier := &Artefact{artefactChunks: artefactChunks}
	return artefactIdentifier
}

func (a *Artefact) ClearArtefactChunks() {
	a.artefactLock.Lock()
	defer a.artefactLock.Unlock()
	var artefactChunks []*generated.ArtefactChunk
	a.artefactChunks = artefactChunks
}

func (a *Artefact) AppendArtefactChunk(artefactChunk *generated.ArtefactChunk) {
	a.artefactLock.Lock()
	defer a.artefactLock.Unlock()
	a.artefactChunks = append(a.artefactChunks, artefactChunk)
}

func (a *Artefact) GetArtefactChunks() []*generated.ArtefactChunk {
	a.artefactLock.Lock()
	defer a.artefactLock.Unlock()
	return a.artefactChunks
}

func (a *Artefact) ReadDataFromFile(fileName string) error {
	a.artefactLock.Lock()
	defer a.artefactLock.Unlock()
	// TODO: implement
	return nil
}

func (a *Artefact) WriteDataToFile(fileName string) error {
	a.artefactLock.Lock()
	defer a.artefactLock.Unlock()
	// TODO: implement
	return nil
}
