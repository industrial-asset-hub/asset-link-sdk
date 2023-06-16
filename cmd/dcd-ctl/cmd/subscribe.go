/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
  Use:   "subscribe",
  Short: "Subscribe to discovery results",
  Long:  `This commands subscribes to the results of an discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("subscribe called")
  },
}

func init() {
  discoveryCmd.AddCommand(subscribeCmd)
  subscribeCmd.Flags().BoolP("output", "o", false, "output format")
}
