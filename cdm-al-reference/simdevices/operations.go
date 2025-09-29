/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package simdevices

import (
	"time"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
)

// ScanResult represents the result of scanning a single device
type ScanResult struct {
	Device SimulatedDevice          `json:"device,omitempty"`
	Error  *generated.DiscoverError `json:"error,omitempty"`
}

// ScanDevicesWithErrors returns both successfully scanned devices and any errors encountered
func ScanDevices(ethInterface, ipRangeFilter string) ([]SimulatedDevice, []*generated.DiscoverError) {
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
	scanErrors := []*generated.DiscoverError{}

	// filter for IP Range (if required) and simulate authentication/connection errors
	for i, device := range interfaceDevices {
		device.DeviceState = StateReading
		handleDeviceChanges(true)

		time.Sleep(2 * time.Second) // simulate reading information from device

		// Simulate authentication/connection errors for some devices
		// This is where you would put real authentication logic
		// Simulate error for every alternate device
		if i%2 == 1 {
			scanError := &generated.DiscoverError{
				ResultCode:  int32(codes.Unavailable),
				Description: "Simulated error for device " + device.GetDeviceName(),
			}
			scanErrors = append(scanErrors, scanError)
			device.DeviceState = StateActive
			handleDeviceChanges(true)
			continue
		}

		device.DeviceState = StateActive
		handleDeviceChanges(true)

		if ipRangeFilter == "" || device.hasIPInRange(ipRangeFilter) {
			filteredDevices = append(filteredDevices, device)
		}
	}

	return filteredDevices, scanErrors
}

// shouldSimulateAuthError simulates authentication errors for demonstration
// In real implementation, this would be actual authentication logic
func shouldSimulateAuthError(device *simulatedDeviceInfo) bool {
	// Simulate auth error for devices with serial numbers ending in "01"
	return len(device.SerialNumber) > 0 && device.SerialNumber[len(device.SerialNumber)-2:] == "01"
}
