/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapabilities(t *testing.T) {
	t.Run("AddCapabilties", func(t *testing.T) {
		m, err := NewDevice("Asset", "device")
		assert.NoError(t, err)

		err = m.AddCapabilities("Capabilties-1", true)
		assert.NoError(t, err)
		err = m.AddCapabilities("Capabilties-2", false)
		assert.NoError(t, err)
		capabilities := m.AssetOperations
		if len(capabilities) != 2 {
			fmt.Printf("Expected 2 added capabilities, got %d\n", len(capabilities))
			t.Fail()
		}
		found := 0
		for _, v := range capabilities {
			if v.OperationName == "Capabilties-1" {
				found++
				assert.True(t, v.ActivationFlag)
			}
			if v.OperationName == "Capabilties-2" {
				found++
				assert.False(t, v.ActivationFlag)
			}
		}

		assert.Equal(t, 2, found)
	})

	t.Run("AddCapabilities_ValidName", func(t *testing.T) {
		m, err := NewDevice("Asset", "device")
		assert.NoError(t, err)
		err = m.AddCapabilities("firmware_update", true)
		assert.NoError(t, err)
	})

	t.Run("AddCapabilities_EmptyName", func(t *testing.T) {
		m, err := NewDevice("Asset", "device")
		assert.NoError(t, err)
		err = m.AddCapabilities("", true)
		assert.Error(t, err)
	})
}
