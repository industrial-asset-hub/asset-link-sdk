/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package cmd

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/dcdconnection"
  "encoding/json"
  "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  "os"
)

var (
  outputFile = ""
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
  Use:   "subscribe",
  Short: "Subscribe to discovery results",
  Long:  `This commands subscribes to the results of an discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
    resp := dcdconnection.Subscribe(dcdEndpoint)

    log.Trace().Str("File", outputFile).Msg("Saving to file")
    f, _ := os.Create(outputFile)
    defer f.Close()

    asJson, _ := json.MarshalIndent(resp, "", "  ")
    _, err := f.Write(asJson)
    if err != nil {
      log.Err(err).Msg("error during writing of the json file")
    }
  },
}

func init() {
  discoveryCmd.AddCommand(subscribeCmd)
  subscribeCmd.Flags().StringVarP(&outputFile, "output-file", "o", "result.json", "output format")
}
