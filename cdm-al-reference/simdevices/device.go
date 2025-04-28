/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package simdevices

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SimulatedDeviceState string

const (
	StateActive  SimulatedDeviceState = "active"
	StateReading SimulatedDeviceState = "reading"
)

type SimulatedDevice interface {
	GetDeviceName() string
	GetManufacturer() string
	GetProductDesignation() string
	GetArticleNumber() string
	GetSerialNumber() string
	GetHardwareVersion() string
	GetFirmwareVersion() string
	IsUpdateSupported() bool
	GetAssetLinkNIC() string
	GetDeviceNIC() string
	GetMacAddress() string
	GetIpDevice() string
	GetIpNetmask() string
	GetIpRoute() string
	GetIDLink() string
}

type simulatedDeviceInfo struct {
	UniqueDeviceID     string               `json:"unique_device_id"` // for internal use only
	DeviceName         string               `json:"device_name"`
	Manufacturer       string               `json:"manufacturer"`
	ProductDesignation string               `json:"product_designation"`
	ArticleNumber      string               `json:"article_number"`
	SerialNumber       string               `json:"serial_number"`
	HardwareVersion    string               `json:"hardware_version"`
	FirmwareVersion    string               `json:"firmware_version"`
	UpdateSupport      bool                 `json:"update_support"`
	AssetLinkNIC       string               `json:"al_nic"`
	DeviceNIC          string               `json:"device_nic"`
	MacAddress         string               `json:"mac_address"`
	IpDevice           string               `json:"ip_device"`
	IpNetmask          string               `json:"ip_netmask"`
	IpRoute            string               `json:"ip_route"`
	DeviceState        SimulatedDeviceState `json:"device_state"`
}

const (
	interfaceEth0 = "eth0"
	interfaceEth1 = "eth1"
)

var (
	visuActive bool

	simulatedDevicesEth0 []*simulatedDeviceInfo
	simulatedDevicesEth1 []*simulatedDeviceInfo

	simLock sync.Mutex
)

func newSimulatedDevice(alNIC, serial, mac, ip, name string) *simulatedDeviceInfo {
	uid := uuid.New()
	return &simulatedDeviceInfo{
		UniqueDeviceID:     uid.String(),
		DeviceName:         name,
		Manufacturer:       "Siemens AG",
		ProductDesignation: "Simulated Device",
		ArticleNumber:      "AN0123456789",
		HardwareVersion:    "3",
		FirmwareVersion:    "1.0.0",
		UpdateSupport:      true,
		AssetLinkNIC:       alNIC,
		DeviceNIC:          "enp0",
		MacAddress:         mac,
		IpDevice:           ip,
		IpNetmask:          "255.255.255.0",
		IpRoute:            "",
		SerialNumber:       serial,
		DeviceState:        StateActive,
	}
}

func StartSimulatedDevices(visuServerAddress string) {
	simulatedDevicesEth0 = append(simulatedDevicesEth0, newSimulatedDevice(interfaceEth0, "SN123450000", "00:16:3e:00:00:00", "192.168.0.10", "Simulated Device A0"))
	simulatedDevicesEth0 = append(simulatedDevicesEth0, newSimulatedDevice(interfaceEth0, "SN123450001", "00:16:3e:00:00:01", "192.168.0.11", "Simulated Device A1"))

	simulatedDevicesEth1 = append(simulatedDevicesEth1, newSimulatedDevice(interfaceEth1, "SN123450100", "00:16:3e:00:01:00", "192.168.1.10", "Simulated Device B0"))
	simulatedDevicesEth1 = append(simulatedDevicesEth1, newSimulatedDevice(interfaceEth1, "SN123450101", "00:16:3e:00:01:01", "192.168.1.11", "Simulated Device B1"))

	visuActive = false
	if visuServerAddress != "" {
		startDeviceVisualization(visuServerAddress)
		visuActive = true
	}
}

func getDeviceListCopy(lockTaken bool) []simulatedDeviceInfo { // for internal use only
	if !lockTaken {
		simLock.Lock()
		defer simLock.Unlock()
	}

	deviceList := make([]simulatedDeviceInfo, 0, len(simulatedDevicesEth0)+len(simulatedDevicesEth1))

	for _, device := range simulatedDevicesEth0 {
		deviceList = append(deviceList, *device)
	}

	for _, device := range simulatedDevicesEth1 {
		deviceList = append(deviceList, *device)
	}

	return deviceList
}

func handleDeviceChanges(lockTaken bool) {
	if visuActive {
		deviceList := getDeviceListCopy(lockTaken)
		broadcastDeviceUpdates(deviceList)
	}
}

func (d *simulatedDeviceInfo) GetDeviceName() string {
	return d.DeviceName
}

func (d *simulatedDeviceInfo) GetManufacturer() string {
	return d.Manufacturer
}

func (d *simulatedDeviceInfo) GetProductDesignation() string {
	return d.ProductDesignation
}

func (d *simulatedDeviceInfo) GetArticleNumber() string {
	return d.ArticleNumber
}

func (d *simulatedDeviceInfo) GetSerialNumber() string {
	return d.SerialNumber
}

func (d *simulatedDeviceInfo) GetHardwareVersion() string {
	return d.HardwareVersion
}

func (d *simulatedDeviceInfo) GetFirmwareVersion() string {
	simLock.Lock()
	defer simLock.Unlock()

	return d.FirmwareVersion
}

func (d *simulatedDeviceInfo) IsUpdateSupported() bool {
	return d.UpdateSupport
}

func (d *simulatedDeviceInfo) GetAssetLinkNIC() string {
	return d.AssetLinkNIC
}

func (d *simulatedDeviceInfo) GetDeviceNIC() string {
	return d.DeviceNIC
}

func (d *simulatedDeviceInfo) GetMacAddress() string {
	return d.MacAddress
}

func (d *simulatedDeviceInfo) GetIpDevice() string {
	return d.IpDevice
}

func (d *simulatedDeviceInfo) GetIpNetmask() string {
	return d.IpNetmask
}

func (d *simulatedDeviceInfo) GetIpRoute() string {
	return d.IpRoute
}

func (d *simulatedDeviceInfo) GetIDLink() string {
	// generate ID link from device data
	idLink := fmt.Sprintf("https://industrial-assets.io?1P=%s&S=%s", d.ArticleNumber, d.SerialNumber)
	return idLink
}

func (d *simulatedDeviceInfo) hasIPInRange(ipRange string) bool {
	deviceIPs := []string{}
	deviceIPs = append(deviceIPs, d.IpDevice) // there is currently only one IP address per device

	return containsIpInRange(ipRange, deviceIPs)
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
