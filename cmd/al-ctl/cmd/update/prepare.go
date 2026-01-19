/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package update

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var prepareParams al.UpdateParams

// UpdatePrepareCommand represents the artefact prepare command
var UpdatePrepareCommand = &cobra.Command{
	Use:   "prepare",
	Short: "Prepare update on device",
	Long:  `Prepares a firmware/software update on the specified device`,
	Run: func(cmd *cobra.Command, args []string) {
		err := al.PrepareUpdate(shared.AssetLinkEndpoint, prepareParams)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to prepare update")
		}
	},
}

func init() {
	UpdatePrepareCommand.Flags().StringVarP(&prepareParams.JobId, "job-id", "j", "", shared.JobIdDesc)
	UpdatePrepareCommand.Flags().StringVarP(&prepareParams.ArtefactFile, "artefact-file", "a", "", "source filename of artefact")
	UpdatePrepareCommand.Flags().StringVarP(&prepareParams.ArtefactType, "artefact-type", "t", "", "provided artefact type (\"backup\", \"configuration\", or \"firmware\")")
	UpdatePrepareCommand.Flags().StringVarP(&prepareParams.DeviceIdentifierFile, "device-identifier-file", "d", "", shared.DeviceIdentifierFileDesc)
	UpdatePrepareCommand.Flags().BoolVarP(&prepareParams.ConvertDeviceIdentifier, "convert-device-identifier", "c", false, shared.ConvertDeviceIdentifierDesc)
	UpdatePrepareCommand.Flags().StringVarP(&prepareParams.DeviceCredentialsFile, "device-credentials-file", "l", "", shared.DeviceCredentialsFileDesc)
	UpdatePrepareCommand.Flags().StringVarP(&prepareParams.ArtefactCredentialsFile, "artefact-credentials-file", "", "", shared.ArtefactCredentialsFileDesc)
}
