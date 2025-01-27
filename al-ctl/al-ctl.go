/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package alctl

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/al-ctl/cmd"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/logging"
)

func main() {
	logging.SetupLogging()
	cmd.SetVersionInfo()
	cmd.Execute()
}
