/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package cmd

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/registry"
	"github.com/spf13/cobra"
)

// InfoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print asset link information",
	Long:  `This command prints information on the asset link.`,
	Run: func(cmd *cobra.Command, args []string) {
		registry.PrintInfo(assetLinkEndpoint)
	},
}

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered asset links",
	Long:  `This command lists all asset links registered in the registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		registry.PrintList(registryEndpoint)
	},
}
