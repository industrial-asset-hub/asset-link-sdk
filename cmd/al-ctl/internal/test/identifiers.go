/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package test

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
)

func TestGetIdentifiers(testConfig TestConfig) bool {
	log.Info().Msg("Running Test for GetIdentifiers")
	identifiers, err := al.GetIdentifiers(shared.AssetLinkEndpoint, testConfig.Credential)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GetIdentifiersRequest from file")
		return false
	}
	if identifiers == nil {
		log.Error().Msg("get-identifiers test failed")
		return false
	}
	discDevice := &generated.DiscoveredDevice{
		Identifiers: identifiers.Identifiers,
	}
	if len(discDevice.Identifiers) == 0 {
		log.Error().Msg("No identifiers found")
		return false
	}
	discResult := &generated.DiscoverResponse{
		Devices: []*generated.DiscoveredDevice{discDevice},
	}

	return createAndValidateDiscoveredAsset(testConfig, []*generated.DiscoverResponse{discResult})
}
