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

type ArtefactMetaData interface {
	GetJobId() string
	GetDeviceIdentifierBlob() []byte
	GetArtefactType() generated.ArtefactType
}

type ArtefactMetaDataImpl struct {
	jobId                string
	deviceIdentifierBlob []byte
	artefactType         generated.ArtefactType
}

func NewArtefactMetaData(jobId string, deviceIdentifierBlob []byte, artefactType generated.ArtefactType) *ArtefactMetaDataImpl {
	artefactIdentifier := &ArtefactMetaDataImpl{
		jobId:                jobId,
		deviceIdentifierBlob: deviceIdentifierBlob,
		artefactType:         artefactType,
	}
	return artefactIdentifier
}

func NewArtefactMetaDataFromInternal(internalMetaData *generated.ArtefactMetaData) (*ArtefactMetaDataImpl, error) {
	if internalMetaData == nil {
		return nil, nil
	}

	deviceIdentifierBlob, err := model.DecodeMetadata(string(internalMetaData.DeviceIdentifier.Blob))
	if err != nil {
		return nil, err
	}

	artefactIdentifier := &ArtefactMetaDataImpl{
		jobId:                internalMetaData.JobIdentifier.JobId,
		deviceIdentifierBlob: deviceIdentifierBlob,
		artefactType:         internalMetaData.ArtefactIdentifier.Type,
	}
	return artefactIdentifier, nil
}

func (am *ArtefactMetaDataImpl) GetJobId() string {
	return am.jobId
}

func (am *ArtefactMetaDataImpl) GetDeviceIdentifierBlob() []byte {
	return am.deviceIdentifierBlob
}

func (am *ArtefactMetaDataImpl) GetArtefactType() generated.ArtefactType {
	return am.artefactType
}
