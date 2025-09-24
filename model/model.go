/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"time"

	"github.com/rs/zerolog/log"
)

const (
	baseSchemaVersion   = "v0.12.0"
	baseSchemaPrefix    = "https://schema.industrial-assets.io/base/" + baseSchemaVersion
	baseSchemaInContext = "https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/" + baseSchemaVersion + "/"

	gateway = "Gateway"
)

// NewDevice Generates a new asset skeleton
func NewDevice(typeOfAsset string, assetName string) *DeviceInfo {

	d := DeviceInfo{}
	if !isNonEmptyValues(typeOfAsset) {
		log.Warn().Msg("Asset type is empty")
		return &d
	}

	d.Type = typeOfAsset
	d.Name = &assetName
	d.Context = getAssetContext()

	d.AddManagementState(ManagementStateValuesUnknown)
	d.addReachabilityState()

	return &d
}

// NewGateway Generates a new gateway skeleton
func NewGateway(gatewayName string) *GatewayInfo {

	d := GatewayInfo{}

	d.Type = gateway
	d.Name = &gatewayName
	d.Context = getAssetContext()

	d.addManagementState(ManagementStateValuesRegarded)

	return &d
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
func (d *DeviceInfo) AddManagementState(stateValue ManagementStateValues) {
	mgmtState := managementStatePtr(stateValue, getAssetCreationTimestamp(d.ManagementState.StateTimestamp))
	if mgmtState == nil {
		return
	}

	d.ManagementState = *mgmtState
}

func (d *GatewayInfo) addManagementState(stateValue ManagementStateValues) {
	mgmtState := managementStatePtr(stateValue, getAssetCreationTimestamp(d.ManagementState.StateTimestamp))
	if mgmtState == nil {
		return
	}

	d.ManagementState = *mgmtState
}

func managementStatePtr(stateValue ManagementStateValues, timestamp time.Time) *ManagementState {
	if !isNonEmptyValues(string(stateValue)) {
		log.Warn().Msg("Management state value is empty")
		return nil
	}
	if stateValue != ManagementStateValuesIgnored && stateValue != ManagementStateValuesRegarded && stateValue != ManagementStateValuesUnknown {
		log.Warn().Msgf("Management state value %s is not valid", stateValue)
		return nil
	}

	mgmtState := ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &stateValue,
	}

	return &mgmtState
}

func (d *DeviceInfo) AddDescription(description string) {

	if !isNonEmptyValues(description) {
		log.Warn().Msg("Description is empty")
		return
	}

	d.Description = &description
}
