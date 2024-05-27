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
)

const (
	baseSchemaPrefix = "https://schema.industrial-assets.io/base/v0.7.5"
)

func NewDevice() *generated.DiscoveredDevice {
	name := "Example Device"
	serialNumber := uuid.New().String()
	articelNumber := "test-article-number"
	var timestamp uint64 = 133344110897340000
	device := generated.DiscoveredDevice{
		Identifiers: []*generated.DeviceIdentifier{
			{
				Value: &generated.DeviceIdentifier_Text{Text: "Siemens AG"},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#product_instance_identifier/manufacturer_product/manufacturer/name", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Children{
					Children: &generated.DeviceIdentifierValueList{
						Value: []*generated.DeviceIdentifier{
							{
								Value: &generated.DeviceIdentifier_Text{
									Text: "30:13:89:1E:C7:61",
								},
								Classifiers: []*generated.SemanticClassifier{
									{
										Type:  "URI",
										Value: fmt.Sprintf("%s/Asset#mac_identifiers/mac_address", baseSchemaPrefix),
									},
								},
							},
						},
					},
				},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#mac_identifiers", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Text{Text: articelNumber},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#product_instance_identifier/manufacturer_product/product_id", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Text{Text: name},
				Classifiers: []*generated.SemanticClassifier{
					{
						Type:  "URI",
						Value: fmt.Sprintf("%s/Asset#name", baseSchemaPrefix),
					},
				},
			},
			{
				Value: &generated.DeviceIdentifier_Text{Text: serialNumber},
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
									Text: "0_Ethernet",
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
									Text: "uuid:40ead537-6faa-4a38-beb3-f55b34578ats",
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
									Text: "EthernetPort",
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
									Text: "30:13:89:1E:C7:72",
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
									Text: "EthernetPort",
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
									Text: "uuid:40ead537-6faa-4a38-beb3-f55b3123456s",
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
		Timestamp:              timestamp,
	}
	return &device
}
