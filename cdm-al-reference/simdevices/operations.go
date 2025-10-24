/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package simdevices

import (
	"encoding/json"
	"errors"
	"os"
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
			deviceAddress := device.GetDeviceAddress()
			filteredDeviceAddresses = append(filteredDeviceAddresses, deviceAddress)
		}
	}

	simulateCostlyOperation(1 * time.Second) // simulate some scan delay before returning
	return filteredDeviceAddresses
}

// Similar to ConnectToDevice but also simulates additional reading of device details (incl. delay and state changes in frontend)
func RetrieveDeviceDetails(deviceAddress SimulatedDeviceAddress, credentials *SimulatedDeviceCredentials) (SimulatedDevice, error) {
	simLock.Lock()
	defer simLock.Unlock()

	device, err := connectToDevice(deviceAddress, credentials)
	if err != nil {
		return nil, err
	}

	device.DeviceState = StateReading
	handleDeviceChanges()

	simulateCostlyOperation(2 * time.Second) // simulate reading information from device

	device.DeviceState = StateActive
	handleDeviceChanges()

	return device, nil
}

func ConnectToDevice(deviceAddress SimulatedDeviceAddress, credentials *SimulatedDeviceCredentials) (SimulatedDevice, error) {
	simLock.Lock()
	defer simLock.Unlock()

	device, err := connectToDevice(deviceAddress, credentials)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (d *simulatedDeviceInfo) UpdateFirmware(firmwareFilename string) error {
	simLock.Lock()
	defer simLock.Unlock()

	if !d.UpdateSupport {
		return errors.New("device does not support updates")
	}

	firmwareData, fileErr := os.ReadFile(firmwareFilename)
	if fileErr != nil {
		return fileErr
	}

	var fwFile firmwareFile
	parseErr := json.Unmarshal(firmwareData, &fwFile)
	if parseErr != nil {
		return errors.New("invalid firmware artefact (overall file format)")
	}

	if fwFile.ArtefactType != "firmware" {
		return errors.New("artefact type mismatch")
	}

	if fwFile.Manufacturer != d.Manufacturer {
		return errors.New("manufacturer mismatch")
	}

	if fwFile.ProductDesignation != d.ProductDesignation {
		return errors.New("product designation mismatch")
	}

	//nolint:gocritic
	if fwFile.FirmwareVersion == d.InstalledFirmwareVersion && fwFile.FirmwareVersion == d.ActiveFirmwareVersion {
		return errors.New("firmware version is already installed")
	}

	d.DeviceState = StateUpdating
	handleDeviceChanges()

	simulateCostlyOperation(3 * time.Second) // simulate updating firmware

	d.InstalledFirmwareVersion = fwFile.FirmwareVersion
	d.DeviceState = StateActive
	handleDeviceChanges()

	return nil
}

func (d *simulatedDeviceInfo) RebootDevice() error {
	simLock.Lock()
	defer simLock.Unlock()

	d.DeviceState = StateBooting
	handleDeviceChanges()

	simulateCostlyOperation(3 * time.Second) // simulate reboot to activate firmware

	d.ActiveFirmwareVersion = d.InstalledFirmwareVersion
	d.DeviceState = StateActive
	handleDeviceChanges()

	return nil
}

// Loads the specified file and stores the contained configuration on the device
func (d *simulatedDeviceInfo) SetConfig(configFilename string) error {
	simLock.Lock()
	defer simLock.Unlock()

	if !d.UpdateSupport {
		return errors.New("device does not support updates")
	}

	configData, fileErr := os.ReadFile(configFilename)
	if fileErr != nil {
		return fileErr
	}

	var configFile configFile
	parseErr := json.Unmarshal(configData, &configFile)
	if parseErr != nil {
		return errors.New("invalid config artefact (overall file format)")
	}

	if configFile.ArtefactType != "config" {
		return errors.New("artefact type mismatch")
	}

	d.DeviceState = StateSetting
	handleDeviceChanges()

	simulateCostlyOperation(3 * time.Second) // simulate storing config on the device

	d.DeviceName = configFile.DeviceName
	d.DeviceState = StateActive
	handleDeviceChanges()

	return nil
}

// Retrieves configuration from the device and saves it to the specified file
func (d *simulatedDeviceInfo) GetConfig(configFilename string) error {
	simLock.Lock()
	defer simLock.Unlock()

	configData := configFile{
		ArtefactType: "config",
		DeviceName:   d.DeviceName,
	}

	fwFile, _ := json.Marshal(configData)

	fileErr := os.WriteFile(configFilename, fwFile, 0644)
	if fileErr != nil {
		return fileErr
	}

	d.DeviceState = StateGetting
	handleDeviceChanges()

	simulateCostlyOperation(3 * time.Second) // simulate retrieving config from the device

	d.DeviceState = StateActive
	handleDeviceChanges()

	return nil
}
