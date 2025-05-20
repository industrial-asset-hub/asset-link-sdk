/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package reference

import (
	"encoding/json"
	"sync"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gd "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
)

// Implements both the Discovery and the Upgrade interface/feature

type ReferenceAssetLink struct {
	driverLock sync.Mutex
}

type DeviceIdentifierOrConnectionData struct {
	AssetLinkNIC string `json:"al_nic"`
	DeviceIP     string `json:"device_ip"`
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

	devicesFound := simdevices.ScanDevices(alInterface, ipRange)

	for _, device := range devicesFound {
		deviceInfo := model.NewDevice("EthernetDevice", device.GetDeviceName())
		deviceInfo.AddNameplate(device.GetManufacturer(), device.GetIDLink(), device.GetArticleNumber(),
			device.GetProductDesignation(), device.GetHardwareVersion(), device.GetSerialNumber())

		nicID := deviceInfo.AddNic(device.GetDeviceNIC(), device.GetMacAddress())
		deviceInfo.AddIPv4(nicID, device.GetIpDevice(), device.GetIpNetmask(), device.GetIpRoute())

		deviceInfo.AddSoftware("Firmware", device.GetActiveFirmwareVersion(), true)
		deviceInfo.AddCapabilities("firmware_update", device.IsUpdateSupported())

		deviceIdentifier := DeviceIdentifierOrConnectionData{
			AssetLinkNIC: device.GetAssetLinkNIC(),
			DeviceIP:     device.GetIpDevice(),
		}
		deviceIdentifierBlob, _ := json.Marshal(deviceIdentifier)
		deviceInfo.AddMetadata(deviceIdentifierBlob) // device ID or device connection data used for artefact uploads/downloads

		discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()

		err := devicePublisher.PublishDevice(discoveredDevice)
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
