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
		m := NewDevice("asset", "device")

		m.AddCapabilities("Capabilties-1", true)
		m.AddCapabilities("Capabilties-2", false)
		capabilities := m.AssetOperations
		if len(capabilities) != 2 {
			fmt.Printf("Expected 2 added capabilities, got %d\n", len(capabilities))
			t.Fail()
		}
		found := 0
		for _, v := range capabilities {
			if *v.OperationName == "Capabilties-1" {
				found++
				assert.True(t, *v.ActivationFlag)
				break
			}
		}

		assert.Equal(t, 1, found)
	})
}
