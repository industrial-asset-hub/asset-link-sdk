/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package dcdconnection

import (
	"encoding/json"
	"io"

	generated "code.siemens.com/common-device-management/shared/cdm-dcd-sdk/v2/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

func StartDiscovery(endpoint string, option string, filter string) []generated.DiscoverResponse {
	log.Trace().Str("Endpoint", endpoint).Str("Option", option).Str("Filter", filter).Msg("Starting discovery job")
	// TODO: Generate option
	parsedOptions := []*generated.ActiveOption{}
	err := json.Unmarshal([]byte(option), &parsedOptions)
	if err != nil {
		log.Err(err).Msg("Parsing of the discovery option returned an error")
		return nil
	}
	log.Trace().Interface("Options", parsedOptions).Msg("Parsed discovery options")

	// TODO: Generate Filter
	parsedFilters := []*generated.ActiveFilter{}
	err = json.Unmarshal([]byte(filter), &parsedFilters)
	if err != nil {
		log.Err(err).Msg("Parsing of the discovery filter returned an error")
		return nil
	}
	log.Trace().Interface("Filters", parsedFilters).Msg("Parsed discovery filter")

	conn := grpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewDeviceDiscoverApiClient(conn)
	ctx := context.Background()
	stream, err := client.DiscoverDevices(ctx, &generated.DiscoverRequest{
		Filters: parsedFilters,
		Options: parsedOptions,
		Target:  nil,
	})

	if err != nil {
		log.Err(err).Msg("StartDeviceDiscovery request returned an error")
		return nil
	}

	devices := make([]generated.DiscoverResponse, 0)
	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Err(err).Msg("SubscribeDiscovery request returned an error")
			return nil
		}
		log.Info().Interface("response", resp.Devices).Msg("Received Response")

		log.Trace().Interface("Devices", resp).Msg("")
		devices = append(devices, *resp)
	}
	return devices
}
