/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package testsuite

import "github.com/rs/zerolog/log"

func ValidateAPI(params ApiValidationParams) []error {
	if params.CancellationValidationRequired {
		log.Info().Msg("Running Test for CancelDiscovery")
		return []error{TestCancelDiscovery(params.TimeoutInSeconds, params.AssetLinkEndpoint, params.DiscoveryFile)}
	}
	return validateDiscovery(params.AssetLinkEndpoint, params.AssetValidationRequired, params.DiscoveryFile, params.AssetValidationParams)
}

func validateDiscovery(assetLinkEndpoint string, assetValidationRequired bool, discoveryFile string, assetValidationParams AssetValidationParams) []error {
	var discoveryErrors []error
	log.Info().Msg("Running Test for DiscoverDevices")
	discoverDevicesValidationError := TestDiscoverDevices(assetValidationRequired, assetValidationParams, assetLinkEndpoint, discoveryFile)
	if discoverDevicesValidationError != nil {
		discoveryErrors = append(discoveryErrors, discoverDevicesValidationError)
	}
	log.Info().Msg("Running Test for GetFilterTypes")
	getFilterTypeValidationError := TestGetFilterTypes(assetLinkEndpoint)
	if getFilterTypeValidationError != nil {
		discoveryErrors = append(discoveryErrors, getFilterTypeValidationError)
	}
	log.Info().Msg("Running Test for GetFilterOptions")
	getFilterOptionsValidationError := TestGetFilterOptions(assetLinkEndpoint)
	if getFilterOptionsValidationError != nil {
		discoveryErrors = append(discoveryErrors, getFilterOptionsValidationError)
	}
	return discoveryErrors
}
