/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package discovery

import (
	"encoding/json"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var discoveryFile string = ""
var outputFile string = ""

// discoverCmd represents the discovery command
var DiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Start discovery job",
	Long:  `This command starts an discovery job and prints the result.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := al.Discover(shared.AssetLinkEndpoint, discoveryFile)
		if err != nil {
			log.Fatal().Err(err).Msg("error during discovery")
		}
		log.Info().Str("File", outputFile).Msg("Response saved to file")
		f, err := os.Create(outputFile)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating file")
		}
		defer f.Close()

		asJson, _ := json.MarshalIndent(resp, "", "  ")
		_, err = f.Write(asJson)
		if err != nil {
			log.Fatal().Err(err).Msg("error during writing of the json file")
		}

	},
}

func init() {
	DiscoverCmd.Flags().StringVarP(&outputFile, "output-file", "o", "result.json", "output file for the discovery result")
	DiscoverCmd.Flags().StringVarP(&discoveryFile, "discovery-file", "d", "", shared.DiscoveryFileDesc)
}
