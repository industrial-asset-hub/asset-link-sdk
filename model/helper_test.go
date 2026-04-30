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

		device.addCustomIdentifier("ExternalID", "id:source-123")
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

		device.addCustomIdentifier("", "id:source-123")
		assert.Empty(t, device.AssetIdentifiers)
	})

	t.Run("invalid value pattern is ignored", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.addCustomIdentifier("ExternalID", "invalid value with spaces")
		assert.Empty(t, device.AssetIdentifiers)
	})
}

func TestAddCertificateIdentifier(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		device, err := NewDevice("Asset", "TestDevice")
		assert.NoError(t, err)

		device.addCertificateIdentifier("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A")
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

		device.addCertificateIdentifier("")
		assert.Empty(t, device.AssetIdentifiers)
	})
}
