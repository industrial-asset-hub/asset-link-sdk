/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package discovery

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/dcdconnection"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/shared"
  "github.com/spf13/cobra"
)

// stopCmd represents the cancel command
var stopCmd = &cobra.Command{
  Use:   "stop",
  Short: "Stop discovery job",
  Long:  `This command stops an discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
    dcdconnection.StopDiscovery(shared.DcdEndpoint)
  },
}

func init() {
  DiscoveryCmd.AddCommand(stopCmd)
}
