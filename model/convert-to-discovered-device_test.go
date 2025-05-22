/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	iah_discovery "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/stretchr/testify/assert"
)

func TestConvertToDiscoveredDevice(t *testing.T) {
	device := generateDevice("Profinet", "Device")
	discoveredDevice := device.ConvertToDiscoveredDevice()
	discoveredDeviceType := fmt.Sprintf("%s/%s", baseSchemaPrefix, "Asset#@type")
	assert.Equal(t, 22, len(discoveredDevice.Identifiers))
	assert.Equal(t, "URI", discoveredDevice.Identifiers[0].Classifiers[0].GetType())
	assert.Equal(t, discoveredDeviceType, discoveredDevice.Identifiers[0].Classifiers[0].GetValue())
}

func TestConvertFromDerivedSchemaToDiscoveredDevice(t *testing.T) {
	schemaUri := "https://schema.industrial-assets.io/sat/v0.8.2"
	device := generateDevice("SatController", "Device")
	discoveredDevice := ConvertFromDerivedSchemaToDiscoveredDevice(device, schemaUri, "SatController")
	assert.Equal(t, 22, len(discoveredDevice.Identifiers))
	assert.Equal(t, "URI", discoveredDevice.Identifiers[0].Classifiers[0].GetType())
	assert.Equal(t, "https://schema.industrial-assets.io/sat/v0.8.2/SatController#@type", discoveredDevice.Identifiers[0].Classifiers[0].GetValue())
}

type DerivedDeviceInfo struct {
	DeviceInfo
	PasswordProtected *bool `json:"password_protected,omitempty"`
}

func TestConvertDerivedSchemaToDiscoveredDevice(t *testing.T) {
	var satDevice *DerivedDeviceInfo
	device := generateDevice("SatController", "Device")
	satDevice = &DerivedDeviceInfo{
		DeviceInfo:        *device,
		PasswordProtected: new(bool),
	}
	*satDevice.PasswordProtected = true

	discoveredDevice := ConvertFromDerivedSchemaToDiscoveredDevice(satDevice, "https://schema.industrial-assets.io/sat/v0.8.2", "SatController")
	assert.Equal(t, 23, len(discoveredDevice.Identifiers))
	passwordProtectedFound := false
	for _, identifier := range discoveredDevice.Identifiers {
		if strings.Contains(identifier.Classifiers[0].GetValue(), "password_protected") {
			passwordProtectedFound = true
			assert.Equal(t, "true", identifier.GetText())
		}
	}
	assert.True(t, passwordProtectedFound)
}

func generateDevice(typeOfAsset string, assetName string) *DeviceInfo {
	device := NewDevice(typeOfAsset, assetName)
	timestamp := device.getAssetCreationTimestamp()
	Name := "Device"
	device.Name = &Name
	product := "test-product"
	version := "1.0.0"
	vendorName := "test-vendor"
	serialNumber := "test"
	vendor := Organization{
		Address:        nil,
		AlternateNames: nil,
		ContactPoint:   nil,
		Id:             "",
		Name:           &vendorName,
	}
	productSerialidentifier := ProductSerialIdentifier{
		IdentifierType:        nil,
		IdentifierUncertainty: nil,
		ManufacturerProduct: &Product{
			Id:             uuid.New().String(),
			Manufacturer:   &vendor,
			Name:           &Name,
			ProductId:      &product,
			ProductVersion: &version,
		},
		SerialNumber: &serialNumber,
	}
	device.ProductInstanceIdentifier = &productSerialidentifier

	randomMacAddress := "12:12:12:12:12:12"
	identifierUncertainty := 1
	device.MacIdentifiers = append(device.MacIdentifiers, MacIdentifier{
		MacAddress:            &randomMacAddress,
		IdentifierUncertainty: &identifierUncertainty,
	})

	connectionPointType := Ipv4ConnectivityConnectionPointTypeIpv4Connectivity
	Ipv4Address := "192.168.0.1"
	Ipv4NetMask := "255.255.255.0"
	connectionPoint := "EthernetPort"
	connectionPointTypeIpv6 := Ipv6ConnectivityConnectionPointTypeIpv6Connectivity
	routerIpv6Address := "fd12:3456:789a::1"
	Ipv6Address := "fd12:3456:789a::1"
	conPoint := "eth0"
	relatedConnectionPoint := RelatedConnectionPoint{
		ConnectionPoint:    &conPoint,
		CustomRelationship: &connectionPoint,
	}
	relatedConnectionPoints := make([]RelatedConnectionPoint, 0)
	relatedConnectionPoints = append(relatedConnectionPoints, relatedConnectionPoint)
	Ipv4Connectivity := Ipv4Connectivity{
		ConnectionPointType:     &connectionPointType,
		Id:                      "1",
		InstanceAnnotations:     nil,
		Ipv4Address:             &Ipv4Address,
		NetworkMask:             &Ipv4NetMask,
		RelatedConnectionPoints: relatedConnectionPoints,
		RouterIpv4Address:       nil,
	}
	device.ConnectionPoints = append(device.ConnectionPoints, Ipv4Connectivity)
	Ipv6Connectivity := Ipv6Connectivity{
		ConnectionPointType:     &connectionPointTypeIpv6,
		Id:                      "2",
		InstanceAnnotations:     nil,
		Ipv6Address:             &Ipv6Address,
		RelatedConnectionPoints: nil,
		RouterIpv6Address:       &routerIpv6Address,
	}
	device.ConnectionPoints = append(device.ConnectionPoints, Ipv6Connectivity)
	ethernetType := EthernetPortConnectionPointTypeEthernetPort
	EthernetPort := EthernetPort{
		Id:                  "3",
		ConnectionPointType: &ethernetType,
		MacAddress:          &randomMacAddress,
	}
	device.ConnectionPoints = append(device.ConnectionPoints, EthernetPort)

	state := ManagementStateValuesUnknown
	State := ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}
	device.ManagementState = State

	reachabilityStateValue := ReachabilityStateValuesReached
	reachabilityState := ReachabilityState{
		StateTimestamp: &timestamp,
		StateValue:     &reachabilityStateValue,
	}
	device.ReachabilityState = &reachabilityState
	return device
}

