/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package reference

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
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

func createDeviceInfo(device simdevices.SimulatedDevice) *model.DeviceInfo {
	deviceInfo := model.NewDevice("EthernetDevice", device.GetDeviceName())
	deviceInfo.AddNameplate(device.GetManufacturer(), device.GetIDLink(), device.GetArticleNumber(),
		device.GetProductDesignation(), device.GetHardwareVersion(), device.GetSerialNumber())

	nicID := deviceInfo.AddNic(device.GetDeviceNIC(), device.GetMacAddress())
	deviceInfo.AddIPv4(nicID, device.GetIpDevice(), device.GetIpNetmask(), device.GetIpRoute())

	deviceInfo.AddSoftware("Firmware", device.GetActiveFirmwareVersion(), true)
	deviceInfo.AddCapabilities("firmware_update", device.IsUpdateSupported())
	deviceInfo.AddDescription(device.GetProductDesignation())
	return deviceInfo
}
