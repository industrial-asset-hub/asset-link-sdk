/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package model

import (
	"encoding/json"
	"os"
	"testing"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/stretchr/testify/assert"
)

func TestDeviceTransformation(t *testing.T) {
	t.Run("ConvertFromDiscoveredDevice when provided a device with identifier value of type text transforms it successfully", func(t *testing.T) {
		testDeviceForText := &generated.DiscoveredDevice{
			Identifiers: []*generated.DeviceIdentifier{{
				Value: &generated.DeviceIdentifier_Text{
					Text: "AUMA Riester GmbH & Co.KG",
				},
				Classifiers: []*generated.SemanticClassifier{{
					Type: "URI",
					Value: baseSchemaPrefix +
						"/Device#product_instance_information/manufacturer_product/manufacturer/name",
				}},
			},
			},
			Timestamp: 1000010000100010,
		}
		expectedType := "URI"
		actualResult := ConvertFromDiscoveredDevice(testDeviceForText, expectedType)
		expectedResult := map[string]interface{}{
			"functional_object_type":       "Device",
			"functional_object_schema_url": FunctionalObjectSchemaUrl,
			"product_instance_information": map[string]interface{}{
				"manufacturer_product": map[string]interface{}{
					"manufacturer": map[string]interface{}{
						"name": "AUMA Riester GmbH & Co.KG",
					},
				},
			},
		}
		assert.Equal(t, expectedResult, actualResult)
	},
	)

	t.Run("ConvertFromDiscoveredDevice when provided a device with identifier value of type int64 transforms it successfully", func(t *testing.T) {
		testDeviceForInt := &generated.DiscoveredDevice{
			Identifiers: []*generated.DeviceIdentifier{{
				Value: &generated.DeviceIdentifier_Int64Value{
					Int64Value: int64(1),
				},
				Classifiers: []*generated.SemanticClassifier{{
					Type:  "URI",
					Value: baseSchemaPrefix + "/Device#test-1",
				}},
			},
			},
			Timestamp: 1000010000100010,
		}
		expectedType := "URI"
		actualResult := ConvertFromDiscoveredDevice(testDeviceForInt, expectedType)
		expectedResult := map[string]interface{}{
			"functional_object_type":       "Device",
			"functional_object_schema_url": FunctionalObjectSchemaUrl,
			"test-1":                       int64(1),
		}
		assert.Equal(t, expectedResult, actualResult)
	},
	)
	t.Run("ConvertFromDiscoveredDevice when provided a device with identifier value of type float64 transforms it successfully", func(t *testing.T) {
		testDeviceForFloat := &generated.DiscoveredDevice{
			Identifiers: []*generated.DeviceIdentifier{{
				Value: &generated.DeviceIdentifier_Float64Value{
					Float64Value: float64(0.1),
				},
				Classifiers: []*generated.SemanticClassifier{{
					Type:  "URI",
					Value: baseSchemaPrefix + "/Device#test-2/A/B",
				}},
			},
			},
			Timestamp: 1000010000100010,
		}
		expectedType := "URI"
		actualResult := ConvertFromDiscoveredDevice(testDeviceForFloat, expectedType)
		expectedResult := map[string]interface{}{
			"functional_object_type":       "Device",
			"functional_object_schema_url": FunctionalObjectSchemaUrl,
			"test-2": map[string]interface{}{
				"A": map[string]interface{}{
					"B": 0.1,
				},
			},
		}
		assert.Equal(t, expectedResult, actualResult)
	},
	)
	t.Run("ConvertFromDiscoveredDevice when provided a device with identifier value of type rawData transforms it successfully", func(t *testing.T) {
		testDeviceForRawData := &generated.DiscoveredDevice{
			Identifiers: []*generated.DeviceIdentifier{{
				Value: &generated.DeviceIdentifier_RawData{
					RawData: []byte{1},
				},
				Classifiers: []*generated.SemanticClassifier{{
					Type:  "URI",
					Value: baseSchemaPrefix + "/Device#test-2/A/B",
				}},
			},
			},
			Timestamp: 1000010000100010,
		}
		expectedType := "URI"
		actualResult := ConvertFromDiscoveredDevice(testDeviceForRawData, expectedType)
		expectedResult := map[string]interface{}{
			"functional_object_type":       "Device",
			"functional_object_schema_url": FunctionalObjectSchemaUrl,
			"test-2": map[string]interface{}{
				"A": map[string]interface{}{
					"B": []byte{1},
				},
			},
		}
		assert.Equal(t, expectedResult, actualResult)
	},
	)
	t.Run("ConvertFromDiscoveredDevice when provided a device with identifier value of type children transforms it successfully", func(t *testing.T) {
		testDeviceForChildren := &generated.DiscoveredDevice{
			Identifiers: []*generated.DeviceIdentifier{{
				Value: &generated.DeviceIdentifier_Children{
					Children: &generated.DeviceIdentifierValueList{
						Value: []*generated.DeviceIdentifier{
							{
								Value: &generated.DeviceIdentifier_Text{
									Text: "test-connection-point",
								},
								Classifiers: []*generated.SemanticClassifier{{
									Type:  "URI",
									Value: baseSchemaPrefix + "/Device#connection_points/related_connection_points/connection_point",
								}},
							},
						},
					},
				},
				Classifiers: []*generated.SemanticClassifier{{
					Type:  "URI",
					Value: baseSchemaPrefix + "/Device#connection_points",
				}},
			},
			},
			Timestamp: 1000010000100010,
		}
		expectedType := "URI"
		actualResult := ConvertFromDiscoveredDevice(testDeviceForChildren, expectedType)
		expectedResult := map[string]interface{}{
			"functional_object_type":       "Device",
			"functional_object_schema_url": FunctionalObjectSchemaUrl,
			"connection_points": []map[string]interface{}{
				{"related_connection_points": map[string]interface{}{"connection_point": "test-connection-point"}}},
		}
		assert.Equal(t, expectedResult, actualResult)
	},
	)
	t.Run("ConvertFromDiscoveredDevice when provided a device with several identifiers transforms it successfully", func(t *testing.T) {
		testDevice := &generated.DiscoveredDevice{
			Identifiers: []*generated.DeviceIdentifier{{
				Value: &generated.DeviceIdentifier_Children{
					Children: &generated.DeviceIdentifierValueList{
						Value: []*generated.DeviceIdentifier{
							{
								Value: &generated.DeviceIdentifier_Text{
									Text: "another-child-value",
								},
								Classifiers: []*generated.SemanticClassifier{{
									Type:  "URI",
									Value: baseSchemaPrefix + "/Device#parent-property/another-child-property",
								}},
							},
						},
					},
				},
				Classifiers: []*generated.SemanticClassifier{{
					Type:  "URI",
					Value: baseSchemaPrefix + "/Device#parent-property",
				}},
			},
				{
					Value: &generated.DeviceIdentifier_Children{
						Children: &generated.DeviceIdentifierValueList{
							Value: []*generated.DeviceIdentifier{
								{
									Value: &generated.DeviceIdentifier_Text{
										Text: "child-value",
									},
									Classifiers: []*generated.SemanticClassifier{{
										Type:  "URI",
										Value: baseSchemaPrefix + "/Device#parent-property/child-property",
									}},
								},
							},
						},
					},
					Classifiers: []*generated.SemanticClassifier{{
						Type:  "URI",
						Value: baseSchemaPrefix + "/Device#parent-property",
					}},
				},
				{
					Value: &generated.DeviceIdentifier_Text{
						Text: "test-serial-number",
					},
					Classifiers: []*generated.SemanticClassifier{{
						Type:  "URI",
						Value: baseSchemaPrefix + "/Device#product_instance_identifier/serial_number",
					}},
				},
			},

			Timestamp: 1000010000100010,
		}
		expectedType := "URI"
		actualResult := ConvertFromDiscoveredDevice(testDevice, expectedType)
		expectedResult := map[string]interface{}{
			"functional_object_type":       "Device",
			"functional_object_schema_url": FunctionalObjectSchemaUrl,
			"parent-property": []map[string]interface{}{
				{"another-child-property": "another-child-value"},
				{"child-property": "child-value"}},
			"product_instance_identifier": map[string]interface{}{
				"serial_number": "test-serial-number",
			},
		}
		assert.Equal(t, expectedResult, actualResult)
	},
	)
	t.Run("ConvertFromDiscoveredDevice should set functional_object_schema_url if not present", func(t *testing.T) {
		testDevice := &generated.DiscoveredDevice{
			Identifiers: []*generated.DeviceIdentifier{},
			Timestamp:   1000010000100010,
		}
		actualResult := ConvertFromDiscoveredDevice(testDevice, "URI")
		assert.Equal(t, FunctionalObjectSchemaUrl, actualResult["functional_object_schema_url"])
	})

	t.Run("ConvertFromDiscoveredDevice should set functional_object_type if not present", func(t *testing.T) {
		testDevice := &generated.DiscoveredDevice{
			Identifiers: []*generated.DeviceIdentifier{},
			Timestamp:   1000010000100010,
		}
		expectedType := "Asset"
		actualResult := ConvertFromDiscoveredDevice(testDevice, "URI")
		assert.Equal(t, expectedType, actualResult["functional_object_type"])
	})
	t.Run("ConvertFromDiscoveredDevice should not set functional_object_schema_url explicitly if it is already present", func(t *testing.T) {
		testDevice := DeviceInfo{}
		testDevice.FunctionalObjectSchemaUrl = "https://example.com/custom-schema.json"
		discoveredDevice := testDevice.ConvertToDiscoveredDevice()
		actualResult := ConvertFromDiscoveredDevice(discoveredDevice, "URI")
		assert.Equal(t, testDevice.FunctionalObjectSchemaUrl, actualResult["functional_object_schema_url"])
	})

	t.Run("ConvertFromDiscoveredDevice should not set functional_object_type explicitly if it is already present", func(t *testing.T) {
		testDevice := DeviceInfo{}
		testDevice.FunctionalObjectType = "test-type"
		expectedType := "test-type"
		actualResult := ConvertFromDiscoveredDevice(testDevice.ConvertToDiscoveredDevice(), "URI")
		assert.Equal(t, expectedType, actualResult["functional_object_type"])
	})
}

