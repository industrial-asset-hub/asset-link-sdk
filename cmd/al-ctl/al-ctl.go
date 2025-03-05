/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package main

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/cmd"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/logging"
)

func main() {
	logging.SetupLogging()
	logging.SetColorForLogLevel()
	cmd.SetVersionInfo()
	cmd.Execute()
}
