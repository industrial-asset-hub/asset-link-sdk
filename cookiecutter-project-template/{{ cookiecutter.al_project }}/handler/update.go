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

func (m *AssetLinkImplementation) HandlePrepareUpdate(artefactReceiver *artefact.ArtefactReceiver) error {
	log.Info().Msg("Handle Prepare Update")

	// Check if a job is already running
	// We currently support only one running job
	if m.driverLock.TryLock() {
		defer m.driverLock.Unlock()
	} else {
		const errMsg string = "Another job is already running"
		log.Error().Msg(errMsg)
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}

	updateMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact meta data")
		return err
	}

	jobId := updateMetaData.GetJobId()
	artefactType := updateMetaData.GetArtefactType()
	deviceIdentifierBlob, err := updateMetaData.GetDeviceIdentifierBlob()
	if err != nil {
		log.Err(err).Msg("Failed to retrieve device identifier blob")
		return err
	}

	log.Info().Str("JobId", jobId).Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("UpdateMetaData")

	err = artefactReceiver.ReceiveArtefactToFile("artefact_file")
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact file")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_DOWNLOAD, ga.ArtefactOperationState_AOS_OK, "Status Message", 100)

	return nil
}

func (m *AssetLinkImplementation) HandleActivateUpdate(artefactReceiver *artefact.ArtefactReceiver) error {
	log.Info().Msg("Handle Activate Update")

	// Check if a job is already running
	// We currently support only one running job
	if m.driverLock.TryLock() {
		defer m.driverLock.Unlock()
	} else {
		const errMsg string = "Another job is already running"
		log.Error().Msg(errMsg)
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}

	updateMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact meta data")
		return err
	}

	jobId := updateMetaData.GetJobId()
	artefactType := updateMetaData.GetArtefactType()
	deviceIdentifierBlob, err := updateMetaData.GetDeviceIdentifierBlob()
	if err != nil {
		log.Err(err).Msg("Failed to retrieve device identifier blob")
		return err
	}

	log.Info().Str("JobId", jobId).Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("UpdateMetaData")

	err = artefactReceiver.ReceiveArtefactToFile("artefact_file")
	if err != nil {
		log.Err(err).Msg("Failed to receive update file")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_DOWNLOAD, ga.ArtefactOperationState_AOS_OK, "Status Message", 100)

	return nil
}
