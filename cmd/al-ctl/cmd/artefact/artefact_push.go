/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefact

import (
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/spf13/cobra"
)

var pushArtefactFile string = ""
var pushDeviceId string = ""

// artefactPushCommand represents the artefact push command
var ArtefactPushCommand = &cobra.Command{
	Use:   "push",
	Short: "Push artefact to device",
	Long:  `Pushes an artefact (e.g., a software update file) to the specified device`,
	Run: func(cmd *cobra.Command, args []string) {
		exitCode := al.PushArtefact(shared.AssetLinkEndpoint, pushArtefactFile, pushDeviceId)
		os.Exit(exitCode)
	},
}

func init() {
	ArtefactPushCommand.Flags().StringVarP(&pushArtefactFile, "artefact-file", "a", "result.json", "source filename of artefact")
	ArtefactPushCommand.Flags().StringVarP(&pushDeviceId, "device-id", "d", "", "device identifier")
}
