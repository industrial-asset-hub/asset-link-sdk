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

var cancelJobId string = ""
var cancelArtefactType string = ""
var cancelDeviceIdentifierFile string = ""
var cancelConvertDeviceIdentifier bool = false

// UpdateActivateCommand represents the artefact prepare command
var UpdateCancelCommand = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel update on device",
	Long:  `Cancels a firmware/software update on the specified device (after the repare step)`,
	Run: func(cmd *cobra.Command, args []string) {
		err := al.CancelUpdate(shared.AssetLinkEndpoint, cancelJobId, cancelArtefactType, cancelDeviceIdentifierFile, cancelConvertDeviceIdentifier)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to cancel update")
		}
	},
}

func init() {
	UpdateCancelCommand.Flags().StringVarP(&cancelJobId, "job-id", "j", "", shared.JobIdDesc)
	UpdateCancelCommand.Flags().StringVarP(&cancelArtefactType, "artefact-type", "t", "", "provided artefact type (\"software\" or \"firmware\")")
	UpdateCancelCommand.Flags().StringVarP(&cancelDeviceIdentifierFile, "device-identifier-file", "d", "", shared.DeviceIdentifierFileDesc)
	UpdateCancelCommand.Flags().BoolVarP(&cancelConvertDeviceIdentifier, "convert-device-identifier", "c", false, shared.ConvertDeviceIdentifierDesc)
}
