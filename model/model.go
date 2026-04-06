/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"time"
)

const (
	baseSchemaVersion   = "v0.12.0"
	baseSchemaPrefix    = "https://schema.industrial-assets.io/base/" + baseSchemaVersion
	baseSchemaInContext = "https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/" + baseSchemaVersion + "/"

	gateway = "Gateway"
)

// NewDevice Generates a new asset skeleton
func NewDevice(typeOfAsset string, assetName string) (*DeviceInfo, error) {

	d := DeviceInfo{}
	if !isNonEmptyValues(typeOfAsset) {
		err := &EmptyError{
			Field:   "Type",
			Message: "Asset type is required and cannot be empty",
			Value:   typeOfAsset,
		}
		return &d, err
	}

	d.Type = typeOfAsset
	d.Name = &assetName
	d.Context = getAssetContext()

	err := d.AddManagementState(ManagementStateValuesUnknown)
	if err != nil {
		return nil, err
	}
	d.addReachabilityState()

	return &d, nil
}

// NewGateway Generates a new gateway skeleton
func NewGateway(gatewayName string) (*GatewayInfo, error) {

	d := GatewayInfo{}

	d.Type = gateway
	d.Name = &gatewayName
	d.Context = getAssetContext()

	err := d.addManagementState(ManagementStateValuesRegarded)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

type DeviceInfo struct {
	Type    string        `json:"@type"`
	Context *AssetContext `json:"@context,omitempty"`
	// Override connection point, since generated base schema does not provide derived types
	ConnectionPoints []any `json:"connection_points,omitempty"`
	Asset
	MacIdentifiers []MacIdentifier `json:"mac_identifiers"`
	// To Be clarified
	SoftwareComponents []any `json:"software_components,omitempty"`
}

type GatewayInfo struct {
	Type    string        `json:"@type"`
	Context *AssetContext `json:"@context,omitempty"`
	Gateway
}

type AssetContext struct {
	Lis       string `json:"lis"`
	Base      string `json:"base"`
	Skos      string `json:"skos"`
	Vocab     string `json:"@vocab"`
	Linkml    string `json:"linkml"`
	SchemaOrg string `json:"schemaorg"`
}

// Add Management state to the asset
func (d *DeviceInfo) AddManagementState(stateValue ManagementStateValues) error {
	mgmtState, err := managementStatePtr(stateValue, getAssetCreationTimestamp(d.ManagementState.StateTimestamp))
	if mgmtState == nil {
		return err
	}
	d.ManagementState = *mgmtState
	return nil
}

func (d *GatewayInfo) addManagementState(stateValue ManagementStateValues) error {
	mgmtState, err := managementStatePtr(stateValue, getAssetCreationTimestamp(d.ManagementState.StateTimestamp))
	if mgmtState == nil {
		return err
	}

	d.ManagementState = *mgmtState
	return nil
}

func managementStatePtr(stateValue ManagementStateValues, timestamp time.Time) (*ManagementState, error) {
	if !isNonEmptyValues(string(stateValue)) {
		err := &EmptyError{
			Field:   "ManagementState",
			Message: "Management state value is empty",
			Value:   stateValue,
		}
		return nil, err
	}
	if stateValue != ManagementStateValuesIgnored && stateValue != ManagementStateValuesRegarded && stateValue != ManagementStateValuesUnknown {
		err := &PermissibleValuesError{
			Field:   "ManagementState",
			Value:   stateValue,
			Allowed: []interface{}{ManagementStateValuesIgnored, ManagementStateValuesRegarded, ManagementStateValuesUnknown},
		}
		return nil, err
	}

	mgmtState := ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &stateValue,
	}

	return &mgmtState, nil
}

func (d *DeviceInfo) AddDescription(description string) error {

	if !isNonEmptyValues(description) {
		err := &EmptyError{
			Field:   "Description",
			Message: "Description is empty",
		}
		return err
	}

	d.Description = &description
	return nil
}
