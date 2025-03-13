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
