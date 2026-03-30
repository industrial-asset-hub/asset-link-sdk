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
	deviceInfo, err := NewDevice("testAsset", "test")
	assert.NoError(t, err)
	_, err = deviceInfo.AddNic("test", "ab:ab:ab:ab:ab:ab")
	assert.NoError(t, err)

	assert.Equal(t, "testAsset", deviceInfo.Type)
	assert.Equal(t, "test", *deviceInfo.Name)
	assert.Equal(t, "ab:ab:ab:ab:ab:ab", *deviceInfo.MacIdentifiers[0].MacAddress)
	assert.Equal(t, 1, *deviceInfo.MacIdentifiers[0].IdentifierUncertainty)
	assert.Equal(t, ReachabilityStateValuesReached, *deviceInfo.ReachabilityState.StateValue)
	assert.Equal(t, ManagementStateValuesUnknown, *deviceInfo.ManagementState.StateValue)
	assert.Equal(t, getAssetContext(), deviceInfo.Context)
}

func TestNewDevice_EmptyType(t *testing.T) {
	deviceInfo, err := NewDevice("", "test")
	assert.Error(t, err)
	assert.NotNil(t, deviceInfo)
	var ee *EmptyError
	if assert.ErrorAs(t, err, &ee) {
		assert.Equal(t, "Type", ee.Field)
		assert.Equal(t, "Asset type is required and cannot be empty", ee.Message)
		assert.Equal(t, "", ee.Value)
	}
}

func TestNewGateway(t *testing.T) {
	gatewayInfo, err := NewGateway("testGateway")
	assert.NoError(t, err)
	assert.Equal(t, "Gateway", gatewayInfo.Type)
	assert.Equal(t, "testGateway", *gatewayInfo.Name)
	assert.Equal(t, ManagementStateValuesRegarded, *gatewayInfo.ManagementState.StateValue)
	assert.Equal(t, getAssetContext(), gatewayInfo.Context)
}

func TestAddManagementStateValidStates(t *testing.T) {
	deviceInfo, err := NewDevice("testAsset", "test")
	assert.NoError(t, err)

	err = deviceInfo.AddManagementState(ManagementStateValuesUnknown)
	assert.NoError(t, err)
	assert.Equal(t, ManagementStateValuesUnknown, *deviceInfo.ManagementState.StateValue)

	err = deviceInfo.AddManagementState(ManagementStateValuesIgnored)
	assert.NoError(t, err)
	assert.Equal(t, ManagementStateValuesIgnored, *deviceInfo.ManagementState.StateValue)

	err = deviceInfo.AddManagementState(ManagementStateValuesRegarded)
	assert.NoError(t, err)
	assert.Equal(t, ManagementStateValuesRegarded, *deviceInfo.ManagementState.StateValue)
}

func TestAddManagementStateEmptyState(t *testing.T) {
	deviceInfo, err := NewDevice("testAsset", "test")
	assert.NoError(t, err)
	prevState := deviceInfo.ManagementState

	err = deviceInfo.AddManagementState("")
	assert.Error(t, err)
	var ee *EmptyError
	if assert.ErrorAs(t, err, &ee) {
		assert.Equal(t, "ManagementState", ee.Field)
		assert.Equal(t, "Management state value is empty", ee.Message)
		assert.Equal(t, ManagementStateValues(""), ee.Value)
	}
	assert.Equal(t, prevState, deviceInfo.ManagementState)
}

