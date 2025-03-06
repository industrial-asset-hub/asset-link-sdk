/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package testsuite

import (
	"github.com/rs/zerolog/log"
)

func RunAssetValidation(params AssetValidationParams) {
	err := ValidateAsset(params)
	if err != nil {
		log.Fatal().Msg("error during asset validation")
	}
	log.Info().Msg("Asset validation successful")
}

func RunAPIValidation(params ApiValidationParams) {
	if params.ServiceName != "discovery" {
		log.Fatal().Msgf("api validation is not supported for %s", params.ServiceName)
	}
	var totalDiscoveryTests int
	if params.CancellationValidationRequired {
		totalDiscoveryTests = 1
	} else {
		totalDiscoveryTests = 3
	}

	errors := ValidateAPI(params)

	if len(errors) > 0 {
		for _, err := range errors {
			log.Error().Err(err).Msg("API validation failed")
		}
		log.Fatal().Msgf("Number of tests passed %d/%d", totalDiscoveryTests-len(errors), totalDiscoveryTests)
	}
	log.Info().Msgf("Number of tests passed %d/%d", totalDiscoveryTests-len(errors), totalDiscoveryTests)
}

func RunRegistrationValidation(params RegistryValidationParams) {
	err := ValidateAssetLinkRegistration(params)
	if err != nil {
		log.Fatal().Err(err).Msg("error during validation of registration")
	}
	log.Info().Msg("Registration validation successful")
}
