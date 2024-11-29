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

	"github.com/industrial-asset-hub/asset-link-sdk/v2/cmd/dcd-ctl/internal/dcd"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/cmd/dcd-ctl/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var filters string
var options string
var outputFile string = ""

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start discovery job",
	Long:  `This command starts an discovery job.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp := dcd.StartDiscovery(shared.AssetLinkEndpoint, options, filters)

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

func init() {
	DiscoveryCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&outputFile, "output-file", "c", "result.json", "output format")
	// TODO: introduce examples
	startCmd.PersistentFlags().StringVarP(&options, "options", "o", "[]",
		shared.DiscoveryOptionsDesc,
	)

	startCmd.PersistentFlags().StringVarP(&filters, "filters", "f", "[]",
		shared.DiscoveryFiltersDesc,
	)
}
