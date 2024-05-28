package model

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertToDiscoveredDevice(t *testing.T) {
	device := NewDevice("Profinet")
	timestamp := CreateTimestamp()

	Name := "Device"
	device.Name = &Name
	product := "test-dcd"
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

	connectionPointType := "Ipv4Connectivity"
	Ipv4Address := "192.168.0.1"
	Ipv4NetMask := "255.255.255.0"
	connectionPoint := "ethernet"
	relatedConnectionPoint := RelatedConnectionPoint{
		ConnectionPoint:    &connectionPoint,
		CustomRelationship: nil,
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
	device.ConnectionPoints = append(device.ConnectionPoints, Ipv4Connectivity)

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
	discoveredDevice := device.ConvertToDiscoveredDevice()
	assert.Equal(t, 17, len(discoveredDevice.Identifiers))
	assert.Equal(t, "URI", discoveredDevice.Identifiers[0].Classifiers[0].GetType())
	assert.Equal(t, "https://schema.industrial-assets.io/base/v0.7.5/Asset#@type", discoveredDevice.Identifiers[0].Classifiers[0].GetValue())
}
