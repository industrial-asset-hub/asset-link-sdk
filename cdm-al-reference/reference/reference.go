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
	"strings"
	"sync"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/config"
	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements both the Discovery and DeviceInfo interface/feature

type ReferenceAssetLink struct {
	discoveryLock sync.Mutex
}

type parameterQuery struct {
	AssetLinkNIC     string            `json:"alNic"`
	DeviceIP         string            `json:"ipAddress"`
	SubDeviceID      int               `json:"subDeviceID"`
	AssetIdentifiers []assetIdentifier `json:"asset_identifiers"`
}

type assetIdentifier struct {
	AssetIdentifierType string `json:"asset_identifier_type"`
	IdLink              string `json:"id_link"`
	MacAddress          string `json:"mac_address"`
	Name                string `json:"name"`
	Version             string `json:"version"`
	Value               string `json:"value"`
}

const (
	identifierTypeIDLink   = "idlinkidentifier"
	identifierTypeMAC      = "macidentifier"
	identifierTypeSoftware = "softwareidentifier"
	identifierTypeCustom   = "customidentifier"

	customIdentifierIPAddress    = "ip_address"
	customIdentifierSerialNumber = "serial_number"
	customIdentifierArticle      = "article_number"
	customIdentifierDeviceName   = "device_name"
)

// referenceSupportedProperties is treated as immutable.
var referenceSupportedProperties = []*generatedDeviceInfo.SupportedProperty{
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

func (m *ReferenceAssetLink) GetPropertyValues(request *generatedDeviceInfo.GetPropertyValuesRequest) (*generatedDeviceInfo.GetPropertyValuesResponse, error) {
	log.Info().Msg("Handle GetPropertyValues Request")
	deviceInfo, err := m.loadDeviceInfo(request.GetDevice())
	if err != nil {
		return nil, err
	}
	propertyResults, err := deviceInfo.ConvertToPropertyValueResults()
	if err != nil {
		return nil, err
	}
	return &generatedDeviceInfo.GetPropertyValuesResponse{PropertyResults: propertyResults}, nil
}

func (m *ReferenceAssetLink) GetSupportedProperties(_ *generatedDeviceInfo.GetSupportedPropertiesRequest) (*generatedDeviceInfo.GetSupportedPropertiesResponse, error) {
	log.Info().Msg("Handle GetSupportedProperties Request")

	return &generatedDeviceInfo.GetSupportedPropertiesResponse{Properties: referenceSupportedProperties}, nil
}

func (m *ReferenceAssetLink) loadDeviceInfo(device *generated.Destination) (*model.DeviceInfo, error) {
	if device == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing device target")
	}

	connectionParameterSet := device.GetConnectionParameterSet()
	if connectionParameterSet == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing connection parameter set")
	}

	credentials, err := parseCredentials(connectionParameterSet.GetCredentials())
	if err != nil {
		log.Error().Err(err).Msg("Could not parse credentials")
		return nil, status.Errorf(codes.InvalidArgument, "Could not parse credentials: %v", err)
	}

	deviceDetails, deviceAddress, err := resolveDeviceDetails(connectionParameterSet.GetParameterJson(), credentials)
	if err != nil {
		log.Error().Err(err).Msgf("Could not retrieve device details for device at IP address %s and NIC %s", deviceAddress.DeviceIP,
			deviceAddress.AssetLinkNIC)
		return nil, status.Errorf(codes.Unavailable, "Could not retrieve device details: %v", err)
	}

	return wrapDeviceInfo(deviceDetails)
}

func parseCredentials(credentials []*generated.ConnectionCredential) ([]*simdevices.SimulatedDeviceCredentials, error) {
	parsedCredentials := make([]*simdevices.SimulatedDeviceCredentials, 0, len(credentials))

	for _, credential := range credentials {
		parsedCredential := &simdevices.SimulatedDeviceCredentials{}
		if err := json.Unmarshal([]byte(credential.GetCredentials()), parsedCredential); err != nil {
			return nil, err
		}

		if parsedCredential.Username == "" && parsedCredential.Password == "" {
			continue
		}

		parsedCredentials = append(parsedCredentials, parsedCredential)
	}

	return parsedCredentials, nil
}

