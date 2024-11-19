/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package info

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v2/cmd/dcd-ctl/internal/dcd"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/cmd/dcd-ctl/internal/shared"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print asset link information",
	Long:  `This command prints information on the asset link.`,
	Run: func(cmd *cobra.Command, args []string) {
		dcd.PrintInfo(shared.AssetLinkEndpoint)
	},
	// Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}
