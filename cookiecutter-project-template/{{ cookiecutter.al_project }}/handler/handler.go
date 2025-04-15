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
	ga "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/artefact-update"
	gd "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements the Discovery interface and feature

type AssetLinkImplementation struct {
	driverLock sync.Mutex
}

func (m *AssetLinkImplementation) Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error {
	log.Info().Msg("Handle Discovery Request")

	// Check if a job is already running
	// We currently support only one running job
	if m.driverLock.TryLock() {
		defer m.driverLock.Unlock()
	} else {
		const errMsg string = "Another job is already running"
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

	// If there are any device-specific errors, they can be communicated via devicePublisher.PublishError(discoverError)
	// discoverError can be created like this:
	// discoverError := &generated.DiscoverError{
	//			ResultCode:  int32(codes.Unavailable),
	//			Description: "Error retrieving device details",
	//		}

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

func (m *AssetLinkImplementation) GetSupportedOptions() []*gd.SupportedOption {
	supportedOptions := make([]*gd.SupportedOption, 0)
	supportedOptions = append(supportedOptions, &gd.SupportedOption{
		Key:      "option",
		Datatype: gd.VariantType_VT_BOOL,
	})
	return supportedOptions
}

func (m *AssetLinkImplementation) GetSupportedFilters() []*gd.SupportedFilter {
	supportedFilters := make([]*gd.SupportedFilter, 0)
	supportedFilters = append(supportedFilters, &gd.SupportedFilter{
		Key:      "filter",
		Datatype: gd.VariantType_VT_STRING,
	})
	return supportedFilters
}

// if the GetIdentifiers() interface is implemented by your asset link,
// uncomment the Identifiers line in main.go to register the interface/feature and
// add "siemens.common.identifiers.v1" to the app_types in the registry.json file
func (m *AssetLinkImplementation) GetIdentifiers(identifiersRequest config.IdentifiersRequest) ([]*generated.DeviceIdentifier, error) {
	log.Info().Msg("Handle Get Identifiers Request")
	// Add your custom logic here to retrieve identifiers based on the provided parameters and credentials
	identifiers := []*generated.DeviceIdentifier{}
	return identifiers, nil
}

func (m *AssetLinkImplementation) HandlePushArtefact(artefactReceiver *artefact.ArtefactReceiver) error {
	log.Info().Msg("Handle Push Artefact by receiving the artefact")

	// Check if a job is already running
	// We currently support only one running job
	if m.driverLock.TryLock() {
		defer m.driverLock.Unlock()
	} else {
		const errMsg string = "Another job is already running"
		log.Error().Msg(errMsg)
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}

	artefactMetaData, err := artefactReceiver.ReceiveArtefactMetaData()
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact meta data")
		return err
	}

	deviceIdentifier := artefactMetaData.GetDeviceIdentifier()
	artefactType := artefactMetaData.GetArtefactType()

	log.Info().Str("DeviceIdentifier", string(deviceIdentifier)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	err = artefactReceiver.ReceiveArtefactToFile("artefact_file")
	if err != nil {
		log.Err(err).Msg("Failed to receive artefact file")
		return err
	}

	_ = artefactReceiver.UpdateStatus(ga.ArtefactUpdateState_AUS_DOWNLOAD, ga.TransferStatus_AS_OK, "Status Message", 100)

	return nil
}

func (m *AssetLinkImplementation) HandlePullArtefact(artefactMetaData *artefact.ArtefactMetaData, artefactTransmitter *artefact.ArtefactTransmitter) error {
	log.Info().Msg("Handle Pull Artefact by transmitting the arefact")

	// Check if a job is already running
	// We currently support only one running job
	if m.driverLock.TryLock() {
		defer m.driverLock.Unlock()
	} else {
		const errMsg string = "Another job is already running"
		log.Error().Msg(errMsg)
		return status.Errorf(codes.ResourceExhausted, errMsg)
	}

	deviceIdentifier := artefactMetaData.GetDeviceIdentifier()
	artefactType := artefactMetaData.GetArtefactType()

	log.Info().Str("DeviceIdentifier", string(deviceIdentifier)).Str("ArtefactType", artefactType.String()).Msg("ArtefactMetaData")

	err := artefactTransmitter.TransmitArtefactFromFile("artefact_file", 1024)
	if err != nil {
		log.Err(err).Msg("Failed to transmit artefact file")
		return err
	}

	_ = artefactTransmitter.UpdateStatus(ga.TransferStatus_AS_OK, "Status Message")

	return nil
}
