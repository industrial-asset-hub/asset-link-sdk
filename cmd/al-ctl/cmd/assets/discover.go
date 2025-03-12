/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assets

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var discoverConfigFile string = ""
var discoverOutputFile string = ""

var DiscoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Trigger discovery",
	Long:  `This command triggers a discovery job and retrieves the assets.`,
	Run: func(cmd *cobra.Command, args []string) {
		discoverResponses, err := al.Discover(shared.AssetLinkEndpoint, discoverConfigFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Error during discovery")
		}

		err = al.WriteDiscoveryResponsesFile(discoverOutputFile, discoverResponses)
		if err != nil {
			log.Fatal().Err(err).Msg("Error writing discovery responses to file")
		}
	},
}

func init() {
	DiscoverCmd.Flags().StringVarP(&discoverConfigFile, "discovery-file", "d", "", shared.DiscoveryFileDesc)
	DiscoverCmd.Flags().StringVarP(&discoverOutputFile, "output", "o", "", "output filename (default stdout)")
}
