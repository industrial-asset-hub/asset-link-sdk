/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: {{cookiecutter.company}}
 *
 */

package handler

import (
	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/iah-discovery"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/model"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

// Implements the features of the DCD.
// see
type AssetLinkImplementation struct {
	discoveryJobCancelationToken chan uint32
	discoveryJobRunning          bool
}

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

	// Thus, this function is executed as Goroutine,
	// and the gRPC Server methods blocks, until the job is started, we assume at this point,
	// that the discover job is started successfully
	err <- nil
	m.discoveryJobRunning = true
	m.discoveryJobCancelationToken = make(chan uint32)
	device := model.NewDevice("Profinet")
	timestamp := model.CreateTimestamp()

	Name := "Device"
	device.Name = &Name
	product := "{{ cookiecutter.al_name }}"
	version := "1.0.0"
	vendorName := "{{ cookiecutter.company }}"
	//serialNumber := uuid.NewString()
	//serialNumber := "sn"
	lastSerialNumber.Add(1)
	serialNumber := fmt.Sprint(lastSerialNumber.Load())
	vendor := model.Organization{
		Address:        nil,
		AlternateNames: nil,
		ContactPoint:   nil,
		Id:             "",
		Name:           &vendorName,
	}
	productSerialidentifier := model.ProductSerialIdentifier{
		IdentifierType:        nil,
		IdentifierUncertainty: nil,
		ManufacturerProduct: &model.Product{
			Id:             "",
			Manufacturer:   &vendor,
			Name:           nil,
			ProductId:      &product,
			ProductVersion: &version,
		},
		SerialNumber: &serialNumber,
	}
	device.ProductInstanceIdentifier = &productSerialidentifier

	randomMacAddress := generateRandomMacAddress()
	identifierUncertainty := 1
	device.MacIdentifiers = append(device.MacIdentifiers, model.MacIdentifier{
		MacAddress:            &randomMacAddress,
		IdentifierUncertainty: &identifierUncertainty,
	})

	connectionPointType := "Ipv4Connectivity"
	Ipv4Address := "192.168.0.1"
	Ipv4NetMask := "255.255.255.0"
	Ipv4Connectivity := model.Ipv4Connectivity{
		ConnectionPointType:     &connectionPointType,
		Id:                      "1",
		InstanceAnnotations:     nil,
		Ipv4Address:             &Ipv4Address,
		NetworkMask:             &Ipv4NetMask,
		RelatedConnectionPoints: nil,
		RouterIpv4Address:       nil,
	}
	device.ConnectionPoints = append(device.ConnectionPoints, Ipv4Connectivity)

	state := model.ManagementStateValuesUnknown
	State := model.ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}
	device.ManagementState = State

	reachabilityStateValue := model.ReachabilityStateValuesReached
	reachabilityState := model.ReachabilityState{
		StateTimestamp: &timestamp,
		StateValue:     &reachabilityStateValue,
	}
	device.ReachabilityState = &reachabilityState

	discoveredDevice := device.ConvertToDiscoveredDevice()
	devices := make([]*generated.DiscoveredDevice, 0)
	devices = append(devices, discoveredDevice)
	deviceChannel <- devices
	m.discoveryJobRunning = false
	log.Debug().
		Msg("Start function exiting")
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

func generateRandomMacAddress() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		0x00, 0x16, 0x3e,
		byte(lastSerialNumber.Load()>>8),
		byte(lastSerialNumber.Load()>>16),
		byte(lastSerialNumber.Load()>>24))
}