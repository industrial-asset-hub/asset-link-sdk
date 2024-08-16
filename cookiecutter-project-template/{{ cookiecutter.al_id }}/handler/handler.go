/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: {{cookiecutter.company}}
 *
 */

package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync/atomic"

	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/iah-discovery"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/model"
	"github.com/rs/zerolog/log"
)

// Implements the features of the DCD.
// see
type AssetLinkImplementation struct {
	discoveryJobRunning bool
}

var lastSerialNumber = atomic.Int64{}

// Implementation of the Discovery feature

// Start implements the function, which is called, with the
// dcdconnection method is executed
func (m *AssetLinkImplementation) Start(jobId uint32, deviceChannel chan []*generated.DiscoveredDevice, err chan error, filters map[string]string) {

	// Check if job is already running
	// We currently support here only one running job
	if m.discoveryJobRunning {
		errMsg := "Discovery job is already running"
		err <- errors.New(errMsg)
	}
	defer close(deviceChannel)

	// Thus, this function is executed as Goroutine,
	// and the gRPC Server methods blocks, until the job is started, we assume at this point,
	// that the discover job is started successfully
	err <- nil
	m.discoveryJobRunning = true

	//
	// At here our custom logic to discover our device and fetch the needed information
	//

	// Fillup the device information
	device := model.NewDevice("EthernetDevice", "My First Ethernet Device")

	product := "{{ cookiecutter.al_name }}"
	orderNumber := "PRODUCT-ONE"
	productVersion := "1.0.0"
	vendorName := "{{ cookiecutter.company }}"
	lastSerialNumber.Add(1)
	serialNumber := fmt.Sprint(lastSerialNumber.Load())
	productUri := fmt.Sprintf("urn:%s/%s", strings.ReplaceAll(vendorName, " ", "_"), strings.ReplaceAll(product, " ", "_"), serialNumber)
	device.AddNameplate(
		vendorName,
		productUri,
		orderNumber,
		product,
		productVersion,
		serialNumber)

	device.AddCapabilities("firmware_update", false)

	randomMacAddress := generateRandomMacAddress()
	id := device.AddNic("eth0", randomMacAddress)
	device.AddIPv4(id, "192.168.0.1", "255.255.255.0", "")

	// Convert and stream device to upstream system
	discoveredDevice := device.ConvertToDiscoveredDevice()
	devices := make([]*generated.DiscoveredDevice, 0)
	devices = append(devices, discoveredDevice)
	deviceChannel <- devices

	m.discoveryJobRunning = false
	log.Debug().
		Msg("Start function exiting")
}

func (m *AssetLinkImplementation) FilterTypes(filterTypesChannel chan []*generated.SupportedFilter) {
	filterTypes := make([]*generated.SupportedFilter, 0)
	filterTypes = append(filterTypes, &generated.SupportedFilter{
		Key:      "type",
		Datatype: generated.VariantType_VT_BYTES,
	})
	filterTypesChannel <- filterTypes
}

func (m *AssetLinkImplementation) FilterOptions(filterOptionsChannel chan []*generated.SupportedOption) {
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