func checkForIdentifierUncertainty(t *testing.T, identifiers []*iah_discovery.DeviceIdentifier) bool {
	for _, identifier := range identifiers {
		for _, classifier := range identifier.GetClassifiers() {
			if strings.Contains(classifier.Value, "identifier_uncertainty") {
				// converting 1 to int64
				assert.Equal(t, int64(1), identifier.GetInt64Value())
				return true
			}
		}

		children := identifier.GetChildren()
		if children != nil {
			childrenIdentifiers := children.GetValue()
			if childrenIdentifiers != nil {
				if checkForIdentifierUncertainty(t, childrenIdentifiers) {
					return true
				}
			}
		}
	}
	return false
}

func TestConvertNumberTypeToDiscoveredDevice(t *testing.T) {
	device := NewDevice("Profinet", "Device")
	device.addIdentifier("ffeffawfafwfw")

	discoveredDevice := device.ConvertToDiscoveredDevice()
	result := checkForIdentifierUncertainty(t, discoveredDevice.GetIdentifiers())

	if !result {
		assert.Fail(t, "identifier_uncertainty not found in identifiers")
	}
}

// Struct containing all Go data types
type AllGoTypes struct {
	IntValue         int        `json:"int_value"`
	Int8Value        int8       `json:"int8_value"`
	Int16Value       int16      `json:"int16_value"`
	Int32Value       int32      `json:"int32_value"`
	Int64Value       int64      `json:"int64_value"`
	UintValue        uint       `json:"uint_value"`
	Uint8Value       uint8      `json:"uint8_value"`
	Uint16Value      uint16     `json:"uint16_value"`
	Uint32Value      uint32     `json:"uint32_value"`
	Uint64Value      uint64     `json:"uint64_value"`
	Float32Value     float32    `json:"float32_value"`
	Float64Value     float64    `json:"float64_value"`
	ByteValue        byte       `json:"byte_value"`
	RuneValue        rune       `json:"rune_value"`
	StringValue      string     `json:"string_value"`
	StringPointer    *string    `json:"string_pointer"`
	ByteSlice        []byte     `json:"raw_data"`
	Timestamp        time.Time  `json:"timestamp"`
	TimestampPointer *time.Time `json:"timestamp_pointer"`
}

