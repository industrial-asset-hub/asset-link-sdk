/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package reference

//TODO: maybe add a webserver that shows all devices and their state

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type SimulatedDevice struct {
	deviceName         string
	manufacturer       string
	productDesignation string
	articleNumber      string
	hardwareVersion    string
	firmwareVersion    string
	updateSupport      bool
	alNIC              string
	deviceNIC          string
	macAddress         string
	ipDevice           string
	ipNetmask          string
	ipRoute            string
	serialNumber       string
}

type firmwareFile struct {
	ArtefactType       string `json:"artefact_type"`
	Manufacturer       string `json:"manufacturer"`
	ProductDesignation string `json:"product_designation"`
	FirmwareVersion    string `json:"firmware_version"`
}

const (
	interfaceEth0 = "eth0"
	interfaceEth1 = "eth1"
)

var simulatedDevicesEth0 []*SimulatedDevice
var simulatedDevicesEth1 []*SimulatedDevice

func newSimulatedDevice(alNIC, serial, mac, ip, name string) *SimulatedDevice {
	return &SimulatedDevice{
		deviceName:         name,
		manufacturer:       "Siemens AG",
		productDesignation: "Simulated Device",
		articleNumber:      "AN0123456789",
		hardwareVersion:    "3",
		firmwareVersion:    "1.0.0",
		updateSupport:      true,
		alNIC:              alNIC,
		deviceNIC:          "enp0",
		macAddress:         mac,
		ipDevice:           ip,
		ipNetmask:          "255.255.255.0",
		ipRoute:            "",
		serialNumber:       serial,
	}
}

func init() {
	simulatedDevicesEth0 = append(simulatedDevicesEth0, newSimulatedDevice(interfaceEth0, "SN123450000", "00:16:3e:00:00:00", "192.168.0.10", "Simulated Device A0"))
	simulatedDevicesEth0 = append(simulatedDevicesEth0, newSimulatedDevice(interfaceEth0, "SN123450001", "00:16:3e:00:00:01", "192.168.0.11", "Simulated Device A1"))

	simulatedDevicesEth1 = append(simulatedDevicesEth1, newSimulatedDevice(interfaceEth1, "SN123450100", "00:16:3e:00:01:00", "192.168.1.10", "Simulated Device B0"))
	simulatedDevicesEth1 = append(simulatedDevicesEth1, newSimulatedDevice(interfaceEth1, "SN123450101", "00:16:3e:00:01:01", "192.168.1.11", "Simulated Device B1"))
}

func (d *SimulatedDevice) GetDeviceName() string {
	return d.deviceName
}

func (d *SimulatedDevice) GetManufacturer() string {
	return d.manufacturer
}

func (d *SimulatedDevice) GetProductDesignation() string {
	return d.productDesignation
}

func (d *SimulatedDevice) GetArticleNumber() string {
	return d.articleNumber
}

func (d *SimulatedDevice) GetHardwareVersion() string {
	return d.hardwareVersion
}

func (d *SimulatedDevice) GetFirmwareVersion() string {
	return d.firmwareVersion
}

func (d *SimulatedDevice) IsUpdateSupported() bool {
	return d.updateSupport
}

func (d *SimulatedDevice) GetALNIC() string {
	return d.alNIC
}

func (d *SimulatedDevice) GetDeviceNIC() string {
	return d.deviceNIC
}

func (d *SimulatedDevice) GetMacAddress() string {
	return d.macAddress
}

func (d *SimulatedDevice) GetIpDevice() string {
	return d.ipDevice
}

func (d *SimulatedDevice) GetIpNetmask() string {
	return d.ipNetmask
}

func (d *SimulatedDevice) GetIpRoute() string {
	return d.ipRoute
}

func (d *SimulatedDevice) GetIDLink() string {
	// generate ID link from device data
	idLink := fmt.Sprintf("https://industrial-assets.io?1P=%s&S=%s", d.GetArticleNumber(), d.GetSerialNumber())
	return idLink
}

func (d *SimulatedDevice) GetSerialNumber() string {
	return d.serialNumber
}

func (d *SimulatedDevice) HasIPInRange(ipRange string) bool {
	deviceIPs := []string{}
	deviceIPs = append(deviceIPs, d.ipDevice) // there is currently only one IP address per device

	return containsIpInRange(ipRange, deviceIPs)
}

