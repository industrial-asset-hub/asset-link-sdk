/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"testing"
)

func TestConvertToJson(t *testing.T) {
	device := NewDevice("DummyDevice", "Dummy Asset")

	device.AddNameplate("Dummy Manufacturer", "http://example.com/idlink", "12345",
		"Dummy Product", "v1.0", "SN123456")
	nicID := device.AddNic("eth0", "00:1A:2B:3C:4D:5E")
	device.AddIPv4(nicID, "192.168.1.100", "255.255.255.0", "192.168.1.1")
	device.AddSoftware("DummySoftware", "1.0.0", true)
	device.AddCapabilities("firmware_update", true)
	firmwareVersionKey := "firmware_version"
	firmwareVersionValue := "1.0.0"

	device.InstanceAnnotations = append(device.InstanceAnnotations, InstanceAnnotation{
		Key:   &firmwareVersionKey,
		Value: &firmwareVersionValue,
	})

	jsonMap, err := device.ConvertToJson()
	if err != nil {
		t.Fatalf("convertToJson failed: %v", err)
	}

	expectedLength := 12
	if len(jsonMap) != expectedLength {
		t.Fatalf("convertToJson should return %d keys, got: %d", expectedLength, len(jsonMap))
	}
	if _, ok := jsonMap["id"]; ok {
		t.Errorf("convertToJson should not return 'id' key")
	}
}

func TestConvertToJsonWithNilDevice(t *testing.T) {
	var device *DeviceInfo
	jsonMap, err := device.ConvertToJson()
	if err == nil {
		t.Fatalf("Expected an error for nil DeviceInfo, but got none")
	}
	if jsonMap != nil {
		t.Errorf("Expected nil map for nil DeviceInfo, but got: %v", jsonMap)
	}
}
