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

var pullJobId string = ""
var pullArtefactFile string = ""
var pullArtefactType string = ""
var pullDeviceIdentifierFile string = ""
var pullConvertDeviceIdentifier bool = false

// artefactPullCommand represents the discovery command
var ArtefactPullCommand = &cobra.Command{
	Use:   "pull",
	Short: "Pull artefact from device",
	Long:  `Pulls an artefact of a specific type from the specified device`,
	Run: func(cmd *cobra.Command, args []string) {
		err := al.PullArtefact(shared.AssetLinkEndpoint, pullJobId, pullArtefactFile, pullArtefactType, pullDeviceIdentifierFile, pullConvertDeviceIdentifier)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to pull artefact")
		}
	},
}

func init() {
	ArtefactPullCommand.Flags().StringVarP(&pullJobId, "job-id", "j", "", shared.JobIdDesc)
	ArtefactPullCommand.Flags().StringVarP(&pullArtefactFile, "artefact-file", "a", "", "destination filename of artefact")
	ArtefactPullCommand.Flags().StringVarP(&pullArtefactType, "artefact-type", "t", "", "requested artefact type (\"configuration\", \"backup\", or \"log\")")
	ArtefactPullCommand.Flags().StringVarP(&pullDeviceIdentifierFile, "device-identifier-file", "d", "", shared.DeviceIdentifierFileDesc)
	ArtefactPullCommand.Flags().BoolVarP(&pullConvertDeviceIdentifier, "convert-device-identifier", "c", false, shared.ConvertDeviceIdentifierDesc)
}
