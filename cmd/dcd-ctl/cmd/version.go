/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
  "fmt"
  "runtime"
)

var (
  // values provided by linker
  version = "dev"
  commit  = "unknown"
  date    = "unknown"
)

func SetVersionInfo() {
  goversion := runtime.Version()
  rootCmd.Version = fmt.Sprintf("%s\nBuild Time: %s\nCommit: %s\nGoVersion: %s", version, date, commit, goversion)
}