func TestFilter(t *testing.T) {
	t.Run("filter when provided with invalid classifier type should filter out classifier", func(t *testing.T) {
		testIdentifier := &generated.DeviceIdentifier{
			Value: &generated.DeviceIdentifier_Text{
				Text: "test-text",
			},
			Classifiers: []*generated.SemanticClassifier{{
				Type:  "test",
				Value: baseSchemaPrefix + "/Device#A",
			}},
		}
		result := filter(testIdentifier, "URI", baseSchemaPrefix)
		assert.Nil(t, result)
	})
	t.Run("filter when provided with valid classifier type should not filter out", func(t *testing.T) {
		testIdentifier := &generated.DeviceIdentifier{
			Value: &generated.DeviceIdentifier_Text{
				Text: "test-text",
			},
			Classifiers: []*generated.SemanticClassifier{{
				Type:  "URI",
				Value: baseSchemaPrefix + "/Device#A",
			}},
		}
		result := filter(testIdentifier, "URI", baseSchemaPrefix)
		assert.NotNil(t, result)
	})
}

func TestTransformKeys(t *testing.T) {
	t.Run("transformKeys when provided with keys it assigns transformed property to device", func(t *testing.T) {
		testKeys := []string{"prop1", "prop2", "prop3", "prop4"}
		testDevice := map[string]interface{}{
			"prop1": map[string]interface{}{},
		}
		expectedDevice := map[string]interface{}{
			"prop1": map[string]interface{}{
				"prop2": map[string]interface{}{
					"prop3": map[string]interface{}{
						"prop4": "test-value",
					},
				},
			},
		}
		transformKeys(testKeys, "test-value", testDevice)
		assert.Equal(t, expectedDevice, testDevice)
	})

}

