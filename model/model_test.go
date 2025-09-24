/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"testing"
	"time"

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

func TestNewGateway(t *testing.T) {
	gatewayInfo := NewGateway("testGateway")
	assert.Equal(t, "Gateway", gatewayInfo.Type)
	assert.Equal(t, "testGateway", *gatewayInfo.Name)
	assert.Equal(t, ManagementStateValuesRegarded, *gatewayInfo.ManagementState.StateValue)
	assert.Equal(t, getAssetContext(), gatewayInfo.Context)
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

func TestAddManagementStateForGateway(t *testing.T) {
	t.Run("addManagementState - Valid Ignored State", func(t *testing.T) {
		g := &GatewayInfo{}
		beforeTime := time.Now().UTC()

		g.addManagementState(ManagementStateValuesIgnored)

		afterTime := time.Now().UTC()

		assert.NotNil(t, g.ManagementState.StateValue)
		assert.Equal(t, ManagementStateValuesIgnored, *g.ManagementState.StateValue)
		assert.NotNil(t, g.ManagementState.StateTimestamp)
		assert.True(t, g.ManagementState.StateTimestamp.After(beforeTime) || g.ManagementState.StateTimestamp.Equal(beforeTime))
		assert.True(t, g.ManagementState.StateTimestamp.Before(afterTime) || g.ManagementState.StateTimestamp.Equal(afterTime))
	})

	t.Run("addManagementState - Valid Regarded State", func(t *testing.T) {
		g := &GatewayInfo{}

		g.addManagementState(ManagementStateValuesRegarded)

		assert.NotNil(t, g.ManagementState.StateValue)
		assert.Equal(t, ManagementStateValuesRegarded, *g.ManagementState.StateValue)
		assert.NotNil(t, g.ManagementState.StateTimestamp)
	})

	t.Run("addManagementState - Valid Unknown State", func(t *testing.T) {
		g := &GatewayInfo{}

		g.addManagementState(ManagementStateValuesUnknown)

		assert.NotNil(t, g.ManagementState.StateValue)
		assert.Equal(t, ManagementStateValuesUnknown, *g.ManagementState.StateValue)
		assert.NotNil(t, g.ManagementState.StateTimestamp)
	})

	t.Run("addManagementState - Empty State Value", func(t *testing.T) {
		g := &GatewayInfo{}
		g.addManagementState(ManagementStateValuesRegarded)
		initialState := g.ManagementState

		g.addManagementState(ManagementStateValues(""))

		assert.Equal(t, initialState, g.ManagementState)
	})

	t.Run("addManagementState - Invalid State Value", func(t *testing.T) {
		g := &GatewayInfo{}
		g.addManagementState(ManagementStateValuesRegarded)
		initialState := g.ManagementState

		g.addManagementState(ManagementStateValues("invalid_state"))

		assert.Equal(t, initialState, g.ManagementState)
	})

	t.Run("addManagementState - Preserves Existing Timestamp When Available", func(t *testing.T) {
		g := &GatewayInfo{}
		initialTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		ignoredState := ManagementStateValuesIgnored
		g.ManagementState = ManagementState{
			StateTimestamp: &initialTime,
			StateValue:     &ignoredState,
		}

		g.addManagementState(ManagementStateValuesRegarded)

		assert.NotNil(t, g.ManagementState.StateValue)
		assert.Equal(t, ManagementStateValuesRegarded, *g.ManagementState.StateValue)
		assert.NotNil(t, g.ManagementState.StateTimestamp)
		assert.Equal(t, initialTime, *g.ManagementState.StateTimestamp)
	})
}
