/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package main

import (
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/cmd"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/logging"
	"runtime"
)

var (
	// values provided by linker
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	logging.SetupLogging()
	SetVersionInfo()
	cmd.Execute()
}
func SetVersionInfo() {
	goversion := runtime.Version()
	cmd.RootCmd.Version = fmt.Sprintf("%s\nBuild Time: %s\nCommit: %s\nGoVersion: %s", version, date, commit, goversion)
}
