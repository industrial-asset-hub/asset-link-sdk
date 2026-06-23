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

	"github.com/industrial-asset-hub/asset-link-sdk/v4/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/publish"
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

	// Example: Create two assets to demonstrate asset relations
	// First, create and publish a parent device
	parentDeviceMAC := "00:16:3e:01:02:04"

	if err := createParentDevice(devicePublisher); err != nil {
		return err
	}

	if err := createAndPublishChildDevice(devicePublisher, parentDeviceMAC); err != nil {
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

// logWarningIfEmpty logs a warning if err is ErrEmpty, otherwise logs a generic warning
func logWarningIfEmpty(err error, emptyMsg, genericMsg string) {
	if errors.Is(err, model.ErrEmpty) {
		log.Warn().Err(err).Msg(emptyMsg)
		return
	}

	log.Warn().Err(err).Msg(genericMsg)
}

// createParentDevice creates and returns a parent device with network configuration
func createParentDevice(devicePublisher publish.DevicePublisher) error {
	parentDeviceName := "Parent Device"
	parentDeviceMAC := "00:16:3e:01:02:04"

	parentDeviceInfo, err := model.NewDevice("Asset", parentDeviceName)
	if err != nil {
		logWarningIfEmpty(err, "one or more required fields for creating parent device are empty", "Could not create parent device info")
		return nil
	}

	parentNicID, err := parentDeviceInfo.AddNic("eth1", parentDeviceMAC)
	if err != nil {
		logWarningIfEmpty(err, "MAC address is empty for parent device", "Could not add NIC to parent device")
		if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("MAC address format is invalid for parent device")
		}
		return nil
	}

	if _, err = parentDeviceInfo.AddIPv4(parentNicID, "192.168.0.5", "255.255.255.0", "192.168.0.1"); err != nil {
		log.Warn().Err(err).Msg("Could not add IPv4 to parent device")
	}

	parentDiscoveredDevice := parentDeviceInfo.ConvertToDiscoveredDevice()
	if err = devicePublisher.PublishDevice(parentDiscoveredDevice); err != nil {
		log.Error().Msgf("Publishing Error for parent device: %v", err)
		return err
	}

	return nil
}

// createAndPublishChildDevice creates a child device with network and module relationships
func createAndPublishChildDevice(devicePublisher publish.DevicePublisher, parentDeviceMAC string) error {
	deviceInfo, err := model.NewDevice("Asset", "Dummy Device 1")
	if err != nil {
		logWarningIfEmpty(err, "one or more required fields for creating device info are empty, cannot create device info for discovered device", "Could not create device info")
		return nil
	}

	vendorName := "{{ cookiecutter.company }}"
	productFamily := "Dummy Product"
	orderNumber := "AN0123456789"
	serialNumber := "SN00012345678900001"
	hardwareVersion := "3"

	productUri := fmt.Sprintf(
		"%s/?1P=%s&S=%s",
		strings.TrimRight("{{ cookiecutter.company_url }}", "/"),
		url.QueryEscape(orderNumber),
		url.QueryEscape(serialNumber),
	)

	if err = deviceInfo.AddNameplate(vendorName, productUri, orderNumber, productFamily, hardwareVersion, serialNumber); err != nil {
		logWarningIfEmpty(err, "One or more nameplate fields are empty, cannot add nameplate information to device info", "Could not add nameplate information to device info")
		return nil
	}

	if err = deviceInfo.AddSoftwareArtifactComponent("Firmware", "1.0.0", true); err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Firmware version is empty, cannot add firmware information to device info")
		}
	}

	if err = deviceInfo.AddAssetRelation(
		"is_module_of",
		model.RelatedAsset{
			AssetIdentifiers: []interface{}{
				model.MacIdentifier{
					AssetIdentifierType: model.MacIdentifierAssetIdentifierTypeMacIdentifier,
					MacAddress:          parentDeviceMAC,
				},
			},
		},
		model.RelationalRoleOfRelatedAssetValuesObject,
		false,
	); err != nil {
		if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("Asset relation format is invalid, cannot add asset relation to device info")
		} else if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Asset relation identifier is empty, cannot add asset relation to device info")
		}
		log.Warn().Err(err).Msg("Could not add asset relation information to device info")
	}

	if err = deviceInfo.AddCapabilities("firmware_update", false); err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Capability name is empty, cannot add capability information to device info")
		}
	}

	if err = deviceInfo.AddDescription("Dummy Device"); err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("Description is empty, cannot add description to device info")
		}
	}

	nicID, err := deviceInfo.AddNic("eth0", "00:16:3e:01:02:03")
	if err != nil {
		logWarningIfEmpty(err, "MAC address is empty, cannot add NIC to device info", "Could not add NIC to device info")
		if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("MAC address format is invalid, cannot add NIC to device info")
		}
		return nil
	}

	if err = addIPv4ConnectionPoint(deviceInfo, nicID, "192.168.0.10", "255.255.255.0", "192.168.0.1"); err != nil {
		return nil
	}

	if err = addIPv4ConnectionPoint(deviceInfo, nicID, "10.0.0.153", "255.255.255.0", "10.0.0.1"); err != nil {
		return nil
	}

	discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()
	if err = devicePublisher.PublishDevice(discoveredDevice); err != nil {
		log.Error().Msgf("Publishing Error: %v", err)
		return err
	}

	return nil
}

// addIPv4ConnectionPoint adds IPv4 to a device with error handling
func addIPv4ConnectionPoint(deviceInfo *model.DeviceInfo, nicID string, ipAddr, netMask, routerIP string) error {
	_, err := deviceInfo.AddIPv4(nicID, ipAddr, netMask, routerIP)
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("IP address or network mask is empty, cannot add IPv4 connectivity to device info")
			return err
		}

		if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("IP address or network mask format is invalid, cannot add IPv4 connectivity to device info")
			return err
		}

		log.Warn().Err(err).Msg("Could not add IPv4 connectivity to device info")
	}
	return err
}
