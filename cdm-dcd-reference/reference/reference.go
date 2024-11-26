/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package reference

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/industrial-asset-hub/asset-link-sdk/v2/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements the Discovery interface and feature

type ReferenceClassDriver struct {
	discoveryLock sync.Mutex
}

var lastSerialNumber = atomic.Int64{}

func (m *ReferenceClassDriver) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error {
	log.Info().
		Msg("Start Discovery")

	// Check if a job is already running
	// We currently support only one running job
	if m.discoveryLock.TryLock() {
		defer m.discoveryLock.Unlock()
	} else {
		errMsg := "Discovery job is already running"
		log.Error().Msg(errMsg)
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}

	alInterface, err := discoveryConfig.GetOptionSettingString("interface", "")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	ipRange, err := discoveryConfig.GetOptionSettingString("ip_range", "")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	if alInterface == "" {
		log.Info().Msg("Scanning for devices on all interfaces")
	} else {
		log.Info().Msg("Scanning for devices on interface " + alInterface)
	}

	if ipRange == "" {
		log.Info().Msg("No Filtering of Devices for IP range")
	} else {
		log.Info().Msg("Filtering of Devices for IP range " + ipRange)
	}

	if alInterface == "" || alInterface == "eth0" {
		//time.Sleep(20 * time.Second)

		// "scan" for devices connected to eth0 ...
		deviceNIC := "enp0"
		deviceIPs := []string{"192.168.0.123", "10.0.0.1"}
		if ContainsIpInRange(ipRange, deviceIPs) {
			// Just provide a static asset
			name := "Device"
			lastSerialNumber.Add(1)
			manufacturer := "Siemens AG"
			serialNumber := fmt.Sprint(lastSerialNumber.Load())
			product := "cdm-reference-dcd-test2"
			deviceInfo := model.NewDevice("EthernetDevice", name)

			uriOfTheProduct := fmt.Sprintf("https://%s/%s-%s", strings.ReplaceAll(manufacturer, " ", "_"), strings.ReplaceAll(product, " ", "_"), serialNumber)
			deviceInfo.AddNameplate(manufacturer, uriOfTheProduct, "MyOrderNumber", product, "1.0.0", serialNumber)

			deviceInfo.AddSoftware("firmware", "1.2.5")
			deviceInfo.AddCapabilities("firmware_update", false)

			randomMacAddress := generateRandomMacAddress()
			id := deviceInfo.AddNic(deviceNIC, randomMacAddress)
			deviceInfo.AddIPv4(id, deviceIPs[0], "255.255.255.0", "")
			deviceInfo.AddIPv4(id, deviceIPs[1], "255.255.255.0", "")
			discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()

			err := devicePublisher.PublishDevice(discoveredDevice)
			if err != nil {
				// discovery request was likely cancelled -> terminate discovery and return error
				log.Error().Msgf("Publishing Error: %v", err)
				return err
			}
		}
	}

	log.Debug().
		Msg("Discover function exiting")
	return nil
}

func (m *ReferenceClassDriver) FilterTypes() []*generated.SupportedFilter {
	filterTypes := make([]*generated.SupportedFilter, 0)
	filterTypes = append(filterTypes, &generated.SupportedFilter{
		Key:      "ip_range",
		Datatype: generated.VariantType_VT_STRING,
	})
	return filterTypes
}

func (m *ReferenceClassDriver) FilterOptions() []*generated.SupportedOption {
	filterOptions := make([]*generated.SupportedOption, 0)
	filterOptions = append(filterOptions, &generated.SupportedOption{
		Key:      "interface",
		Datatype: generated.VariantType_VT_STRING,
	})
	return filterOptions
}

func generateRandomMacAddress() string {
	r := rand.Uint64()
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		0x00, 0x16, 0x3e,
		byte(r>>8),
		byte(r>>16),
		byte(r>>24))
}
