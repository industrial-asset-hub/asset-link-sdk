/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

const (
	baseSchemaPrefix = "https://schema.industrial-assets.io/base/v0.7.5"
)

func NewDevice(typeOfAsset string) *DeviceInfo {
	d := DeviceInfo{}
	d.Type = typeOfAsset
	d.Context = map[string]interface{}{
		"base":      "https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/v0.7.5/",
		"linkml":    "https://w3id.org/linkml/",
		"lis":       "http://rds.posccaesar.org/ontology/lis14/rdl/",
		"schemaorg": "https://schema.org/",
		"skos":      "http://www.w3.org/2004/02/skos/core#",
		"@vocab":    "https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/v0.7.5/",
	}
	return &d
}

type DeviceInfo struct {
	Type    string                 `json:"@type"`
	Context map[string]interface{} `json:"@context"`
	// Override connection point, since generated base schema does not provide derived types
	ConnectionPoints []Ipv4Connectivity `json:"connection_points,omitempty"`
	Asset
}
