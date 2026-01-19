/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package client

import (
	"errors"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
)

type ArtefactMessageHandler interface {
	HandleStatusUpdate(statusUpdate *generated.ArtefactOperationStatus)
	HandleError(err error)
}

func ArtefactCreateArtefactIdentifier(artefactType string) (*generated.ArtefactIdentifier, error) {
	artefactIdentifier := generated.ArtefactIdentifier{Type: generated.ArtefactType_AT_CONFIGURATION}

	switch artefactType {
	case "software":
		artefactIdentifier.Type = generated.ArtefactType_AT_SOFTWARE
	case "firmware":
		artefactIdentifier.Type = generated.ArtefactType_AT_FIRMWARE
	case "configuration":
		artefactIdentifier.Type = generated.ArtefactType_AT_CONFIGURATION
	case "backup":
		artefactIdentifier.Type = generated.ArtefactType_AT_BACKUP
	case "logfile":
		artefactIdentifier.Type = generated.ArtefactType_AT_LOGFILE
	default:
		return nil, errors.New("invalid artefact type")
	}

	return &artefactIdentifier, nil
}

func ArtefactCreateMetadata(jobId string, deviceIdentifierBlob []byte, artefactType string, deviceCredentials *generated.DeviceCredentials, artefactCredentials *generated.ArtefactCredentials) (*generated.ArtefactMetaData, error) {
	jobIdentifier := generated.JobIdentifier{JobId: jobId}

	artefactIdentifier, err := ArtefactCreateArtefactIdentifier(artefactType)
	if err != nil {
		return nil, err
	}

	deviceIdentifier := generated.DeviceIdentifier{Blob: deviceIdentifierBlob}

	artefactMetaData := &generated.ArtefactMetaData{
		JobIdentifier:       &jobIdentifier,
		DeviceIdentifier:    &deviceIdentifier,
		DeviceCredentials:   deviceCredentials,
		ArtefactIdentifier:  artefactIdentifier,
		ArtefactCredentials: artefactCredentials,
	}

	return artefactMetaData, nil
}
