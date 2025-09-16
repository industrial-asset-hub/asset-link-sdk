/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"github.com/rs/zerolog/log"
)

const (
	baseSchemaVersion   = "v0.10.0"
	baseSchemaPrefix    = "https://schema.industrial-assets.io/base/" + baseSchemaVersion
	baseSchemaInContext = "https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/" + baseSchemaVersion + "/"
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

	if !isNonEmptyValues(string(stateValue)) {
		log.Warn().Msg("Management state value is empty")
		return
	}
	if stateValue != ManagementStateValuesIgnored && stateValue != ManagementStateValuesRegarded && stateValue != ManagementStateValuesUnknown {
		log.Warn().Msgf("Management state value %s is not valid", stateValue)
		return
	}
	timestamp := d.getAssetCreationTimestamp()
	state := stateValue

	mgmtState := ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}

	d.ManagementState = mgmtState
}
