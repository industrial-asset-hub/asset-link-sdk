/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package reference

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"

	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/iah-discovery"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/model"
	"github.com/rs/zerolog/log"
)

// Implements the features of the DCD.
// see
type ReferenceClassDriver struct {
	discoveryJobCancelationToken chan uint32
	discoveryJobRunning          bool
}

var lastSerialNumber = atomic.Int64{}

// Implementation of the Discovery feature

// Start implements the function, which is called, with the
// grpc method is executed
func (m *ReferenceClassDriver) Start(jobId uint32, deviceChannel chan []*generated.DiscoveredDevice, err chan error, filters map[string]string) {
	log.Info().
		Msg("Start Discovery")

	log.Debug().
		Bool("running", m.discoveryJobRunning).
		Msg("Discovery running?")

	defer close(deviceChannel)
	// Check if job is already running
	// We currently support here only one running job
	if m.discoveryJobRunning {
		errMsg := "Discovery job is already running"
		err <- errors.New(errMsg)
	}

	// Thus, this function is executed as Goroutine,
	// and the gRPC Server methods blocks, until the job is started, we assume at this point,
	// that the discover job is started successfully
	err <- nil

	m.discoveryJobRunning = true
	m.discoveryJobCancelationToken = make(chan uint32)

	// Just provide a static asset
	name := "Device"
	lastSerialNumber.Add(1)
	serialNumber := fmt.Sprint(lastSerialNumber.Load())

	deviceInfo := model.NewDevice("EthernetDevice", name)
	deviceInfo.AddNameplate(
		"Siemens AG",
		"MyOrderNumber",
		"cdm-reference-dcd-test2",
		"1.0.0",
		serialNumber)

	deviceInfo.AddSoftware("firmware", "1.2.5")
	randomMacAddress := generateRandomMacAddress()
	id := deviceInfo.AddNic("eth0", randomMacAddress)
	deviceInfo.AddIPv4(id, "192.168.0.1", "255.255.255.0", "")
	deviceInfo.AddIPv4(id, "10.0.0.1", "255.255.255.0", "")

	discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()
	devices := make([]*generated.DiscoveredDevice, 0)
	devices = append(devices, discoveredDevice)
	deviceChannel <- devices

	m.discoveryJobRunning = false
	log.Debug().
		Msg("Start function exiting")

}

func (m *ReferenceClassDriver) Cancel(jobId uint32) error {
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

func (m *ReferenceClassDriver) FilterTypes(filterTypesChannel chan []*generated.SupportedFilter) {
	filterTypes := make([]*generated.SupportedFilter, 0)
	filterTypes = append(filterTypes, &generated.SupportedFilter{
		Key:      "type",
		Datatype: generated.VariantType_VT_BYTES,
	})
	filterTypesChannel <- filterTypes
}

func (m *ReferenceClassDriver) FilterOptions(filterOptionsChannel chan []*generated.SupportedOption) {
	filterOptions := make([]*generated.SupportedOption, 0)
	filterOptions = append(filterOptions, &generated.SupportedOption{
		Key:      "option",
		Datatype: generated.VariantType_VT_BOOL,
	})
	filterOptionsChannel <- filterOptions
}

func generateRandomMacAddress() string {
	r := rand.Uint64()
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		0x00, 0x16, 0x3e,
		byte(r>>8),
		byte(r>>16),
		byte(r>>24))
}
