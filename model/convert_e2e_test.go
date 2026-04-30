/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackAndForthConversion(t *testing.T) {
	testDevice := getTestDevice()
	discoveredDevice := testDevice.ConvertToDiscoveredDevice()
	transformedDevice := ConvertFromDiscoveredDevice(discoveredDevice, "URI")

	assert.IsType(t, map[string]interface{}{}, transformedDevice)

	assert.Equal(t, "Asset", transformedDevice["functional_object_type"])
	assert.Equal(t, FunctionalObjectSchemaUrl, transformedDevice["functional_object_schema_url"])
	assert.Equal(t, "TestDevice", transformedDevice["name"])
	connectionPoint := transformedDevice["connection_points"].([]map[string]interface{})[0]
	assetIdentifiers := transformedDevice["asset_identifiers"].([]map[string]interface{})
	macIdentifier := findMapWithKey(assetIdentifiers, "mac_address")
	if !assert.NotNil(t, macIdentifier) {
		return
	}
	softwareComponent := transformedDevice["software_components"].([]map[string]interface{})[0]
	productInstanceInformation := transformedDevice["product_instance_information"].(map[string]interface{})
	manufacturerProduct := productInstanceInformation["manufacturer_product"].(map[string]interface{})
	manufacturer := manufacturerProduct["manufacturer"].(map[string]interface{})
	assert.Equal(t, "Siemens AG", manufacturer["name"])
	assert.Equal(t, testIDLink, manufacturerProduct["product_link"])
	assert.Equal(t, "MyOrderNumber", manufacturerProduct["product_id"])
	assert.Equal(t, "1.0.0", manufacturerProduct["product_version"])

	softwareIdentifier := softwareComponent["asset_identifiers"].([]map[string]interface{})[0]
	swName := softwareIdentifier["name"]
	swVersion := softwareIdentifier["version"]
	isFirmwareBool, isFirmwareErr := strconv.ParseBool(softwareComponent["is_firmware"].(string))
	assert.Equal(t, "Firmware", swName)
	assert.Equal(t, "1.2.5", swVersion)
	assert.Nil(t, isFirmwareErr)
	assert.True(t, isFirmwareBool)

	assert.Equal(t, "EthernetPort", connectionPoint["connection_point_type"])
	assert.Equal(t, "enp0", connectionPoint["name"])
	assert.Equal(t, "00:00:00:00:00:00", macIdentifier["mac_address"])
	assert.Equal(t, "123456", productInstanceInformation["serial_number"])
}

func findMapWithKey(values []map[string]interface{}, key string) map[string]interface{} {
	for _, value := range values {
		if _, ok := value[key]; ok {
			return value
		}
	}
	return nil
}

func getTestDevice() *DeviceInfo {
	manufacturer := "Siemens AG"
	product := "TestDevice"
	serialNumber := "123456"
	deviceInfo, err := NewDevice("Asset", "TestDevice")
	if err != nil {
		panic(err)
	}

	uriOfTheProduct := testIDLink
	err = deviceInfo.AddNameplate(manufacturer, uriOfTheProduct, "MyOrderNumber", product, "1.0.0", serialNumber)
	if err != nil {
		panic(err)
	}
	err = deviceInfo.AddSoftware("Firmware", "1.2.5", true)
	if err != nil {
		panic(err)
	}
	macAddress := "00:00:00:00:00:00"
	deviceNIC := "enp0"
	_, _ = deviceInfo.AddNic(deviceNIC, macAddress)
	return deviceInfo
}
