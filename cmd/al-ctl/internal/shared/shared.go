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

	TimeoutSeconds float64
)

const (
	DiscoveryFileDesc           string = "discovery configuration filename allows the configuration of discovery filters and options (see discovery.json for an example)"
	DeviceIdentifierFileDesc    string = "file with device identifier/connection blob (see device_address.json for an example)"
	ConvertDeviceIdentifierDesc string = "convert device identifier in base64 encoding (required if the device identifier is not yet encoded)"
	DeviceCredentialsFileDesc   string = "file with device credentials (see device_credentials_admin.json and device_credentials_user.json for examples)"
	ArtefactCredentialsFileDesc string = "file with artefact credentials"
	JobIdDesc                   string = "unique ID that identifiers the job"
)
