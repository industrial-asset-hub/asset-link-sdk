/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/publish"
	"github.com/rs/zerolog/log"
)

// Implements the Discovery interface and feature

type AssetLinkImplementation struct {
	discoveryLock sync.Mutex
}

var lastSerialNumber = atomic.Int64{}

func (m *AssetLinkImplementation) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) {
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

	//
	// Add your custom logic here to discover devices and publish them
	//

	// Fillup the device information
	deviceInfo := model.NewDevice("EthernetDevice", "My First Ethernet Device")

	product := "{{ cookiecutter.al_name }}"
	orderNumber := "PRODUCT-ONE"
	productVersion := "1.0.0"
	vendorName := "{{ cookiecutter.company }}"
	lastSerialNumber.Add(1)
	serialNumber := fmt.Sprint(lastSerialNumber.Load())
	productUri := fmt.Sprintf("urn:%s/%s", strings.ReplaceAll(vendorName, " ", "_"), strings.ReplaceAll(product, " ", "_"), serialNumber)
	deviceInfo.AddNameplate(
		vendorName,
		productUri,
		orderNumber,
		product,
		productVersion,
		serialNumber)

	deviceInfo.AddCapabilities("firmware_update", false)

	randomMacAddress := generateRandomMacAddress()
	id := deviceInfo.AddNic("eth0", randomMacAddress)
	deviceInfo.AddIPv4(id, "192.168.0.1", "255.255.255.0", "")

	// Convert and stream device to upstream system
	discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()

	err := devicePublisher.PublishDevice(discoveredDevice)
	if err != nil {
		// discovery request was likely cancelled -> terminate discovery and return error
		log.Error().Msgf("Publishing Error: %v", err)
		return err
	}

	log.Debug().
		Msg("Discover function exiting")
	return nil
}

func (m *AssetLinkImplementation) FilterTypes() []*generated.SupportedFilter {
	filterTypes := make([]*generated.SupportedFilter, 0)
	filterTypes = append(filterTypes, &generated.SupportedFilter{
		Key:      "type",
		Datatype: generated.VariantType_VT_BYTES,
	})
	return filterTypes
}

func (m *AssetLinkImplementation) FilterOptions() []*generated.SupportedOption {
	filterOptions := make([]*generated.SupportedOption, 0)
	filterOptions = append(filterOptions, &generated.SupportedOption{
		Key:      "option",
		Datatype: generated.VariantType_VT_BOOL,
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
