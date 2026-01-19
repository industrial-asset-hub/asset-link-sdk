/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package reference

//TODO: make sure that an activation actually activates the expected version and that no other preparation has taken place in the mean time
//TODO: support explicit cancellation of an update after preparation

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (m *ReferenceAssetLink) HandlePrepareUpdate(artefactMetaData artefact.ArtefactMetaData, artefactReceiver artefact.ArtefactReceiver) error {
	log.Info().Str("JobID", artefactMetaData.GetJobId()).Msg("Handle Prepare Update")

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
	deviceCredentials := artefactMetaData.GetDeviceCredentials()
	log.Info().Str("JobId", jobId).Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Interface("DeviceCredentials", deviceCredentials).Msg("ArtefactMetaData")

	// Perform checks
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_PREPARE, ga.ArtefactOperationState_AOS_OK, "Performing checks", 0)

	if artefactType != ga.ArtefactType_AT_FIRMWARE {
		err := errors.New("artefact type not supported")
		log.Err(err).Msg("Failed to handle push artefact")
		return err
	}

	// Receiving new firmware
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_DOWNLOAD, ga.ArtefactOperationState_AOS_OK, "Receiving new firmware", 0)

	tempDir, err := os.MkdirTemp("", "firmware_update")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	firmwareFilename := path.Join(tempDir, "firmware.fwu")
	err = artefactReceiver.ReceiveArtefactToFile(firmwareFilename)
	if err != nil {
		log.Err(err).Msg("Failed to receive update file")
		return err
	}

	time.Sleep(2 * time.Second)

	// Verify new firmware, connect to device, and install new firmware
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Verifying new firmware", 0)

	var deviceAddress simdevices.SimulatedDeviceAddress
	err = json.Unmarshal(deviceIdentifierBlob, &deviceAddress)
	if err != nil {
		log.Err(err).Msg("Failed to parse connection blob")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Connecting to device", 0)

	var dc *simdevices.SimulatedDeviceCredentials = nil
	if deviceCredentials != nil {
		dc = &simdevices.SimulatedDeviceCredentials{
			Username: deviceCredentials.Username,
			Password: deviceCredentials.Password,
		}
	}

	device, err := simdevices.ConnectToDevice(deviceAddress, dc)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

	oldFirmwareVersion := device.GetInstalledFirmwareVersion()

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Installing new firmware on device", 0)

	err = device.UpdateFirmware(firmwareFilename)
	if err != nil {
		log.Err(err).Msg("Failed to update device firmware")
		return err
	}

	newFirmwareVersion := device.GetInstalledFirmwareVersion()

	finalMessage := fmt.Sprintf("New firmware installed (new version %s, old version %s)", newFirmwareVersion, oldFirmwareVersion)
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, finalMessage, 100)

	return nil
}

func (m *ReferenceAssetLink) HandleActivateUpdate(artefactMetaData artefact.ArtefactMetaData, artefactReceiver artefact.ArtefactReceiver) error {
	log.Info().Str("JobID", artefactMetaData.GetJobId()).Msg("Handle Activate Update")

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
	deviceCredentials := artefactMetaData.GetDeviceCredentials()
	log.Info().Str("JobId", jobId).Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Interface("DeviceCredentials", deviceCredentials).Msg("ArtefactMetaData")

	// Connect to device and activate new firmware
	var deviceAddress simdevices.SimulatedDeviceAddress
	err := json.Unmarshal(deviceIdentifierBlob, &deviceAddress)
	if err != nil {
		log.Err(err).Msg("Failed to parse connection blob")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_ACTIVATION, ga.ArtefactOperationState_AOS_OK, "Connecting to device", 0)

	var dc *simdevices.SimulatedDeviceCredentials = nil
	if deviceCredentials != nil {
		dc = &simdevices.SimulatedDeviceCredentials{
			Username: deviceCredentials.Username,
			Password: deviceCredentials.Password,
		}
	}

	device, err := simdevices.ConnectToDevice(deviceAddress, dc)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

	if device.GetActiveFirmwareVersion() == device.GetInstalledFirmwareVersion() {
		err = errors.New("installed version and active version are the same")
		log.Err(err).Msg("Failed to activate new firmware")
		return err
	}

	oldFirmwareVersion := device.GetActiveFirmwareVersion()

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_ACTIVATION, ga.ArtefactOperationState_AOS_OK, "Activating new firmware on device", 0)

	err = device.RebootDevice()
	if err != nil {
		log.Err(err).Msg("Failed to reboot device")
		return err
	}

	newFirmwareVersion := device.GetActiveFirmwareVersion()

	finalMessage := fmt.Sprintf("New firmware activated (new version %s, old version %s)", newFirmwareVersion, oldFirmwareVersion)
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_ACTIVATION, ga.ArtefactOperationState_AOS_OK, finalMessage, 100)

	return nil
}

func (m *ReferenceAssetLink) HandleCancelUpdate(artefactMetaData artefact.ArtefactMetaData, statusTransmitter artefact.StatusTransmitter) error {
	log.Info().Str("JobID", artefactMetaData.GetJobId()).Msg("Handle Cancel Update")

	_ = statusTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_CANCELLATION, ga.ArtefactOperationState_AOS_OK, "Update cancelled", 100)

	return nil
}
