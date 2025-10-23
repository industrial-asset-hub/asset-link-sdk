/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/artefact"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements the Discovery interface and feature

type AssetLinkImplementation struct {
	discoveryLock sync.Mutex
}

func (m *AssetLinkImplementation) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error {
	log.Info().Msg("Handle Discovery Request")

	// Check if a job is already running
	// We currently support only one running job
	if m.discoveryLock.TryLock() {
		defer m.discoveryLock.Unlock()
	} else {
		const errMsg string = "Another discovery job is already running"
		log.Error().Msg(errMsg)
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}

	//
	// Add your custom logic here to retrieve discovery config, discover devices, and publish them
	//

	optionSetting, optionErr := discoveryConfig.GetOptionSettingBool("option", false)
	if optionErr != nil {
		log.Error().Err(optionErr)
		return optionErr
	}

	filterSetting, filterErr := discoveryConfig.GetFilterSettingString("filter", "default")
	if filterErr != nil {
		log.Error().Err(filterErr)
		return filterErr
	}

	_ = optionSetting
	_ = filterSetting

	// Fillup the device information
	assetName := "Dummy Device 1"
	vendorName := "{{ cookiecutter.company }}"
	productName := "Dummy Product"
	orderNumber := "AN0123456789"
	serialNumber := "SN00012345678900001"
	hardwareVersion := "3"
	firmwareVersion := "1.0.0"

	productUri := fmt.Sprintf("urn:%s/%s/%s", strings.ReplaceAll(vendorName, " ", "_"), strings.ReplaceAll(productName, " ", "_"), serialNumber)

	deviceInfo := model.NewDevice("EthernetDevice", assetName)
	deviceInfo.AddNameplate(vendorName, productUri, orderNumber, productName, hardwareVersion, serialNumber)
	deviceInfo.AddSoftware("Firmware", firmwareVersion, true)
	deviceInfo.AddCapabilities("firmware_update", false)
	deviceInfo.AddDescription("Dummy Device")

	deviceInfo.AddMetadata("DEVICE-ID") // device ID or device connection data used for artefact uploads/downloads

	nicID := deviceInfo.AddNic("eth0", "00:16:3e:01:02:03") // random mac address
	deviceInfo.AddIPv4(nicID, "192.168.0.10", "255.255.255.0", "")
	deviceInfo.AddIPv4(nicID, "10.0.0.153", "255.255.255.0", "")

	// Convert and publish device
	discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()

	err := devicePublisher.PublishDevice(discoveredDevice)
	if err != nil {
		// discovery request was likely cancelled -> terminate discovery and return error
		log.Error().Msgf("Publishing Error: %v", err)
		return err
	}

	return nil
}

func (m *AssetLinkImplementation) GetSupportedOptions() []*generated.SupportedOption {
	supportedOptions := make([]*generated.SupportedOption, 0)
	supportedOptions = append(supportedOptions, &generated.SupportedOption{
		Key:      "option",
		Datatype: generated.VariantType_VT_BOOL,
	})
	return supportedOptions
}

func (m *AssetLinkImplementation) GetSupportedFilters() []*generated.SupportedFilter {
	supportedFilters := make([]*generated.SupportedFilter, 0)
	supportedFilters = append(supportedFilters, &generated.SupportedFilter{
		Key:      "filter",
		Datatype: generated.VariantType_VT_STRING,
	})
	return supportedFilters
}

func (m *AssetLinkImplementation) HandlePushArtefact(artefactReceiver *artefact.ArtefactReceiver) error {
	log.Info().Msg("Handle Push Artefact by receiving the artefact")

	artefactMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact meta data")
		return err
	}

	deviceIdentifier := string(artefactMetaData.GetDeviceIdentifier())
	artefactType := artefactMetaData.GetArtefactType().String()

	log.Info().Str("DeviceIdentifier", deviceIdentifier).Str("ArtefactType", artefactType).Msg("ArtefactMetaData")

	err = artefactReceiver.ReceiveArtefactToFile("artefact_file")
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact file")
		return err
	}

	return nil
}

func (m *AssetLinkImplementation) HandlePullArtefact(artefactMetaData *artefact.ArtefactMetaData, artefactTransmitter *artefact.ArtefactTransmitter) error {
	log.Info().Msg("Handle Pull Artefact by transmitting the arefact")

	deviceIdentifier := string(artefactMetaData.GetDeviceIdentifier())
	artefactType := artefactMetaData.GetArtefactType().String()

	log.Info().Str("DeviceIdentifier", deviceIdentifier).Str("ArtefactType", artefactType).Msg("ArtefactMetaData")

	err := artefactTransmitter.TransmitArtefactFromFile("artefact_file", 1024)
	if err != nil {
		log.Err(err).Msg("Failed to transmit artefact file")
		return err
	}

	return nil
}
