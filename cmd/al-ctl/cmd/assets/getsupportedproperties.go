/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
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

var supportedPropertiesRequestFilePath = ""
var supportedPropertiesOutputFile = ""

var GetSupportedPropertiesCmd = &cobra.Command{
	Use:   "getsupportedproperties",
	Short: "Trigger get supported properties",
	Long:  `This command triggers a get supported properties request and retrieves supported device info properties.`,
	Run: func(cmd *cobra.Command, args []string) {
		supportedProperties, err := al.GetSupportedDeviceInfoProperties(shared.AssetLinkEndpoint, supportedPropertiesRequestFilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Error during get supported properties")
		}

		err = al.WriteSupportedPropertyResponsesFile(supportedPropertiesOutputFile, supportedProperties)
		if err != nil {
			log.Fatal().Err(err).Msg("Error writing supported properties to file")
		}
	},
}

func init() {
	GetSupportedPropertiesCmd.Flags().StringVarP(&supportedPropertiesRequestFilePath, "property-request-file-path", "p", "", "Path to the GetSupportedProperties request file")
	GetSupportedPropertiesCmd.Flags().StringVarP(&supportedPropertiesOutputFile, "output", "o", "", "Output filename (default stdout)")
	_ = GetSupportedPropertiesCmd.MarkFlagRequired("property-request-file-path")
}
