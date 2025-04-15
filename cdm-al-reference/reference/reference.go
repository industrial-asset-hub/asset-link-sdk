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
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/artefact"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	ga "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	gd "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements the Discovery interface and feature

type ReferenceClassDriver struct {
	driverLock sync.Mutex
}

type DeviceIdentifierOrConnectionData struct {
	AssetLinkNIC string `json:"al_nic"`
	DeviceIP     string `json:"device_ip"`
}

func (m *ReferenceClassDriver) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error {
	log.Info().Msg("Handle Discovery Request")

	// Check if a job is already running
	// We currently support only one running job
	if m.driverLock.TryLock() {
		defer m.driverLock.Unlock()
	} else {
		const errMsg string = "Another job is already running"
		log.Error().Msg(errMsg)
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}

	alInterface, err := discoveryConfig.GetOptionSettingString("interface", "")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	ipRange, err := discoveryConfig.GetOptionSettingString("ip_range", "")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	devicesFound := ScanDevices(alInterface, ipRange)

	for _, device := range devicesFound {
		deviceInfo := model.NewDevice("EthernetDevice", device.GetDeviceName())
		deviceInfo.AddNameplate(device.GetManufacturer(), device.GetIDLink(), device.GetArticleNumber(),
			device.GetProductDesignation(), device.GetHardwareVersion(), device.GetSerialNumber())

		nicID := deviceInfo.AddNic(device.GetDeviceNIC(), device.GetMacAddress())
		deviceInfo.AddIPv4(nicID, device.GetIpDevice(), device.GetIpNetmask(), device.GetIpRoute())

		deviceInfo.AddSoftware("firmware", device.GetFirmwareVersion())
		deviceInfo.AddCapabilities("firmware_update", device.IsUpdateSupported())

		deviceIdentifier := DeviceIdentifierOrConnectionData{
			AssetLinkNIC: device.GetALNIC(),
			DeviceIP:     device.GetIpDevice(),
		}
		deviceIdentifierBlob, _ := json.Marshal(deviceIdentifier)
		deviceInfo.AddMetadata(string(deviceIdentifierBlob)) // device ID or device connection data used for artefact uploads/downloads

		discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()

		time.Sleep(1 * time.Second)

		err := devicePublisher.PublishDevice(discoveredDevice)
		if err != nil {
			// discovery request was likely cancelled -> terminate discovery and return error
			log.Error().Msgf("Publishing Error: %v", err)
			return err
		}
	}

	return nil
}

func (m *ReferenceClassDriver) GetSupportedFilters() []*gd.SupportedFilter {
	supportedFilters := make([]*gd.SupportedFilter, 0)
	supportedFilters = append(supportedFilters, &gd.SupportedFilter{
		Key:      "ip_range",
		Datatype: gd.VariantType_VT_STRING,
	})
	return supportedFilters
}

func (m *ReferenceClassDriver) GetSupportedOptions() []*gd.SupportedOption {
	supportedOptions := make([]*gd.SupportedOption, 0)
	supportedOptions = append(supportedOptions, &gd.SupportedOption{
		Key:      "interface",
		Datatype: gd.VariantType_VT_STRING,
	})
	return supportedOptions
}

