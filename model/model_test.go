/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDevice(t *testing.T) {
	deviceInfo, err := NewDevice("testAsset", "test")
	assert.NoError(t, err)
	_, err = deviceInfo.AddNic("test", "ab:ab:ab:ab:ab:ab")
	assert.NoError(t, err)

	assert.Equal(t, "testAsset", deviceInfo.FunctionalObjectType)
	assert.Equal(t, "test", *deviceInfo.Name)
	if assert.Len(t, deviceInfo.AssetIdentifiers, 1) {
		macIdentifier, ok := deviceInfo.AssetIdentifiers[0].(MacIdentifier)
		if assert.True(t, ok) {
			assert.Equal(t, "ab:ab:ab:ab:ab:ab", macIdentifier.MacAddress)
			if assert.NotNil(t, macIdentifier.IdentifierUncertainty) {
				assert.Equal(t, 1, *macIdentifier.IdentifierUncertainty)
			}
		}
	}
	assert.Equal(t, FunctionalObjectSchemaUrl, deviceInfo.FunctionalObjectSchemaUrl)
}

func TestNewDevice_EmptyType(t *testing.T) {
	deviceInfo, err := NewDevice("", "test")
	assert.Error(t, err)
	assert.NotNil(t, deviceInfo)
	var ee *EmptyError
	if assert.ErrorAs(t, err, &ee) {
		assert.Equal(t, "FunctionalObjectType", ee.Field)
		assert.Equal(t, "Functional object type is required and cannot be empty", ee.Message)
		assert.Equal(t, "", ee.Value)
	}
}

func TestAddDescriptionValidDescription(t *testing.T) {
	deviceInfo, err := NewDevice("testAsset", "test")
	assert.NoError(t, err)
	description := "This is a test device"
	err = deviceInfo.AddDescription(description)
	assert.NoError(t, err)
	assert.NotNil(t, deviceInfo.Description)
	assert.Equal(t, description, *deviceInfo.Description)
}

func TestAddDescriptionEmptyDescription(t *testing.T) {
	deviceInfo, err := NewDevice("testAsset", "test")
	assert.NoError(t, err)
	// Set an initial description to check it doesn't get overwritten
	initialDescription := "Initial description"
	deviceInfo.Description = &initialDescription

	err = deviceInfo.AddDescription("")
	assert.Error(t, err)

	// Description should remain unchanged
	assert.NotNil(t, deviceInfo.Description)
	assert.Equal(t, initialDescription, *deviceInfo.Description)
}
