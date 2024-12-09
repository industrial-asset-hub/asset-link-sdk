/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package info

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/registry"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered asset links",
	Long:  `This command lists all asset links registered in the registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		registry.PrintList(shared.RegistryEndpoint)
	},
	//Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}