func TestRetrieveAssetTypeFromDiscoveredDeviceForValidClassifier(t *testing.T) {
	assetTypeToCheck := "Device"
	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{{
			Value: &generated.DeviceIdentifier_Text{
				Text: "test-text",
			},
			Classifiers: []*generated.SemanticClassifier{{
				Type:  "URI",
				Value: baseSchemaPrefix + "/Device#A",
			}},
		}},
	}
	assetType := retrieveAssetTypeFromDiscoveredDevice(&discoveredDevice)
	assert.Equal(t, assetTypeToCheck, assetType)
}

func TestRetrieveAssetTypeFromDiscoveredDeviceWithoutAnyIdentifier(t *testing.T) {
	assetTypeToCheck := "Asset"
	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{},
	}
	assetType := retrieveAssetTypeFromDiscoveredDevice(&discoveredDevice)
	assert.Equal(t, assetTypeToCheck, assetType)
}

func TestRetrieveAssetTypeFromDiscoveredDeviceWithUnsupportedClassifier(t *testing.T) {
	assetTypeToCheck := "Asset"
	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{{
			Value: &generated.DeviceIdentifier_Text{
				Text: "test-text",
			},
			Classifiers: []*generated.SemanticClassifier{{
				Type:  "any-other-test-type-that-is-not-supported",
				Value: "unsupported-value",
			}},
		}},
	}
	assetType := retrieveAssetTypeFromDiscoveredDevice(&discoveredDevice)
	assert.Equal(t, assetTypeToCheck, assetType)
}

