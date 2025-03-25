/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefact

import (
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
)

type ArtefactMetaData struct {
	deviceIdentifier []byte
	artefactType     *generated.ArtefactType
}

func NewArtefactMetaData(deviceIdentifier []byte, artefactType *generated.ArtefactType) *ArtefactMetaData {
	artefactIdentifier := &ArtefactMetaData{deviceIdentifier: deviceIdentifier, artefactType: artefactType}
	return artefactIdentifier
}

func (ai *ArtefactMetaData) GetDeviceIdentifier() []byte {
	return ai.deviceIdentifier
}

func (ai *ArtefactMetaData) GetArtefactType() *generated.ArtefactType {
	return ai.artefactType
}
