/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package testsuite

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/al-ctl/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/al-ctl/shared"
	iah_discovery "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model/conversion"
	"github.com/rs/zerolog/log"
	"os"
)

type testFunction func(string, string) interface{}

type Test struct {
	name     string
	function testFunction
}

func RunApiMockTests(address, discoveryFile string) {
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
		if test.name == "TestDiscoverDevices" && shared.AssetValidationRequired {
			createAssetFileFromDiscoveryResponse(data)
		}
		fmt.Printf("Total tests passed: %d/%d, failed: %d\n", testPassed, len(allTests), len(allTests)-testPassed)
	}

}

func createAssetFileFromDiscoveryResponse(data interface{}) {
	discoveryResponse := data.([]iah_discovery.DiscoverResponse)
	baseSchemaVersion, err := GetBaseSchemaVersionFromExtendedSchema()
	if err != nil {
		log.Err(err).Msg("failed to get base schema baseSchemaVersion from extended schema defaulting to v0.9.0")
		baseSchemaVersion = "v0.9.0"
	}
	for i := range discoveryResponse {
		for _, discoveredDevice := range discoveryResponse[i].Devices {
			transformedDevice := conversion.TransformDevice(discoveredDevice, "URI", baseSchemaVersion)
			// Add a unique id to the transformed device
			transformedDevice["id"] = uuid.New().String()
			jsonDevice, err := json.Marshal(transformedDevice)
			if err != nil {
				log.Err(err).Msg("failed to marshal transformed device")
			}
			// Create a Test asset file
			assetFile, err := os.Create("Test.json")
			if err != nil {
				log.Err(err).Msg("failed to create Test asset file")
			}
			_, err = assetFile.Write(jsonDevice)
			if err != nil {
				log.Err(err).Msg("failed to write Test asset file")
			}
			// Set the asset path to the Test asset file
			shared.AssetJsonPath = "Test.json"
		}
	}
}

func TestDiscoverDevices(address, discoveryFile string) interface{} {
	fmt.Println("Running Test for StartDiscovery")
	data := al.Discover(address, discoveryFile)
	return data
}

func TestGetFilterTypes(address, discoveryFile string) interface{} {
	fmt.Println("Running Test for GetFilterTypes")
	data := al.GetFilterTypes(address)
	return data
}

func TestGetFilterOptions(address, discoveryFile string) interface{} {
	fmt.Println("Running Test for GetFilterOptions")
	data := al.GetFilterOptions(address)
	return data
}