func TestAddManagementStateInvalidState(t *testing.T) {
	deviceInfo, err := NewDevice("testAsset", "test")
	assert.NoError(t, err)
	prevState := deviceInfo.ManagementState

	err = deviceInfo.AddManagementState("invalid_state")
	assert.Error(t, err)
	var pe *PermissibleValuesError
	if assert.ErrorAs(t, err, &pe) {
		assert.Equal(t, "ManagementState", pe.Field)
		assert.Equal(t, ManagementStateValues("invalid_state"), pe.Value)
		assert.Equal(t, []interface{}{ManagementStateValuesIgnored, ManagementStateValuesRegarded, ManagementStateValuesUnknown}, pe.Allowed)
	}
	assert.Equal(t, prevState, deviceInfo.ManagementState)
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

func TestAddManagementStateForGateway(t *testing.T) {
	t.Run("addManagementState - Valid Ignored State", func(t *testing.T) {
		g := &GatewayInfo{}
		beforeTime := time.Now().UTC()

		err := g.addManagementState(ManagementStateValuesIgnored)
		assert.NoError(t, err)

		afterTime := time.Now().UTC()

		assert.NotNil(t, g.ManagementState.StateValue)
		assert.Equal(t, ManagementStateValuesIgnored, *g.ManagementState.StateValue)
		assert.NotNil(t, g.ManagementState.StateTimestamp)
		assert.True(t, g.ManagementState.StateTimestamp.After(beforeTime) || g.ManagementState.StateTimestamp.Equal(beforeTime))
		assert.True(t, g.ManagementState.StateTimestamp.Before(afterTime) || g.ManagementState.StateTimestamp.Equal(afterTime))
	})

	t.Run("addManagementState - Valid Regarded State", func(t *testing.T) {
		g := &GatewayInfo{}

		err := g.addManagementState(ManagementStateValuesRegarded)
		assert.NoError(t, err)

		assert.NotNil(t, g.ManagementState.StateValue)
		assert.Equal(t, ManagementStateValuesRegarded, *g.ManagementState.StateValue)
		assert.NotNil(t, g.ManagementState.StateTimestamp)
	})

	t.Run("addManagementState - Valid Unknown State", func(t *testing.T) {
		g := &GatewayInfo{}

		err := g.addManagementState(ManagementStateValuesUnknown)
		assert.NoError(t, err)

		assert.NotNil(t, g.ManagementState.StateValue)
		assert.Equal(t, ManagementStateValuesUnknown, *g.ManagementState.StateValue)
		assert.NotNil(t, g.ManagementState.StateTimestamp)
	})

	t.Run("addManagementState - Empty State Value", func(t *testing.T) {
		g := &GatewayInfo{}
		err := g.addManagementState(ManagementStateValuesRegarded)
		assert.NoError(t, err)
		initialState := g.ManagementState

		err = g.addManagementState(ManagementStateValues(""))
		assert.Error(t, err)
		var ee *EmptyError
		if assert.ErrorAs(t, err, &ee) {
			assert.Equal(t, "ManagementState", ee.Field)
			assert.Equal(t, "Management state value is empty", ee.Message)
			assert.Equal(t, ManagementStateValues(""), ee.Value)
		}

		assert.Equal(t, initialState, g.ManagementState)
	})

	t.Run("addManagementState - Invalid State Value", func(t *testing.T) {
		g := &GatewayInfo{}
		err := g.addManagementState(ManagementStateValuesRegarded)
		assert.NoError(t, err)
		initialState := g.ManagementState

		err = g.addManagementState(ManagementStateValues("invalid_state"))
		assert.Error(t, err)
		var pe *PermissibleValuesError
		if assert.ErrorAs(t, err, &pe) {
			assert.Equal(t, "ManagementState", pe.Field)
			assert.Equal(t, ManagementStateValues("invalid_state"), pe.Value)
			assert.Equal(t, []interface{}{ManagementStateValuesIgnored, ManagementStateValuesRegarded, ManagementStateValuesUnknown}, pe.Allowed)
		}

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

		err := g.addManagementState(ManagementStateValuesRegarded)
		assert.NoError(t, err)

		assert.NotNil(t, g.ManagementState.StateValue)
		assert.Equal(t, ManagementStateValuesRegarded, *g.ManagementState.StateValue)
		assert.NotNil(t, g.ManagementState.StateTimestamp)
		assert.Equal(t, initialTime, *g.ManagementState.StateTimestamp)
	})
}
