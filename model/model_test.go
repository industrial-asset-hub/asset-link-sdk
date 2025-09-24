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
	deviceInfo := NewDevice("testAsset", "test")
	deviceInfo.AddNic("test", "ab:ab:ab:ab:ab:ab")

	assert.Equal(t, "testAsset", deviceInfo.Type)
	assert.Equal(t, "test", *deviceInfo.Name)
	assert.Equal(t, "ab:ab:ab:ab:ab:ab", *deviceInfo.MacIdentifiers[0].MacAddress)
	assert.Equal(t, 1, *deviceInfo.MacIdentifiers[0].IdentifierUncertainty)
	assert.Equal(t, ReachabilityStateValuesReached, *deviceInfo.ReachabilityState.StateValue)
	assert.Equal(t, ManagementStateValuesUnknown, *deviceInfo.ManagementState.StateValue)
	assert.Equal(t, getAssetContext(), deviceInfo.Context)
}
func TestAddManagementStateValidStates(t *testing.T) {
	deviceInfo := NewDevice("testAsset", "test")

	deviceInfo.AddManagementState(ManagementStateValuesUnknown)
	assert.Equal(t, ManagementStateValuesUnknown, *deviceInfo.ManagementState.StateValue)

	deviceInfo.AddManagementState(ManagementStateValuesIgnored)
	assert.Equal(t, ManagementStateValuesIgnored, *deviceInfo.ManagementState.StateValue)

	deviceInfo.AddManagementState(ManagementStateValuesRegarded)
	assert.Equal(t, ManagementStateValuesRegarded, *deviceInfo.ManagementState.StateValue)
}

func TestAddManagementStateEmptyState(t *testing.T) {
	deviceInfo := NewDevice("testAsset", "test")
	// Save previous state
	prevState := deviceInfo.ManagementState

	deviceInfo.AddManagementState("")
	// Should not update state
	assert.Equal(t, prevState, deviceInfo.ManagementState)
}

func TestAddManagementStateInvalidState(t *testing.T) {
	deviceInfo := NewDevice("testAsset", "test")
	// Save previous state
	prevState := deviceInfo.ManagementState

	deviceInfo.AddManagementState("invalid_state")
	// Should not update state
	assert.Equal(t, prevState, deviceInfo.ManagementState)
}
func TestAddDescriptionValidDescription(t *testing.T) {
	deviceInfo := NewDevice("testAsset", "test")
	description := "This is a test device"
	deviceInfo.AddDescription(description)

	assert.NotNil(t, deviceInfo.Description)
	assert.Equal(t, description, *deviceInfo.Description)
}

func TestAddDescriptionEmptyDescription(t *testing.T) {
	deviceInfo := NewDevice("testAsset", "test")
	// Set an initial description to check it doesn't get overwritten
	initialDescription := "Initial description"
	deviceInfo.Description = &initialDescription

	deviceInfo.AddDescription("")

	// Description should remain unchanged
	assert.NotNil(t, deviceInfo.Description)
	assert.Equal(t, initialDescription, *deviceInfo.Description)
}
