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
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	ga "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	gd "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements both the Discovery and the Identifiers interface/feature

type ReferenceAssetLink struct {
	driverLock sync.Mutex
}

func (m *ReferenceAssetLink) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error {
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

	ipRange, err := discoveryConfig.GetFilterSettingString("ip_range", "")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	devicesAddressesFound := simdevices.ScanForDevices(alInterface, ipRange)

	for _, address := range devicesAddressesFound {
		// connect to device and retrieve its details
		device, err := simdevices.RetrieveDeviceDetails(address, nil) // provide no credentials
		if err != nil {
			log.Error().Err(err).Msg("Could not retrieve device details")
			discoverError := createDiscoverError(err, address)
			pubErr := devicePublisher.PublishError(discoverError)
			if pubErr != nil {
				log.Error().Err(pubErr).Msg("Failed to publish device error")
				return pubErr
			}
			continue // try next device
		}

		deviceInfo := createDeviceInfo(device)
		discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()
		err = devicePublisher.PublishDevice(discoveredDevice)
		if err != nil {
			// discovery request was likely cancelled -> terminate discovery and return error
			log.Error().Msgf("Publishing Error: %v", err)
			return err
		}
	}

	return nil
}

func (m *ReferenceAssetLink) GetSupportedOptions() []*gd.SupportedOption {
	supportedOptions := make([]*gd.SupportedOption, 0)
	supportedOptions = append(supportedOptions, &gd.SupportedOption{
		Key:      "interface",
		Datatype: gd.VariantType_VT_STRING,
	})
	return supportedOptions
}

func (m *ReferenceAssetLink) GetSupportedFilters() []*gd.SupportedFilter {
	supportedFilters := make([]*gd.SupportedFilter, 0)
	supportedFilters = append(supportedFilters, &gd.SupportedFilter{
		Key:      "ip_range",
		Datatype: gd.VariantType_VT_STRING,
	})
	return supportedFilters
}

func (m *ReferenceAssetLink) GetIdentifiers(identifiersRequest config.IdentifiersRequest) ([]*gd.DeviceIdentifier, error) {
	log.Info().Msg("Handle GetIdentifiers Request")
	parameterJson := identifiersRequest.GetParameterJson()
	var deviceAddress simdevices.SimulatedDeviceAddress
	err := json.Unmarshal([]byte(parameterJson), &deviceAddress)
	if err != nil {
		log.Error().Err(err).Msg("Could not parse parameterJson")
		return nil, status.Errorf(codes.InvalidArgument, "Could not parse parameterJson: %v", err)
	}
	var deviceCredentials simdevices.SimulatedDeviceCredentials
	credentials := identifiersRequest.GetCredentials()
	for _, credential := range credentials {
		credMap := credential.GetCredentials()
		err = json.Unmarshal([]byte(credMap), &deviceCredentials)
		if err != nil {
			log.Error().Err(err).Msg("Could not parse credentials")
			return nil, status.Errorf(codes.InvalidArgument, "Could not parse credentials: %v", err)
		}
		deviceDetails, err := simdevices.RetrieveDeviceDetails(deviceAddress, &deviceCredentials)
		if err != nil {
			log.Warn().Err(err).Msgf("Could not retrieve device details with provided credentials for device at IP address %s and NIC %s",
				deviceAddress.DeviceIP, deviceAddress.AssetLinkNIC)
			continue // try next credentials
		}
		// if device can be reached with provided credentials, return its identifiers
		// otherwise try next credentials
		deviceInfo := createDeviceInfo(deviceDetails)
		discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()
		if discoveredDevice != nil {
			return discoveredDevice.GetIdentifiers(), nil
		}
	}
	// if device can not be reached with any provided credentials, try without credentials
	deviceDetails, err := simdevices.RetrieveDeviceDetails(deviceAddress, nil)
	if err != nil {
		log.Error().Err(err).Msgf("Could not retrieve device details for device at IP address %s and NIC %s", deviceAddress.DeviceIP,
			deviceAddress.AssetLinkNIC)
		return nil, status.Errorf(codes.Unavailable, "Could not retrieve device details: %v", err)
	}
	deviceInfo := createDeviceInfo(deviceDetails)
	discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()
	if discoveredDevice != nil {
		return discoveredDevice.GetIdentifiers(), nil
	}
	return nil, status.Errorf(codes.NotFound, "Could not retrieve device details for device at IP address %s and NIC %s", deviceAddress.DeviceIP,
		deviceAddress.AssetLinkNIC)
}

func createDiscoverError(sdError error, sdAddress simdevices.SimulatedDeviceAddress) *gd.DiscoverError {
	var resultCode codes.Code
	switch sdError {
	case simdevices.ErrInvalidInterface:
		resultCode = codes.InvalidArgument
	case simdevices.ErrDeviceNotFound:
		resultCode = codes.NotFound
	case simdevices.ErrSubDeviceNotFound:
		resultCode = codes.NotFound
	case simdevices.ErrUnauthenticated:
		resultCode = codes.Unauthenticated
	default:
		resultCode = codes.Unknown
	}
	sdAddressBytes, _ := json.Marshal(sdAddress)
	description := fmt.Sprintf("Error retrieving device details (%s)", sdError.Error())
	discoverError := &gd.DiscoverError{
		ResultCode:  int32(resultCode),
		Description: description,
		Source: &gd.DiscoverError_Device{
			Device: &gd.Destination{Target: &gd.Destination_ConnectionParameterSet{ConnectionParameterSet: &gd.ConnectionParameterSet{
				ParameterJson: string(sdAddressBytes),
			}}},
		},
	}
	return discoverError
}

