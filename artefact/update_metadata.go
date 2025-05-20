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

type UpdateMetaData struct {
	jobId                string
	deviceIdentifierBlob []byte
	artefactType         generated.ArtefactType
}

func NewUpdateMetaData(jobId string, deviceIdentifierBlob []byte, artefactType generated.ArtefactType) *UpdateMetaData {
	artefactIdentifier := &UpdateMetaData{
		jobId:                jobId,
		deviceIdentifierBlob: deviceIdentifierBlob,
		artefactType:         artefactType,
	}
	return artefactIdentifier
}

func (um *UpdateMetaData) GetJobId() string {
	return um.jobId
}

func (um *UpdateMetaData) GetDeviceIdentifierBlob() ([]byte, error) {
	return model.DecodeMetadata(string(um.deviceIdentifierBlob))
}

func (um *UpdateMetaData) GetArtefactType() generated.ArtefactType {
	return um.artefactType
}
