/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package reference

import (
	"fmt"
	"testing"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/config"
	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/publish"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDiscovery(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	t.Run("discoverySucceeds", func(t *testing.T) {

		devicePublisher := &publish.DevicePublisherMock{}

		discoveryConfig := config.NewDiscoveryConfigWithDefaults()

		driver := &ReferenceAssetLink{}

		assert.NoError(t, driver.Discover(discoveryConfig, devicePublisher))
		assert.NotEmpty(t, devicePublisher.GetDevices())
	})

	t.Run("discoveryFails", func(t *testing.T) {
		devicePublisher := &publish.DevicePublisherMock{}

		err := status.Errorf(codes.Canceled, "Discovery was Canceled")
		devicePublisher.SetError(err)

		discoveryConfig := config.NewDiscoveryConfigWithDefaults()

		driver := &ReferenceAssetLink{}

		assert.Error(t, err, driver.Discover(discoveryConfig, devicePublisher))
		assert.Empty(t, devicePublisher.GetDevices())
	})

	t.Run("discoveryWithDeviceDetailsError", func(t *testing.T) {
		devicePublisher := &publish.DevicePublisherMock{}

		discoveryConfig := config.NewDiscoveryConfigWithDefaults()

		driver := &ReferenceAssetLink{}

		assert.NoError(t, driver.Discover(discoveryConfig, devicePublisher))

		existingErrors := devicePublisher.GetErrors()
		initialErrorCount := len(existingErrors)

		deviceDetailsError := &generated.DiscoverError{
			ResultCode:  int32(codes.Unavailable),
			Description: "Failed to retrieve device details for discovered device",
		}
		assert.NoError(t, devicePublisher.PublishError(deviceDetailsError))

		allErrors := devicePublisher.GetErrors()
		assert.NotEmpty(t, allErrors)
		assert.Equal(t, initialErrorCount+1, len(allErrors))

		foundTestError := false
		for _, err := range allErrors {
			if err.ResultCode == int32(codes.Unavailable) &&
				err.Description == "Failed to retrieve device details for discovered device" {
				foundTestError = true
				break
			}
		}
		assert.True(t, foundTestError, "Should find the test error we specifically published")
		assert.True(t, len(allErrors) > 0, "Should have device detail errors published")
	})
}

func TestConfig(t *testing.T) {
	t.Run("requestSupportedFilters", func(t *testing.T) {
		driver := &ReferenceAssetLink{}

		driver.GetSupportedFilters()
	})

	t.Run("requestSupportedOptions", func(t *testing.T) {
		driver := &ReferenceAssetLink{}

		driver.GetSupportedOptions()
	})
}
func TestGetPropertyValues(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	driver := &ReferenceAssetLink{}
	credential := `{"username":"user","password":"user_password"}`
	credentials := []*generated.ConnectionCredential{{Credentials: credential}}
	parameters := `{"alNic":"eth0","ipAddress":"192.168.0.12","subDeviceID":-1}`
	target := &generated.Destination{Target: &generated.Destination_ConnectionParameterSet{ConnectionParameterSet: &generated.ConnectionParameterSet{Credentials: credentials, ParameterJson: parameters}}}

	propertyResp, err := driver.GetPropertyValues(&generatedDeviceInfo.GetPropertyValuesRequest{Device: target})
	assert.NoError(t, err)
	assert.NotNil(t, propertyResp)
	assert.NotEmpty(t, propertyResp.GetPropertyResults())

	resultsByKey := make(map[string]*generated.Variant, len(propertyResp.GetPropertyResults()))
	for _, result := range propertyResp.GetPropertyResults() {
		property := result.GetProperty()
		if property == nil {
			continue
		}
		resultsByKey[property.GetKey()] = property.GetValue()
	}

	connectionPoints := resultsByKey["connection_points"]
	if assert.NotNil(t, connectionPoints) {
		arrayValue := connectionPoints.GetArrayValue()
		if assert.NotNil(t, arrayValue) && assert.NotEmpty(t, arrayValue.GetValues()) {
			assert.NotNil(t, arrayValue.GetValues()[0].GetStructValue())
		}
	}

	productInfo := resultsByKey["product_instance_information"]
	if assert.NotNil(t, productInfo) {
		assert.NotNil(t, productInfo.GetStructValue())
	}

	for _, result := range propertyResp.GetPropertyResults() {
		property := result.GetProperty()
		if property == nil {
			continue
		}
		assert.NotContains(t, property.GetKey(), ".", fmt.Sprintf("unexpected flattened property key %q", property.GetKey()))
		assert.NotContains(t, property.GetKey(), "[", fmt.Sprintf("unexpected flattened property key %q", property.GetKey()))
	}
}

