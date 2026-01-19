/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package update

import (
	"github.com/rs/zerolog/log"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/spf13/cobra"
)

var cancelParams al.UpdateParams

// UpdateActivateCommand represents the artefact prepare command
var UpdateCancelCommand = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel update on device",
	Long:  `Cancels a firmware/software update on the specified device (after the repare step)`,
	Run: func(cmd *cobra.Command, args []string) {
		err := al.CancelUpdate(shared.AssetLinkEndpoint, cancelParams)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to cancel update")
		}
	},
}

func init() {
	UpdateCancelCommand.Flags().StringVarP(&cancelParams.JobId, "job-id", "j", "", shared.JobIdDesc)
	UpdateCancelCommand.Flags().StringVarP(&cancelParams.ArtefactType, "artefact-type", "t", "", "provided artefact type (\"software\" or \"firmware\")")
	UpdateCancelCommand.Flags().StringVarP(&cancelParams.DeviceIdentifierFile, "device-identifier-file", "d", "", shared.DeviceIdentifierFileDesc)
	UpdateCancelCommand.Flags().BoolVarP(&cancelParams.ConvertDeviceIdentifier, "convert-device-identifier", "c", false, shared.ConvertDeviceIdentifierDesc)
	UpdateCancelCommand.Flags().StringVarP(&cancelParams.DeviceCredentialsFile, "device-credentials-file", "l", "", shared.DeviceCredentialsFileDesc)
	UpdateCancelCommand.Flags().StringVarP(&cancelParams.ArtefactCredentialsFile, "artefact-credentials-file", "", "", shared.ArtefactCredentialsFileDesc)
}
