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
	GetDeviceCredentials() *generated.DeviceCredentials
	GetArtefactCredentials() *generated.ArtefactCredentials
}

type ArtefactMetaDataImpl struct {
	jobId                string
	deviceIdentifierBlob []byte
	artefactType         generated.ArtefactType
	deviceCredentials    *generated.DeviceCredentials
	artefactCredentials  *generated.ArtefactCredentials
}

func NewArtefactMetaData(jobId string, deviceIdentifierBlob []byte, artefactType generated.ArtefactType, deviceCredentials *generated.DeviceCredentials, artefactCredentials *generated.ArtefactCredentials) *ArtefactMetaDataImpl {
	artefactIdentifier := &ArtefactMetaDataImpl{
		jobId:                jobId,
		deviceIdentifierBlob: deviceIdentifierBlob,
		artefactType:         artefactType,
		deviceCredentials:    deviceCredentials,
		artefactCredentials:  artefactCredentials,
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
		deviceCredentials:    internalMetaData.DeviceCredentials,
		artefactCredentials:  internalMetaData.ArtefactCredentials,
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

func (am *ArtefactMetaDataImpl) GetDeviceCredentials() *generated.DeviceCredentials {
	return am.deviceCredentials
}

func (am *ArtefactMetaDataImpl) GetArtefactCredentials() *generated.ArtefactCredentials {
	return am.artefactCredentials
}
