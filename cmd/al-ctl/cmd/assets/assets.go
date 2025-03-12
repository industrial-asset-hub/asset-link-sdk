/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assets

import (
	"github.com/spf13/cobra"
)

var AssetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Asset commands",
	Long:  `Commands related to asset discovery.`,
}

func init() {
	AssetsCmd.AddCommand(DiscoverCmd)
	AssetsCmd.AddCommand(ConvertCmd)
}
