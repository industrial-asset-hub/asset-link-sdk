/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCustomIdentifier(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddCustomIdentifier("ExternalID", "id:source-123")
		if assert.Len(t, device.AssetIdentifiers, 1) {
			identifier, ok := device.AssetIdentifiers[0].(CustomIdentifier)
			if assert.True(t, ok) {
				assert.Equal(t, CustomIdentifierAssetIdentifierTypeCustomIdentifier, identifier.AssetIdentifierType)
				assert.Equal(t, "ExternalID", identifier.Name)
				assert.Equal(t, "id:source-123", identifier.Value)
			}
		}
	})

	t.Run("empty name is ignored", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddCustomIdentifier("", "id:source-123")
		assert.Empty(t, device.AssetIdentifiers)
	})

	t.Run("invalid value pattern is ignored", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddCustomIdentifier("ExternalID", "invalid value with spaces")
		assert.Empty(t, device.AssetIdentifiers)
	})
}

func TestAddCertificateIdentifier(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddCertificateIdentifier("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A")
		if assert.Len(t, device.AssetIdentifiers, 1) {
			identifier, ok := device.AssetIdentifiers[0].(CertificateIdentifier)
			if assert.True(t, ok) {
				assert.Equal(t, CertificateIdentifierAssetIdentifierTypeCertificateIdentifier, identifier.AssetIdentifierType)
				assert.Equal(t, "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A", identifier.CertificateId)
			}
		}
	})

	t.Run("empty certificate id is ignored", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddCertificateIdentifier("")
		assert.Empty(t, device.AssetIdentifiers)
	})
}

func TestAddHostBasedSoftwareIdentifier(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddHostBasedSoftwareIdentifier("MyApp", "1.2.3", "host-001")
		if assert.Len(t, device.AssetIdentifiers, 1) {
			id, ok := device.AssetIdentifiers[0].(HostBasedSoftwareIdentifier)
			if assert.True(t, ok) {
				assert.Equal(t, HostBasedSoftwareIdentifierAssetIdentifierTypeHostBasedSoftwareIdentifier, id.AssetIdentifierType)
				assert.Equal(t, "MyApp", id.Name)
				assert.Equal(t, "1.2.3", id.Version)
				assert.Equal(t, "host-001", id.HostIdentifier)
			}
		}
	})

	t.Run("empty name is ignored", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddHostBasedSoftwareIdentifier("", "1.2.3", "host-001")
		assert.Empty(t, device.AssetIdentifiers)
	})

	t.Run("empty version is ignored", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddHostBasedSoftwareIdentifier("MyApp", "", "host-001")
		assert.Empty(t, device.AssetIdentifiers)
	})

	t.Run("empty host identifier is ignored", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddHostBasedSoftwareIdentifier("MyApp", "1.2.3", "")
		assert.Empty(t, device.AssetIdentifiers)
	})
}

func TestAddProductInstanceIdentifier(t *testing.T) {
	t.Run("success with all fields", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddProductInstanceIdentifier("Siemens AG", "6GK5208-0HA00-2AS6", "SN-12345")
		if assert.Len(t, device.AssetIdentifiers, 1) {
			id, ok := device.AssetIdentifiers[0].(ProductInstanceIdentifier)
			if assert.True(t, ok) {
				assert.Equal(t, ProductInstanceIdentifierAssetIdentifierTypeProductInstanceIdentifier, id.AssetIdentifierType)
				assert.Equal(t, "Siemens AG", id.Vendor)
				assert.Equal(t, "6GK5208-0HA00-2AS6", id.ArticleNumber)
				assert.Equal(t, "SN-12345", id.SerialNumber)
			}
		}
	})

	t.Run("empty vendor is accepted", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddProductInstanceIdentifier("", "6GK5208-0HA00-2AS6", "SN-12345")
		assert.Len(t, device.AssetIdentifiers, 1)
	})

	t.Run("empty article number is accepted", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddProductInstanceIdentifier("Siemens AG", "", "SN-12345")
		assert.Len(t, device.AssetIdentifiers, 1)
	})

	t.Run("empty serial number is accepted", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddProductInstanceIdentifier("Siemens AG", "6GK5208-0HA00-2AS6", "")
		assert.Len(t, device.AssetIdentifiers, 1)
	})

	t.Run("all empty is ignored", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddProductInstanceIdentifier("", "", "")
		assert.Empty(t, device.AssetIdentifiers)
	})
}

func intPtr(v int) *int { return &v }

