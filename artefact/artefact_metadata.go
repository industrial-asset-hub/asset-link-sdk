/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artefact

import (
	"encoding/base64"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ArtefactMetaData identifies an update job and its target device.
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
	return &ArtefactMetaDataImpl{jobId, deviceIdentifierBlob, artefactType, deviceCredentials, artefactCredentials}
}

// NewArtefactMetaDataFromInternal validates wire metadata and decodes the
// base64-encoded device identifier used by the update protocol.
func NewArtefactMetaDataFromInternal(metadata *generated.ArtefactMetaData) (*ArtefactMetaDataImpl, error) {
	if metadata == nil || metadata.JobIdentifier == nil || metadata.ArtefactIdentifier == nil || metadata.DeviceIdentifier == nil {
		return nil, status.Error(codes.InvalidArgument, "update metadata requires job, artefact, and device identifiers")
	}
	if _, ok := generated.ArtefactType_name[int32(metadata.ArtefactIdentifier.Type)]; !ok {
		return nil, status.Error(codes.InvalidArgument, "unknown artefact type")
	}
	blob, err := base64.StdEncoding.DecodeString(string(metadata.DeviceIdentifier.Blob))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "device identifier must be base64-encoded")
	}
	return NewArtefactMetaData(metadata.JobIdentifier.JobId, blob, metadata.ArtefactIdentifier.Type, metadata.DeviceCredentials, metadata.ArtefactCredentials), nil
}

func (am *ArtefactMetaDataImpl) GetJobId() string                        { return am.jobId }
func (am *ArtefactMetaDataImpl) GetDeviceIdentifierBlob() []byte         { return am.deviceIdentifierBlob }
func (am *ArtefactMetaDataImpl) GetArtefactType() generated.ArtefactType { return am.artefactType }
func (am *ArtefactMetaDataImpl) GetDeviceCredentials() *generated.DeviceCredentials {
	return am.deviceCredentials
}
func (am *ArtefactMetaDataImpl) GetArtefactCredentials() *generated.ArtefactCredentials {
	return am.artefactCredentials
}
