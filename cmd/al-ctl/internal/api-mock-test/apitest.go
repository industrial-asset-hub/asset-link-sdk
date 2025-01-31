/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package apimock

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	iah_discovery "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/rs/zerolog/log"
	"os"
)

type testFunction func(string, string) interface{}

type Test struct {
	name     string
	function testFunction
}

func RunApiMockTests(address string, discoveryFile string, assetJsonPath string, assetValidationRequired bool) (numberOfAssetsToValidate int) {
	allTests := []Test{
		{"TestDiscoverDevices", TestDiscoverDevices},
		{"TestGetFilterTypes", TestGetFilterTypes},
		{"TestGetFilterOptions", TestGetFilterOptions},
	}
	testPassed := 0
	for _, test := range allTests {
		data := test.function(address, discoveryFile)
		if data == nil {
			fmt.Println("Test failed")
			continue
		}
		fmt.Println("Test passed")
		testPassed++
		if test.name == "TestDiscoverDevices" && assetValidationRequired {
			numberOfAssetsToValidate = createAssetFileFromDiscoveryResponse(data)
		}
		fmt.Printf("Total tests passed: %d/%d, failed: %d\n", testPassed, len(allTests), len(allTests)-testPassed)
	}

	return numberOfAssetsToValidate
}

func createAssetFileFromDiscoveryResponse(data interface{}) (numberOfAssetsToValidate int) {
	discoveryResponse := data.([]*iah_discovery.DiscoverResponse)
	// to store all created files
	var assetFiles []*os.File
	for discoveryResponseIndex := range discoveryResponse {
		for discoveredDeviceIndex, discoveredDevice := range discoveryResponse[discoveryResponseIndex].Devices {

			// Convert the discovered device to a transformed device
			transformedDevice := model.ConvertFromDiscoveredDevice(discoveredDevice, "URI")

			// Add a unique id to the transformed device
			transformedDevice["id"] = uuid.New().String()

			// Marshal the transformed device
			jsonDevice, err := json.Marshal(transformedDevice)
			if err != nil {
				log.Err(err).Msg("failed to marshal transformed device")
				return numberOfAssetsToValidate
			}

			// Create a Test asset file
			assetFileName := fmt.Sprintf("Test-%d.json", discoveredDeviceIndex)
			assetFile, err := os.Create(assetFileName)
			if err != nil {
				log.Err(err).Msg("failed to create Test asset file")
				return numberOfAssetsToValidate
			}
			numberOfAssetsToValidate++

			// Append the asset file to the list of asset files
			assetFiles = append(assetFiles, assetFile)
			_, err = assetFile.Write(jsonDevice)
			if err != nil {
				log.Err(err).Msg("failed to write Test asset file")
				return numberOfAssetsToValidate
			}
		}
	}
	for _, assetFile := range assetFiles {
		err := assetFile.Close()
		if err != nil {
			log.Err(err).Msg("failed to close asset file")
		}
	}
	return numberOfAssetsToValidate
}
