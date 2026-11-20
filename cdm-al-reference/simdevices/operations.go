/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package simdevices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

// Similar to ConnectToDevice but only requires read permissions and also simulates additional reading of device details (incl. delay and state changes in frontend)
func RetrieveDeviceDetails(deviceAddress SimulatedDeviceAddress, credentials *SimulatedDeviceCredentials) (SimulatedDeviceInfo, error) {
	simLock.Lock()
	defer simLock.Unlock()

	device, err := connectToDevice(deviceAddress, credentials, true)
	if err != nil {
		return nil, err
	}

	device.DeviceState = StateReading
	handleDeviceChanges()

	simulateCostlyOperation(1 * time.Second) // simulate reading information from device

	device.DeviceState = StateActive
	handleDeviceChanges()

	return device, nil
}

func ConnectToDevice(deviceAddress SimulatedDeviceAddress, credentials *SimulatedDeviceCredentials) (SimulatedDevice, error) {
	simLock.Lock()
	defer simLock.Unlock()

	device, err := connectToDevice(deviceAddress, credentials, false)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (d *simulatedDeviceInfo) UpdateFirmware(firmwareFilename string) error {
	firmware, err := d.PrepareFirmware(firmwareFilename)
	if err != nil {
		return err
	}
	return d.InstallFirmware(context.Background(), firmware)
}

func (d *simulatedDeviceInfo) PrepareFirmware(firmwareFilename string) (PreparedFirmware, error) {
	data, err := os.ReadFile(firmwareFilename)
	if err != nil {
		return PreparedFirmware{}, err
	}

	var fwFile firmwareFile
	if err := json.Unmarshal(data, &fwFile); err != nil {
		return PreparedFirmware{}, err
	}

	simLock.Lock()
	defer simLock.Unlock()
	if err := d.validateFirmware(fwFile); err != nil {
		return PreparedFirmware{}, err
	}
	return PreparedFirmware{file: fwFile}, nil
}

func (d *simulatedDeviceInfo) validateFirmware(fwFile firmwareFile) error {
	if !d.UpdateSupport {
		return errors.New("device does not support updates")
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
	if fwFile.FirmwareVersion == "" {
		return errors.New("firmware version is required")
	}

	//nolint:gocritic
	if fwFile.FirmwareVersion == d.InstalledFirmwareVersion && fwFile.FirmwareVersion == d.ActiveFirmwareVersion {
		return errors.New("firmware version is already installed")
	}
	return nil
}

func (d *simulatedDeviceInfo) InstallFirmware(ctx context.Context, firmware PreparedFirmware) error {
	return d.installFirmware(ctx, firmware, simulationDelay(3*time.Second))
}

func (d *simulatedDeviceInfo) installFirmware(ctx context.Context, firmware PreparedFirmware, delay time.Duration) error {
	simLock.Lock()
	defer simLock.Unlock()
	if err := d.validateFirmware(firmware.file); err != nil {
		return fmt.Errorf("%w: %w", ErrFirmwareRejected, err)
	}
	return d.runFirmwareStep(ctx, StateUpdating, delay, func() {
		d.InstalledFirmwareVersion = firmware.Version()
	})
}

func (d *simulatedDeviceInfo) RebootDevice() error {
	return d.RebootDeviceContext(context.Background())
}

func (d *simulatedDeviceInfo) RebootDeviceContext(ctx context.Context) error {
	return d.rebootDevice(ctx, simulationDelay(3*time.Second))
}

func (d *simulatedDeviceInfo) rebootDevice(ctx context.Context, delay time.Duration) error {
	simLock.Lock()
	defer simLock.Unlock()
	return d.runFirmwareStep(ctx, StateBooting, delay, func() {
		d.ActiveFirmwareVersion = d.InstalledFirmwareVersion
	})
}

// runFirmwareStep is called with simLock held. The simulator applies the
// change only once the step completes; real devices may continue on their own.
func (d *simulatedDeviceInfo) runFirmwareStep(ctx context.Context, state SimulatedDeviceState, delay time.Duration, apply func()) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	d.DeviceState = state
	handleDeviceChanges()
	defer func() {
		d.DeviceState = StateActive
		handleDeviceChanges()
	}()
	if err := waitForFirmwareStep(ctx, delay); err != nil {
		return err
	}
	apply()
	return nil
}

func waitForFirmwareStep(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
	return ctx.Err()
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