func TestRetrieveAssetTypeFromDiscoveredDeviceWithInvalidClassifier(t *testing.T) {
	assetTypeToCheck := "Asset"
	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{{
			Value: &generated.DeviceIdentifier_Text{
				Text: "test-text",
			},
			Classifiers: []*generated.SemanticClassifier{{
				Type:  "URI",
				Value: "https://schema.industrial-assets.io/un-supported-uri-format",
			}},
		}},
	}
	assetType := retrieveAssetTypeFromDiscoveredDevice(&discoveredDevice)
	assert.Equal(t, assetTypeToCheck, assetType)
}

func TestMapManyArrayElementsViaMultipleArrayContainersIntoIahDevice(t *testing.T) {
	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{
			{
				Value: &generated.DeviceIdentifier_Children{
					Children: &generated.DeviceIdentifierValueList{
						Value: []*generated.DeviceIdentifier{
							{
								Value:       &generated.DeviceIdentifier_Text{Text: "element-1"},
								Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array/id"}},
							},
							{
								Value:       &generated.DeviceIdentifier_Text{Text: "element-1-name"},
								Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array/name"}},
							},
						},
					},
				},
				Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array"}},
			},
			{
				Value: &generated.DeviceIdentifier_Children{
					Children: &generated.DeviceIdentifierValueList{
						Value: []*generated.DeviceIdentifier{
							{
								Value:       &generated.DeviceIdentifier_Text{Text: "element-2"},
								Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array/id"}},
							},
							{
								Value:       &generated.DeviceIdentifier_Text{Text: "element-2-name"},
								Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array/name"}},
							},
						},
					},
				},
				Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array"}},
			},
		},
	}
	iahDevice := ConvertFromDiscoveredDevice(&discoveredDevice, "URI")
	assert.Equal(t, len(iahDevice["array"].([]map[string]interface{})), 2)
	assert.Equal(t, len(iahDevice["array"].([]map[string]interface{})[0]), 2)
	assert.Equal(t, len(iahDevice["array"].([]map[string]interface{})[1]), 2)
}

func TestMapManyArrayElementsIntoIahDevice(t *testing.T) {
	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{{
			Value: &generated.DeviceIdentifier_Children{
				Children: &generated.DeviceIdentifierValueList{
					Value: []*generated.DeviceIdentifier{
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "element-1"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array%5B0%5D%2Fid"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "some-name-1"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array%5B0%5D%2Fname"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "element-2"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array%5B1%5D%2Fid"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "some-name-2"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array%5B1%5D%2Fname"}},
						},
					},
				},
			},
			Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array"}},
		}}}
	iahDevice := ConvertFromDiscoveredDevice(&discoveredDevice, "URI")

	arrayProperty, ok := iahDevice["array"].([]map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(arrayProperty))
	assert.Equal(t, 2, len(arrayProperty[0]))
	assert.Equal(t, 2, len(arrayProperty[1]))
}

