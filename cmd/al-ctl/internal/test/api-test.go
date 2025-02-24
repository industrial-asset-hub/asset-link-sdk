/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	iah_discovery "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/rs/zerolog/log"
)

type TestConfig struct {
	DiscoveryFile           string
	AssetValidationRequired bool
	LinkMLSupported         bool
	AssetValidationParams   AssetValidationParams
}

type testFunction func(testConfig TestConfig) bool

type Test struct {
	name     string
	function testFunction
}

func RunApiTests(serviceName string, cancelValidationRequired bool, testConfig TestConfig) {
	allTests := []Test{}
	if serviceName == "discovery" {
		allTests = []Test{
			{"TestDiscoverDevices", TestDiscoverDevices},
			{"TestGetFilterTypes", TestGetFilterTypes},
			{"TestGetFilterOptions", TestGetFilterOptions},
		}
		if cancelValidationRequired {
			allTests = []Test{
				{"TestCancelDiscovery", TestCancelDiscovery},
			}
		}
	} else {
		log.Fatal().Msgf("Service \"%s\"is not supported", serviceName)
	}

	testPassed := 0
	for _, test := range allTests {
		result := test.function(testConfig)
		if result {
			testPassed++
		} else {
			log.Error().Str("test-name", test.name).Msg("test failed")
		}
	}

	log.Info().Msgf("Total tests passed: %d/%d, failed: %d\n", testPassed, len(allTests), len(allTests)-testPassed)

	if testPassed < len(allTests) {
		os.Exit(1)
	}
}

func createAssetFilesFromDiscoveryResponse(data interface{}) (numberOfAssetsToValidate int, errOccurred bool) {
	discoveryResponse := data.([]*iah_discovery.DiscoverResponse)
	// to store all created files
	for discoveryResponseIndex := range discoveryResponse {
		for discoveredDeviceIndex, discoveredDevice := range discoveryResponse[discoveryResponseIndex].Devices {
			// increment for each discovered device
			numberOfAssetsToValidate++
			// Convert the discovered device to a transformed device
			transformedDevice := model.ConvertFromDiscoveredDevice(discoveredDevice, "URI")

			// Add a unique id to the transformed device
			transformedDevice["id"] = uuid.New().String()

			// Marshal the transformed device
			jsonDevice, err := json.Marshal(transformedDevice)
			if err != nil {
				errOccurred = true
				log.Err(err).Msg("failed to marshal transformed device")
				continue
			}

			// Create a Test asset file
			assetFileName := fmt.Sprintf("Test-%d.json", discoveredDeviceIndex)
			assetFile, err := os.Create(assetFileName)
			if err != nil {
				errOccurred = true
				log.Err(err).Msg("failed to create Test asset file")
				continue
			}
			_, err = assetFile.Write(jsonDevice)
			if err != nil {
				errOccurred = true
				log.Err(err).Msg("failed to write Test asset file")
				assetFile.Close()
				continue
			}
			assetFile.Close()
		}
	}
	return numberOfAssetsToValidate, errOccurred
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
