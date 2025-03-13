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
	baseSchemaVersion   = "v0.9.0"
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

	d.addManagementState()
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

func createTimestamp() time.Time {
	currentTime := time.Now().UTC()
	return currentTime
}

func getAssetContext() *AssetContext {
	return &AssetContext{
		Lis:       "http://rds.posccaesar.org/ontology/lis14/rdl/",
		Base:      baseSchemaInContext,
		Skos:      "http://www.w3.org/2004/02/skos/core#",
		Vocab:     baseSchemaInContext,
		Linkml:    "https://w3id.org/linkml/",
		SchemaOrg: "https://schema.org/",
	}
}