func resolveDeviceDetails(parameterJSON string, credentials []*simdevices.SimulatedDeviceCredentials) (simdevices.SimulatedDeviceInfo, simdevices.SimulatedDeviceAddress, error) {
	query := parameterQuery{}
	if err := json.Unmarshal([]byte(parameterJSON), &query); err != nil {
		return nil, simdevices.SimulatedDeviceAddress{}, status.Errorf(codes.InvalidArgument, "Could not parse parameterJson: %v", err)
	}

	deviceAddress := simdevices.SimulatedDeviceAddress{
		AssetLinkNIC: query.AssetLinkNIC,
		DeviceIP:     query.DeviceIP,
		SubDeviceID:  query.SubDeviceID,
	}

	if deviceAddress.AssetLinkNIC != "" && deviceAddress.DeviceIP != "" {
		deviceDetails, err := retrieveDeviceDetailsWithCredentials(deviceAddress, credentials)
		return deviceDetails, deviceAddress, err
	}

	if len(query.AssetIdentifiers) == 0 {
		return nil, simdevices.SimulatedDeviceAddress{}, status.Errorf(codes.InvalidArgument, "parameterJson must contain either alNic/ipAddress or asset_identifiers")
	}

	for _, address := range simdevices.ScanForDevices("", "") {
		if details, err := retrieveDeviceDetailsWithCredentials(address, credentials); err == nil {
			if matchesAssetIdentifiers(details, query.AssetIdentifiers) {
				return details, address, nil
			}
		}
	}

	return nil, simdevices.SimulatedDeviceAddress{}, simdevices.ErrDeviceNotFound
}

func retrieveDeviceDetailsWithCredentials(deviceAddress simdevices.SimulatedDeviceAddress, credentials []*simdevices.SimulatedDeviceCredentials) (simdevices.SimulatedDeviceInfo, error) {
	for _, credential := range credentials {
		deviceDetails, err := simdevices.RetrieveDeviceDetails(deviceAddress, credential)
		if err != nil {
			log.Warn().Err(err).Msgf("Could not retrieve device details with provided credentials for device at IP address %s and NIC %s",
				deviceAddress.DeviceIP, deviceAddress.AssetLinkNIC)
			continue
		}

		return deviceDetails, nil
	}

	return simdevices.RetrieveDeviceDetails(deviceAddress, nil)
}

func matchesAssetIdentifiers(device simdevices.SimulatedDeviceInfo, identifiers []assetIdentifier) bool {
	recognizedIdentifiers := 0

	for _, identifier := range identifiers {
		switch normalizeIdentifierType(identifier.AssetIdentifierType) {
		case identifierTypeIDLink:
			if identifier.IdLink == "" {
				continue
			}
			recognizedIdentifiers++
			if !strings.EqualFold(device.GetIDLink(), identifier.IdLink) {
				return false
			}
		case identifierTypeMAC:
			if identifier.MacAddress == "" {
				continue
			}
			recognizedIdentifiers++
			if !strings.EqualFold(device.GetMacAddress(), identifier.MacAddress) {
				return false
			}
		case identifierTypeSoftware:
			if identifier.Name == "" || identifier.Version == "" {
				continue
			}
			recognizedIdentifiers++
			if !strings.EqualFold(identifier.Name, "Firmware") {
				return false
			}
			if !strings.EqualFold(device.GetActiveFirmwareVersion(), identifier.Version) {
				return false
			}
		case identifierTypeCustom:
			if identifier.Name == "" || identifier.Value == "" {
				continue
			}
			recognizedIdentifiers++
			if !matchesCustomIdentifier(device, identifier.Name, identifier.Value) {
				return false
			}
		}
	}

	return recognizedIdentifiers > 0
}

