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
	device, err := NewDevice("Device", "Dummy Asset")
	if err != nil {
		t.Fatalf("NewDevice failed: %v", err)
	}

	err = device.AddNameplate("Dummy Manufacturer", testIDLink, "12345",
		"Dummy Product", "v1.0", "SN123456")
	if err != nil {
		t.Fatalf("AddNameplate failed: %v", err)
	}
	nicID, err := device.AddNic("eth0", "00:1A:2B:3C:4D:5E")
	if err != nil {
		t.Fatalf("AddNic failed: %v", err)
	}
	_, err = device.AddIPv4(nicID, "192.168.1.100", "255.255.255.0", "192.168.1.1")
	if err != nil {
		t.Fatalf("AddIPv4 failed: %v", err)
	}
	err = device.AddSoftwareArtifactComponent("DummySoftware", "1.0.0", true)
	if err != nil {
		t.Fatalf("AddSoftwareArtifactComponent failed: %v", err)
	}
	err = device.AddCapabilities("firmware_update", true)
	if err != nil {
		t.Fatalf("AddCapabilities failed: %v", err)
	}
	firmwareVersionKey := "firmware_version"
	firmwareVersionValue := "1.0.0"

	device.InstanceAnnotations = append(device.InstanceAnnotations, InstanceAnnotation{
		Key:   &firmwareVersionKey,
		Value: &firmwareVersionValue,
	})

	jsonMap, err := device.ConvertToJson()
	if err != nil {
		t.Fatalf("ConvertToJson failed: %v", err)
	}

	expectedLength := 9
	if len(jsonMap) != expectedLength {
		t.Fatalf("ConvertToJson should return %d keys, got: %d", expectedLength, len(jsonMap))
	}
	if _, ok := jsonMap["id"]; ok {
		t.Errorf("ConvertToJson should not return 'id' key")
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
