/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import (
	"github.com/google/uuid"
)

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
	timestamp := CreateTimestamp()
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

	timestamp := CreateTimestamp()
	state := ManagementStateValuesUnknown

	mgmtState := ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}

	d.ManagementState = mgmtState
}

func (d *DeviceInfo) addSoftwareComponent(assetName string) {

	softwareAsset := SoftwareAsset{
		AssetOperations:           nil,
		ConnectionPoints:          nil,
		CustomUiProperties:        nil,
		FunctionalParts:           nil,
		Id:                        uuid.New().String(),
		InstanceAnnotations:       []InstanceAnnotation{},
		ManagementState:           ManagementState{},
		Name:                      &assetName,
		OtherStates:               nil,
		ProductInstanceIdentifier: nil,
		ReachabilityState:         nil,
		SoftwareComponents:        nil,
	}
	softwareAsset.AddManagementState()
	softwareAsset.AddInstanceAnnotation("description", "IAH description")
	d.SoftwareComponents = append(d.SoftwareComponents, softwareAsset)
}

// Add Management state to the software component
// Only used internal
func (s *SoftwareAsset) AddManagementState() {
	timestamp := CreateTimestamp()
	state := ManagementStateValuesRegarded
	mgmtState := ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}
	s.ManagementState = mgmtState
}

func (s *SoftwareAsset) AddInstanceAnnotation(key string, value string) {
	instanceAnnotation := InstanceAnnotation{
		Key:   &key,
		Value: &value,
	}
	s.InstanceAnnotations = append(s.InstanceAnnotations, instanceAnnotation)
}

func (d *DeviceInfo) AddAssetOperations(operationName string, activationFlag bool) {
	assetOperation := AssetOperation{
		OperationName:  &operationName,
		ActivationFlag: &activationFlag,
	}
	d.AssetOperations = append(d.AssetOperations, assetOperation)
}

func (d *DeviceInfo) AddMetadata(key string, value string) {
	instanceAnnotation := InstanceAnnotation{
		Key:   &key,
		Value: &value,
	}
	d.InstanceAnnotations = append(d.InstanceAnnotations, instanceAnnotation)
}
