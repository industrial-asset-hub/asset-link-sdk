/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package test

import (
	"context"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
)

func TestGetIdentifiers(testConfig TestConfig) bool {
	log.Info().Msg("Running Test for GetIdentifiers")
	identifiersReq := createIdentifiersRequestFromCredential(testConfig.Credential)
	identifiers := GetIdentifiers(identifiersReq, shared.AssetLinkEndpoint)
	if identifiers == nil {
		log.Error().Msg("get-identifiers test failed")
		return false
	}
	log.Info().Msgf("Identifiers: %v\n", identifiers)

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

func GetIdentifiers(identifiers *generated.GetIdentifiersRequest, endpoint string) *generated.GetIdentifiersResponse {
	log.Trace().Str("Endpoint", endpoint).Msg("Getting Identifiers")
	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewIdentifiersApiClient(conn)
	resp, err := client.GetIdentifiers(context.Background(), identifiers)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get identifiers")
		return nil
	}

	return resp
}

func createIdentifiersRequestFromCredential(filePath string) *generated.GetIdentifiersRequest {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &generated.GetIdentifiersRequest{}
	}

	identifiers := &generated.GetIdentifiersRequest{
		Target: &generated.Destination{
			Target: &generated.Destination_ConnectionParameterSet{
				ConnectionParameterSet: &generated.ConnectionParameterSet{
					Credentials: []*generated.ConnectionCredential{{
						Credentials: string(data),
					}},
				},
			},
		},
	}

	return identifiers
}
