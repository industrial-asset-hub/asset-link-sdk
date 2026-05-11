/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package main

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/cmd"
)

func main() {
	cmd.SetVersionInfo()
	cmd.Execute()
}
