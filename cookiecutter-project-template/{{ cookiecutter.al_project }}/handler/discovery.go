/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gd "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
)

// Implements both the Discovery and the Upgrade interface/feature

type AssetLinkImplementation struct {
	driverLock sync.Mutex
}

func (m *AssetLinkImplementation) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error {
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

	//
	// Add your custom logic here to retrieve discovery config, discover devices, and publish them
	//

	optionSetting, optionErr := discoveryConfig.GetOptionSettingBool("option", false)
	if optionErr != nil {
		log.Error().Err(optionErr)
		return optionErr
	}

	filterSetting, filterErr := discoveryConfig.GetFilterSettingString("filter", "default")
	if filterErr != nil {
		log.Error().Err(filterErr)
		return filterErr
	}

	_ = optionSetting
	_ = filterSetting

	// Fillup the device information
	assetName := "Dummy Device 1"
	vendorName := "{{ cookiecutter.company }}"
	productName := "Dummy Product"
	orderNumber := "AN0123456789"
	serialNumber := "SN00012345678900001"
	hardwareVersion := "3"
	firmwareVersion := "1.0.0"
	deviceIdentifierBlob := []byte("dummy1")

	productUri := fmt.Sprintf("urn:%s/%s/%s", strings.ReplaceAll(vendorName, " ", "_"), strings.ReplaceAll(productName, " ", "_"), serialNumber)

	deviceInfo := model.NewDevice("EthernetDevice", assetName)
	deviceInfo.AddNameplate(vendorName, productUri, orderNumber, productName, hardwareVersion, serialNumber)
	deviceInfo.AddSoftware("Firmware", firmwareVersion, true)
	deviceInfo.AddCapabilities("firmware_update", false)
	deviceInfo.AddDescription("Dummy Device")

	deviceInfo.AddMetadata(deviceIdentifierBlob) // device ID or device connection data used for artefact uploads/downloads

	nicID := deviceInfo.AddNic("eth0", "00:16:3e:01:02:03") // random mac address
	deviceInfo.AddIPv4(nicID, "192.168.0.10", "255.255.255.0", "")
	deviceInfo.AddIPv4(nicID, "10.0.0.153", "255.255.255.0", "")

	// Convert and publish device
	discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()

	err := devicePublisher.PublishDevice(discoveredDevice)
	if err != nil {
		// discovery request was likely cancelled -> terminate discovery and return error
		log.Error().Msgf("Publishing Error: %v", err)
		return err
	}

	return nil
}

func (m *AssetLinkImplementation) GetSupportedOptions() []*gd.SupportedOption {
	supportedOptions := make([]*gd.SupportedOption, 0)
	supportedOptions = append(supportedOptions, &gd.SupportedOption{
		Key:      "option",
		Datatype: gd.VariantType_VT_BOOL,
	})
	return supportedOptions
}

func (m *AssetLinkImplementation) GetSupportedFilters() []*gd.SupportedFilter {
	supportedFilters := make([]*gd.SupportedFilter, 0)
	supportedFilters = append(supportedFilters, &gd.SupportedFilter{
		Key:      "filter",
		Datatype: gd.VariantType_VT_STRING,
	})
	return supportedFilters
}
