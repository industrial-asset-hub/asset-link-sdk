/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package info

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v2/cmd/dcd-ctl/internal/registry"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/cmd/dcd-ctl/internal/shared"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered asset links",
	Long:  `This command lists all asset links registered in the registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		registry.GetList(shared.RegistryEndpoint)
	},
	//Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}
