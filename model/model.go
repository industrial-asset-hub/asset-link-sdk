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
	baseSchemaPrefix = "https://schema.industrial-assets.io/base/v0.9.0"
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

	d.addManagementState()
	d.addReachabilityState()

	return &d
}

type DeviceInfo struct {
	Type string `json:"@type"`
	// Override connection point, since generated base schema does not provide derived types
	ConnectionPoints []any `json:"connection_points,omitempty"`
	Asset
	MacIdentifiers []MacIdentifier `json:"mac_identifiers"`
	// To Be clarified
	SoftwareComponents []any `json:"software_components,omitempty"`
}

func createTimestamp() time.Time {
	currentTime := time.Now().UTC()
	return currentTime
}