func TestMapManyMixedArrayElementsIntoIahDevice(t *testing.T) {
	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{{
			Value: &generated.DeviceIdentifier_Children{
				Children: &generated.DeviceIdentifierValueList{
					Value: []*generated.DeviceIdentifier{
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "element-1"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array%5B0%5D%2Fid"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "some-name-2"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array%5B1%5D%2Fname"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "some-name-1"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array%5B0%5D%2Fname"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "element-2"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array%5B1%5D%2Fid"}},
						},
					},
				},
			},
			Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#array"}},
		}}}
	iahDevice := ConvertFromDiscoveredDevice(&discoveredDevice, "URI")
	arrayProperty, ok := iahDevice["array"].([]map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(arrayProperty))
	assert.Equal(t, 2, len(arrayProperty[0]))
	assert.Equal(t, 2, len(arrayProperty[1]))
	assert.Equal(t, "some-name-1", arrayProperty[0]["name"])
	assert.Equal(t, "some-name-2", arrayProperty[1]["name"])
}

func TestMapManyDeepArrayElementsIntoIahDevice(t *testing.T) {
	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{{
			Value: &generated.DeviceIdentifier_Children{
				Children: &generated.DeviceIdentifierValueList{
					Value: []*generated.DeviceIdentifier{
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "element-1"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fid"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "some-name-1"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fname"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "element-2"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fid"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "some-name-2"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fname"}},
						},
					},
				},
			},
			Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray"}},
		}}}
	iahDevice := ConvertFromDiscoveredDevice(&discoveredDevice, "URI")
	arrayProperty, ok := iahDevice["a"].(map[string]interface{})["deeper"].(map[string]interface{})["array"].([]map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(arrayProperty))
	assert.Equal(t, 2, len(arrayProperty[0]))
	assert.Equal(t, 2, len(arrayProperty[1]))
	assert.Equal(t, "element-1", arrayProperty[0]["id"])
	assert.Equal(t, "element-2", arrayProperty[1]["id"])
}

func TestMapManyDeepArrayElementsWithDeepPathsIntoIahDevice(t *testing.T) {
	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{{
			Value: &generated.DeviceIdentifier_Children{
				Children: &generated.DeviceIdentifierValueList{
					Value: []*generated.DeviceIdentifier{
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "element-1"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fid"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "some-name-1"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fname"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "element-2"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fid"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "some-name-2"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fname"}},
						},
					},
				},
			},
			Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray"}},
		}}}
	iahDevice := ConvertFromDiscoveredDevice(&discoveredDevice, "URI")

	arrayProperty, ok := iahDevice["a"].(map[string]interface{})["deeper"].(map[string]interface{})["array"].([]map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(arrayProperty))
	assert.Equal(t, "element-1", arrayProperty[0]["some_object"].(map[string]interface{})["id"])
	assert.Equal(t, "element-2", arrayProperty[1]["some_object"].(map[string]interface{})["id"])
}

