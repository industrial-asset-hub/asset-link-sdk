/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"testing"

	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestConvertToPropertyValueResultsCompleteAsset(t *testing.T) {
	device := generateDevice("Device", "Profinet")

	actual, err := device.ConvertToPropertyValueResults()
	require.NoError(t, err)
	require.NotEmpty(t, actual)

	expected := expectedPropertyValueResultsForGeneratedDevice()
	assertPropertyResultsEqual(t, expected, actual)
}

func TestConvertToPropertyValueResults(t *testing.T) {
	device := generateDevice("Device", "Profinet")

	results, err := device.ConvertToPropertyValueResults()
	require.NoError(t, err)
	require.NotEmpty(t, results)

	resultsByKey := propertyValuesByKey(results)
	assertPropertyValueStructure(t, resultsByKey)
}

func TestConvertToPropertyValueResultsNilDevice(t *testing.T) {
	var device *DeviceInfo

	results, err := device.ConvertToPropertyValueResults()
	require.Error(t, err)
	require.Nil(t, results)
}

func expectedPropertyValueResultsForGeneratedDevice() []*generatedDeviceInfo.PropertyValueResult {
	return []*generatedDeviceInfo.PropertyValueResult{
		propertyResult("asset_identifiers", variantArray(
			variantStruct(map[string]*generated.Variant{
				"asset_identifier_type":  textVariant("MacIdentifier"),
				"identifier_uncertainty": int64Variant(1),
				"mac_address":            textVariant("12:12:12:12:12:12"),
			}),
		)),
		propertyResult("asset_operations", variantArray(
			variantStruct(map[string]*generated.Variant{
				"activation_flag": boolVariant(true),
				"operation_name":  textVariant("Firmware Update"),
			}),
		)),
		propertyResult("connection_points", variantArray(
			variantStruct(map[string]*generated.Variant{
				"connection_point_type": textVariant("Ipv4Connectivity"),
				"id":                    textVariant("1"),
				"ipv4_address":          textVariant("192.168.0.1"),
				"network_mask":          textVariant("255.255.255.0"),
				"related_connection_points": variantArray(
					variantStruct(map[string]*generated.Variant{
						"connection_point_id": textVariant("eth0"),
						"custom_relationship": textVariant("EthernetPort"),
					}),
				),
			}),
			variantStruct(map[string]*generated.Variant{
				"connection_point_type": textVariant("Ipv6Connectivity"),
				"id":                    textVariant("2"),
				"ipv6_address":          textVariant("fd12:3456:789a::1"),
				"router_ipv6_address":   textVariant("fd12:3456:789a::1"),
			}),
			variantStruct(map[string]*generated.Variant{
				"connection_point_type": textVariant("EthernetPort"),
				"id":                    textVariant("3"),
				"mac_address":           textVariant("12:12:12:12:12:12"),
			}),
		)),
		propertyResult("functional_object_schema_url", textVariant("https://industrial-assets.io/schemas/iah/base-schema/released/v1/iah-base.json")),
		propertyResult("functional_object_type", textVariant("Device")),
		propertyResult("name", textVariant("Device")),
		propertyResult("product_instance_information", variantStruct(map[string]*generated.Variant{
			"manufacturer_product": variantStruct(map[string]*generated.Variant{
				"manufacturer": variantStruct(map[string]*generated.Variant{
					"name": textVariant("test-vendor"),
				}),
				"product_id":      textVariant("test-product"),
				"product_version": textVariant("1.0.0"),
			}),
			"serial_number": textVariant("test"),
		})),
	}
}

func propertyResult(key string, value *generated.Variant) *generatedDeviceInfo.PropertyValueResult {
	return &generatedDeviceInfo.PropertyValueResult{
		Result: &generatedDeviceInfo.PropertyValueResult_Property{
			Property: &generatedDeviceInfo.PropertyKeyValuePair{
				Key:   key,
				Value: value,
			},
		},
	}
}

func assertPropertyResultsEqual(t *testing.T, expected, actual []*generatedDeviceInfo.PropertyValueResult) {
	t.Helper()
	require.Len(t, actual, len(expected))

	for index := range expected {
		require.Truef(t, proto.Equal(expected[index], actual[index]), "property result %d mismatch", index)
	}
}

func assertPropertyValueStructure(t *testing.T, resultsByKey map[string]*generated.Variant) {
	t.Helper()

	require.Contains(t, resultsByKey, "connection_points")
	require.Contains(t, resultsByKey, "product_instance_information")
	require.NotContains(t, resultsByKey, "connection_points.")
	require.NotContains(t, resultsByKey, "connection_points[")
	require.NotNil(t, resultsByKey["connection_points"].GetArrayValue())
	require.NotNil(t, resultsByKey["product_instance_information"].GetStructValue())
}

func propertyValuesByKey(results []*generatedDeviceInfo.PropertyValueResult) map[string]*generated.Variant {
	out := make(map[string]*generated.Variant, len(results))
	for _, result := range results {
		property := result.GetProperty()
		if property == nil {
			continue
		}
		out[property.GetKey()] = property.GetValue()
	}
	return out
}

func variantStruct(fields map[string]*generated.Variant) *generated.Variant {
	return &generated.Variant{Value: &generated.Variant_StructValue{StructValue: &generated.VariantStruct{Fields: fields}}}
}

func variantArray(values ...*generated.Variant) *generated.Variant {
	return &generated.Variant{Value: &generated.Variant_ArrayValue{ArrayValue: &generated.VariantArray{Values: values}}}
}

func textVariant(value string) *generated.Variant {
	return &generated.Variant{Value: &generated.Variant_Text{Text: value}}
}

func boolVariant(value bool) *generated.Variant {
	return &generated.Variant{Value: &generated.Variant_BoolValue{BoolValue: value}}
}

func int64Variant(value int64) *generated.Variant {
	return &generated.Variant{Value: &generated.Variant_Int64Value{Int64Value: value}}
}
