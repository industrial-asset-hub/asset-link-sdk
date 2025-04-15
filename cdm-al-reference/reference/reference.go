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

	"github.com/industrial-asset-hub/asset-link-sdk/v3/artefact"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generatedArtefact "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements the Discovery interface and feature

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
		deviceInfo.AddMetadata(string(deviceIdentifierBlob)) // device ID or device connection data used for artefact uploads/downloads

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

func (m *ReferenceAssetLink) GetSupportedOptions() []*generated.SupportedOption {
	supportedOptions := make([]*generated.SupportedOption, 0)
	supportedOptions = append(supportedOptions, &generated.SupportedOption{
		Key:      "interface",
		Datatype: generated.VariantType_VT_STRING,
	})
	return supportedOptions
}

func (m *ReferenceAssetLink) GetSupportedFilters() []*generated.SupportedFilter {
	supportedFilters := make([]*generated.SupportedFilter, 0)
	supportedFilters = append(supportedFilters, &generated.SupportedFilter{
		Key:      "ip_range",
		Datatype: generated.VariantType_VT_STRING,
	})
	return supportedFilters
}

func createDiscoverError(sdError error, sdAddress simdevices.SimulatedDeviceAddress) *generated.DiscoverError {
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
	_ = sdAddress //TODO: add device address to discoverError
	description := fmt.Sprintf("Error retrieving device details (%s)", sdError.Error())
	discoverError := &generated.DiscoverError{
		ResultCode:  int32(resultCode),
		Description: description,
	}
	return discoverError
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

	artefactMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact meta data")
		return err
	}

	deviceIdentifierBlob := artefactMetaData.GetDeviceIdentifier()
	artefactType := artefactMetaData.GetArtefactType()

	log.Info().Str("DeviceIdentifier", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	if artefactType != generatedArtefact.ArtefactType_AT_FIRMWARE {
		err = errors.New("artefact type not supported")
		log.Err(err).Msg("Failed to handle push artefact")
		return err
	}

	tempDir, err := os.MkdirTemp("", "artefact_pull")
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

	var deviceAddress simdevices.SimulatedDeviceAddress
	err = json.Unmarshal(deviceIdentifierBlob, &deviceAddress)
	if err != nil {
		log.Err(err).Msg("Failed to parse connection blob")
		return err
	}

	device, err := simdevices.ConnectToDevice(deviceAddress, nil)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

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

	deviceIdentifierBlob := artefactMetaData.GetDeviceIdentifier()
	artefactType := artefactMetaData.GetArtefactType()

	log.Info().Str("DeviceIdentifier", string(deviceIdentifierBlob)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	if artefactType != generatedArtefact.ArtefactType_AT_CONFIGURATION {
		err := errors.New("artefact type not supported")
		log.Err(err).Msg("Failed to handle pull artefact")
		return err
	}

	var deviceAddress simdevices.SimulatedDeviceAddress
	err := json.Unmarshal(deviceIdentifierBlob, &deviceAddress)
	if err != nil {
		log.Err(err).Msg("Failed to parse connection blob")
		return err
	}

	device, err := simdevices.ConnectToDevice(deviceAddress, nil)
	if err != nil {
		log.Err(err).Msg("Failed to connect to device")
		return err
	}

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

	err = artefactTransmitter.TransmitArtefactFromFile(artefactFilename, 1024)
	if err != nil {
		log.Err(err).Msg("Failed to transmit artefact file")
		return err
	}

	return nil
}
