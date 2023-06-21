/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package status

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/dcdconnection"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/shared"
  "github.com/spf13/cobra"
)

// statusCmd represents the status command
var StatusCmd = &cobra.Command{
  Use:   "status",
  Short: "Get Status of the DCD",
  Long:  `This command gathers the status of the command.`,
  Run: func(cmd *cobra.Command, args []string) {
    dcdconnection.GetHealth(shared.DcdEndpoint)
    dcdconnection.GetVersion(shared.DcdEndpoint)
  },
  //Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}
