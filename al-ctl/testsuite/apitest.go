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
	"strings"
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
	discoveryResponse := data.([]*iah_discovery.DiscoverResponse)
	baseSchemaVersion, err := GetBaseSchemaVersionFromExtendedSchema()
	if err != nil {
		log.Err(err).Msg("failed to get base schema baseSchemaVersion from extended schema defaulting to v0.9.0")
		baseSchemaVersion = "v0.9.0"
	}
	for deviceIndex, device := range discoveryResponse[0].Devices {
		transformedDevice := conversion.TransformDevice(device, "URI", baseSchemaVersion)
		// Add a unique id to the transformed device
		transformedDevice["id"] = uuid.New().String()
		file, err := os.Create(fmt.Sprintf("testdevice-%d.json", deviceIndex))
		if err != nil {
			fmt.Printf("Error creating testsuite.json file: %v", err)
		}
		defer file.Close() // Ensure the file is closed when done

		arr := strings.Split(shared.AssetJsonPath, "/")
		newArr := arr[:len(arr)-1]
		shared.AssetJsonPath = strings.Join(newArr, "/")
		shared.AssetJsonPath = shared.AssetJsonPath + fmt.Sprintf("testdevice-%d.json", deviceIndex)
		// Write the transformed asset to a file
		jsonWriter := json.NewEncoder(file)
		if err := jsonWriter.Encode(transformedDevice); err != nil {
			log.Err(err).Msg("failed to write transformed asset to file")
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
