/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefacts

import (
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/spf13/cobra"
)

var pushArtefactFile string = ""
var pushArtefactType string = ""
var pushDeviceIdentifierFile string = ""
var pushConvertDeviceIdentifier bool = false

// artefactPushCommand represents the artefact push command
var ArtefactPushCommand = &cobra.Command{
	Use:   "push",
	Short: "Push artefact to device",
	Long:  `Pushes an artefact (e.g., a software update file) to the specified device`,
	Run: func(cmd *cobra.Command, args []string) {
		exitCode := al.PushArtefact(shared.AssetLinkEndpoint, pushArtefactFile, pushArtefactType, pushDeviceIdentifierFile, pushConvertDeviceIdentifier)
		os.Exit(exitCode)
	},
}

func init() {
	ArtefactPushCommand.Flags().StringVarP(&pushArtefactFile, "artefact-file", "a", "", "source filename of artefact")
	ArtefactPushCommand.Flags().StringVarP(&pushArtefactType, "artefact-type", "t", "", "provided artefact type (\"backup\", \"configuration\", or \"firmware\")")
	ArtefactPushCommand.Flags().StringVarP(&pushDeviceIdentifierFile, "device-identifier-file", "d", "", shared.DeviceIdentifierFileDesc)
	ArtefactPushCommand.Flags().BoolVarP(&pushConvertDeviceIdentifier, "convert-device-identifier", "c", false, shared.ConvertDeviceIdentifierDesc)
}
