/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assets

import (
	"encoding/json"
	"fmt"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/dataio"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/shared"
	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/model"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var convertInputFile string = ""
var convertOutputFile string = ""

var ConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert payload to actual assets",
	Long:  `This command converts a discover payload or property response payload to actual assets.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Trace().Str("Input", convertInputFile).Str("Output", convertOutputFile).Msg("Starting conversion")

		buffer, err := dataio.ReadInput(convertInputFile)
		if err != nil {
			log.Fatal().Err(err).Str("file-path", convertInputFile).Msg("Error reading input")
		}

		convertedDevices, err := convertPayloadToAssets(buffer)
		if err != nil {
			log.Fatal().Err(err).Msg("Converting input payload to assets failed")
		}

		transformedDevicesAsJson, err := json.MarshalIndent(convertedDevices, "", "  ")
		if err != nil {
			log.Fatal().Err(err).Msg("Marshalling of transformed devices failed")
		}

		if err := dataio.WriteOutput(convertOutputFile, transformedDevicesAsJson); err != nil {
			log.Fatal().Err(err).Str("file-path", convertOutputFile).Msg("Error writing output")
		}
	},
}

func convertPayloadToAssets(buffer []byte) ([]map[string]interface{}, error) {
	var propertyResponse generatedDeviceInfo.GetPropertyValuesResponse
	if err := protojson.Unmarshal(buffer, &propertyResponse); err == nil && len(propertyResponse.GetPropertyResults()) > 0 {
		return convertPropertyResponsePayload(&propertyResponse)
	}

	discoveryResponses, err := parseDiscoveryResponses(buffer)
	if err == nil {
		if converted, convErr := convertDiscoveryPayload(discoveryResponses); convErr == nil {
			return converted, nil
		}
	}

	// If input is already an array of assets, keep it as-is.
	var assets []map[string]interface{}
	if err := json.Unmarshal(buffer, &assets); err == nil && len(assets) > 0 {
		return assets, nil
	}

	return nil, fmt.Errorf("input format is not supported: expected discovery response array or GetPropertyValuesResponse")
}

func parseDiscoveryResponses(buffer []byte) ([]*generated.DiscoverResponse, error) {
	var rawResponses []json.RawMessage
	if err := json.Unmarshal(buffer, &rawResponses); err != nil {
		return nil, err
	}

	if len(rawResponses) == 0 {
		return nil, fmt.Errorf("discovery response array is empty")
	}

	responses := make([]*generated.DiscoverResponse, 0, len(rawResponses))
	for i, raw := range rawResponses {
		protoMsg := generated.DiscoverResponse{}
		if err := protojson.Unmarshal(raw, &protoMsg); err != nil {
			return nil, fmt.Errorf("unmarshal discovery response at index %d: %w", i, err)
		}
		responses = append(responses, &protoMsg)
	}

	return responses, nil
}

func convertDiscoveryPayload(discoveryResponses []*generated.DiscoverResponse) ([]map[string]interface{}, error) {
	var convertedDevices []map[string]interface{}
	for _, response := range discoveryResponses {
		for _, device := range response.Devices {
			device := model.ConvertFromDiscoveredDevice(device, "URI")
			convertedDevices = append(convertedDevices, device)
		}

		for _, respError := range response.Errors {
			log.Warn().Int32("Result Code", respError.ResultCode).
				Str("Description", respError.Description).
				Interface("Source", respError.Source).
				Msg("Dropped discovery error")
		}
	}

	return convertedDevices, nil
}

func convertPropertyResponsePayload(propertyResponse *generatedDeviceInfo.GetPropertyValuesResponse) ([]map[string]interface{}, error) {
	if propertyResponse == nil {
		return nil, fmt.Errorf("property response is nil")
	}

	asset, propertyCount, errorEntries, err := shared.BuildAssetFromPropertyResults(propertyResponse.GetPropertyResults(), "property_results")
	if err != nil {
		return nil, err
	}

	for _, propertyErr := range errorEntries {
		log.Warn().Str("Key", propertyErr.GetKey()).
			Str("Description", propertyErr.GetDescription()).
			Msg("Dropped property error")
	}

	if propertyCount == 0 {
		return nil, fmt.Errorf("no successful properties found in property response")
	}

	return []map[string]interface{}{asset}, nil
}
func init() {
	ConvertCmd.Flags().StringVarP(&convertInputFile, "input", "i", "", "input filename (default stdin)")
	ConvertCmd.Flags().StringVarP(&convertOutputFile, "output", "o", "", "output filename (default stdout)")
}