func TestMapManyDeepArrayElementsWithDeepPathsThatContainArraysIntoIahDevice(t *testing.T) {
	expectedResult := map[string]interface{}{"functional_object_schema_url": FunctionalObjectSchemaUrl, "functional_object_type": "device", "a": map[string]interface{}{"deeper": map[string]interface{}{"array": []map[string]interface{}{{"some_object": map[string]interface{}{"connection_points": []map[string]interface{}{{"id": "array-0-connection-point-0", "related_connection_points": []map[string]interface{}{{"id": "array-0-con-point-0-related-connection-point-0"}, {"id": "array-0-con-point-0-related-connection-point-1"}}}, {}, {"related_connection_points": []map[string]interface{}{{}, {}, {"id": "array-0-con-point-2-related-connection-point-2"}, {"id": "array-0-con-point-2-related-connection-point-3"}}}}, "name": "array-0-name"}}, {"some_object": map[string]interface{}{"connection_points": []map[string]interface{}{{"id": "array-1-connection-point-0", "related_connection_points": []map[string]interface{}{{"id": "array-1-con-point-0-related-connection-point-0"}, {"id": "array-1-con-point-0-related-connection-point-1"}}}, {"id": "array-1-connection-point-1", "related_connection_points": []map[string]interface{}{{}, {}, {"id": "array-1-con-point-1-related-connection-point-2"}, {"id": "array-1-con-point-1-related-connection-point-3"}}}}, "id": "array-1-id", "name": "array-1-name"}}, {"some_object": map[string]interface{}{"id": "array-2-id"}}}}}}

	discoveredDevice := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{{
			Value: &generated.DeviceIdentifier_Children{
				Children: &generated.DeviceIdentifierValueList{
					Value: []*generated.DeviceIdentifier{
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "array-1-id"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fid"}},
						},
						{
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fconnection_points"}},
							Value: &generated.DeviceIdentifier_Children{
								Children: &generated.DeviceIdentifierValueList{
									Value: []*generated.DeviceIdentifier{
										{
											Value:       &generated.DeviceIdentifier_Text{Text: "array-0-connection-point-0"},
											Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fconnection_points%5B0%5D%2Fid"}},
										},
										{
											Value: &generated.DeviceIdentifier_Children{
												Children: &generated.DeviceIdentifierValueList{
													Value: []*generated.DeviceIdentifier{
														{
															Value:       &generated.DeviceIdentifier_Text{Text: "array-0-con-point-0-related-connection-point-0"},
															Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fconnection_points%5B0%5D%2Frelated_connection_points%5B0%5D%2Fid"}},
														},
														{
															Value:       &generated.DeviceIdentifier_Text{Text: "array-0-con-point-0-related-connection-point-1"},
															Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fconnection_points%5B0%5D%2Frelated_connection_points%5B1%5D%2Fid"}},
														},
													},
												},
											},
											Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fconnection_points%5B0%5D%2Frelated_connection_points"}},
										},
										{
											Value: &generated.DeviceIdentifier_Children{
												Children: &generated.DeviceIdentifierValueList{
													Value: []*generated.DeviceIdentifier{
														{
															Value:       &generated.DeviceIdentifier_Text{Text: "array-0-con-point-2-related-connection-point-2"},
															Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fconnection_points%5B2%5D%2Frelated_connection_points%5B2%5D%2Fid"}},
														},
														{
															Value:       &generated.DeviceIdentifier_Text{Text: "array-0-con-point-2-related-connection-point-3"},
															Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fconnection_points%5B2%5D%2Frelated_connection_points%5B3%5D%2Fid"}},
														},
													},
												},
											},
											Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fconnection_points%5B2%5D%2Frelated_connection_points"}},
										},
									},
								},
							},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "array-0-name"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B0%5D%2Fsome_object%2Fname"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "array-2-id"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B2%5D%2Fsome_object%2Fid"}},
						},
						{
							Value:       &generated.DeviceIdentifier_Text{Text: "array-1-name"},
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fname"}},
						},
						{
							Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fconnection_points"}},
							Value: &generated.DeviceIdentifier_Children{
								Children: &generated.DeviceIdentifierValueList{
									Value: []*generated.DeviceIdentifier{
										{
											Value:       &generated.DeviceIdentifier_Text{Text: "array-1-connection-point-0"},
											Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fconnection_points%5B0%5D%2Fid"}},
										},
										{
											Value: &generated.DeviceIdentifier_Children{
												Children: &generated.DeviceIdentifierValueList{
													Value: []*generated.DeviceIdentifier{
														{
															Value:       &generated.DeviceIdentifier_Text{Text: "array-1-con-point-0-related-connection-point-0"},
															Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fconnection_points%5B0%5D%2Frelated_connection_points%5B0%5D%2Fid"}},
														},
														{
															Value:       &generated.DeviceIdentifier_Text{Text: "array-1-con-point-0-related-connection-point-1"},
															Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fconnection_points%5B0%5D%2Frelated_connection_points%5B1%5D%2Fid"}},
														},
													},
												},
											},
											Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fconnection_points%5B0%5D%2Frelated_connection_points"}},
										},
										{
											Value: &generated.DeviceIdentifier_Children{
												Children: &generated.DeviceIdentifierValueList{
													Value: []*generated.DeviceIdentifier{
														{
															Value:       &generated.DeviceIdentifier_Text{Text: "array-1-con-point-1-related-connection-point-2"},
															Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fconnection_points%5B1%5D%2Frelated_connection_points%5B2%5D%2Fid"}},
														},
														{
															Value:       &generated.DeviceIdentifier_Text{Text: "array-1-con-point-1-related-connection-point-3"},
															Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fconnection_points%5B1%5D%2Frelated_connection_points%5B3%5D%2Fid"}},
														},
													},
												},
											},
											Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fconnection_points%5B1%5D%2Frelated_connection_points"}},
										},
										{
											Value:       &generated.DeviceIdentifier_Text{Text: "array-1-connection-point-1"},
											Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray%5B1%5D%2Fsome_object%2Fconnection_points%5B1%5D%2Fid"}},
										},
									},
								},
							},
						},
					},
				},
			},
			Classifiers: []*generated.SemanticClassifier{{Type: "URI", Value: baseSchemaPrefix + "/device#a%2Fdeeper%2Farray"}},
		}}, Timestamp: 1000010000100010}
	iahDevice := ConvertFromDiscoveredDevice(&discoveredDevice, "URI")

	assert.Equal(t, expectedResult, iahDevice)
}
func TestPanics(t *testing.T) {

	t.Run("transformKeys panics if empty key is found", func(t *testing.T) {
		testKeys := []string{}
		testDevice := map[string]interface{}{}
		assert.NotPanics(t, func() { transformKeys(testKeys, "", testDevice) })
	})

	t.Run("transformKeys panics if ip v6 connection point is provided", func(t *testing.T) {
		testKeys := []string{"prop1", "prop2", "prop3", "prop4"}

		testDevice := make(map[string]interface{})
		assert.NotPanics(t, func() { transformKeys(testKeys, "test-value", testDevice) })
	})
}

