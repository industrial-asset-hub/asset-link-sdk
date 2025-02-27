/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package reference

import (
	"sync"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/artefact"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements the Discovery interface and feature

type ReferenceAssetLink struct {
	discoveryLock sync.Mutex
}

func (m *ReferenceAssetLink) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error {
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

	alInterface, err := discoveryConfig.GetOptionSettingString("interface", "")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	ipRange, err := discoveryConfig.GetFilterSettingString("ip_range", "")
	if err != nil {
		log.Error().Err(err)
		return err
	}

	devicesFound := simdevices.ScanDevices(alInterface, ipRange)

	for _, device := range devicesFound {
		deviceInfo := model.NewDevice("EthernetDevice", device.GetDeviceName())
		deviceInfo.AddNameplate(device.GetManufacturer(), device.GetIDLink(), device.GetArticleNumber(),
			device.GetProductDesignation(), device.GetHardwareVersion(), device.GetSerialNumber())

		nicID := deviceInfo.AddNic(device.GetDeviceNIC(), device.GetMacAddress())
		deviceInfo.AddIPv4(nicID, device.GetIpDevice(), device.GetIpNetmask(), device.GetIpRoute())

		deviceInfo.AddSoftware("firmware", device.GetFirmwareVersion())
		deviceInfo.AddCapabilities("firmware_update", device.IsUpdateSupported())

		discoveredDevice := deviceInfo.ConvertToDiscoveredDevice()

		err := devicePublisher.PublishDevice(discoveredDevice)
		if err != nil {
			// discovery request was likely cancelled -> terminate discovery and return error
			log.Error().Msgf("Publishing Error: %v", err)
			return err
		}
	}

	return nil
}

func (m *ReferenceAssetLink) GetSupportedOptions() []*generated.SupportedOption {
	supportedOptions := make([]*generated.SupportedOption, 0)
	supportedOptions = append(supportedOptions, &generated.SupportedOption{
		Key:      "interface",
		Datatype: generated.VariantType_VT_STRING,
	})
	return supportedOptions
}

func (m *ReferenceAssetLink) GetSupportedFilters() []*generated.SupportedFilter {
	supportedFilters := make([]*generated.SupportedFilter, 0)
	supportedFilters = append(supportedFilters, &generated.SupportedFilter{
		Key:      "ip_range",
		Datatype: generated.VariantType_VT_STRING,
	})
	return supportedFilters
}

func (m *ReferenceAssetLink) HandlePushArtefact(artefactReceiver *artefact.ArtefactReceiver) error {
	log.Info().Msg("Handle Push Artefact by receiving the artefact")

	artefactMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact meta data")
		return err
	}

	deviceConnectionInformation := artefactMetaData.GetDeviceConnectionInformation()
	deviceIdentifier := string(deviceConnectionInformation[:])

	log.Info().Str("DeviceIdentifier", deviceIdentifier).Msg("ArtefactMetaData")

	err = artefactReceiver.ReceiveArtefactToFile("artefact_file")
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact file")
		return err
	}

	return nil
}

func (m *ReferenceAssetLink) HandlePullArtefact(artefactIdentifier *artefact.ArtefactIdentifier, artefactTransmitter *artefact.ArtefactTransmitter) error {
	log.Info().Msg("Handle Pull Artefact by transmitting the arefact")

	err := artefactTransmitter.TransmitArtefactFromFile("artefact_file", 1024)
	if err != nil {
		log.Err(err).Msg("Failed to transmit artefact file")
		return err
	}

	return nil
}
