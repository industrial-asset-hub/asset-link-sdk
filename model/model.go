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

func NewDevice(typeOfAsset string) *DeviceInfo {
	d := DeviceInfo{}
	d.Type = typeOfAsset
	return &d
}

type DeviceInfo struct {
	Type string `json:"@type"`
	// Override connection point, since generated base schema does not provide derived types
	ConnectionPoints []Ipv4Connectivity `json:"connection_points,omitempty"`
	Asset
	MacIdentifiers []MacIdentifier `json:"mac_identifiers"`
}

func CreateTimestamp() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC3339Nano)
}
