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
	"sync/atomic"
	"testing"
)

var lastSerialNumber = atomic.Int64{}

func TestBackAndForthConversion(t *testing.T) {
	testDevice := getTestDevice()
	discoveredDevice := testDevice.ConvertToDiscoveredDevice()
	transformedDevice := TransformDevice(discoveredDevice, "URI")
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
}

func getTestDevice() *DeviceInfo {
	manufacturer := "Siemens AG"
	product := "TestDevice"
	serialNumber := fmt.Sprint(lastSerialNumber.Load())
	deviceInfo := NewDevice("Asset", "TestDevice")

	uriOfTheProduct := fmt.Sprintf("https://%s/%s-%s", strings.ReplaceAll(manufacturer, " ", "_"), strings.ReplaceAll(product, " ", "_"), serialNumber)
	deviceInfo.AddNameplate(manufacturer, uriOfTheProduct, "MyOrderNumber", product, "1.0.0", serialNumber)
	deviceInfo.AddSoftware("firmware", "1.2.5")
	macAddress := "00:00:00:00:00:00"
	deviceNIC := "enp0"
	_ = deviceInfo.AddNic(deviceNIC, macAddress)
	return deviceInfo
}
