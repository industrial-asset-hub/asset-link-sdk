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
	"sync"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements both the Discovery and the Identifiers interface/feature

type ReferenceAssetLink struct {
	discoveryLock sync.Mutex
}

func (m *ReferenceAssetLink) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error {
	log.Info().Msg("Handle Discovery Request")

	// Check if a job is already running
	// We currently support only one running job
	if m.discoveryLock.TryLock() {
		defer m.discoveryLock.Unlock()
	} else {
		const errMsg string = "Another discovery job is already running"
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
		deviceInfo, err := createDeviceInfo(device)
		if err != nil {
			log.Error().Err(err).Msg("Could not create device info")
			continue
		}
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

func (m *ReferenceAssetLink) GetIdentifiers(identifiersRequest config.IdentifiersRequest) ([]*generated.DeviceIdentifier, error) {
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
		deviceInfo, err := createDeviceInfo(deviceDetails)
		if err != nil {
			log.Error().Err(err).Msg("Could not create device info")
			return nil, status.Errorf(codes.Internal, "Could not create device info: %v", err)
		}
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
	deviceInfo, err := createDeviceInfo(deviceDetails)
	if err != nil {
		log.Error().Err(err).Msg("Could not create device info")
		return nil, status.Errorf(codes.Internal, "Could not create device info: %v", err)
	}
	discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()
	if discoveredDevice != nil {
		return discoveredDevice.GetIdentifiers(), nil
	}
	return nil, status.Errorf(codes.NotFound, "Could not retrieve device details for device at IP address %s and NIC %s", deviceAddress.DeviceIP,
		deviceAddress.AssetLinkNIC)
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
	sdAddressBytes, _ := json.Marshal(sdAddress)
	description := fmt.Sprintf("Error retrieving device details (%s)", sdError.Error())
	discoverError := &generated.DiscoverError{
		ResultCode:  int32(resultCode),
		Description: description,
		Source: &generated.DiscoverError_Device{
			Device: &generated.Destination{Target: &generated.Destination_ConnectionParameterSet{ConnectionParameterSet: &generated.ConnectionParameterSet{
				ParameterJson: string(sdAddressBytes),
			}}},
		},
	}
	return discoverError
}

func createDeviceInfo(device simdevices.SimulatedDeviceInfo) (*model.DeviceInfo, error) {
	deviceInfo, err := model.NewDevice("EthernetDevice", device.GetDeviceName())
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("one or more required fields for creating device info are empty, cannot create device info for discovered device")
			return deviceInfo, nil
		}
		log.Warn().Err(err).Msg("Could not create device info")
		return deviceInfo, nil
	}
	err = deviceInfo.AddNameplate(device.GetManufacturer(), device.GetIDLink(), device.GetArticleNumber(),
		device.GetProductDesignation(), device.GetHardwareVersion(), device.GetSerialNumber())
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("One or more nameplate fields are empty, cannot add nameplate information to device info")
			return deviceInfo, nil
		}
	}

	nicID, err := deviceInfo.AddNic(device.GetDeviceNIC(), device.GetMacAddress())
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("MAC address is empty, cannot add NIC to device info")
			return deviceInfo, nil
		} else if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("MAC address format is invalid, cannot add NIC to device info")
			return deviceInfo, nil // return device info without NIC -ask?
		}
		log.Warn().Err(err).Msg("Could not add NIC to device info")
		return deviceInfo, nil
	}
	_, err = deviceInfo.AddIPv4(nicID, device.GetIpDevice(), device.GetIpNetmask(), device.GetIpRoute())
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("IP address or network mask is empty, cannot add IPv4 connectivity to device info")
			return deviceInfo, nil // return device info without IPv4 connectivity
		} else if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("IP address or network mask format is invalid, cannot add IPv4 connectivity to device info")
			return deviceInfo, nil // return device info without IPv4 connectivity
		}
		log.Warn().Err(err).Msg("Could not add IPv4 connectivity to device info")
		return deviceInfo, nil
	}

	err = deviceInfo.AddSoftware("Firmware", device.GetActiveFirmwareVersion(), true)
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Firmware version is empty, cannot add firmware information to device info")
			return deviceInfo, nil
		}
	}
	err = deviceInfo.AddCapabilities("firmware_update", device.IsUpdateSupported())
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Capability name is empty, cannot add capability information to device info")
			return deviceInfo, nil
		}
	}
	err = deviceInfo.AddDescription(device.GetProductDesignation())
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Description is empty, cannot add description to device info")
		}
	}
	return deviceInfo, nil
}
