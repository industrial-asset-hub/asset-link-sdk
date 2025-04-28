/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package simdevices

import (
	"time"

	"github.com/rs/zerolog/log"
)

func ScanDevices(ethInterface, ipRangeFilter string) []SimulatedDevice {
	simLock.Lock()
	defer simLock.Unlock()

	if ethInterface == "" {
		log.Info().Msg("Scanning for devices on all interfaces")
	} else {
		log.Info().Msg("Scanning for devices on interface " + ethInterface)
	}

	if ipRangeFilter == "" {
		log.Info().Msg("No Filtering of Devices for IP range")
	} else {
		log.Info().Msg("Filtering of Devices for IP range " + ipRangeFilter)
	}

	interfaceDevices := []*simulatedDeviceInfo{}

	if ethInterface == interfaceEth0 || ethInterface == "" {
		interfaceDevices = append(interfaceDevices, simulatedDevicesEth0...)
	}

	if ethInterface == interfaceEth1 || ethInterface == "" {
		interfaceDevices = append(interfaceDevices, simulatedDevicesEth1...)
	}

	filteredDevices := []SimulatedDevice{}
	// filter for IP Range (if required)
	for _, device := range interfaceDevices {
		device.DeviceState = StateReading
		handleDeviceChanges(true)

		time.Sleep(1 * time.Second) // simulate reading information from device

		device.DeviceState = StateActive
		handleDeviceChanges(true)

		if ipRangeFilter == "" || device.hasIPInRange(ipRangeFilter) {
			filteredDevices = append(filteredDevices, device)
		}
	}

	return filteredDevices
}
