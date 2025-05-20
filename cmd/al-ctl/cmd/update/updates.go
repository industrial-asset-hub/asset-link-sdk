/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package update

import (
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update commands",
	Long:  `Commands related to update management (e.g., firmware updates, software updates). An update consists of two steps: prepare and activate.`,
}

func init() {
	UpdateCmd.AddCommand(UpdatePrepareCommand)
	UpdateCmd.AddCommand(UpdateActivateCommand)
}
