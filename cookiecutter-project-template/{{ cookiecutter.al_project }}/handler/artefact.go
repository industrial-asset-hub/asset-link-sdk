/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/artefact"
	ga "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
)

func (m *AssetLinkImplementation) HandlePushArtefact(artefactMetaData artefact.ArtefactMetaData, artefactReceiver artefact.ArtefactReceiver) error {
	log.Info().Msg("Handle Push Artefact by receiving the artefact")

	// Check if a job is already running
	// We currently support only one running job
	if m.driverLock.TryLock() {
		defer m.driverLock.Unlock()
	} else {
		const errMsg string = "Another job is already running"
		log.Error().Msg(errMsg)
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}

	jobId := artefactMetaData.GetJobId()
	artefactType := artefactMetaData.GetArtefactType()
	deviceIdentifierBlob := artefactMetaData.GetDeviceIdentifierBlob()
	log.Info().Str("JobId", jobId).Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	err := artefactReceiver.ReceiveArtefactToFile("artefact_file")
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact file")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_DOWNLOAD, ga.ArtefactOperationState_AOS_OK, "Status Message", 100)

	return nil
}

func (m *AssetLinkImplementation) HandlePullArtefact(artefactMetaData artefact.ArtefactMetaData, artefactTransmitter artefact.ArtefactTransmitter) error {
	log.Info().Msg("Handle Pull Artefact by transmitting the arefact")

	// Check if a job is already running
	// We currently support only one running job
	if m.driverLock.TryLock() {
		defer m.driverLock.Unlock()
	} else {
		const errMsg string = "Another job is already running"
		log.Error().Msg(errMsg)
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}

	jobId := artefactMetaData.GetJobId()
	artefactType := artefactMetaData.GetArtefactType()
	deviceIdentifierBlob := artefactMetaData.GetDeviceIdentifierBlob()
	log.Info().Str("JobId", jobId).Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	err := artefactTransmitter.TransmitArtefactFromFile("artefact_file", 1024)
	if err != nil {
		log.Err(err).Msg("Failed to transmit artefact file")
		return err
	}

	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_UPLOAD, ga.ArtefactOperationState_AOS_OK, "Status Message", 100)

	return nil
}
