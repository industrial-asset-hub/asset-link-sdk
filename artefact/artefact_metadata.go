/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefact

import (
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
)

type ArtefactMetaData struct {
	deviceIdentifierBlob []byte
	artefactType         generated.ArtefactType
}

func NewArtefactMetaData(deviceIdentifierBlob []byte, artefactType generated.ArtefactType) *ArtefactMetaData {
	artefactIdentifier := &ArtefactMetaData{deviceIdentifierBlob: deviceIdentifierBlob, artefactType: artefactType}
	return artefactIdentifier
}

func (ai *ArtefactMetaData) GetDeviceIdentifierBlob() ([]byte, error) {
	return model.DecodeMetadata(string(ai.deviceIdentifierBlob))
}

func (ai *ArtefactMetaData) GetArtefactType() generated.ArtefactType {
	return ai.artefactType
}
