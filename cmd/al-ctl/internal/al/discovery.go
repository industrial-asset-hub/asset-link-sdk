/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package al

import (
	"encoding/json"
	"io"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/dataio"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/fileformat"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

func Discover(endpoint string, discoveryFile string) ([]*generated.DiscoverResponse, error) {
	log.Info().Str("Endpoint", endpoint).Str("Discovery Request Config File", discoveryFile).Msg("Starting discovery job")

	discoveryConfig := config.NewDiscoveryConfigWithDefaults()

	if discoveryFile != "" {
		var configError error
		discoveryConfig, configError = config.NewDiscoveryConfigFromFile(discoveryFile)
		if configError != nil {
			log.Err(configError).Msg("Failed to read config file")
			return nil, configError
		}
	}

	discoveryRequest := discoveryConfig.GetDiscoveryRequest()

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewDeviceDiscoverApiClient(conn)
	ctx := context.Background()
	if shared.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		go func() {
			duration := time.Duration(shared.TimeoutSeconds * float64(time.Second))
			time.Sleep(duration)
			cancel()
		}()
	}
	stream, err := client.DiscoverDevices(ctx, discoveryRequest)

	if err != nil {
		log.Err(err).Msg("StartDeviceDiscovery request returned an error")
		return nil, err
	}

	discoveryResponses := make([]*generated.DiscoverResponse, 0)
	deviceCount := 0
	errorCount := 0
	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Err(err).Msg("SubscribeDiscovery request returned an error")
			return nil, err
		}

		log.Trace().Interface("response", resp).Msg("")

		discoveryResponses = append(discoveryResponses, resp)

		errorCount += len(resp.Errors)
		deviceCount += len(resp.Devices)

		log.Debug().Int("Assets", len(resp.Devices)).
			Int("Discovery errors", len(resp.Errors)).
			Msg("Received response")

		for _, respError := range resp.Errors {
			log.Warn().Int32("Result Code", respError.ResultCode).
				Str("Description", respError.Description).
				Interface("Source", respError.Source).
				Msg("Received discovery error")
		}
	}
	log.Info().Int("Discovery responses", len(discoveryResponses)).
		Int("Assets", deviceCount).
		Int("Discovery errors", errorCount).
		Msg("Received all responses")
	return discoveryResponses, nil
}

func GetFilterTypes(endpoint string) *generated.FilterTypesResponse {
	log.Trace().Str("Endpoint", endpoint).Msg("Getting filter types")
	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewDeviceDiscoverApiClient(conn)
	ctx := context.Background()
	resp, err := client.GetFilterTypes(ctx, &generated.FilterTypesRequest{})
	if err != nil {
		log.Err(err).Msg("GetFilterTypes request returned an error")
		return nil
	}
	log.Trace().Interface("DiscoveryResponse", resp).Msg("Received DiscoveryResponse")
	return resp
}

func GetFilterOptions(endpoint string) *generated.FilterOptionsResponse {
	log.Trace().Str("Endpoint", endpoint).Msg("Getting filter options")
	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewDeviceDiscoverApiClient(conn)
	ctx := context.Background()
	resp, err := client.GetFilterOptions(ctx, &generated.FilterOptionsRequest{})
	if err != nil {
		log.Err(err).Msg("GetFilterOptions request returned an error")
		return nil
	}
	log.Trace().Interface("DiscoveryResponse", resp).Msg("Received DiscoveryResponse")
	return resp
}

func WriteDiscoveryResponsesFile(discoverOutputFile string, discoverResponses []*generated.DiscoverResponse) error {
	var discoveryResponsesInFile fileformat.DiscoveryResponsesInFile
	for _, discoverResponse := range discoverResponses {

		// marshals the discovery discoveryResponsesInFile to json
		message, err := protojson.Marshal(discoverResponse)
		if err != nil {
			log.Err(err).Msg("Marshalling of discovery responses failed")
			return err
		}

		discoveryResponsesInFile = append(discoveryResponsesInFile, fileformat.DiscoveryResponseInFile{DiscoveryResponse: message})
	}

	// marshals the array of discovery discoveryResponsesInFile to json
	discoveryResponsesJson, err := json.MarshalIndent(discoveryResponsesInFile, "", "	")
	if err != nil {
		log.Err(err).Msg("Marshalling to array of discovery responses failed")
		return err
	}

	if err := dataio.WriteOutput(discoverOutputFile, discoveryResponsesJson); err != nil {
		log.Err(err).Str("file-path", discoverOutputFile).Msg("Error writing output")
		return err
	}
	return nil
}
