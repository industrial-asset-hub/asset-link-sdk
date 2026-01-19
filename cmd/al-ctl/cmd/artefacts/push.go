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

var pushParams al.ArtefactParams

// artefactPushCommand represents the artefact push command
var ArtefactPushCommand = &cobra.Command{
	Use:   "push",
	Short: "Push artefact to device",
	Long:  `Pushes an artefact of a specific type to the specified device`,
	Run: func(cmd *cobra.Command, args []string) {
		err := al.PushArtefact(shared.AssetLinkEndpoint, pushParams)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to push artefact")
		}
	},
}

func init() {
	ArtefactPushCommand.Flags().StringVarP(&pushParams.JobId, "job-id", "j", "", shared.JobIdDesc)
	ArtefactPushCommand.Flags().StringVarP(&pushParams.ArtefactFile, "artefact-file", "a", "", "source filename of artefact")
	ArtefactPushCommand.Flags().StringVarP(&pushParams.ArtefactType, "artefact-type", "t", "", "provided artefact type (\"configuration\", \"backup\", or \"log\")")
	ArtefactPushCommand.Flags().StringVarP(&pushParams.DeviceIdentifierFile, "device-identifier-file", "d", "", shared.DeviceIdentifierFileDesc)
	ArtefactPushCommand.Flags().BoolVarP(&pushParams.ConvertDeviceIdentifier, "convert-device-identifier", "c", false, shared.ConvertDeviceIdentifierDesc)
	ArtefactPushCommand.Flags().StringVarP(&pushParams.DeviceCredentialsFile, "device-credentials-file", "l", "", shared.DeviceCredentialsFileDesc)
	ArtefactPushCommand.Flags().StringVarP(&pushParams.ArtefactCredentialsFile, "artefact-credentials-file", "", "", shared.ArtefactCredentialsFileDesc)
}
