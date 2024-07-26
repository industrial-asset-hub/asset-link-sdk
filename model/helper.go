/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import "github.com/google/uuid"

func (d *DeviceInfo) addIdentifier(mac string) {

	identifierUncertainty := 1
	identifier := MacIdentifier{
		IdentifierType:        nil,
		IdentifierUncertainty: &identifierUncertainty,
		MacAddress:            &mac,
	}
	d.MacIdentifiers = append(d.MacIdentifiers, identifier)
}

// Add reachability state to the asset
func (d *DeviceInfo) addReachabilityState() {
	timestamp := createTimestamp()
	state := ReachabilityStateValuesReached

	reachabilityState := ReachabilityState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}

	d.ReachabilityState = &reachabilityState
}

// Add Management state to the asset
// Only used internal
func (d *DeviceInfo) addManagementState() {

	timestamp := createTimestamp()
	state := ManagementStateValuesUnknown

	mgmtState := ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}

	d.ManagementState = mgmtState
}

func (d *DeviceInfo) addSoftwareComponent(annotations map[string]string) {
	softwareAsset := SoftwareAsset{
		AssetOperations:           nil,
		ConnectionPoints:          nil,
		CustomUiProperties:        nil,
		FunctionalParts:           nil,
		Id:                        uuid.New().String(),
		InstanceAnnotations:       nil,
		ManagementState:           ManagementState{},
		Name:                      nil,
		OtherStates:               nil,
		ProductInstanceIdentifier: nil,
		ReachabilityState:         nil,
		SoftwareComponents:        nil,
	}
	d.SoftwareComponents = append(d.SoftwareComponents, softwareAsset)
}

func (d *DeviceInfo) AddOneMetadata(key string, value string) {
	instanceAnnotation := InstanceAnnotation{
		Key:   &key,
		Value: &value,
	}
	d.InstanceAnnotations = append(d.InstanceAnnotations, instanceAnnotation)
}
