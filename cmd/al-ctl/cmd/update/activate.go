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

var activateParams al.UpdateParams

// UpdateActivateCommand represents the artefact prepare command
var UpdateActivateCommand = &cobra.Command{
	Use:   "activate",
	Short: "Activate update on device",
	Long:  `Activates a firmware/software update on the specified device`,
	Run: func(cmd *cobra.Command, args []string) {
		err := al.ActivateUpdate(shared.AssetLinkEndpoint, activateParams)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to activate update")
		}
	},
}

func init() {
	UpdateActivateCommand.Flags().StringVarP(&activateParams.JobId, "job-id", "j", "", shared.JobIdDesc)
	UpdateActivateCommand.Flags().StringVarP(&activateParams.ArtefactFile, "artefact-file", "a", "", "source filename of artefact")
	UpdateActivateCommand.Flags().StringVarP(&activateParams.ArtefactType, "artefact-type", "t", "", "provided artefact type (\"software\" or \"firmware\")")
	UpdateActivateCommand.Flags().StringVarP(&activateParams.DeviceIdentifierFile, "device-identifier-file", "d", "", shared.DeviceIdentifierFileDesc)
	UpdateActivateCommand.Flags().BoolVarP(&activateParams.ConvertDeviceIdentifier, "convert-device-identifier", "c", false, shared.ConvertDeviceIdentifierDesc)
	UpdateActivateCommand.Flags().StringVarP(&activateParams.DeviceCredentialsFile, "device-credentials-file", "l", "", shared.DeviceCredentialsFileDesc)
	UpdateActivateCommand.Flags().StringVarP(&activateParams.ArtefactCredentialsFile, "artefact-credentials-file", "", "", shared.ArtefactCredentialsFileDesc)
}
