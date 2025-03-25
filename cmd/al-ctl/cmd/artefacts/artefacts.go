/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package artefacts

import (
	"github.com/spf13/cobra"
)

var ArtefactsCmd = &cobra.Command{
	Use:   "artefacts",
	Short: "Artefact commands",
	Long:  `Commands related to artefact management (e.g., firmware updates, backups).`,
}

func init() {
	ArtefactsCmd.AddCommand(ArtefactPullCommand)
	ArtefactsCmd.AddCommand(ArtefactPushCommand)
}
