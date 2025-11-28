/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package simdevices

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SimulatedDeviceState string

const (
	// regular operation
	StateActive SimulatedDeviceState = "active" // device state during regular operation

	// device discovery
	StateReading SimulatedDeviceState = "reading" // reading device information (e.g., during discovery)

	// update management
	StateUpdating SimulatedDeviceState = "updating" // installing update on a device
	StateBooting  SimulatedDeviceState = "booting"  // rebooting device (e.g., to activate a prior update)

	// artifact management
	StateGetting SimulatedDeviceState = "getting" // getting artifact (e.g., configuration) from device
	StateSetting SimulatedDeviceState = "setting" // setting artifact (e.g., configuration) on device
)

var (
	ErrInvalidInterface  = errors.New("invalid interface")
	ErrDeviceNotFound    = errors.New("device not found")
	ErrSubDeviceNotFound = errors.New("sub-device not found")
	ErrUnauthenticated   = errors.New("invalid credentials")
)

type SimulatedDevice interface {
	GetDeviceName() string
	GetManufacturer() string
	GetProductDesignation() string
	GetArticleNumber() string
	GetSerialNumber() string
	GetHardwareVersion() string
	GetActiveFirmwareVersion() string
	GetInstalledFirmwareVersion() string
	IsUpdateSupported() bool
	GetAssetLinkNIC() string
	GetDeviceNIC() string
	GetMacAddress() string
	GetIpDevice() string
	GetIpNetmask() string
	GetIpRoute() string
	GetIDLink() string

	GetDeviceAddress() SimulatedDeviceAddress

	GetSubDeviceID() int
	GetSubDevices() []SimulatedDevice
	GetParentDevice() SimulatedDevice

	UpdateFirmware(artefactFilename string) error
	RebootDevice() error

	GetConfig(configFilename string) error
	SetConfig(configFilename string) error
}

type SimulatedDeviceAddress struct {
	AssetLinkNIC string `json:"alNic"`
	DeviceIP     string `json:"ipAddress"`
	SubDeviceID  int    `json:"subDeviceID"` // -1 for top-level devices
}

type SimulatedDeviceCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type simulatedDeviceInfo struct {
	UniqueDeviceID           string                      `json:"unique_device_id"` // for internal use only
	DeviceName               string                      `json:"device_name"`
	Manufacturer             string                      `json:"manufacturer"`
	ProductDesignation       string                      `json:"product_designation"`
	ArticleNumber            string                      `json:"article_number"`
	SerialNumber             string                      `json:"serial_number"`
	HardwareVersion          string                      `json:"hardware_version"`
	ActiveFirmwareVersion    string                      `json:"active_firmware_version"`
	InstalledFirmwareVersion string                      `json:"installed_firmware_version"`
	UpdateSupport            bool                        `json:"update_support"`
	AssetLinkNIC             string                      `json:"al_nic"`
	DeviceNIC                string                      `json:"device_nic"`
	MacAddress               string                      `json:"mac_address"`
	IpDevice                 string                      `json:"ip_device"`
	IpNetmask                string                      `json:"ip_netmask"`
	IpRoute                  string                      `json:"ip_route"`
	DeviceState              SimulatedDeviceState        `json:"device_state"`
	SubDeviceID              int                         `json:"sub_device_id"` // -1 for top-level devices
	SubDevices               []*simulatedDeviceInfo      `json:"sub_devices"`
	ParentDevice             *simulatedDeviceInfo        `json:"-"` // parent device is not serialized
	credentials              *SimulatedDeviceCredentials `json:"-"` // credentials are not serialized
}

type firmwareFile struct {
	ArtefactType       string `json:"artefact_type"`
	Manufacturer       string `json:"manufacturer"`
	ProductDesignation string `json:"product_designation"`
	FirmwareVersion    string `json:"firmware_version"`
}

