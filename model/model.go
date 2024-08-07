/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import (
	"time"
)

const (
	baseSchemaPrefix = "https://schema.industrial-assets.io/base/v0.7.5"
)

// NewDevice Generates a new asset skeleton
func NewDevice(typeOfAsset string, assetName string) *DeviceInfo {
	d := DeviceInfo{}
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

func createTimestamp() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC3339Nano)
}
