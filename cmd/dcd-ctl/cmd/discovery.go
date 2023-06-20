/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package cmd

import (
  "github.com/spf13/cobra"
)

// discoveryCmd represents the discovery command
var discoveryCmd = &cobra.Command{
  Use:   "discovery",
  Short: "Use discovery feature of an DCD",
  Long: `This command allows to start/stop and receive the results of an
discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
  },
  Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}

func init() {
  rootCmd.AddCommand(discoveryCmd)
}
