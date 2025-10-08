/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package simdevices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimDevices(t *testing.T) {
	StartSimulatedDevices("") // start without visualization web server

	t.Run("discoveryWithCredentials", func(t *testing.T) {
		deviceAddresses := ScanForDevices("", "") // scan without filters
		assert.Len(t, deviceAddresses, 6)

		for _, addr := range deviceAddresses {
			device, err := RetrieveDeviceDetails(addr, "admin", "admin") // use correct credentials
			assert.NoError(t, err)
			assert.NotNil(t, device)
		}
	})

	t.Run("discoveryWithoutCredentials", func(t *testing.T) {
		deviceAddresses := ScanForDevices("", "") // scan without filters
		assert.Len(t, deviceAddresses, 6)

		failed := 0
		successful := 0
		for _, addr := range deviceAddresses {
			device, err := RetrieveDeviceDetails(addr, "", "") // do not use credentials
			if err != nil || device == nil {
				failed++
			} else {
				successful++
			}
		}

		assert.Equal(t, 4, successful)
		assert.Equal(t, 2, failed)
	})
}
