/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConvertToDiscoveredDevice(t *testing.T) {
	device := generateDevice("Profinet", "Device")
	discoveredDevice := device.ConvertToDiscoveredDevice()
	discoveredDeviceType := fmt.Sprintf("%s/%s", baseSchemaPrefix, "Asset#@type")
	assert.Equal(t, 16, len(discoveredDevice.Identifiers))
	assert.Equal(t, "URI", discoveredDevice.Identifiers[0].Classifiers[0].GetType())
	assert.Equal(t, discoveredDeviceType, discoveredDevice.Identifiers[0].Classifiers[0].GetValue())
}

func TestConvertFromDerivedSchemaToDiscoveredDevice(t *testing.T) {
	schemaUri := "https://schema.industrial-assets.io/sat/v0.8.2"
	device := generateDevice("SatController", "Device")
	discoveredDevice := ConvertFromDerivedSchemaToDiscoveredDevice(device, schemaUri, "SatController")
	assert.Equal(t, 16, len(discoveredDevice.Identifiers))
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
	assert.Equal(t, 17, len(discoveredDevice.Identifiers))
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
	timestamp := createTimestamp()
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
