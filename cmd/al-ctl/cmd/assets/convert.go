/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assets

import (
	"encoding/json"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/dataio"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/fileformat"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/model"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var convertInputFile string = ""
var convertOutputFile string = ""

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

		transformedDevicesAsJson, err := json.MarshalIndent(convertedDevices, "", "  ")
		if err != nil {
			log.Fatal().Err(err).Msg("Marshalling of transformed devices failed")
		}

		if err := dataio.WriteOutput(convertOutputFile, transformedDevicesAsJson); err != nil {
			log.Fatal().Err(err).Str("file-path", convertOutputFile).Msg("Error writing output")
		}
	},
}

func init() {
	ConvertCmd.Flags().StringVarP(&convertInputFile, "input", "i", "", "input filename (default stdin)")
	ConvertCmd.Flags().StringVarP(&convertOutputFile, "output", "o", "", "output filename (default stdout)")
}
