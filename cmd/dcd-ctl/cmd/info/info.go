/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package info

import (
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/cmd/dcd-ctl/internal/dcdconnection"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/cmd/dcd-ctl/internal/shared"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information of the AssetLink",
	Long:  `This command gathers the info of the command.`,
	Run: func(cmd *cobra.Command, args []string) {
		dcdconnection.GetInfo(shared.AssetLinkEndpoint)
	},
	//Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}
