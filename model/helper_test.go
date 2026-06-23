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
