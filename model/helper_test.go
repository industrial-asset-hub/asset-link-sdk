/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddIdentifier(t *testing.T) {
	deviceInfo := NewDevice("testAsset", "test")
	mac := "test-mac"
	identifierUncertainty := 1
	deviceInfo.addIdentifier(mac)
	assert.Equal(t, len(deviceInfo.MacIdentifiers), 1)
	assert.Equal(t, deviceInfo.MacIdentifiers[0], MacIdentifier{
		IdentifierType:        nil,
		IdentifierUncertainty: &identifierUncertainty,
		MacAddress:            &mac,
	})
}

func TestAddReachabilityState(t *testing.T) {
	deviceInfo := NewDevice("testAsset", "test")
	state := ReachabilityStateValuesReached
	deviceInfo.addReachabilityState()
	assert.Equal(t, deviceInfo.ReachabilityState.StateValue, &state)
}

func TestAddManagementState(t *testing.T) {
	deviceInfo := NewDevice("testAsset", "test")
	state := ManagementStateValuesUnknown
	deviceInfo.addReachabilityState()
	assert.Equal(t, deviceInfo.ManagementState.StateValue, &state)
}
