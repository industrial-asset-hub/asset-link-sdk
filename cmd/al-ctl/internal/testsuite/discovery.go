/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package testsuite

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/discovery"
	iah_discovery "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

func TestDiscoverDevices(assetValidationRequired bool, assetValidationParams AssetValidationParams, assetLinkEndpoint string, discoveryFile string) error {
	discoveryResponses, err := discovery.Discover(assetLinkEndpoint, discoveryFile, 0)
	if len(discoveryResponses) == 0 || err != nil {
		return err
	}
	if assetValidationRequired {
		assetValidationPassed := true
		for _, discoveryResponse := range discoveryResponses {
			assetFileNames := createAssetFilesFromDiscoveryResponse(discoveryResponse)
			for _, assetFileName := range assetFileNames {
				if fileExists(assetFileName) {
					assetValidationParams.AssetJsonPath = assetFileName
					err := ValidateAsset(assetValidationParams)
					if err != nil {
						assetValidationPassed = false
						log.Error().Err(err).Msgf("error during asset validation for asset %s", assetFileName)
					}
				} else {
					assetValidationPassed = false
					log.Error().Msgf("asset file %s does not exist", assetFileName)
				}
			}
		}
		if !assetValidationPassed {
			return errors.New("asset validation failed")
		}
	}

	return nil
}

func TestGetFilterTypes(assetLinkEndpoint string) error {
	data := discovery.GetFilterTypes(assetLinkEndpoint)
	if data == nil {
		return errors.New("get-filter-types test failed")
	}
	return nil
}

func TestGetFilterOptions(assetLinkEndpoint string) error {
	fmt.Println("Running Test for GetFilterOptions")
	data := discovery.GetFilterOptions(assetLinkEndpoint)
	if data == nil {
		return errors.New("get-filter-options test failed")
	}
	return nil
}

func TestCancelDiscovery(timeoutInSeconds uint, assetlinkEndpoint string, discoveryFile string) error {
	if timeoutInSeconds == 0 {
		return errors.New("CancelDiscovery can only be used with a specified timeout")
	}

	_, err := discovery.Discover(assetlinkEndpoint, discoveryFile, 0)
	if err == nil {
		return errors.New("failed to cancel discovery job")
	}
	st, ok := status.FromError(err)
	if ok && st.Code() == codes.Canceled {
		return nil
	} else {
		return errors.New("failed to cancel discovery job with error: " + err.Error())
	}
}

func createAssetFilesFromDiscoveryResponse(discoveryResponse *iah_discovery.DiscoverResponse) (fileNames []string) {
	assetFileNames := make([]string, 0)
	for discoveredDeviceIndex, discoveredDevice := range discoveryResponse.Devices {
		// Convert the discovered device to a transformed device
		transformedDevice := model.ConvertFromDiscoveredDevice(discoveredDevice, "URI")

		// Add a unique id to the transformed device
		transformedDevice["id"] = uuid.New().String()

		// Marshal the transformed device
		jsonDevice, err := json.Marshal(transformedDevice)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal device")
		}

		assetFileName := fmt.Sprintf("Asset-%d.json", discoveredDeviceIndex)
		assetFile, err := os.Create(assetFileName)
		if err != nil {
			log.Error().Err(err).Msg("failed to create asset file")
		}
		_, err = assetFile.Write(jsonDevice)
		if err != nil {
			assetFile.Close()
			log.Error().Err(err).Msg("failed to write to asset file")
		}
		assetFile.Close()
		assetFileNames = append(assetFileNames, assetFileName)
	}
	return assetFileNames
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		return true
	}
	log.Warn().Msgf("File %s does not exist", filename)
	return false
}
