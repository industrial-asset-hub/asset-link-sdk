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
	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
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

// templateSupportedProperties is treated as immutable.
var templateSupportedProperties = []*generatedDeviceInfo.SupportedProperty{
	{Key: "name", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_STRING}},
	{Key: "description", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_STRING}},
	{Key: "functional_object_type", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_STRING}},
	{Key: "functional_object_schema_url", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_STRING}},
	{Key: "asset_identifiers", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_ARRAY}},
	{Key: "asset_relations", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_ARRAY}},
	{Key: "connection_points", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_ARRAY}},
	{Key: "software_components", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_ARRAY}},
	{Key: "asset_operations", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_ARRAY}},
	{Key: "product_instance_information", Type: &generatedDeviceInfo.SupportedProperty_Datatype{Datatype: generated.VariantType_VT_STRUCT}},
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

	// Minimal example: publish one device with an asset identifier (MAC) only.
	if err := createAndPublishMinimalDevice(devicePublisher); err != nil {
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

// if your asset link implements the DeviceInfo interface,
// uncomment the DeviceInfo line in main.go to register the interface/feature and
// add "siemens.connectivitysuite.deviceinfo.v1" to the app_types in the registry.json file
func (m *AssetLinkImplementation) GetPropertyValues(request *generatedDeviceInfo.GetPropertyValuesRequest) (*generatedDeviceInfo.GetPropertyValuesResponse, error) {
	log.Info().Msg("Handle GetPropertyValues Request")
	_ = request

	deviceInfo, err := buildTemplateDeviceInfo()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not build device info: %v", err)
	}

	propertyResults, err := deviceInfo.ConvertToPropertyValueResults()
	if err != nil {
		return nil, err
	}

	return &generatedDeviceInfo.GetPropertyValuesResponse{PropertyResults: propertyResults}, nil
}

func (m *AssetLinkImplementation) GetSupportedProperties(_ *generatedDeviceInfo.GetSupportedPropertiesRequest) (*generatedDeviceInfo.GetSupportedPropertiesResponse, error) {
	log.Info().Msg("Handle GetSupportedProperties Request")

	return &generatedDeviceInfo.GetSupportedPropertiesResponse{Properties: templateSupportedProperties}, nil
}

func buildTemplateDeviceInfo() (*model.DeviceInfo, error) {
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

	deviceInfo, err := model.NewDevice("Asset", "Dummy Device 1")
	if err != nil {
		return nil, err
	}

	if err = deviceInfo.AddNameplate(vendorName, productUri, orderNumber, productFamily, hardwareVersion, serialNumber); err != nil {
		return nil, err
	}

	if err = deviceInfo.AddSoftwareArtifactComponent("Firmware", "1.0.0", true); err != nil {
		return nil, err
	}

	if err = deviceInfo.AddCapabilities("firmware_update", false); err != nil {
		return nil, err
	}

	if err = deviceInfo.AddDescription("Dummy Device"); err != nil {
		return nil, err
	}

	nicID, err := deviceInfo.AddNic("eth0", "00:16:3e:01:02:03")
	if err != nil {
		return nil, err
	}

	if _, err = deviceInfo.AddIPv4(nicID, "192.168.0.10", "255.255.255.0", "192.168.0.1"); err != nil {
		return nil, err
	}

	if _, err = deviceInfo.AddIPv4(nicID, "10.0.0.153", "255.255.255.0", "10.0.0.1"); err != nil {
		return nil, err
	}

	return deviceInfo, nil
}

// logWarningIfEmpty logs a warning if err is ErrEmpty, otherwise logs a generic warning
func logWarningIfEmpty(err error, emptyMsg, genericMsg string) {
	if errors.Is(err, model.ErrEmpty) {
		log.Warn().Err(err).Msg(emptyMsg)
		return
	}

	log.Warn().Err(err).Msg(genericMsg)
}

// createAndPublishMinimalDevice creates one device with just an asset identifier.
func createAndPublishMinimalDevice(devicePublisher publish.DevicePublisher) error {
	deviceInfo, err := model.NewDevice("Asset", "Dummy Device 1")
	if err != nil {
		logWarningIfEmpty(err, "one or more required fields for creating device info are empty", "Could not create device info")
		return err
	}

	if _, err = deviceInfo.AddNic("eth0", "00:16:3e:01:02:03"); err != nil {
		logWarningIfEmpty(err, "MAC address is empty, cannot add NIC to device info", "Could not add NIC to device info")
		if errors.Is(err, model.ErrValidation) {
			log.Warn().Err(err).Msg("MAC address format is invalid, cannot add NIC to device info")
		}
		return err
	}

	if err = devicePublisher.PublishDevice(deviceInfo.ConvertToDiscoveredDevice()); err != nil {
		log.Error().Msgf("Publishing Error: %v", err)
		return err
	}

	return nil
}
