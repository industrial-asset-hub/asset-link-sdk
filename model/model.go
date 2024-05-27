/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import (
	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/iah-discovery"
	"fmt"
	"github.com/google/uuid"
	"time"
)

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

func CreateTimestamp() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC3339Nano)
}

func (d *DeviceInfo) ConvertToDiscoveredDevice() *generated.DiscoveredDevice {
	device := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{
			{
				Value: &generated.DeviceIdentifier_Text{Text: *d.ProductInstanceIdentifier.ManufacturerProduct.Name},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#product_instance_identifier/manufacturer_product/manufacturer/name", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Text{Text: *d.ProductInstanceIdentifier.ManufacturerProduct.ProductId},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#product_instance_identifier/manufacturer_product/product_id", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Text{Text: *d.Name},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#name", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Text{Text: uuid.New().String()},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#product_instance_identifier/serial_number", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Children{
					Children: &generated.DeviceIdentifierValueList{
						Value: []*generated.DeviceIdentifier{
							{
								Value: &generated.DeviceIdentifier_Text{
									Text: *d.ConnectionPoints[0].RelatedConnectionPoints[0].ConnectionPoint,
								},
								Classifiers: []*generated.SemanticClassifier{
									{
										Type:  "URI",
										Value: fmt.Sprintf("%s/Asset#connection_points/related_connection_points/connection_point", baseSchemaPrefix),
									},
								},
							},
						},
					},
				},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#connection_points", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Children{
					Children: &generated.DeviceIdentifierValueList{
						Value: []*generated.DeviceIdentifier{
							{
								Value: &generated.DeviceIdentifier_Text{
									Text: uuid.New().String(),
								},
								Classifiers: []*generated.SemanticClassifier{
									{
										Type:  "URI",
										Value: fmt.Sprintf("%s/Asset#connection_points/id", baseSchemaPrefix),
									},
								},
							},
							{
								Value: &generated.DeviceIdentifier_Text{
									Text: *d.ConnectionPoints[0].ConnectionPointType,
								},
								Classifiers: []*generated.SemanticClassifier{
									{
										Type:  "URI",
										Value: fmt.Sprintf("%s/Asset#connection_points/connection_point_type", baseSchemaPrefix),
									},
								},
							},
						},
					},
				},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#connection_points", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Children{
					Children: &generated.DeviceIdentifierValueList{
						Value: []*generated.DeviceIdentifier{
							{
								Value: &generated.DeviceIdentifier_Text{
									Text: *d.ConnectionPoints[0].Ipv4Address,
								},
								Classifiers: []*generated.SemanticClassifier{
									{
										Type:  "URI",
										Value: fmt.Sprintf("%s/Asset#connection_points/mac_address", baseSchemaPrefix),
									},
								},
							},
							{
								Value: &generated.DeviceIdentifier_Text{
									Text: *d.ConnectionPoints[0].ConnectionPointType,
								},
								Classifiers: []*generated.SemanticClassifier{
									{
										Type:  "URI",
										Value: fmt.Sprintf("%s/Asset#connection_points/connection_point_type", baseSchemaPrefix),
									},
								},
							},
							{
								Value: &generated.DeviceIdentifier_Text{
									Text: uuid.New().String(),
								},
								Classifiers: []*generated.SemanticClassifier{
									{
										Type:  "URI",
										Value: fmt.Sprintf("%s/Asset#connection_points/id", baseSchemaPrefix),
									},
								},
							},
						},
					},
				},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#connection_points", baseSchemaPrefix),
					},
				},
			},
		},
		ConnectionParameterSet: nil,
		Timestamp:              19347439483904,
	}
	return &device

}
