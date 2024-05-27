/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package discovery

import (
	"github.com/spf13/cobra"
)

var (
	jobId uint32 = 1
)

// discoveryCmd represents the discovery command
var DiscoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Use discovery feature of an AssetLink",
	Long: `This command allows to start/stop and receive the results of an
discovery job.

Example workflow:
  - Start discovery job
  - Subscribe to receive discovered devices as stream
`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}
