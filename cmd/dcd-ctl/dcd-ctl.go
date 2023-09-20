/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package main

import (
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/cmd"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/logging"
)

func main() {
	logging.SetupLogging()
	cmd.SetVersionInfo()
	cmd.Execute()
}
