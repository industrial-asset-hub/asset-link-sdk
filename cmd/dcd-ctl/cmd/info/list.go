/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package info

import (
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/cmd/dcd-ctl/internal/registry"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/cmd/dcd-ctl/internal/shared"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered asset links",
	Long:  `List all asset links registered in the registry`,
	Run: func(cmd *cobra.Command, args []string) {
		registry.GetList(shared.RegistryEndpoint)
	},
	//Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}
