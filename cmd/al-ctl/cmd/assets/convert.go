/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assets

import (
	"bytes"
	"encoding/json"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/dataio"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/fileformat"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var convertInputFile string = ""
var convertOutputFile string = ""
var convertOutputFormat string = ""

var ConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert payload to actual assets",
	Long:  `This command converts a discover payload to actual assets.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Trace().Str("Input", convertInputFile).Str("Output", convertOutputFile).Msg("Starting conversion")

		buffer, err := dataio.ReadInput(convertInputFile)
		if err != nil {
			log.Fatal().Err(err).Str("file-path", convertInputFile).Msg("Error reading input")
		}

		var discoveryResponsesInFile fileformat.DiscoveryResponsesInFile
		if err := json.Unmarshal(buffer, &discoveryResponsesInFile); err != nil {
			log.Fatal().Err(err).Msg("Unmarshalling to array of discovery responses failed")
		}

		var convertedDevices []map[string]interface{}
		for _, discoveryResponse := range discoveryResponsesInFile {
			protoMsg := generated.DiscoverResponse{}
			if err := protojson.Unmarshal(discoveryResponse.DiscoveryResponse, &protoMsg); err != nil {
				log.Fatal().Err(err).Msg("Unmarshalling of discovery responses failed")
			}

			for _, device := range protoMsg.Devices {
				device := model.ConvertFromDiscoveredDevice(device, "URI")
				convertedDevices = append(convertedDevices, device)
			}
		}

		var outputDevices []byte
		switch convertOutputFormat {
		case "json":
			outputDevices, err = json.MarshalIndent(convertedDevices, "", "  ")
			if err != nil {
				log.Fatal().Err(err).Msg("Marshalling of transformed devices failed")
			}
		case "text":
			for _, device := range convertedDevices {
				assetOperations, assetOperationsOK := device["asset_operations"].([]map[string]interface{})
				if assetOperationsOK {
					// bools are not supported by CS (activation_flag)
					for _, assetOperation := range assetOperations {
						activationFlag, activationFlagOK := assetOperation["activation_flag"].(string)
						if activationFlagOK {
							if activationFlag == "true" {
								assetOperation["activation_flag"] = true
							} else {
								assetOperation["activation_flag"] = false
							}
						}
					}
				}

				delete(device, "reachability_state")
				delete(device, "management_state")
			}

			var devices []model.Asset
			err := mapstructure.Decode(convertedDevices, &devices)
			if err != nil {
				log.Fatal().Err(err).Msg("Decoding of transformed devices failed")
			}

			var buffer bytes.Buffer
			numDevices := len(devices)
			for n, device := range devices {
				buffer.WriteString(device.String())
				if n < numDevices-1 {
					buffer.WriteString("----------------------------------------\n")
				}
			}
			outputDevices = buffer.Bytes()
		default:
			log.Fatal().Str("format", convertOutputFormat).Msg("Invalid output format specified")
		}

		if err := dataio.WriteOutput(convertOutputFile, outputDevices); err != nil {
			log.Fatal().Err(err).Str("file-path", convertOutputFile).Msg("Error writing output")
		}
	},
}

func init() {
	ConvertCmd.Flags().StringVarP(&convertInputFile, "input", "i", "", "input filename (default stdin)")
	ConvertCmd.Flags().StringVarP(&convertOutputFile, "output", "o", "", "output filename (default stdout)")
	ConvertCmd.Flags().StringVarP(&convertOutputFormat, "format", "f", "json", "output format: json (default) or text")
}
