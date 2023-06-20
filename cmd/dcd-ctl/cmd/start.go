/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package cmd

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/dcdconnection"
  "github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
  Use:   "start",
  Short: "Start Discovery job",
  Long:  `This command starts an discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
    dcdconnection.StartDiscovery(dcdEndpoint)
  },
}

func init() {
  discoveryCmd.AddCommand(startCmd)
}
