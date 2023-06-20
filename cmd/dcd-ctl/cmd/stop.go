/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/dcdconnection"
  "github.com/spf13/cobra"
)

// stopCmd represents the cancel command
var stopCmd = &cobra.Command{
  Use:   "stop",
  Short: "Stop discovery job",
  Long:  `This command stops an discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
    dcdconnection.StopDiscovery(dcdEndpoint)
  },
}

func init() {
  discoveryCmd.AddCommand(stopCmd)
}
