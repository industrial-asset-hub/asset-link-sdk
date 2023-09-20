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

var filters string
var options string

// startCmd represents the start command
var startCmd = &cobra.Command{
  Use:   "start",
  Short: "Start discovery job",
  Long:  `This command starts an discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
    dcdconnection.StartDiscovery(shared.DcdEndpoint, options, filters)
  },
}

func init() {
  DiscoveryCmd.AddCommand(startCmd)
  // TODO: introduce examples
  startCmd.PersistentFlags().StringVarP(&options, "options", "o", "[]",
    `Discovery options

Key/Value: TODO Description

Operator:
    EQUAL = 0
    NOT_EQUAL = 1
    GREATER_THAN = 2
    GREATER_THAN_OR_EQUAL_TO = 3
    LESS_THAN = 4
    LESS_THAN_OR_EQUAL_TO = 5

Please be aware to use quotes on our commandline

Example options:
  - [{"key": "test", "value": "value", "operator": 1}]`,
  )

  startCmd.PersistentFlags().StringVarP(&filters, "filters", "f", "[]",
    `Discovery filters

Key/Value: TODO Description

Operator:
    EQUAL = 0
    NOT_EQUAL = 1
    GREATER_THAN = 2
    GREATER_THAN_OR_EQUAL_TO = 3
    LESS_THAN = 4
    LESS_THAN_OR_EQUAL_TO = 5

Please be aware to use quotes on our commandline

Example filters:
  - [{"key": "test", "value": "value", "operator": 1}]`,
  )
}
