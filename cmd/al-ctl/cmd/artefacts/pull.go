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

var pullArtefactFile string = ""
var pullArtefactType string = ""
var pullDeviceId string = ""

// artefactPullCommand represents the discovery command
var ArtefactPullCommand = &cobra.Command{
	Use:   "pull",
	Short: "Pull artefact from device",
	Long:  `Pulls an artefact of a specific type from the specified device`,
	Run: func(cmd *cobra.Command, args []string) {
		exitCode := al.PullArtefact(shared.AssetLinkEndpoint, pullArtefactFile, pullDeviceId, pullArtefactType)
		os.Exit(exitCode)
	},
}

func init() {
	ArtefactPullCommand.Flags().StringVarP(&pullArtefactFile, "artefact-file", "a", "result.json", "destination filename of artefact")
	ArtefactPullCommand.Flags().StringVarP(&pullArtefactType, "artefact-type", "t", "", "requested artefact type (\"backup\", \"configuration\", or \"firmware\")")
	ArtefactPullCommand.Flags().StringVarP(&pullDeviceId, "device-id", "d", "", "device identifier")
}
