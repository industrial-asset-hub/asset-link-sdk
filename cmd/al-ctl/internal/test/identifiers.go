/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package test

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
)

func TestGetIdentifiers(testConfig TestConfig) bool {
	log.Info().Msg("Running Test for GetIdentifiers")
	identifiersReq, err := createIdentifiersRequestFromInputFile(testConfig.Credential)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GetIdentifiersRequest from file")
		return false
	}
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

func createIdentifiersRequestFromInputFile(filePath string) (*generated.GetIdentifiersRequest, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &generated.GetIdentifiersRequest{}, nil
	}

	fmt.Println(string(data))
	var getIdentifiersReq generated.GetIdentifiersRequest
	err = protojson.Unmarshal(data, &getIdentifiersReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal GetIdentifiersRequest from file")
	}
	fmt.Println(getIdentifiersReq)
	return &getIdentifiersReq, err
}
