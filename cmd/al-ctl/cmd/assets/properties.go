/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assets

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var requestFilePath = ""
var outputFile = ""

var PropertyCmd = &cobra.Command{
	Use:   "properties",
	Short: "Trigger get device properties",
	Long:  `This command triggers a get property values job and retrieves device info properties.`,
	Run: func(cmd *cobra.Command, args []string) {
		properties, err := al.GetDeviceInfoProperties(shared.AssetLinkEndpoint, requestFilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Error during get property values")
		}

		err = al.WritePropertyResponsesFile(outputFile, properties)
		if err != nil {
			log.Fatal().Err(err).Msg("Error writing properties to file")
		}
	},
}

func init() {
	PropertyCmd.Flags().StringVarP(&requestFilePath, "property-request-file-path", "p", "", "Path to the GetPropertyValues request file")
	PropertyCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output filename (default stdout)")
	_ = PropertyCmd.MarkFlagRequired("property-request-file-path")
}
