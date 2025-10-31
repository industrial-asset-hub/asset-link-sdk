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

	"encoding/json"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/dataio"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/fileformat"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func GetIdentifiers(endpoint string, credentialFilePath string) (*generated.GetIdentifiersResponse, error) {
	log.Info().Msg("Running Test for GetIdentifiers")
	identifiersReq, err := createIdentifiersRequestFromInputFile(credentialFilePath)
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

func GetIdentifiersRequest(identifiers *generated.GetIdentifiersRequest, endpoint string) *generated.GetIdentifiersResponse {
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
	var identifierResponsesInFile fileformat.DiscoveryResponseInFile
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

	// marshals the identifierResponsesInFile to json
	message, err := protojson.Marshal(discResult)
	if err != nil {
		log.Err(err).Msg("Marshalling of identifier responses failed")
		return err
	}

	identifierResponsesInFile = fileformat.DiscoveryResponseInFile{DiscoveryResponse: message}

	// marshals the array of identifierResponsesInFile to json
	identifierResponsesJson, err := json.MarshalIndent(identifierResponsesInFile, "", "	")
	if err != nil {
		log.Err(err).Msg("Marshalling to array of identifier responses failed")
		return err
	}

	if err := dataio.WriteOutput(identifierOutputFile, identifierResponsesJson); err != nil {
		log.Err(err).Str("file-path", identifierOutputFile).Msg("Error writing output")
		return err
	}
	return nil
}
