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

		foundSubdevices := false
		failed := 0
		successful := 0
		for _, addr := range deviceAddresses {
			device, err := RetrieveDeviceDetails(addr, "admin", "admin") // use correct credentials
			if err != nil || device == nil {
				failed++
				continue
			} else {
				successful++
			}

			subdevices := device.GetSubDevices()
			if device.GetDeviceName() == "Simulated Device B1" {
				assert.Len(t, subdevices, 3) // device B1 has 3 subdevices
				foundSubdevices = true
			} else {
				assert.Empty(t, subdevices)
			}
		}

		assert.Equal(t, 6, successful)
		assert.Equal(t, 0, failed)
		assert.True(t, foundSubdevices)
	})

	t.Run("discoveryWithoutCredentials", func(t *testing.T) {
		deviceAddresses := ScanForDevices("", "") // scan without filters
		assert.Len(t, deviceAddresses, 6)

		foundSubdevices := false
		failed := 0
		successful := 0
		for _, addr := range deviceAddresses {
			device, err := RetrieveDeviceDetails(addr, "", "") // do not use credentials
			if err != nil || device == nil {
				failed++
				continue
			} else {
				successful++
			}

			subdevices := device.GetSubDevices()
			if device.GetDeviceName() == "Simulated Device B1" {
				assert.Len(t, subdevices, 3) // device B1 has 3 subdevices
				foundSubdevices = true
			} else {
				assert.Empty(t, subdevices)
			}
		}

		assert.Equal(t, 4, successful)
		assert.Equal(t, 2, failed)
		assert.True(t, foundSubdevices)
	})
}
