/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package shared

var (
	RegistryEndpoint  string
	AssetLinkEndpoint string

	TimeoutSeconds uint
)

const (
	DiscoveryFileDesc           string = "discovery configuration filename allows the configuration of discovery filters and options (see discovery.json for an example)"
	DeviceIdentifierFileDesc    string = "file with device identifier/connection blob (see device-identifier.json for an example)"
	ConvertDeviceIdentifierDesc string = "convert device identifier in base64 encoding (required if the device identifier is not yet encoded)"
)
