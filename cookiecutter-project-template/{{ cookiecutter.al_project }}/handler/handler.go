/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements the Discovery interface and feature

type AssetLinkImplementation struct {
	discoveryLock sync.Mutex
}

func (m *AssetLinkImplementation) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error {
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

	// If there are any device-specific errors, they can be communicated via devicePublisher.PublishError(discoverError)
	// discoverError can be created like this:
	// discoverError := &generated.DiscoverError{
	//			ResultCode:  int32(codes.Unavailable),
	//			Description: "Error retrieving device details",
	//		}

	// Fillup the device information
	assetName := "Dummy Device 1"
	vendorName := "{{ cookiecutter.company }}"
	productName := "Dummy Product"
	orderNumber := "AN0123456789"
	serialNumber := "SN00012345678900001"
	hardwareVersion := "3"
	firmwareVersion := "1.0.0"

	productUri := fmt.Sprintf(
		"%s/?1P=%s&S=%s",
		strings.TrimRight("{{ cookiecutter.company_url }}", "/"),
		url.QueryEscape(orderNumber),
		url.QueryEscape(serialNumber),
	)

	deviceInfo, err := model.NewDevice("Asset", assetName)
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("one or more required fields for creating device info are empty, cannot create device info for discovered device")
			return nil
		}
		log.Warn().Err(err).Msg("Could not create device info")
		return nil
	}
	if err = deviceInfo.AddNameplate(vendorName, productUri, orderNumber, productName, hardwareVersion, serialNumber); err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("One or more nameplate fields are empty, cannot add nameplate information to device info")
			return nil
		}
		log.Warn().Err(err).Msg("Could not add nameplate information to device info")
		return nil
	}
	if err = deviceInfo.AddSoftware("Firmware", firmwareVersion, true); err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Firmware version is empty, cannot add firmware information to device info")
			return nil
		}
	}
	if err = deviceInfo.AddCapabilities("firmware_update", false); err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Capability name is empty, cannot add capability information to device info")
			return nil
		}
	}
	if err = deviceInfo.AddDescription("Dummy Device"); err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Description is empty, cannot add description to device info")
			return nil
		}
	}

	nicID, err := deviceInfo.AddNic("eth0", "00:16:3e:01:02:03") // random mac address
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("MAC address is empty, cannot add NIC to device info")
			return nil
		} else if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("MAC address format is invalid, cannot add NIC to device info")
			return nil // return device info without NIC -ask?
		}
		log.Warn().Err(err).Msg("Could not add NIC to device info")
		return nil
	}
	if _, err = deviceInfo.AddIPv4(nicID, "192.168.0.10", "255.255.255.0", "192.168.0.1"); err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("IP address or network mask is empty, cannot add IPv4 connectivity to device info")
			return nil // return device info without IPv4 connectivity
		} else if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("IP address or network mask format is invalid, cannot add IPv4 connectivity to device info")
			return nil // return device info without IPv4 connectivity
		}
		log.Warn().Err(err).Msg("Could not add IPv4 connectivity to device info")
		return nil
	}
	if _, err = deviceInfo.AddIPv4(nicID, "10.0.0.153", "255.255.255.0", "10.0.0.1"); err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("IP address or network mask is empty, cannot add IPv4 connectivity to device info")
			return nil // return device info without IPv4 connectivity
		} else if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("IP address or network mask format is invalid, cannot add IPv4 connectivity to device info")
			return nil // return device info without IPv4 connectivity
		}
		log.Warn().Err(err).Msg("Could not add IPv4 connectivity to device info")
		return nil
	}

	// Convert and publish device
	discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()

	err = devicePublisher.PublishDevice(discoveredDevice)
	if err != nil {
		// discovery request was likely cancelled -> terminate discovery and return error
		log.Error().Msgf("Publishing Error: %v", err)
		return err
	}

	return nil
}

func (m *AssetLinkImplementation) GetSupportedOptions() []*generated.SupportedOption {
	supportedOptions := make([]*generated.SupportedOption, 0)
	supportedOptions = append(supportedOptions, &generated.SupportedOption{
		Key:      "option",
		Datatype: generated.VariantType_VT_BOOL,
	})
	return supportedOptions
}

func (m *AssetLinkImplementation) GetSupportedFilters() []*generated.SupportedFilter {
	supportedFilters := make([]*generated.SupportedFilter, 0)
	supportedFilters = append(supportedFilters, &generated.SupportedFilter{
		Key:      "filter",
		Datatype: generated.VariantType_VT_STRING,
	})
	return supportedFilters
}

// if the GetIdentifiers() interface is implemented by your asset link,
// uncomment the Identifiers line in main.go to register the interface/feature and
// add "siemens.common.identifiers.v1" to the app_types in the registry.json file
func (m *AssetLinkImplementation) GetIdentifiers(identifiersRequest config.IdentifiersRequest) ([]*generated.DeviceIdentifier, error) {
	log.Info().Msg("Handle Get Identifiers Request")
	// Add your custom logic here to retrieve identifiers based on the provided parameters and credentials
	identifiers := []*generated.DeviceIdentifier{}
	return identifiers, nil
}
