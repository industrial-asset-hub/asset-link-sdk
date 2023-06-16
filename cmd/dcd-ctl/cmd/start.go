/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
  Use:   "start",
  Short: "A brief description of your command",
  Long:  `This command starts an discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("start called")
  },
}

func init() {
  discoveryCmd.AddCommand(startCmd)
}
