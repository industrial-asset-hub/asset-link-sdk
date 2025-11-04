/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package al

import (
	"context"
	"errors"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func GetIdentifiers(endpoint string, requestFilePath string) (*generated.GetIdentifiersResponse, error) {
	log.Info().Msg("Running Test for GetIdentifiers")
	identifiersReq, err := createIdentifiersRequestFromInputFile(requestFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GetIdentifiersRequest from file")
		return nil, err
	}
	identifiers := GetIdentifiersRequest(identifiersReq, shared.AssetLinkEndpoint)
	if identifiers == nil {
		log.Error().Msg("get-identifiers test failed")
		return nil, err
	}
	log.Info().Msgf("Identifiers: %v\n", identifiers)

	return identifiers, nil
}

func GetIdentifiersRequest(request *generated.GetIdentifiersRequest, endpoint string) *generated.GetIdentifiersResponse {
	log.Trace().Str("Endpoint", endpoint).Msg("Getting Identifiers")
	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewIdentifiersApiClient(conn)
	resp, err := client.GetIdentifiers(context.Background(), request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get identifiers")
		return nil
	}

	return resp
}

func createIdentifiersRequestFromInputFile(filePath string) (*generated.GetIdentifiersRequest, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &generated.GetIdentifiersRequest{}, err
	}

	var getIdentifiersReq generated.GetIdentifiersRequest
	err = protojson.Unmarshal(data, &getIdentifiersReq)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal GetIdentifiersRequest from file")
		return &generated.GetIdentifiersRequest{}, err
	}

	return &getIdentifiersReq, nil
}

func WriteIdentifierResponsesFile(identifierOutputFile string, identifierResponses *generated.GetIdentifiersResponse) error {
	if identifierResponses == nil {
		return errors.New("No identifier responses to write")
	}
	discDevice := &generated.DiscoveredDevice{
		Identifiers: identifierResponses.Identifiers,
	}
	if len(discDevice.Identifiers) == 0 {
		// return no identifier found error
		err := errors.New("No identifiers found")
		return err
	}
	discResult := &generated.DiscoverResponse{
		Devices: []*generated.DiscoveredDevice{discDevice},
	}

	err := WriteDiscoveryResponsesFile(identifierOutputFile, []*generated.DiscoverResponse{discResult})
	if err != nil {
		log.Err(err).Msg("Error writing discovery responses to file")
		return err
	}
	return nil
}
