/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package cmd

import (
	"encoding/json"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/al-ctl/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/al-ctl/registry"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/al-ctl/shared"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var discoveryFile string = ""
var outputFile string = ""

// DiscoverCmd represents the discovery command
var DiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Start discovery job",
	Long:  `This command starts an discovery job and prints the result.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := al.Discover(shared.AssetLinkEndpoint, discoveryFile)

		log.Trace().Str("File", outputFile).Msg("Saving to file")
		f, _ := os.Create(outputFile)
		defer f.Close()

		asJson, _ := json.MarshalIndent(resp, "", "  ")
		_, err := f.Write(asJson)
		if err != nil {
			log.Err(err).Msg("error during writing of the json file")
		}

	},
}

// InfoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print asset link information",
	Long:  `This command prints information on the asset link.`,
	Run: func(cmd *cobra.Command, args []string) {
		registry.PrintInfo(shared.AssetLinkEndpoint)
	},
}

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered asset links",
	Long:  `This command lists all asset links registered in the registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		registry.PrintList(shared.RegistryEndpoint)
	},
}

func init() {
	DiscoverCmd.Flags().StringVarP(&outputFile, "output-file", "o", "result.json", "output file")
	DiscoverCmd.Flags().StringVarP(&discoveryFile, "discovery-file", "d", "", shared.DiscoveryFileDesc)
}