func TestAddParentRelativeIdentifier(t *testing.T) {
	t.Run("success with MacIdentifier parent", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		parent := MacIdentifier{
			AssetIdentifierType: MacIdentifierAssetIdentifierTypeMacIdentifier,
			MacAddress:          "AA:BB:CC:DD:EE:FF",
		}
		device.AddParentRelativeIdentifier(parent, intPtr(1), intPtr(2))
		if assert.Len(t, device.AssetIdentifiers, 1) {
			id, ok := device.AssetIdentifiers[0].(ParentRelativeIdentifier)
			if assert.True(t, ok) {
				assert.Equal(t, ParentRelativeIdentifierAssetIdentifierTypeParentRelativeIdentifier, id.AssetIdentifierType)
				assert.Equal(t, parent, id.ParentIdentifier)
				assert.Equal(t, 1, id.Slot)
				assert.Equal(t, 2, id.Subslot)
			}
		}
	})

	t.Run("success with IdLinkIdentifier parent", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		parent := IdLinkIdentifier{
			AssetIdentifierType: IdLinkIdentifierAssetIdentifierTypeIdLinkIdentifier,
			IdLink:              "https://example.com/product/123",
		}
		device.AddParentRelativeIdentifier(parent, intPtr(0), intPtr(0))
		if assert.Len(t, device.AssetIdentifiers, 1) {
			id, ok := device.AssetIdentifiers[0].(ParentRelativeIdentifier)
			if assert.True(t, ok) {
				assert.Equal(t, parent, id.ParentIdentifier)
				assert.Equal(t, 0, id.Slot)
				assert.Equal(t, 0, id.Subslot)
			}
		}
	})

	t.Run("nil slot and subslot default to 0", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		parent := MacIdentifier{
			AssetIdentifierType: MacIdentifierAssetIdentifierTypeMacIdentifier,
			MacAddress:          "AA:BB:CC:DD:EE:FF",
		}
		device.AddParentRelativeIdentifier(parent, nil, nil)
		if assert.Len(t, device.AssetIdentifiers, 1) {
			id, ok := device.AssetIdentifiers[0].(ParentRelativeIdentifier)
			if assert.True(t, ok) {
				assert.Equal(t, 0, id.Slot)
				assert.Equal(t, 0, id.Subslot)
			}
		}
	})

	t.Run("nil parent is ignored", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.AddParentRelativeIdentifier(nil, intPtr(1), intPtr(2))
		assert.Empty(t, device.AssetIdentifiers)
	})
}

