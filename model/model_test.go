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
	deviceInfo := NewDevice("testAsset", "testsuite")
	deviceInfo.AddNic("testsuite", "ab:ab:ab:ab:ab:ab")

	assert.Equal(t, "testAsset", deviceInfo.Type)
	assert.Equal(t, "testsuite", *deviceInfo.Name)
	assert.Equal(t, "ab:ab:ab:ab:ab:ab", *deviceInfo.MacIdentifiers[0].MacAddress)
	assert.Equal(t, 1, *deviceInfo.MacIdentifiers[0].IdentifierUncertainty)
	assert.Equal(t, ReachabilityStateValuesReached, *deviceInfo.ReachabilityState.StateValue)
	assert.Equal(t, ManagementStateValuesUnknown, *deviceInfo.ManagementState.StateValue)
}
