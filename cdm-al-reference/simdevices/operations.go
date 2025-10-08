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

func ScanForDevices(ethInterface, ipRangeFilter string) []SimulatedDeviceAddress {
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

	filteredDeviceAddresses := []SimulatedDeviceAddress{}
	// filter for IP Range (if required)
	for _, device := range interfaceDevices {
		if ipRangeFilter == "" || device.hasIPInRange(ipRangeFilter) {
			deviceAddress := device.getDeviceAddress()
			filteredDeviceAddresses = append(filteredDeviceAddresses, deviceAddress)
		}
	}

	simulateCostlyOperation(1 * time.Second) // simulate some scan delay before returning
	return filteredDeviceAddresses
}

func RetrieveDeviceDetails(deviceAddress SimulatedDeviceAddress, username, password string) (SimulatedDevice, error) {
	simLock.Lock()
	defer simLock.Unlock()

	device, err := connectToDevice(deviceAddress, username, password)
	if err != nil {
		return nil, err
	}

	device.DeviceState = StateReading
	handleDeviceChanges(true)

	simulateCostlyOperation(2 * time.Second) // simulate reading information from device

	device.DeviceState = StateActive
	handleDeviceChanges(true)

	return device, nil
}
