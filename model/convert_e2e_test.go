/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBackAndForthConversion(t *testing.T) {
	testDevice := getTestDevice()
	discoveredDevice := testDevice.ConvertToDiscoveredDevice()
	transformedDevice := ConvertFromDiscoveredDevice(discoveredDevice, "URI")
	assert.IsType(t, map[string]interface{}{}, transformedDevice)
	assert.Equal(t, "Asset", transformedDevice["@type"])
	assert.Equal(t, "TestDevice", transformedDevice["name"])
	connectionPoint := transformedDevice["connection_points"].([]map[string]interface{})[0]
	macIdentifier := transformedDevice["mac_identifiers"].([]map[string]interface{})[0]
	softwareComponent := transformedDevice["software_components"].([]map[string]interface{})[0]
	assert.Equal(t, "Siemens AG", transformedDevice["product_instance_identifier"].(map[string]interface{})["manufacturer_product"].(map[string]interface{})["manufacturer"].(map[string]interface{})["name"])
	assert.Equal(t, "TestDevice", transformedDevice["product_instance_identifier"].(map[string]interface{})["manufacturer_product"].(map[string]interface{})["name"])
	assert.Equal(t, "MyOrderNumber", transformedDevice["product_instance_identifier"].(map[string]interface{})["manufacturer_product"].(map[string]interface{})["product_id"])
	assert.Equal(t, "1.0.0", transformedDevice["product_instance_identifier"].(map[string]interface{})["manufacturer_product"].(map[string]interface{})["product_version"])
	name := softwareComponent["artifact"].(map[string]interface{})["software_identifier"].(map[string]interface{})["name"]
	version := softwareComponent["artifact"].(map[string]interface{})["software_identifier"].(map[string]interface{})["version"]
	assert.Equal(t, "EthernetPort", connectionPoint["connection_point_type"])
	assert.Equal(t, "enp0", connectionPoint["instance_annotations"].([]map[string]interface{})[0]["value"])
	assert.Equal(t, "00:00:00:00:00:00", macIdentifier["mac_address"])
	assert.Equal(t, "firmware", name)
	assert.Equal(t, "1.2.5", version)
	assert.Equal(t, "123456", transformedDevice["product_instance_identifier"].(map[string]interface{})["serial_number"])
	assert.Equal(t, "http://rds.posccaesar.org/ontology/lis14/rdl/", transformedDevice["@context"].(map[string]interface{})["lis"])
	assert.Equal(t, "https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/v0.9.0/",
		transformedDevice["@context"].(map[string]interface{})["@vocab"])
	assert.Equal(t, "https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/v0.9.0/",
		transformedDevice["@context"].(map[string]interface{})["base"])
	assert.Equal(t, "https://w3id.org/linkml/", transformedDevice["@context"].(map[string]interface{})["linkml"])
	assert.Equal(t, "http://www.w3.org/2004/02/skos/core#", transformedDevice["@context"].(map[string]interface{})["skos"])
	assert.Equal(t, "https://schema.org/", transformedDevice["@context"].(map[string]interface{})["schemaorg"])
}

func getTestDevice() *DeviceInfo {
	manufacturer := "Siemens AG"
	product := "TestDevice"
	serialNumber := "123456"
	deviceInfo := NewDevice("Asset", "TestDevice")

	uriOfTheProduct := fmt.Sprintf("https://%s/%s-%s", strings.ReplaceAll(manufacturer, " ", "_"), strings.ReplaceAll(product, " ", "_"), serialNumber)
	deviceInfo.AddNameplate(manufacturer, uriOfTheProduct, "MyOrderNumber", product, "1.0.0", serialNumber)
	deviceInfo.AddSoftware("firmware", "1.2.5")
	macAddress := "00:00:00:00:00:00"
	deviceNIC := "enp0"
	_ = deviceInfo.AddNic(deviceNIC, macAddress)
	return deviceInfo
}
