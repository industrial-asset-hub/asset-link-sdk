/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package discovery

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/dcdconnection"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/shared"
  "github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
  Use:   "start",
  Short: "Start discovery job",
  Long:  `This command starts an discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
    dcdconnection.StartDiscovery(shared.DcdEndpoint)
  },
}

func init() {
  DiscoveryCmd.AddCommand(startCmd)
}
