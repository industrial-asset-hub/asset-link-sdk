/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package test

import (
	"context"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
)

func TestGetIdentifiers(_ TestConfig) bool {
	log.Info().Msg("Running Test for GetIdentifiers")
	identifiers := GetIdentifiers(shared.AssetLinkEndpoint)
	if identifiers == nil {
		log.Error().Msg("get-identifiers test failed")
		return false
	}
	log.Info().Msgf("Identifiers: %v\n", identifiers)
	return true
}

func GetIdentifiers(endpoint string) *generated.GetIdentifiersResponse {
	log.Trace().Str("Endpoint", endpoint).Msg("Getting Identifiers")
	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewIdentifiersApiClient(conn)
	resp, err := client.GetIdentifiers(context.Background(), &generated.GetIdentifiersRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get identifiers")
		return nil
	}
	return resp
}
