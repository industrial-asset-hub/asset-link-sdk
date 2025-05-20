/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package reference

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/artefact"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cdm-al-reference/simdevices"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ga "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
)

func (m *ReferenceAssetLink) HandlePushArtefact(artefactReceiver *artefact.ArtefactReceiver) error {
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

	// Retrieve meta data
	artefactMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact meta data")
		return err
	}

	jobId := artefactMetaData.GetJobId()
	artefactType := artefactMetaData.GetArtefactType()
	deviceIdentifierBlob, err := artefactMetaData.GetDeviceIdentifierBlob()
	if err != nil {
		log.Err(err).Msg("Failed to retrieve device identifier blob")
		return err
	}

	log.Info().Str("JobId", jobId).Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	// Perform checks
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_PREPARE, ga.ArtefactOperationState_AOS_OK, "Performing checks", 0)

	if artefactType != ga.ArtefactType_AT_CONFIGURATION {
		err = errors.New("artefact type not supported")
		log.Err(err).Msg("Failed to handle push artefact")
		return err
	}

	// Receiving new configuration
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_DOWNLOAD, ga.ArtefactOperationState_AOS_OK, "Receiving new configuration", 0)

	tempDir, err := os.MkdirTemp("", "artefact_push")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	artefactFilename := path.Join(tempDir, "artefact_file_in.fwu")
	err = artefactReceiver.ReceiveArtefactToFile(artefactFilename)
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact file")
		return err
	}

	time.Sleep(2 * time.Second)

	// Verify artefact, connect to device, and update device configuration
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Verifying new configuration", 0)

	var deviceAddress simdevices.SimulatedDeviceAddress
	err = json.Unmarshal(deviceIdentifierBlob, &deviceAddress)
	if err != nil {
		log.Err(err).Msg("Failed to parse connection blob")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Connecting to device", 0)

	device, err := simdevices.ConnectToDevice(deviceAddress, nil)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Installing new configuration on device", 0)

	err = device.SetConfig(artefactFilename)
	if err != nil {
		log.Err(err).Msg("Failed to update device configuration")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Installed new configuration on device", 100)

	return nil
}

func (m *ReferenceAssetLink) HandlePullArtefact(artefactMetaData *artefact.ArtefactMetaData, artefactTransmitter *artefact.ArtefactTransmitter) error {
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

	// Retrieve meta data
	jobId := artefactMetaData.GetJobId()
	artefactType := artefactMetaData.GetArtefactType()
	deviceIdentifierBlob, err := artefactMetaData.GetDeviceIdentifierBlob()
	if err != nil {
		log.Err(err).Msg("Failed to retrieve device identifier blob")
		return err
	}

	log.Info().Str("JobId", jobId).Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	// Perform checks
	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_PREPARE, ga.ArtefactOperationState_AOS_OK, "Performing checks", 0)

	if artefactType != ga.ArtefactType_AT_CONFIGURATION {
		err := errors.New("artefact type not supported")
		log.Err(err).Msg("Failed to handle pull artefact")
		return err
	}

	var deviceAddress simdevices.SimulatedDeviceAddress
	err = json.Unmarshal(deviceIdentifierBlob, &deviceAddress)
	if err != nil {
		log.Err(err).Msg("Failed to parse connection blob")
		return err
	}

	// Connect to device and retrieve current configuration from device
	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_ARCHIVE, ga.ArtefactOperationState_AOS_OK, "Connecting to device", 0)

	device, err := simdevices.ConnectToDevice(deviceAddress, nil)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_ARCHIVE, ga.ArtefactOperationState_AOS_OK, "Retrieving current configuration from device", 0)

	tempDir, err := os.MkdirTemp("", "artefact_pull")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	artefactFilename := path.Join(tempDir, "artefact_file_out.fwu")
	err = device.GetConfig(artefactFilename)
	if err != nil {
		log.Err(err).Msg("Failed to retrieve device configuration")
		return err
	}

	// Transmit current configuration
	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_UPLOAD, ga.ArtefactOperationState_AOS_OK, "Transmitting current configuration", 0)

	err = artefactTransmitter.TransmitArtefactFromFile(artefactFilename, 1024)
	if err != nil {
		log.Err(err).Msg("Failed to transmit artefact file")
		return err
	}

	time.Sleep(2 * time.Second)

	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_UPLOAD, ga.ArtefactOperationState_AOS_OK, "Transmitted current configuration", 100)

	return nil
}
