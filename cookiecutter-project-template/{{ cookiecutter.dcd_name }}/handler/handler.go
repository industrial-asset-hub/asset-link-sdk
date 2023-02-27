/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: {{cookiecutter.company}}
 *
 */

package handler

import (
	"errors"
	"strconv"
	"time"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/model"
	softwareUpdate "code.siemens.com/common-device-management/utils/go-modules/firmwareupdate.git/pkg/firmware-update"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Implements the features of the DCD.
// see
type DCDImplementation struct {
	discoveryJobCancelationToken chan uint32
	discoveryJobRunning          bool
}

// Implementation of the Discovery feature

// Start implements the function, which is called, with the
// grpc method is executed
func (m *DCDImplementation) Start(jobId uint32, deviceInfoReply chan model.DeviceInfo) error {
	log.Info().
		Msg("Start Discovery")

	log.Debug().
		Bool("running", m.discoveryJobRunning).
		Msg("Discovery running?")
	defer close(deviceInfoReply)

	// Check if job is already running
	// We currently support here only one running job
	if m.discoveryJobRunning {
		errMsg := "Discovery job is already running"
		return errors.New(errMsg)
	}

	m.discoveryJobRunning = true
	m.discoveryJobCancelationToken = make(chan uint32)

	// For loop, just to simulation some time for an discovery job
	select {
	case cancelationJobId := <-m.discoveryJobCancelationToken:
		log.Debug().
			Uint32("Job Id", cancelationJobId).
			Msg("Received cancel request")
		m.discoveryJobRunning = false
		return nil
	default:

		// Default Device Information structure
		deviceInfo := model.DeviceInfo{}
		deviceInfo.Vendor = "{{ cookiecutter.company }}"
		deviceInfo.DeviceFamily = "{{ cookiecutter.dcd_name }}"
		// Exact Device Type e.g. CPU 1516-3 PN/DP
		deviceInfo.DeviceDescription = "{{ cookiecutter.company }} - DeviceOne"
		deviceInfo.ArticleNumber = ""
		deviceInfo.DeviceHwVersion = "1.0.1"
		deviceInfo.DeviceSwVersion = "0.2.0"
		deviceInfo.PasswordProtected = false
		deviceInfo.SerialNumber = model.SerialNumber(uuid.NewString())
		stationName := "Hello {{ cookiecutter.author_name }}, our first device"
		deviceType := model.PropertiesDeviceTypeNative
		connectivityType := model.PropertiesConnectivityStatusOnline
		opMode := model.PropertiesOperatingModeRun

		// Network Interfaces
		nic1Name := model.NicIdentifier("eth0")
		nic1MAC := "ab:cd:ef:ab:cd:ef"
		nic1IP := "192.168.0.99"
		nic1Netmask := "192.168.0.99"
		nic1DefaultGateway := ""
		nic2Name := model.NicIdentifier("eth2")
		nic2MAC := "ab:cd:ef:ab:cd:ef"
		nic2IP := "10.0.0.2"
		nic2Netmask := "255.255.0.0"
		nic2DefaultGateway := "10.0.0.1"

		nic1 := model.DeviceInfoNicsElem{MacAddress: &nic1MAC, NicIdentifier: &nic1Name}
		nic2 := model.DeviceInfoNicsElem{MacAddress: &nic2MAC, NicIdentifier: &nic2Name}

		deviceInfo.Nics = append(deviceInfo.Nics, nic1, nic2)
		ipSettingsNic1 := model.PropertiesIpConfigurationsElemIpSettingsElem{
			IpAddress:      nic1IP,
			SubnetMask:     nic1Netmask,
			DefaultGateway: nic1DefaultGateway}
		ipSettingsNic2 := model.PropertiesIpConfigurationsElemIpSettingsElem{
			IpAddress:      nic2IP,
			SubnetMask:     nic2Netmask,
			DefaultGateway: nic2DefaultGateway}

		var IPSettings1 []model.PropertiesIpConfigurationsElemIpSettingsElem
		IPSettings1 = append(IPSettings1, ipSettingsNic1)
		var IPSettings2 []model.PropertiesIpConfigurationsElemIpSettingsElem
		IPSettings2 = append(IPSettings2, ipSettingsNic2)
		IPInfoNic1 := model.PropertiesIpConfigurationsElem{NicIdentifier: &nic1Name, IpSettings: IPSettings1}
		IPInfoNic2 := model.PropertiesIpConfigurationsElem{NicIdentifier: &nic2Name, IpSettings: IPSettings2}
		var IPInfo []model.PropertiesIpConfigurationsElem
		IPInfo = append(IPInfo, IPInfoNic1, IPInfoNic2)

		// Capabilities of ot the device
		var ptrTrue bool = true
		var ptrFalse bool = false
		cap := model.Capabilities{
			FirmwareUpdate: &ptrTrue,
			ProgramUpdate:  &ptrTrue,
			Backup:         &ptrFalse,
			Restore:        &ptrFalse,
			ResetToFactory: &ptrTrue,
		}
		deviceInfo.Capabilities = &cap

		var connectedTo []model.PropertiesConnectedToElem
		connectedTo = append(connectedTo, model.PropertiesConnectedToElem{
			Name:          "Connection to Profibus",
			InterfaceType: model.PropertiesConnectedToElemInterfaceTypeProfibus,
			Devices:       []string{"device-1", "device-2"},
		})

		// Capabilities of ot the device
		module1Cap := model.Capabilities{
			FirmwareUpdate: &ptrTrue,
			ProgramUpdate:  &ptrFalse,
			Backup:         &ptrFalse,
			Restore:        &ptrFalse,
			ResetToFactory: &ptrFalse,
		}

		module1 := model.Module{
			ArticleNumber:   "Module 1 Article number",
			Capabilities:    module1Cap,
			Description:     "Module 1 description",
			DeviceHwVersion: "Module 1 HW version",
			DeviceSwVersion: "Module 1 SW version",
			Name:            "Module 1 Name",
			SerialNumber:    "Module 1 Serial Number",
			Slot:            model.Slot(1),
			SubSlot:         model.SubSlot(1),
		}

		// Values located under properties
		runtimeMode := model.PropertiesRuntimeModeNormal
		properties := model.Properties{
			ConnectivityStatus: connectivityType,
			RuntimeMode:        &runtimeMode,
			IpConfigurations:   IPInfo,
			ConnectedTo:        connectedTo,
			StationName:        stationName,
			ProfinetName:       &stationName,
			DeviceType:         deviceType,
			OperatingMode:      &opMode,
		}
		deviceInfo.Properties = &properties
		deviceInfo.Modules = append(deviceInfo.Modules, module1)

		deviceInfoReply <- deviceInfo
		time.Sleep(1000 * time.Millisecond)
	}

	// Close channel, to signal that no more data is to be transfered
	m.discoveryJobRunning = false
	log.Debug().
		Msg("Start function exiting")

	return nil
}

func (m *DCDImplementation) Cancel(jobId uint32) error {
	log.Info().
		Uint32("Job Id", jobId).
		Msg("Cancel Discovery")

	if m.discoveryJobRunning {
		log.Info().
			Msg("Cancel received. Sending notification to stop current discovery job")
		m.discoveryJobCancelationToken <- jobId
	} else {
		log.Warn().
			Msg("Cancel received, but no discovery is running")
	}
	log.Debug().
		Msg("Cancel function exiting")
	return nil

}

// Implementation of the Software Update feature
func (m *DCDImplementation) Update(jobId string, deviceId string, metaData []*softwareUpdate.FirmwareMetaData, statusChannel chan *softwareUpdate.FirmwareReply) error {

	log.Info().
		Str("Job Id", jobId).
		Str("Device Id", deviceId).
		Msg("Firmware Update Implementation")

	for _, d := range metaData {
		log.Debug().
			Str("Metadata", d.String()).
			Msg("Metadata received")
	}

	// Emulate an Software Update
	for i := 0; i <= 100; i += 25 {
		progressInfo := softwareUpdate.ProgressInfo{
			Operation:  softwareUpdate.FirmwareOperation_DOWNLOAD,
			Percentage: strconv.Itoa(i)}

		UpdateStatus := softwareUpdate.FirmwareReply{
			DeviceId:       deviceId,
			JobId:          jobId,
			ProgressStatus: &progressInfo,
			Status:         softwareUpdate.FirmwareStatus_IN_PROGRESS,
			ErrorMsg:       ""}

		// Report success after the last iteration
		if i >= 100 {
			progressInfo.Percentage = "100"
			progressInfo.Operation = softwareUpdate.FirmwareOperation_INSTALL
			UpdateStatus.Status = softwareUpdate.FirmwareStatus_SUCCESS
		}

		statusChannel <- &UpdateStatus

		// Wait until next iteration
		time.Sleep(1 * time.Second)
	}
	defer close(statusChannel)
	log.Debug().
		Msg("Update function exiting")
	return nil

}
