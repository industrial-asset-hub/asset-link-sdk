/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package update

import (
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/spf13/cobra"
)

var activateJobId string = ""
var activateArtefactFile string = ""
var activateArtefactType string = ""
var activateDeviceIdentifierFile string = ""
var activateConvertDeviceIdentifier bool = false

// UpdateActivateCommand represents the artefact prepare command
var UpdateActivateCommand = &cobra.Command{
	Use:   "activate",
	Short: "Activate update on device",
	Long:  `Activates a firmware/software update on the specified device`,
	Run: func(cmd *cobra.Command, args []string) {
		exitCode := al.ActivateUpdate(shared.AssetLinkEndpoint, activateJobId, activateArtefactFile, activateArtefactType, activateDeviceIdentifierFile, activateConvertDeviceIdentifier)
		os.Exit(exitCode)
	},
}

func init() {
	UpdateActivateCommand.Flags().StringVarP(&activateJobId, "job-id", "j", "", shared.JobIdDesc)
	UpdateActivateCommand.Flags().StringVarP(&activateArtefactFile, "artefact-file", "a", "", "source filename of artefact")
	UpdateActivateCommand.Flags().StringVarP(&activateArtefactType, "artefact-type", "t", "", "provided artefact type (\"software\" or \"firmware\")")
	UpdateActivateCommand.Flags().StringVarP(&activateDeviceIdentifierFile, "device-identifier-file", "d", "", shared.DeviceIdentifierFileDesc)
	UpdateActivateCommand.Flags().BoolVarP(&activateConvertDeviceIdentifier, "convert-device-identifier", "c", false, shared.ConvertDeviceIdentifierDesc)
}
