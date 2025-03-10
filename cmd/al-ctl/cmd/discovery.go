/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package cmd

import (
	"encoding/json"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/discovery"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

// discoverCmd represents the discovery command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Start discovery job",
	Long:  `This command starts an discovery job and prints the result.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := discovery.Discover(assetLinkEndpoint, discoveryFile, timeoutInSeconds)
		if err != nil {
			log.Fatal().Err(err).Msg("error during discovery")
		}
		log.Trace().Str("File", outputFile).Msg("Saving to file")
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
	discoverCmd.Flags().StringVarP(&outputFile, "output-file", "o", "result.json", "output file")
	discoverCmd.Flags().StringVarP(&discoveryFile, "discovery-file", "d", "", DiscoveryFileDesc)
}