func assignValuesToAllGoTypes(timestamp time.Time) AllGoTypes {
	var allTypes AllGoTypes
	allTypes.IntValue = 42

	allTypes.Int8Value = 8
	allTypes.Int16Value = 16
	allTypes.Int32Value = 32
	allTypes.Int64Value = 64
	allTypes.UintValue = 42
	allTypes.Uint8Value = 8
	allTypes.Uint16Value = 16
	allTypes.Uint32Value = 32
	allTypes.Uint64Value = 64
	allTypes.Float64Value = 6.28
	allTypes.ByteValue = byte('a')
	allTypes.RuneValue = 'b'
	allTypes.StringValue = "test"
	allTypes.StringPointer = &allTypes.StringValue
	allTypes.ByteSlice = []byte{'a', 'b', 'c'}
	allTypes.Timestamp = timestamp
	allTypes.TimestampPointer = &timestamp

	return allTypes
}

func TestConversionOfAllTypes(t *testing.T) {
	timestampString := "2025-02-18T08:08:18.970618Z"
	parseTime, parseErr := time.Parse(time.RFC3339Nano, timestampString)
	if parseErr != nil {
		assert.Error(t, parseErr, "timestamp initialization failed")
	}

	allTypes := assignValuesToAllGoTypes(parseTime)
	discoveredDevice := ConvertFromDerivedSchemaToDiscoveredDevice(&allTypes, "https://schema.industrial-assets.io/test/v0.0.1", "Test")

	assert.Equal(t, 19, len(discoveredDevice.Identifiers))

	for _, identifier := range discoveredDevice.Identifiers {
		classifierValue := identifier.Classifiers[0].GetValue()
		propertyName := strings.Split(classifierValue, "#")[1]
		switch propertyName {
		case "int_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.IntValue), extractValue(identifier.String()), "Property: int_value")
		case "int8_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Int8Value), extractValue(identifier.String()), "Property: int8_value")
		case "int16_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Int16Value), extractValue(identifier.String()), "Property: int16_value")
		case "int32_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Int32Value), extractValue(identifier.String()), "Property: int32_value")
		case "int64_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Int64Value), extractValue(identifier.String()), "Property: int64_value")
		case "uint_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.UintValue), extractValue(identifier.String()), "Property: uint_value")
		case "uint8_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Uint8Value), extractValue(identifier.String()), "Property: uint8_value")
		case "uint16_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Uint16Value), extractValue(identifier.String()), "Property: uint16_value")
		case "uint32_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Uint32Value), extractValue(identifier.String()), "Property: uint32_value")
		case "uint64_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Uint64Value), extractValue(identifier.String()), "Property: uint64_value")
		case "float32_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Float32Value), extractValue(identifier.String()), "Property: float32_value")
		case "float64_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.Float64Value), extractValue(identifier.String()), "Property: float64_value")
		case "byte_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.ByteValue), extractValue(identifier.String()), "Property: byte_value")
		case "rune_value":
			assert.Equal(t, fmt.Sprintf("%v", allTypes.RuneValue), extractValue(identifier.String()), "Property: rune_value")
		case "string_value":
			assert.Equal(t, allTypes.StringValue, extractValue(identifier.String()), "Property: string_value")
		case "string_pointer":
			assert.Equal(t, *allTypes.StringPointer, extractValue(identifier.String()), "Property: string_pointer")
		case "raw_data":
			assert.Equal(t, "abc", extractValue(identifier.String()), "Property: raw_data")
		case "timestamp":
			assert.Equal(t, timestampString, extractValue(identifier.String()), "Property: timestamp")
		case "timestamp_pointer":
			assert.Equal(t, timestampString, extractValue(identifier.String()), "Property: timestamp_pointer")
		}
	}
}

func TestConvertToDeviceIdentifiers_IgnoredIdentifier(t *testing.T) {
	device := generateDevice("Profinet", "Device")
	testCases := []struct {
		fieldValue interface{}
		uri        string
	}{
		{device.ProductInstanceIdentifier.IdentifierUncertainty, "https://schema.industrial-assets.io/base/v0.9.0/Asset#asset_identifiers/identifier_uncertainty"},
		{device.InstanceAnnotations, "https://schema.industrial-assets.io/base/v0.9.0/Asset#instance_annotations"},
	}

	for _, testCase := range testCases {
		identifiers := convertToDeviceIdentifiers(reflect.ValueOf(testCase.fieldValue), testCase.uri, 0)
		assert.Empty(t, identifiers, fmt.Sprintf("Expected no identifiers for %s", testCase.uri))
	}
}

func extractValue(input string) string {
	parts := strings.SplitN(input, ":", 2)
	typeValue := strings.Split(parts[1], " ")[0]
	return strings.Trim(typeValue, `"`)
}
