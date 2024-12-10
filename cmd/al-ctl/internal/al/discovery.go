/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package al

import (
	"fmt"
	"io"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

func Discover(endpoint string, discoveryFile string) []*generated.DiscoverResponse {
	log.Info().Str("Endpoint", endpoint).Str("Discovery Request Config File", discoveryFile).Msg("Starting discovery job")

	discoveryConfig := config.NewDiscoveryConfigWithDefaults()

	if discoveryFile != "" {
		var configError error
		discoveryConfig, configError = config.NewDiscoveryConfigFromFile(discoveryFile)
		if configError != nil {
			log.Err(configError).Msg("Failed to read config file")
			return nil
		}
	}

	discoveryRequest := discoveryConfig.GetDiscoveryRequest()

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewDeviceDiscoverApiClient(conn)
	ctx := context.Background()
	stream, err := client.DiscoverDevices(ctx, discoveryRequest)

	if err != nil {
		log.Err(err).Msg("StartDeviceDiscovery request returned an error")
		return nil
	}

	devices := make([]*generated.DiscoverResponse, 0)
	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Err(err).Msg("SubscribeDiscovery request returned an error")
			return nil
		}

		fmt.Printf("%+v\n", resp.Devices)

		log.Trace().Interface("Devices", resp).Msg("")
		devices = append(devices, resp)
	}
	return devices
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
	log.Trace().Interface("Response", resp).Msg("Received Response")
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
	log.Trace().Interface("Response", resp).Msg("Received Response")
	return resp
}
