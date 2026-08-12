/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"encoding/json"
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

// TestConvertToJsonDeepNesting verifies that nested arrays and structs at multiple
// levels are fully preserved (not flattened or dropped).
func TestConvertToJsonDeepNesting(t *testing.T) {
	device, err := NewDevice("Device", "Deep Test")
	if err != nil {
		t.Fatalf("NewDevice failed: %v", err)
	}

	// 3-level struct: product_instance_information → manufacturer_product → manufacturer.name
	err = device.AddNameplate("Siemens AG", testIDLink, "S7-1500", "PLC", "2.0", "SN-42")
	if err != nil {
		t.Fatalf("AddNameplate failed: %v", err)
	}

	// array of structs with a nested array: connection_points → related_connection_points
	nicID, err := device.AddNic("eth0", "AA:BB:CC:DD:EE:FF")
	if err != nil {
		t.Fatalf("AddNic failed: %v", err)
	}
	_, err = device.AddIPv4(nicID, "10.0.0.1", "255.0.0.0", "10.0.0.254")
	if err != nil {
		t.Fatalf("AddIPv4 failed: %v", err)
	}

	m, err := device.ConvertToJson()
	if err != nil {
		t.Fatalf("ConvertToJson failed: %v", err)
	}

	// --- 3-level struct depth ---
	pii, ok := m["product_instance_information"].(map[string]interface{})
	if !ok {
		t.Fatalf("product_instance_information is not a map, got %T", m["product_instance_information"])
	}
	mp, ok := pii["manufacturer_product"].(map[string]interface{})
	if !ok {
		t.Fatalf("manufacturer_product is not a map, got %T", pii["manufacturer_product"])
	}
	mfr, ok := mp["manufacturer"].(map[string]interface{})
	if !ok {
		t.Fatalf("manufacturer is not a map, got %T", mp["manufacturer"])
	}
	if mfr["name"] != "Siemens AG" {
		t.Errorf("manufacturer.name: want %q, got %v", "Siemens AG", mfr["name"])
	}

	// --- array → struct → nested array ---
	cps, ok := m["connection_points"].([]interface{})
	if !ok {
		t.Fatalf("connection_points is not a []interface{}, got %T", m["connection_points"])
	}
	var ipv4 map[string]interface{}
	for _, cp := range cps {
		entry, ok := cp.(map[string]interface{})
		if ok && entry["connection_point_type"] == "Ipv4Connectivity" {
			ipv4 = entry
			break
		}
	}
	if ipv4 == nil {
		t.Fatal("no Ipv4Connectivity found in connection_points")
	}
	if ipv4["ipv4_address"] != "10.0.0.1" {
		t.Errorf("ipv4_address: want %q, got %v", "10.0.0.1", ipv4["ipv4_address"])
	}
	relatedRaw, ok := ipv4["related_connection_points"].([]interface{})
	if !ok || len(relatedRaw) == 0 {
		t.Fatalf("related_connection_points missing or empty, got %T %v", ipv4["related_connection_points"], ipv4["related_connection_points"])
	}
	rel, ok := relatedRaw[0].(map[string]interface{})
	if !ok {
		t.Fatalf("related_connection_points[0] is not a map, got %T", relatedRaw[0])
	}
	if rel["connection_point_id"] != nicID {
		t.Errorf("connection_point_id: want %q, got %v", nicID, rel["connection_point_id"])
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

func TestMarshalJSON(t *testing.T) {
	device, err := NewDevice("Asset", "Dummy Asset")
	if err != nil {
		t.Fatalf("NewDevice failed: %v", err)
	}

	err = device.AddDescription("Remote I/O module for distributed field device integration")
	if err != nil {
		t.Fatalf("AddDescription failed: %v", err)
	}

	jsonBytes, err := json.Marshal(device)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if payload["name"] != "Dummy Asset" {
		t.Fatalf("expected name to be marshaled, got: %v", payload["name"])
	}
	if payload["description"] != "Remote I/O module for distributed field device integration" {
		t.Fatalf("expected description to be marshaled, got: %v", payload["description"])
	}
}
