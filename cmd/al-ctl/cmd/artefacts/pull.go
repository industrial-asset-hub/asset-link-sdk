/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefacts

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var pullParams al.ArtefactParams

// artefactPullCommand represents the discovery command
var ArtefactPullCommand = &cobra.Command{
	Use:   "pull",
	Short: "Pull artefact from device",
	Long:  `Pulls an artefact of a specific type from the specified device`,
	Run: func(cmd *cobra.Command, args []string) {
		err := al.PullArtefact(shared.AssetLinkEndpoint, pullParams)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to pull artefact")
		}
	},
}

func init() {
	ArtefactPullCommand.Flags().StringVarP(&pullParams.JobId, "job-id", "j", "", shared.JobIdDesc)
	ArtefactPullCommand.Flags().StringVarP(&pullParams.ArtefactFile, "artefact-file", "a", "", "destination filename of artefact")
	ArtefactPullCommand.Flags().StringVarP(&pullParams.ArtefactType, "artefact-type", "t", "", "requested artefact type (\"configuration\", \"backup\", or \"log\")")
	ArtefactPullCommand.Flags().StringVarP(&pullParams.DeviceIdentifierFile, "device-identifier-file", "d", "", shared.DeviceIdentifierFileDesc)
	ArtefactPullCommand.Flags().BoolVarP(&pullParams.ConvertDeviceIdentifier, "convert-device-identifier", "c", false, shared.ConvertDeviceIdentifierDesc)
	ArtefactPullCommand.Flags().StringVarP(&pullParams.DeviceCredentialsFile, "device-credentials-file", "l", "", shared.DeviceCredentialsFileDesc)
	ArtefactPullCommand.Flags().StringVarP(&pullParams.ArtefactCredentialsFile, "artefact-credentials-file", "", "", shared.ArtefactCredentialsFileDesc)
}