type configFile struct {
	ArtefactType string `json:"artefact_type"`
	DeviceName   string `json:"device_name"`
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

func assertLocked() {
	if simLock.TryLock() {
		panic("assertion failed: simLock mutex not locked")
	}
}

func newSimulatedDevice(alNIC, serial, mac, ip, name string, credentials *SimulatedDeviceCredentials) *simulatedDeviceInfo {
	uid := uuid.New()
	return &simulatedDeviceInfo{
		UniqueDeviceID:           uid.String(),
		DeviceName:               name,
		Manufacturer:             "Siemens AG",
		ProductDesignation:       "Simulated Device",
		ArticleNumber:            "AN0123456789",
		HardwareVersion:          "3",
		ActiveFirmwareVersion:    "1.0.0",
		InstalledFirmwareVersion: "1.0.0",
		UpdateSupport:            true,
		AssetLinkNIC:             alNIC,
		DeviceNIC:                "enp0",
		MacAddress:               mac,
		IpDevice:                 ip,
		IpNetmask:                "255.255.255.0",
		IpRoute:                  "",
		SerialNumber:             serial,
		DeviceState:              StateActive,
		SubDevices:               []*simulatedDeviceInfo{},
		SubDeviceID:              -1,
		ParentDevice:             nil, // parent device is nil for top-level devices
		credentials:              credentials,
	}
}

func newSimulatedSubDevice(name, serial, mac, ip string) *simulatedDeviceInfo {
	uid := uuid.New()
	return &simulatedDeviceInfo{
		UniqueDeviceID:           uid.String(),
		DeviceName:               name,
		Manufacturer:             "Siemens AG",
		ProductDesignation:       "Simulated Sub Device",
		ArticleNumber:            "AN9876543210",
		HardwareVersion:          "5",
		ActiveFirmwareVersion:    "1.2.3",
		InstalledFirmwareVersion: "1.2.3",
		UpdateSupport:            false,
		AssetLinkNIC:             "",
		DeviceNIC:                "enp0",
		MacAddress:               mac,
		IpDevice:                 ip,
		IpNetmask:                "255.255.255.0",
		IpRoute:                  "",
		SerialNumber:             serial,
		DeviceState:              StateActive,
		SubDevices:               []*simulatedDeviceInfo{},
		SubDeviceID:              -1,
		ParentDevice:             nil, // parent device will be set when appending to a parent
		credentials:              nil, // sub-devices do not have credentials
	}
}

func (d *simulatedDeviceInfo) appendSimulatedSubDevice(subDevice *simulatedDeviceInfo) {
	if d.ParentDevice != nil {
		log.Error().Msg("Cannot append sub-device to a sub-device")
		return
	}

	subDevice.SubDeviceID = len(d.SubDevices) // assign sub-device ID
	subDevice.ParentDevice = d                // set the parent device for the sub-device
	d.SubDevices = append(d.SubDevices, subDevice)
}

func StartSimulatedDevices(visuServerAddress string) {
	credentials := &SimulatedDeviceCredentials{
		Username: "admin",
		Password: "admin", // storing the password in plain text is obviously not secure, but this is just a simulation
	}

	simulatedDevicesEth0 = []*simulatedDeviceInfo{}
	simulatedDevicesEth1 = []*simulatedDeviceInfo{}

	// Eth 0 (device A0)
	deviceA0 := newSimulatedDevice(interfaceEth0, "SN123450000", "00:16:3e:00:00:00", "192.168.0.10", "Simulated Device A0", nil)
	simulatedDevicesEth0 = append(simulatedDevicesEth0, deviceA0)

	// Eth 0 (device A1)
	deviceA1 := newSimulatedDevice(interfaceEth0, "SN123450001", "00:16:3e:00:00:01", "192.168.0.11", "Simulated Device A1", nil)
	simulatedDevicesEth0 = append(simulatedDevicesEth0, deviceA1)

	// Eth 0 (device A2)
	deviceA2 := newSimulatedDevice(interfaceEth0, "SN123450002", "00:16:3e:00:00:02", "192.168.0.12", "Simulated Device A2", credentials)
	simulatedDevicesEth0 = append(simulatedDevicesEth0, deviceA2)

	// Eth 1 (device B0)
	deviceB0 := newSimulatedDevice(interfaceEth1, "SN123450100", "00:16:3e:00:01:00", "192.168.1.10", "Simulated Device B0", nil)
	simulatedDevicesEth1 = append(simulatedDevicesEth1, deviceB0)

	// Eth 1 (device B1 and subdevices)
	deviceB1 := newSimulatedDevice(interfaceEth1, "SN123450101", "00:16:3e:00:01:01", "192.168.1.11", "Simulated Device B1", nil)
	subDeviceB10 := newSimulatedSubDevice("Simulated Sub Device B1-0", "SN123450101-0", "00:16:3e:00:01:10", "192.168.1.12")
	deviceB1.appendSimulatedSubDevice(subDeviceB10)
	subDeviceB11 := newSimulatedSubDevice("Simulated Sub Device B1-1", "SN123450101-1", "00:16:3e:00:01:11", "192.168.1.13")
	deviceB1.appendSimulatedSubDevice(subDeviceB11)
	subDeviceB12 := newSimulatedSubDevice("Simulated Sub Device B1-2", "SN123450101-2", "00:16:3e:00:01:12", "192.168.1.14")
	deviceB1.appendSimulatedSubDevice(subDeviceB12)
	simulatedDevicesEth1 = append(simulatedDevicesEth1, deviceB1)

	// Eth 1 (device B2)
	deviceB2 := newSimulatedDevice(interfaceEth1, "SN123450102", "00:16:3e:00:01:02", "192.168.1.12", "Simulated Device B2", credentials)
	simulatedDevicesEth1 = append(simulatedDevicesEth1, deviceB2)

	visuActive = false
	if visuServerAddress != "" && !testing.Testing() {
		startDeviceVisualization(visuServerAddress)
		visuActive = true
	}
}

func getDeviceListCopy(lockTaken bool) []simulatedDeviceInfo { // for internal use only
	if lockTaken {
		assertLocked()
	} else {
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

func handleDeviceChanges() {
	assertLocked()

	if visuActive {
		deviceList := getDeviceListCopy(true)
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

func (d *simulatedDeviceInfo) GetActiveFirmwareVersion() string {
	simLock.Lock()
	defer simLock.Unlock()

	return d.ActiveFirmwareVersion
}

func (d *simulatedDeviceInfo) GetInstalledFirmwareVersion() string {
	simLock.Lock()
	defer simLock.Unlock()

	return d.InstalledFirmwareVersion
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

func (d *simulatedDeviceInfo) GetSubDeviceID() int {
	return d.SubDeviceID
}

func (d *simulatedDeviceInfo) GetSubDevices() []SimulatedDevice {
	simLock.Lock()
	defer simLock.Unlock()

	subDevices := make([]SimulatedDevice, len(d.SubDevices))
	for i, subDevice := range d.SubDevices {
		subDevice.DeviceState = StateReading
		handleDeviceChanges()

		simulateCostlyOperation(2 * time.Second) // simulate reading information from device

		subDevice.DeviceState = StateActive
		handleDeviceChanges()

		subDevices[i] = subDevice
	}
	return subDevices
}

func (d *simulatedDeviceInfo) GetParentDevice() SimulatedDevice {
	return d.ParentDevice
}

func (d *simulatedDeviceInfo) GetDeviceAddress() SimulatedDeviceAddress {
	if d.ParentDevice != nil {
		// handle sub-devices
		deviceAddress := SimulatedDeviceAddress{
			AssetLinkNIC: d.ParentDevice.AssetLinkNIC,
			DeviceIP:     d.ParentDevice.IpDevice,
			SubDeviceID:  d.SubDeviceID,
		}
		return deviceAddress
	}

	// handle top-level devices
	deviceAddress := SimulatedDeviceAddress{
		AssetLinkNIC: d.AssetLinkNIC,
		DeviceIP:     d.IpDevice,
		SubDeviceID:  -1,
	}
	return deviceAddress
}

func (d *simulatedDeviceInfo) hasIPInRange(ipRange string) bool {
	deviceIPs := []string{}
	deviceIPs = append(deviceIPs, d.IpDevice) // there is currently only one IP address per device

	return containsIpInRange(ipRange, deviceIPs)
}

func (d *simulatedDeviceInfo) checkCredentials(credentials *SimulatedDeviceCredentials) bool {
	if d.credentials == nil {
		return true
	}

	if d.credentials.Username == "" && d.credentials.Password == "" {
		return true
	}

	if credentials == nil {
		return false
	}

	if d.credentials.Username == credentials.Username && d.credentials.Password == credentials.Password {
		return true
	}

	return false
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

func simulateCostlyOperation(duration time.Duration) {
	if testing.Testing() {
		return
	}

	ciValue, ciExists := os.LookupEnv("CI")
	if ciExists && ciValue == "true" {
		return
	}

	time.Sleep(duration)
}

func connectToDevice(deviceAddress SimulatedDeviceAddress, credentials *SimulatedDeviceCredentials) (*simulatedDeviceInfo, error) {
	assertLocked()

	var simulatedDevices []*simulatedDeviceInfo

	switch deviceAddress.AssetLinkNIC {
	case interfaceEth0:
		simulatedDevices = simulatedDevicesEth0
	case interfaceEth1:
		simulatedDevices = simulatedDevicesEth1
	default:
		return nil, ErrInvalidInterface
	}

	simulateCostlyOperation(1 * time.Second) // simulate connecting to device

	for _, device := range simulatedDevices {
		if device.IpDevice == deviceAddress.DeviceIP {
			if !device.checkCredentials(credentials) {
				return nil, ErrUnauthenticated
			}

			if deviceAddress.SubDeviceID < 0 {
				return device, nil // return top-level device
			}

			if deviceAddress.SubDeviceID < len(device.SubDevices) {
				return device.SubDevices[deviceAddress.SubDeviceID], nil // return sub-device
			}

			return nil, ErrSubDeviceNotFound
		}
	}

	return nil, ErrDeviceNotFound
}
