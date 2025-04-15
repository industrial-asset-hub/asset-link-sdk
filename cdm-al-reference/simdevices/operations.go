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

func ConnectToDevice(alNIC, deviceIP string) (SimulatedDevice, error) {
	simLock.Lock()
	defer simLock.Unlock()

	var simulatedDevices []*simulatedDeviceInfo

	switch alNIC {
	case interfaceEth0:
		simulatedDevices = simulatedDevicesEth0
	case interfaceEth1:
		simulatedDevices = simulatedDevicesEth1
	default:
		return nil, errors.New("invalid asset link interface")
	}

	time.Sleep(1 * time.Second) // simulate connecting to device

	for _, device := range simulatedDevices {
		if device.IpDevice == deviceIP {
			return device, nil
		}
	}

	return nil, errors.New("device not found")
}

func (d *simulatedDeviceInfo) UpdateFirmware(artefactFilename string) error {
	simLock.Lock()
	defer simLock.Unlock()

	if !d.UpdateSupport {
		return errors.New("device does not support updates")
	}

	firmwareData, fileErr := os.ReadFile(artefactFilename)
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

	if fwFile.FirmwareVersion == d.FirmwareVersion {
		return errors.New("firmware version is already installed")
	}

	d.DeviceState = StateUpdating
	handleDeviceChanges(true)

	time.Sleep(3 * time.Second) // simulate updating firmware

	d.FirmwareVersion = fwFile.FirmwareVersion
	d.DeviceState = StateActive
	handleDeviceChanges(true)

	return nil
}

func (d *simulatedDeviceInfo) RetrieveFirmware(artefactFilename string) error {
	simLock.Lock()
	defer simLock.Unlock()

	firmwareData := firmwareFile{
		ArtefactType:       "firmware",
		Manufacturer:       d.Manufacturer,
		ProductDesignation: d.ProductDesignation,
		FirmwareVersion:    d.FirmwareVersion,
	}

	fwFile, _ := json.Marshal(firmwareData)

	fileErr := os.WriteFile(artefactFilename, fwFile, 0644)
	if fileErr != nil {
		return fileErr
	}

	d.DeviceState = StateRetrieving
	handleDeviceChanges(true)

	time.Sleep(3 * time.Second) // simulate retrieving firmware

	d.DeviceState = StateActive
	handleDeviceChanges(true)

	return nil
}