func (m *ReferenceClassDriver) HandlePushArtefact(artefactReceiver *artefact.ArtefactReceiver) error {
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

	deviceIdentifierBlob := artefactMetaData.GetDeviceIdentifier()
	artefactType := artefactMetaData.GetArtefactType()

	log.Info().Str("DeviceIdentifier", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	// Perform checks
	_ = artefactReceiver.UpdateStatus(ga.ArtefactUpdateState_AUS_DOWNLOAD, ga.TransferStatus_AS_OK, "Performing checks", 10)

	if artefactType != ga.ArtefactType_AT_FIRMWARE {
		err = errors.New("artefact type not supported")
		log.Err(err).Msg("Failed to handle push artefact")
		return err
	}

	// Receiving new firmware
	_ = artefactReceiver.UpdateStatus(ga.ArtefactUpdateState_AUS_DOWNLOAD, ga.TransferStatus_AS_OK, "Receiving new firmware", 20)

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

	// Verify new firmware
	_ = artefactReceiver.UpdateStatus(ga.ArtefactUpdateState_AUS_INSTALLATION, ga.TransferStatus_AS_OK, "Verifying new firmware", 50)

	var deviceIdentifier DeviceIdentifierOrConnectionData
	err = json.Unmarshal(deviceIdentifierBlob, &deviceIdentifier)
	if err != nil {
		log.Err(err).Msg("Failed to parse connection blob")
		return err
	}

	// Connect to device
	_ = artefactReceiver.UpdateStatus(ga.ArtefactUpdateState_AUS_INSTALLATION, ga.TransferStatus_AS_OK, "Connecting to device", 60)

	device, err := ConnectToDevice(deviceIdentifier.AssetLinkNIC, deviceIdentifier.DeviceIP)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

	oldFirmwareVersion := device.GetFirmwareVersion()

	// Installing new firmware on device
	_ = artefactReceiver.UpdateStatus(ga.ArtefactUpdateState_AUS_INSTALLATION, ga.TransferStatus_AS_OK, "Installing new firmware on device", 70)

	err = device.UpdateFirmware(artefactFilename)
	if err != nil {
		log.Err(err).Msg("Failed to update device firmware")
		return err
	}

	newFirmwareVersion := device.GetFirmwareVersion()

	// Report successful activation
	finalMessage := fmt.Sprintf("New firmware activated (new version %s, old version %s)", newFirmwareVersion, oldFirmwareVersion)
	_ = artefactReceiver.UpdateStatus(ga.ArtefactUpdateState_AUS_ACTIVATION, ga.TransferStatus_AS_OK, finalMessage, 100)

	return nil
}

func (m *ReferenceClassDriver) HandlePullArtefact(artefactMetaData *artefact.ArtefactMetaData, artefactTransmitter *artefact.ArtefactTransmitter) error {
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
	deviceIdentifierBlob := artefactMetaData.GetDeviceIdentifier()
	artefactType := artefactMetaData.GetArtefactType()
	log.Info().Str("DeviceIdentifier", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	// Perform checks
	_ = artefactTransmitter.UpdateStatus(ga.TransferStatus_AS_OK, "Performing checks")

	if artefactType != ga.ArtefactType_AT_FIRMWARE {
		err := errors.New("artefact type not supported")
		log.Err(err).Msg("Failed to handle pull artefact")
		return err
	}

	var deviceIdentifier DeviceIdentifierOrConnectionData
	err := json.Unmarshal(deviceIdentifierBlob, &deviceIdentifier)
	if err != nil {
		log.Err(err).Msg("Failed to parse connection blob")
		return err
	}

	// Connect to device
	_ = artefactTransmitter.UpdateStatus(ga.TransferStatus_AS_OK, "Connecting to device")

	device, err := ConnectToDevice(deviceIdentifier.AssetLinkNIC, deviceIdentifier.DeviceIP)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

	// Retrieve current firmware from device
	_ = artefactTransmitter.UpdateStatus(ga.TransferStatus_AS_OK, "Retrieving current firmware from device")

	tempDir, err := os.MkdirTemp("", "artefact_pull")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	artefactFilename := path.Join(tempDir, "artefact_file_out.fwu")
	err = device.RetrieveFirmware(artefactFilename)
	if err != nil {
		log.Err(err).Msg("Failed to retrieve device firmware")
		return err
	}

	// Transmit current firmware
	_ = artefactTransmitter.UpdateStatus(ga.TransferStatus_AS_OK, "Transmitting current firmware")

	err = artefactTransmitter.TransmitArtefactFromFile(artefactFilename, 1024)
	if err != nil {
		log.Err(err).Msg("Failed to transmit artefact file")
		return err
	}

	time.Sleep(2 * time.Second)

	// Report successful transmission
	finalMessage := fmt.Sprintf("Firmware transmission complete (version %s)", device.GetFirmwareVersion())
	_ = artefactTransmitter.UpdateStatus(ga.TransferStatus_AS_OK, finalMessage)

	return nil
}