// matchesCustomIdentifier maps well-known custom identifier names to the
// corresponding simulated-device property. Callers can use any of these names
// to identify a device without knowing its internal NIC/IP address:
//
//	"ip_address"     - device IP address
//	"serial_number"  - device serial number
//	"article_number" - device article number
//	"device_name"    - device name
func matchesCustomIdentifier(device simdevices.SimulatedDeviceInfo, name, value string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case customIdentifierIPAddress:
		return strings.EqualFold(device.GetIpDevice(), value)
	case customIdentifierSerialNumber:
		return strings.EqualFold(device.GetSerialNumber(), value)
	case customIdentifierArticle:
		return strings.EqualFold(device.GetArticleNumber(), value)
	case customIdentifierDeviceName:
		return strings.EqualFold(device.GetDeviceName(), value)
	}
	return false
}

func normalizeIdentifierType(identifierType string) string {
	return strings.ToLower(strings.TrimSpace(identifierType))
}

// wrapDeviceInfo calls createDeviceInfo and wraps any error with an Internal gRPC status.
func wrapDeviceInfo(deviceDetails simdevices.SimulatedDeviceInfo) (*model.DeviceInfo, error) {
	deviceInfo, err := createDeviceInfo(deviceDetails)
	if err != nil {
		log.Error().Err(err).Msg("Could not create device info")
		return nil, status.Errorf(codes.Internal, "Could not create device info: %v", err)
	}
	return deviceInfo, nil
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
	deviceInfo, err := model.NewDevice("Device", device.GetDeviceName())
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("one or more required fields for creating device info are empty, cannot create device info for discovered device")
			return deviceInfo, nil
		}
		log.Warn().Err(err).Msg("Could not create device info")
		return deviceInfo, nil
	}

	// Keep mandatory base fields non-empty for GetPropertyValues/GetSupportedProperties.
	functionalObjectType, typeOk := deviceInfo.FunctionalObjectType.(string)
	if !typeOk || strings.TrimSpace(functionalObjectType) == "" {
		deviceInfo.FunctionalObjectType = string(model.DeviceFunctionalObjectTypeDevice)
	}

	functionalObjectSchemaURL, schemaURLOk := deviceInfo.FunctionalObjectSchemaUrl.(string)
	if !schemaURLOk || strings.TrimSpace(functionalObjectSchemaURL) == "" {
		deviceInfo.FunctionalObjectSchemaUrl = model.FunctionalObjectSchemaUrl
	}

	err = deviceInfo.AddNameplate(device.GetManufacturer(), device.GetIDLink(), device.GetArticleNumber(),
		device.GetProductDesignation(), device.GetHardwareVersion(), device.GetSerialNumber())
	if err != nil {
		if errors.Is(err, model.ErrEmpty) {
			log.Warn().Err(err).Msg("One or more nameplate fields are empty, cannot add nameplate information to device info")
			return deviceInfo, nil
		}
	}
	// Add asset relations for meaningful architectural relationships
	// For sub-devices, establish the module relationship to the parent device
	if device.GetSubDeviceID() >= 0 {
		parentDevice := device.GetParentDevice()
		err := deviceInfo.AddAssetRelation(
			"is_module_of",
			model.RelatedAsset{
				AssetIdentifiers: []interface{}{
					model.MacIdentifier{
						AssetIdentifierType: model.MacIdentifierAssetIdentifierTypeMacIdentifier,
						MacAddress:          parentDevice.GetMacAddress(),
					},
				},
			},
			model.RelationalRoleOfRelatedAssetValuesObject,
			false,
		)
		if err != nil {
			if errors.Is(err, model.ErrValidation) {
				log.Warn().Err(err).Msg("Asset relation format is invalid, cannot add asset relation to device info")
				return deviceInfo, nil
			} else if errors.Is(err, model.ErrEmpty) {
				log.Warn().Err(err).Msg("Asset relation identifier is empty, cannot add asset relation to device info")
				return deviceInfo, nil
			}
			log.Warn().Err(err).Msg("Cannot add asset relation to device info")
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

	err = deviceInfo.AddSoftwareArtifactComponent("Firmware", device.GetActiveFirmwareVersion(), true)
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
