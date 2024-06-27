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
	deviceInfo := model.NewDevice("Profinet")
	timestamp := model.CreateTimestamp()

	Name := "Device"
	deviceInfo.Name = &Name
	product := "cdm-reference-dcd"
	version := "1.0.0"
	vendorName := "Siemens AG"
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
			Name:           &Name,
			ProductId:      &product,
			ProductVersion: &version,
		},
		SerialNumber: &serialNumber,
	}
	deviceInfo.ProductInstanceIdentifier = &productSerialidentifier
	randomMacAddress := generateRandomMacAddress()
	identifierUncertainty := 1
	deviceInfo.MacIdentifiers = append(deviceInfo.MacIdentifiers, model.MacIdentifier{
		MacAddress:            &randomMacAddress,
		IdentifierUncertainty: &identifierUncertainty,
	})
	// connectionPoint := "ethernet"
	relatedConnectionPoint := model.RelatedConnectionPoint{
		ConnectionPoint:    nil,
		CustomRelationship: nil,
	}
	relatedConnectionPoints := make([]model.RelatedConnectionPoint, 0)
	relatedConnectionPoints = append(relatedConnectionPoints, relatedConnectionPoint)
	// not using it due to unexpected build issue
	//connectionPointType := "Ipv4Connectivtestity"
	Ipv4Address := "192.168.0.1"
	Ipv4NetMask := "255.255.255.0"
	routerIpv6Address := []string{"fd12:3456:789a::1"}
	Ipv6Address := []string{"fd12:3456:789a::1", "fd12:3456:789a::2"}
	Ipv4Connectivity := model.Ipv4Connectivity{
		ConnectionPointType:     nil,
		Id:                      "1",
		InstanceAnnotations:     nil,
		Ipv4Address:             &Ipv4Address,
		NetworkMask:             &Ipv4NetMask,
		RelatedConnectionPoints: relatedConnectionPoints,
		RouterIpv4Address:       nil,
	}
	deviceInfo.ConnectionPoints = append(deviceInfo.ConnectionPoints, model.Connection{
		Ipv4Connectivity: Ipv4Connectivity,
	})
	Ipv6Connectivity := model.Ipv6Connectivity{
		ConnectionPointType:     nil,
		Id:                      "2",
		InstanceAnnotations:     nil,
		Ipv6Address:             Ipv6Address,
		RelatedConnectionPoints: nil,
		RouterIpv6Address:       routerIpv6Address,
	}
	deviceInfo.ConnectionPoints = append(deviceInfo.ConnectionPoints, model.Connection{
		Ipv6Connectivity: Ipv6Connectivity,
	})
	EthernetPort := model.EthernetPort{
		Id:                  "3",
		ConnectionPointType: nil,
		MacAddress:          &randomMacAddress,
	}
	deviceInfo.ConnectionPoints = append(deviceInfo.ConnectionPoints, model.Connection{
		EthernetPort: EthernetPort,
	})

	state := model.ManagementStateValuesUnknown
	State := model.ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}
	deviceInfo.ManagementState = State

	reachabilityStateValue := model.ReachabilityStateValuesReached
	reachabilityState := model.ReachabilityState{
		StateTimestamp: &timestamp,
		StateValue:     &reachabilityStateValue,
	}
	deviceInfo.ReachabilityState = &reachabilityState

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