func (d *SimulatedDevice) UpdateFirmware(artefactFilename string) error {
	if !d.IsUpdateSupported() {
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

	if fwFile.Manufacturer != d.GetManufacturer() {
		return errors.New("manufacturer mismatch")
	}

	if fwFile.ProductDesignation != d.GetProductDesignation() {
		return errors.New("product designation mismatch")
	}

	if fwFile.FirmwareVersion == d.firmwareVersion {
		return errors.New("firmware version is already installed")
	}

	time.Sleep(2 * time.Second)

	d.firmwareVersion = fwFile.FirmwareVersion
	return nil
}

func (d *SimulatedDevice) RetrieveFirmware(artefactFilename string) error {
	firmwareData := firmwareFile{
		ArtefactType:       "firmware",
		Manufacturer:       d.GetManufacturer(),
		ProductDesignation: d.GetProductDesignation(),
		FirmwareVersion:    d.GetFirmwareVersion(),
	}

	fwFile, _ := json.Marshal(firmwareData)

	fileErr := os.WriteFile(artefactFilename, fwFile, 0644)
	if fileErr != nil {
		return fileErr
	}

	time.Sleep(2 * time.Second)

	return nil
}

func ScanDevices(ethInterface, ipRangeFilter string) []*SimulatedDevice {
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

	interfaceDevices := []*SimulatedDevice{}

	if ethInterface == interfaceEth0 || ethInterface == "" {
		interfaceDevices = append(interfaceDevices, simulatedDevicesEth0...)
	}

	if ethInterface == interfaceEth1 || ethInterface == "" {
		interfaceDevices = append(interfaceDevices, simulatedDevicesEth1...)
	}

	filteredDevices := interfaceDevices

	if ipRangeFilter != "" {
		// filter for IP Range (if required)
		filteredDevices = []*SimulatedDevice{}
		for _, device := range interfaceDevices {
			if device.HasIPInRange(ipRangeFilter) {
				filteredDevices = append(filteredDevices, device)
			}
		}
	}

	return filteredDevices
}

func ConnectToDevice(alNIC, deviceIP string) (*SimulatedDevice, error) {
	var simulatedDevices []*SimulatedDevice

	switch alNIC {
	case interfaceEth0:
		simulatedDevices = simulatedDevicesEth0
	case interfaceEth1:
		simulatedDevices = simulatedDevicesEth1
	default:
		return nil, errors.New("invalid asset link interface")
	}

	for _, device := range simulatedDevices {
		if device.GetIpDevice() == deviceIP {
			return device, nil
		}
	}

	return nil, errors.New("device not found")
}

func containsIpInRange(ipRange string, actualIPs []string) bool {
	ipRangeParts := strings.Split(ipRange, "-")
	if len(ipRangeParts) != 2 {
		log.Warn().Msg("Invalid IP range (format)")
		return true
	}

	ipRangeBeginString := ipRangeParts[0]
	ipRangeEndString := ipRangeParts[1]

	ipRangeBegin := net.ParseIP(ipRangeBeginString)
	if ipRangeBegin == nil {
		log.Warn().Msg("Invalid IP range (IP range begin)")
		return true
	}

	ipRangeEnd := net.ParseIP(ipRangeEndString)
	if ipRangeEnd == nil {
		log.Warn().Msg("Invalid IP range (IP range end)")
		return true
	}

	if len(actualIPs) == 0 {
		return false
	}

	for _, actualIPString := range actualIPs {
		actualIP := net.ParseIP(actualIPString)
		if actualIP == nil {
			log.Warn().Msg("Invalid device IP")
			return true
		}

		ipRangeBegin16 := ipRangeBegin.To16()
		ipRangeEnd16 := ipRangeEnd.To16()
		actualIP16 := actualIP.To16()
		if actualIP16 == nil || ipRangeBegin16 == nil || ipRangeEnd16 == nil {
			log.Warn().Msg("IP conversion failed")
			return true
		}

		if bytes.Compare(actualIP16, ipRangeBegin16) >= 0 && bytes.Compare(actualIP16, ipRangeEnd16) <= 0 {
			return true
		}
	}

	return false
}