func createDeviceInfo(device simdevices.SimulatedDeviceInfo) *model.DeviceInfo {
	deviceInfo := model.NewDevice("EthernetDevice", device.GetDeviceName())
	deviceInfo.AddNameplate(device.GetManufacturer(), device.GetIDLink(), device.GetArticleNumber(),
		device.GetProductDesignation(), device.GetHardwareVersion(), device.GetSerialNumber())

	nicID := deviceInfo.AddNic(device.GetDeviceNIC(), device.GetMacAddress())
	deviceInfo.AddIPv4(nicID, device.GetIpDevice(), device.GetIpNetmask(), device.GetIpRoute())

	deviceInfo.AddSoftware("Firmware", device.GetActiveFirmwareVersion(), true)
	deviceInfo.AddCapabilities("firmware_update", device.IsUpdateSupported())
	deviceInfo.AddDescription(device.GetProductDesignation())

	deviceAddress := device.GetDeviceAddress()
	deviceIdentifierBlob, _ := json.Marshal(deviceAddress)
	deviceInfo.AddMetadata(deviceIdentifierBlob) // device ID or device connection data used for artefact uploads/downloads

	return deviceInfo
}

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

	artefactType := artefactMetaData.GetArtefactType()
	deviceIdentifierBlob, err := artefactMetaData.GetDeviceIdentifierBlob()
	if err != nil {
		log.Err(err).Msg("Failed to retrieve device identifier blob")
		return err
	}

	log.Info().Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	// Perform checks
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_PREPARE, ga.ArtefactOperationState_AOS_OK, "Performing checks", 0)

	if artefactType != ga.ArtefactType_AT_FIRMWARE {
		err = errors.New("artefact type not supported")
		log.Err(err).Msg("Failed to handle push artefact")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_PREPARE, ga.ArtefactOperationState_AOS_OK, "Performed checks", 100)

	// Receiving new firmware
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_DOWNLOAD, ga.ArtefactOperationState_AOS_OK, "Receiving new firmware", 0)

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

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_DOWNLOAD, ga.ArtefactOperationState_AOS_OK, "Received new firmware", 100)

	// Verify new firmware, connect to device, and install new firmware
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Verifying new firmware", 0)

	var deviceAddress simdevices.SimulatedDeviceAddress
	err = json.Unmarshal(deviceIdentifierBlob, &deviceAddress)
	if err != nil {
		log.Err(err).Msg("Failed to parse connection blob")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Connecting to device", 10)

	device, err := simdevices.ConnectToDevice(deviceAddress, nil)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

	oldFirmwareVersion := device.GetActiveFirmwareVersion()

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Installing new firmware on device", 30)

	err = device.UpdateFirmware(artefactFilename)
	if err != nil {
		log.Err(err).Msg("Failed to update device firmware")
		return err
	}

	err = device.RebootDevice()
	if err != nil {
		log.Err(err).Msg("Failed to reboot device after firmware update")
		return err
	}

	newFirmwareVersion := device.GetActiveFirmwareVersion()

	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_INSTALLATION, ga.ArtefactOperationState_AOS_OK, "Installed new firmware on device", 100)

	// Activating new firmware
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_ACTIVATION, ga.ArtefactOperationState_AOS_OK, "Activating new firmware", 0)

	finalMessage := fmt.Sprintf("New firmware activated (new version %s, old version %s)", newFirmwareVersion, oldFirmwareVersion)
	_ = artefactReceiver.UpdateStatus(ga.ArtefactOperationPhase_AOP_ACTIVATION, ga.ArtefactOperationState_AOS_OK, finalMessage, 100)

	return nil
}

func (m *ReferenceAssetLink) HandlePullArtefact(artefactMetaData *artefact.ArtefactMetaData, artefactTransmitter *artefact.ArtefactTransmitter) error {
	log.Info().Msg("Handle Pull Artefact by transmitting the artefact")

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
	artefactType := artefactMetaData.GetArtefactType()
	deviceIdentifierBlob, err := artefactMetaData.GetDeviceIdentifierBlob()
	if err != nil {
		log.Err(err).Msg("Failed to retrieve device identifier blob")
		return err
	}

	log.Info().Str("DeviceIdentifierBlob", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

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

	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_PREPARE, ga.ArtefactOperationState_AOS_OK, "Performed checks", 100)

	// Connect to device and retrieve configuration from device
	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_ARCHIVE, ga.ArtefactOperationState_AOS_OK, "Connecting to device", 0)

	device, err := simdevices.ConnectToDevice(deviceAddress, nil)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_ARCHIVE, ga.ArtefactOperationState_AOS_OK, "Retrieving configuration from device", 10)

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

	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_ARCHIVE, ga.ArtefactOperationState_AOS_OK, "Retrieved configuration from device", 100)

	// Transmit configuration
	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_UPLOAD, ga.ArtefactOperationState_AOS_OK, "Transmitting configuration", 0)

	err = artefactTransmitter.TransmitArtefactFromFile(artefactFilename, 1024)
	if err != nil {
		log.Err(err).Msg("Failed to transmit artefact file")
		return err
	}

	time.Sleep(2 * time.Second)

	_ = artefactTransmitter.UpdateStatus(ga.ArtefactOperationPhase_AOP_UPLOAD, ga.ArtefactOperationState_AOS_OK, "Configuration transmission complete", 100)

	return nil
}