func TestCheckConversionForFile(t *testing.T) {
	t.Run("validate PNAS output against IAH schema", func(t *testing.T) {
		t.Skip("Skipping test as it is only used to get an end to end example going manually")
		var testDevice map[string]interface{}
		var scannedPnasResponse generated.DiscoverResponse
		// open file for reading
		file, err := os.Open("./example_dummy_device.json")
		if err != nil {
			t.Fatalf("failed to open input file: %v", err)
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				t.Fatalf("failed to close input file: %v", closeErr)
			}
		}()

		// create result file for writing
		resultFile, err := os.Create("./mapped_device_after_refactor.json")
		if err != nil {
			t.Fatalf("failed to create output file: %v", err)
		}
		defer func() {
			if closeErr := resultFile.Close(); closeErr != nil {
				t.Fatalf("failed to close output file: %v", closeErr)
			}
		}()

		// get file size for creating a matching buffer
		fileInfo, _ := file.Stat()
		byteBuffer := make([]byte, fileInfo.Size())
		_, _ = file.Read(byteBuffer)
		unmarshalOptions := protojson.UnmarshalOptions{
			DiscardUnknown: true,
			AllowPartial:   true,
		}
		if err := unmarshalOptions.Unmarshal(byteBuffer, &scannedPnasResponse); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}
		// do the actual schema transformation
		testDevice = ConvertFromDiscoveredDevice(scannedPnasResponse.Devices[0], "URI")

		// Write the result to the output file
		jsonWriter := json.NewEncoder(resultFile)
		if err := jsonWriter.Encode(testDevice); err != nil {
			t.Fatalf("failed to write result to output file: %v", err)
		}

		// Perform validation against IAH schema
		// ...
		// linkml-validate --include-range-class-descendants --target-class=ProfinetDevice -s ./profinet.yaml mapped_iah_device.json
		// cmd := exec.Command("linkml-validate", "--target-class=ProfinetDevice", "-s", "./profinet.yaml", "./basic_mapped_array_device.json")

		// if err = cmd.Run(); err != nil {
		//  t.Log(err)
		// }
	})
}