func TestDeviceInfoAddAssetRelation(t *testing.T) {
	for _, functionalType := range []string{"Asset", "Device", "Gateway", "SoftwareArtifact"} {
		functionalType := functionalType
		t.Run(functionalType+" success non-bidirectional with MacIdentifier", func(t *testing.T) {
			deviceInfo, err := NewDevice(functionalType, "TestDevice")
			assert.NoError(t, err)

			relatedAsset := RelatedAsset{AssetIdentifiers: []interface{}{
				MacIdentifier{
					AssetIdentifierType: MacIdentifierAssetIdentifierTypeMacIdentifier,
					MacAddress:          "AA:BB:CC:DD:EE:FF",
				},
			}}
			err = deviceInfo.AddAssetRelation("is_part_of", relatedAsset, RelationalRoleOfRelatedAssetValuesObject, false)
			assert.NoError(t, err)

			if assert.Len(t, deviceInfo.AssetRelations, 1) {
				relation := deviceInfo.AssetRelations[0]
				assert.Equal(t, "is_part_of", relation.Predicate)
				assert.Equal(t, relatedAsset, relation.RelatedAsset)
				assert.Equal(t, RelationalRoleOfRelatedAssetValuesObject, relation.RelationalRoleOfRelatedAsset)
				assert.Nil(t, relation.IsBidirectional)
			}
		})

		t.Run(functionalType+" success bidirectional with CustomIdentifier", func(t *testing.T) {
			deviceInfo, err := NewDevice(functionalType, "TestDevice")
			assert.NoError(t, err)

			relatedAsset := RelatedAsset{AssetIdentifiers: []interface{}{
				CustomIdentifier{
					AssetIdentifierType: CustomIdentifierAssetIdentifierTypeCustomIdentifier,
					Name:                "CustomID",
					Value:               "custom-value-123",
				},
			}}
			err = deviceInfo.AddAssetRelation("is_connected_to", relatedAsset, RelationalRoleOfRelatedAssetValuesSubject, true)
			assert.NoError(t, err)

			if assert.Len(t, deviceInfo.AssetRelations, 1) {
				relation := deviceInfo.AssetRelations[0]
				assert.Equal(t, "is_connected_to", relation.Predicate)
				assert.Equal(t, relatedAsset, relation.RelatedAsset)
				assert.Equal(t, RelationalRoleOfRelatedAssetValuesSubject, relation.RelationalRoleOfRelatedAsset)
				if assert.NotNil(t, relation.IsBidirectional) {
					assert.True(t, *relation.IsBidirectional)
				}
			}
		})

		t.Run(functionalType+" empty predicate returns error", func(t *testing.T) {
			deviceInfo, err := NewDevice(functionalType, "TestDevice")
			assert.NoError(t, err)

			relatedAsset := RelatedAsset{AssetIdentifiers: []interface{}{
				MacIdentifier{
					AssetIdentifierType: MacIdentifierAssetIdentifierTypeMacIdentifier,
					MacAddress:          "AA:BB:CC:DD:EE:FF",
				},
			}}
			err = deviceInfo.AddAssetRelation("", relatedAsset, RelationalRoleOfRelatedAssetValuesObject, false)
			assert.Error(t, err)
			assert.Empty(t, deviceInfo.AssetRelations)
		})

		t.Run(functionalType+" invalid predicate format returns error", func(t *testing.T) {
			deviceInfo, err := NewDevice(functionalType, "TestDevice")
			assert.NoError(t, err)

			relatedAsset := RelatedAsset{AssetIdentifiers: []interface{}{
				MacIdentifier{
					AssetIdentifierType: MacIdentifierAssetIdentifierTypeMacIdentifier,
					MacAddress:          "AA:BB:CC:DD:EE:FF",
				},
			}}
			err = deviceInfo.AddAssetRelation("InvalidPredicate", relatedAsset, RelationalRoleOfRelatedAssetValuesObject, false)
			assert.Error(t, err)
			assert.Empty(t, deviceInfo.AssetRelations)
		})

		t.Run(functionalType+" invalid MAC address in identifier returns error", func(t *testing.T) {
			deviceInfo, err := NewDevice(functionalType, "TestDevice")
			assert.NoError(t, err)

			relatedAsset := RelatedAsset{AssetIdentifiers: []interface{}{
				MacIdentifier{
					AssetIdentifierType: MacIdentifierAssetIdentifierTypeMacIdentifier,
					MacAddress:          "INVALID_MAC",
				},
			}}
			err = deviceInfo.AddAssetRelation("is_part_of", relatedAsset, RelationalRoleOfRelatedAssetValuesObject, false)
			assert.Error(t, err)
			assert.Empty(t, deviceInfo.AssetRelations)
		})

		t.Run(functionalType+" empty MAC address in identifier returns error", func(t *testing.T) {
			deviceInfo, err := NewDevice(functionalType, "TestDevice")
			assert.NoError(t, err)

			relatedAsset := RelatedAsset{AssetIdentifiers: []interface{}{
				MacIdentifier{
					AssetIdentifierType: MacIdentifierAssetIdentifierTypeMacIdentifier,
					MacAddress:          "",
				},
			}}
			err = deviceInfo.AddAssetRelation("is_part_of", relatedAsset, RelationalRoleOfRelatedAssetValuesObject, false)
			assert.Error(t, err)
			assert.Empty(t, deviceInfo.AssetRelations)
		})

		t.Run(functionalType+" invalid custom identifier value returns error", func(t *testing.T) {
			deviceInfo, err := NewDevice(functionalType, "TestDevice")
			assert.NoError(t, err)

			// CustomIdentifierValuePattern: ^[A-Za-z0-9._~!$&'()*+,;=:/?@%-]{1,256}$
			// Use invalid characters like spaces or brackets
			relatedAsset := RelatedAsset{AssetIdentifiers: []interface{}{
				CustomIdentifier{
					AssetIdentifierType: CustomIdentifierAssetIdentifierTypeCustomIdentifier,
					Name:                "CustomID",
					Value:               "invalid value with spaces",
				},
			}}
			err = deviceInfo.AddAssetRelation("is_part_of", relatedAsset, RelationalRoleOfRelatedAssetValuesObject, false)
			assert.Error(t, err)
			assert.Empty(t, deviceInfo.AssetRelations)
		})

		t.Run(functionalType+" empty identifiers list returns error", func(t *testing.T) {
			deviceInfo, err := NewDevice(functionalType, "TestDevice")
			assert.NoError(t, err)

			relatedAsset := RelatedAsset{AssetIdentifiers: []interface{}{}}
			err = deviceInfo.AddAssetRelation("is_part_of", relatedAsset, RelationalRoleOfRelatedAssetValuesObject, false)
			assert.Error(t, err)
			assert.Empty(t, deviceInfo.AssetRelations)
		})
	}
}
