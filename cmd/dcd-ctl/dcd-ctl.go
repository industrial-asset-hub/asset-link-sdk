/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package main

import (
	"code.siemens.com/common-device-management/shared/cdm-dcd-sdk/v2/cmd/dcd-ctl/cmd"
	"code.siemens.com/common-device-management/shared/cdm-dcd-sdk/v2/logging"
)

func main() {
	logging.SetupLogging()
	cmd.SetVersionInfo()
	cmd.Execute()
}
