/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package simdevices

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getFilenamePath(filename string) string {
	_, srcFilename, _, _ := runtime.Caller(0)
	srcDir := path.Dir(srcFilename)
	srcRoot := path.Join(srcDir, "../..")
	filenamePath := path.Join(srcRoot, filename)
	return filenamePath
}

func TestSimDevices(t *testing.T) {
	t.Run("discoveryWithCredentials", func(t *testing.T) {
		StartSimulatedDevices("")

		deviceAddresses := ScanForDevices("", "") // scan without filters
		assert.Len(t, deviceAddresses, 6)

		credentials := SimulatedDeviceCredentials{
			Username: "user",
			Password: "user_password",
		}

		foundSubdevices := false
		failed := 0
		successful := 0
		for _, addr := range deviceAddresses {
			device, err := RetrieveDeviceDetails(addr, &credentials) // use correct credentials
			if err != nil || device == nil {
				failed++
				continue
			} else {
				successful++
			}

			subdevices := device.GetSubDevices()
			if device.GetDeviceName() == "Simulated Device B0" {
				for _, subdevice := range subdevices {
					sdAddr := subdevice.GetDeviceAddress()
					assert.GreaterOrEqual(t, sdAddr.SubDeviceID, 0)
					sdViaAddr, err := RetrieveDeviceDetails(sdAddr, nil)
					assert.NoError(t, err)
					assert.NotNil(t, sdViaAddr)
					assert.Equal(t, sdViaAddr.GetDeviceName(), subdevice.GetDeviceName())
				}

				assert.Len(t, subdevices, 3) // device B0 has 3 subdevices
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
		StartSimulatedDevices("")

		deviceAddresses := ScanForDevices("", "") // scan without filters
		assert.Len(t, deviceAddresses, 6)

		foundSubdevices := false
		failed := 0
		successful := 0
		for _, addr := range deviceAddresses {
			device, err := RetrieveDeviceDetails(addr, nil) // do not use credentials
			if err != nil || device == nil {
				failed++
				continue
			} else {
				successful++
			}

			subdevices := device.GetSubDevices()
			if device.GetDeviceName() == "Simulated Device B0" {
				for _, subdevice := range subdevices {
					sdAddr := subdevice.GetDeviceAddress()
					assert.GreaterOrEqual(t, sdAddr.SubDeviceID, 0)
					sdViaAddr, err := RetrieveDeviceDetails(sdAddr, nil)
					assert.NoError(t, err)
					assert.NotNil(t, sdViaAddr)
					assert.Equal(t, sdViaAddr.GetDeviceName(), subdevice.GetDeviceName())
				}

				assert.Len(t, subdevices, 3) // device B0 has 3 subdevices
				foundSubdevices = true
			} else {
				assert.Empty(t, subdevices)
			}
		}

		assert.Equal(t, 4, successful)
		assert.Equal(t, 2, failed)
		assert.True(t, foundSubdevices)
	})

	t.Run("updateWithCredentials", func(t *testing.T) {
		StartSimulatedDevices("")

		deviceAddresses := ScanForDevices("", "") // scan without filters
		assert.Len(t, deviceAddresses, 6)

		credentials := SimulatedDeviceCredentials{
			Username: "admin",
			Password: "admin_password",
		}

		failed := 0
		successful := 0
		for _, addr := range deviceAddresses {
			device, err := ConnectToDevice(addr, &credentials) // use correct credentials
			if err != nil || device == nil {
				failed++
				continue
			}

			assert.Equal(t, "1.0.0", device.GetInstalledFirmwareVersion())
			assert.Equal(t, "1.0.0", device.GetActiveFirmwareVersion())

			firmwareFile := getFilenamePath("misc/simulated_device_firmware_2.0.0.fwu")
			err = device.UpdateFirmware(firmwareFile)
			assert.NoError(t, err)

			assert.Equal(t, "2.0.0", device.GetInstalledFirmwareVersion())
			assert.Equal(t, "1.0.0", device.GetActiveFirmwareVersion())

			err = device.RebootDevice()
			assert.NoError(t, err)

			assert.Equal(t, "2.0.0", device.GetInstalledFirmwareVersion())
			assert.Equal(t, "2.0.0", device.GetActiveFirmwareVersion())

			successful++
		}

		assert.Equal(t, 6, successful)
		assert.Equal(t, 0, failed)
	})

	t.Run("updateWithoutCredentials", func(t *testing.T) {
		StartSimulatedDevices("")

		deviceAddresses := ScanForDevices("", "") // scan without filters
		assert.Len(t, deviceAddresses, 6)

		failed := 0
		successful := 0
		for _, addr := range deviceAddresses {
			device, err := ConnectToDevice(addr, nil) // do not use credentials
			if err != nil || device == nil {
				failed++
				continue
			}

			assert.Equal(t, "1.0.0", device.GetInstalledFirmwareVersion())
			assert.Equal(t, "1.0.0", device.GetActiveFirmwareVersion())

			firmwareFile := getFilenamePath("misc/simulated_device_firmware_3.0.0.fwu")
			err = device.UpdateFirmware(firmwareFile)
			assert.NoError(t, err)
			if err != nil {
				failed++
				continue
			}

			assert.Equal(t, "3.0.0", device.GetInstalledFirmwareVersion())
			assert.Equal(t, "1.0.0", device.GetActiveFirmwareVersion())

			err = device.RebootDevice()
			assert.NoError(t, err)
			if err != nil {
				failed++
				continue
			}

			assert.Equal(t, "3.0.0", device.GetInstalledFirmwareVersion())
			assert.Equal(t, "3.0.0", device.GetActiveFirmwareVersion())

			successful++
		}

		assert.Equal(t, 2, successful)
		assert.Equal(t, 4, failed)
	})

	t.Run("deviceConfig", func(t *testing.T) {
		StartSimulatedDevices("")

		deviceAddress := SimulatedDeviceAddress{
			AssetLinkNIC: "eth0",
			DeviceIP:     "192.168.0.10",
			SubDeviceID:  -1,
		}

		device, err := ConnectToDevice(deviceAddress, nil) // do not use credentials
		assert.NoError(t, err)
		assert.NotNil(t, device)

		assert.Equal(t, "Simulated Device A0", device.GetDeviceName())

		storeFilename := getFilenamePath("misc/simulated_device_cfg_old.cfg")
		defer os.Remove(storeFilename) // clean up after test run
		err = device.GetConfig(storeFilename)
		assert.NoError(t, err)
		assert.FileExists(t, storeFilename)

		data, err := os.ReadFile(storeFilename)
		assert.NoError(t, err)

		expectedJSON := "{\"artefact_type\":\"config\",\"device_name\":\"Simulated Device A0\"}"
		assert.Equal(t, expectedJSON, string(data))

		loadFilename := getFilenamePath("misc/simulated_device_cfg_1.cfg")
		err = device.SetConfig(loadFilename)
		assert.NoError(t, err)

		assert.Equal(t, "Simulated Device - CFG 1", device.GetDeviceName())
	})
}
