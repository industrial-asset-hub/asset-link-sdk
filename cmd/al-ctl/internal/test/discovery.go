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

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDiscoverDevices(testConfig TestConfig) bool {
	fmt.Println("Running Test for StartDiscovery")
	data, err := al.Discover(shared.AssetLinkEndpoint, testConfig.DiscoveryFile)
	if err != nil {
		log.Error().Msg("discovery test failed")
		return false
	}

	log.Trace().Str("File", shared.OutputFile).Msg("Saving to file")
	f, _ := os.Create(shared.OutputFile)
	defer f.Close()

	asJson, _ := json.MarshalIndent(data, "", "  ")
	_, err = f.Write(asJson)
	if err != nil {
		log.Err(err).Msg("error during writing of the json file")
	}

	if testConfig.AssetValidationRequired {
		numberOfAssetsToValidate, errOccurredDuringValidation := createAssetFilesFromDiscoveryResponse(data)
		for i := range numberOfAssetsToValidate {
			assetFileName := fmt.Sprintf("Test-%d.json", i)
			if fileExists(assetFileName) {
				testConfig.AssetValidationParams.AssetJsonPath = assetFileName
				err := ValidateAsset(testConfig.AssetValidationParams, testConfig.LinkMLSupported)
				if err != nil {
					errOccurredDuringValidation = true
					log.Err(err).Str("asset-file-name", assetFileName).Msg("failed to validate asset against schema")
				}
			}
		}

		if errOccurredDuringValidation {
			return false
		}
	}

	return true
}

func TestCancelDiscovery(testConfig TestConfig) bool {
	fmt.Println("Running Test for CancelDiscovery")

	if shared.TimeoutSeconds == 0 {
		log.Fatal().Msg("CancelDiscovery can only be used with a specified timeout")
	}

	_, err := al.Discover(shared.AssetLinkEndpoint, testConfig.DiscoveryFile)
	if err == nil { //TODO: IMPORTANT: has this ever been tested ?!?
		log.Error().Err(err).Msg("Failed to cancel discovery job")
		return false
	}
	st, ok := status.FromError(err)
	if ok && st.Code() == codes.Canceled { //TODO: IMPORTANT: has this ever been tested ?!?
		log.Info().Msg("Discovery job was successfully cancelled")
		return true
	} else {
		log.Error().Err(err).Msg("Failed to cancel discovery job")
		return false
	}
}

func TestGetFilterTypes(testConfig TestConfig) bool {
	fmt.Println("Running Test for GetFilterTypes")
	data := al.GetFilterTypes(shared.AssetLinkEndpoint)
	if data == nil {
		log.Error().Msg("get-filter-types test failed")
		return false
	}
	return true
}

func TestGetFilterOptions(testConfig TestConfig) bool {
	fmt.Println("Running Test for GetFilterOptions")
	data := al.GetFilterOptions(shared.AssetLinkEndpoint)
	if data == nil {
		log.Error().Msg("get-filter-options test failed")
		return false
	}
	return true
}