func TestGetPropertyValuesWithAssetIdentifiers(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	driver := &ReferenceAssetLink{}
	identifierParameters := `{"asset_identifiers":[{"asset_identifier_type":"IdLinkIdentifier","id_link":"https://industrial-assets.io/?1P=AN0123456789&S=SN123450000"},{"asset_identifier_type":"MacIdentifier","mac_address":"00:16:3e:00:00:00"},{"asset_identifier_type":"SoftwareIdentifier","name":"Firmware","version":"1.0.0"}]}`
	target := &generated.Destination{Target: &generated.Destination_ConnectionParameterSet{ConnectionParameterSet: &generated.ConnectionParameterSet{ParameterJson: identifierParameters}}}

	propertyResp, err := driver.GetPropertyValues(&generatedDeviceInfo.GetPropertyValuesRequest{Device: target, Keys: []string{"name", "asset_identifiers"}})
	assert.NoError(t, err)
	assert.NotNil(t, propertyResp)
	assert.NotEmpty(t, propertyResp.GetPropertyResults())
}

func TestGetPropertyValuesIncludesMandatoryFunctionalFields(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	driver := &ReferenceAssetLink{}
	parameters := `{"alNic":"eth0","ipAddress":"192.168.0.10","subDeviceID":-1}`
	target := &generated.Destination{Target: &generated.Destination_ConnectionParameterSet{ConnectionParameterSet: &generated.ConnectionParameterSet{ParameterJson: parameters}}}

	propertyResp, err := driver.GetPropertyValues(&generatedDeviceInfo.GetPropertyValuesRequest{
		Device: target,
		Keys:   []string{"functional_object_type", "functional_object_schema_url"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, propertyResp)

	resultsByKey := make(map[string]*generated.Variant, len(propertyResp.GetPropertyResults()))
	for _, result := range propertyResp.GetPropertyResults() {
		property := result.GetProperty()
		if property == nil {
			continue
		}
		resultsByKey[property.GetKey()] = property.GetValue()
	}

	if assert.Contains(t, resultsByKey, "functional_object_type") {
		assert.NotEmpty(t, resultsByKey["functional_object_type"].GetText())
		assert.Equal(t, "Device", resultsByKey["functional_object_type"].GetText())
	}

	if assert.Contains(t, resultsByKey, "functional_object_schema_url") {
		assert.NotEmpty(t, resultsByKey["functional_object_schema_url"].GetText())
		assert.Equal(t, model.FunctionalObjectSchemaUrl, resultsByKey["functional_object_schema_url"].GetText())
	}
}

func TestGetSupportedProperties(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	driver := &ReferenceAssetLink{}
	credential := `{"username":"user","password":"user_password"}`
	credentials := []*generated.ConnectionCredential{{Credentials: credential}}
	parameters := `{"alNic":"eth0","ipAddress":"192.168.0.12","subDeviceID":-1}`
	target := &generated.Destination{Target: &generated.Destination_ConnectionParameterSet{ConnectionParameterSet: &generated.ConnectionParameterSet{Credentials: credentials, ParameterJson: parameters}}}

	propertyResp, err := driver.GetSupportedProperties(&generatedDeviceInfo.GetSupportedPropertiesRequest{Device: target})
	assert.NoError(t, err)
	assert.NotNil(t, propertyResp)
	assert.NotEmpty(t, propertyResp.GetProperties())

	propertiesByKey := make(map[string]*generatedDeviceInfo.SupportedProperty, len(propertyResp.GetProperties()))
	for _, property := range propertyResp.GetProperties() {
		propertiesByKey[property.GetKey()] = property
		assert.NotContains(t, property.GetKey(), ".")
		assert.NotContains(t, property.GetKey(), "[")
	}

	if assert.Contains(t, propertiesByKey, "connection_points") {
		assert.Equal(t, generated.VariantType_VT_ARRAY, propertiesByKey["connection_points"].GetDatatype())
	}
	if assert.Contains(t, propertiesByKey, "product_instance_information") {
		assert.Equal(t, generated.VariantType_VT_STRUCT, propertiesByKey["product_instance_information"].GetDatatype())
	}
}

func TestCreateDeviceInfoAssetRelations(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	t.Run("adds gateway relation for top-level device", func(t *testing.T) {
		address := simdevices.SimulatedDeviceAddress{
			AssetLinkNIC: "eth0",
			DeviceIP:     "192.168.0.10",
			SubDeviceID:  -1,
		}

		device, err := simdevices.RetrieveDeviceDetails(address, nil)
		assert.NoError(t, err)
		if !assert.NotNil(t, device) {
			return
		}

		deviceInfo, err := createDeviceInfo(device)
		assert.NoError(t, err)
		if !assert.NotNil(t, deviceInfo) {
			return
		}

	})

	t.Run("adds module relation for sub-device", func(t *testing.T) {
		address := simdevices.SimulatedDeviceAddress{
			AssetLinkNIC: "eth1",
			DeviceIP:     "192.168.1.10",
			SubDeviceID:  0,
		}

		device, err := simdevices.RetrieveDeviceDetails(address, nil)
		assert.NoError(t, err)
		if !assert.NotNil(t, device) {
			return
		}

		deviceInfo, err := createDeviceInfo(device)
		assert.NoError(t, err)
		if !assert.NotNil(t, deviceInfo) {
			return
		}

		if assert.Len(t, deviceInfo.AssetRelations, 1) {
			moduleRelation := deviceInfo.AssetRelations[0]
			assert.Equal(t, "is_module_of", moduleRelation.Predicate)
			assert.Equal(t, model.RelationalRoleOfRelatedAssetValuesObject, moduleRelation.RelationalRoleOfRelatedAsset)
			if assert.Len(t, moduleRelation.RelatedAsset.AssetIdentifiers, 1) {
				identifier, ok := moduleRelation.RelatedAsset.AssetIdentifiers[0].(model.MacIdentifier)
				if assert.True(t, ok) {
					assert.Equal(t, model.MacIdentifierAssetIdentifierTypeMacIdentifier, identifier.AssetIdentifierType)
					assert.Equal(t, "00:16:3e:01:00:00", identifier.MacAddress)
				}
			}
		}
	})
}

func TestGetSupportedPropertiesWithAssetIdentifiers(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	driver := &ReferenceAssetLink{}
	identifierParameters := `{"asset_identifiers":[{"asset_identifier_type":"MacIdentifier","mac_address":"00:16:3e:00:00:00"}]}`
	target := &generated.Destination{Target: &generated.Destination_ConnectionParameterSet{ConnectionParameterSet: &generated.ConnectionParameterSet{ParameterJson: identifierParameters}}}

	propertyResp, err := driver.GetSupportedProperties(&generatedDeviceInfo.GetSupportedPropertiesRequest{Device: target})
	assert.NoError(t, err)
	assert.NotNil(t, propertyResp)
	assert.NotEmpty(t, propertyResp.GetProperties())
}

func TestGetPropertyValuesWithSoftwareIdentifierOnly(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	driver := &ReferenceAssetLink{}
	identifierParameters := `{"asset_identifiers":[{"asset_identifier_type":"SoftwareIdentifier","name":"Firmware","version":"1.0.0"}]}`
	target := &generated.Destination{Target: &generated.Destination_ConnectionParameterSet{ConnectionParameterSet: &generated.ConnectionParameterSet{ParameterJson: identifierParameters}}}

	propertyResp, err := driver.GetPropertyValues(&generatedDeviceInfo.GetPropertyValuesRequest{Device: target, Keys: []string{"name", "software_components"}})
	assert.NoError(t, err)
	assert.NotNil(t, propertyResp)
	assert.NotEmpty(t, propertyResp.GetPropertyResults())
}

func TestGetPropertyValuesWithCustomIdentifierIP(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	driver := &ReferenceAssetLink{}
	// CustomIdentifier with name=ip_address lets callers identify a device by IP without knowing its NIC
	identifierParameters := `{"asset_identifiers":[{"asset_identifier_type":"CustomIdentifier","name":"ip_address","value":"192.168.0.10"}]}`
	target := &generated.Destination{Target: &generated.Destination_ConnectionParameterSet{ConnectionParameterSet: &generated.ConnectionParameterSet{ParameterJson: identifierParameters}}}

	propertyResp, err := driver.GetPropertyValues(&generatedDeviceInfo.GetPropertyValuesRequest{Device: target, Keys: []string{"name"}})
	assert.NoError(t, err)
	assert.NotNil(t, propertyResp)
	assert.NotEmpty(t, propertyResp.GetPropertyResults())
}

func TestGetPropertyValuesWithCustomIdentifierSerialNumber(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	driver := &ReferenceAssetLink{}
	identifierParameters := `{"asset_identifiers":[{"asset_identifier_type":"CustomIdentifier","name":"serial_number","value":"SN123450000"}]}`
	target := &generated.Destination{Target: &generated.Destination_ConnectionParameterSet{ConnectionParameterSet: &generated.ConnectionParameterSet{ParameterJson: identifierParameters}}}

	propertyResp, err := driver.GetPropertyValues(&generatedDeviceInfo.GetPropertyValuesRequest{Device: target, Keys: []string{"name"}})
	assert.NoError(t, err)
	assert.NotNil(t, propertyResp)
	assert.NotEmpty(t, propertyResp.GetPropertyResults())
}
