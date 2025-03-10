/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package cmd

import (
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "al-ctl",
	Short: "command line interface to interact with Asset Links",
	Long: `This command line interfaces allows to interact with the so called Asset Links (AL).

This can be useful for validation purposes inside CI/CD pipelines or just
to ease development efforts.`,
}

func init() {
	cobra.OnInitialize(initHandlers)
	RootCmd.AddCommand(discoverCmd)
	RootCmd.AddCommand(InfoCmd)
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(TestCmd)
	RootCmd.PersistentFlags().StringVarP(&registryEndpoint, "registry", "r", "localhost:50051", "gRPC Server Address of the Registry")
	RootCmd.PersistentFlags().StringVarP(&assetLinkEndpoint, "endpoint", "e", "localhost:8081", "gRPC Server Address of the AssetLink")
	RootCmd.PersistentFlags().UintVarP(&timeoutInSeconds, "timeout", "n", 0, "timeout in seconds (default none)")
	RootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "info",
		fmt.Sprintf("set log level. one of: %s,%s,%s,%s,%s,%s,%s",
			zerolog.TraceLevel.String(),
			zerolog.DebugLevel.String(),
			zerolog.InfoLevel.String(),
			zerolog.WarnLevel.String(),
			zerolog.ErrorLevel.String(),
			zerolog.FatalLevel.String(),
			zerolog.PanicLevel.String()))
}

func initHandlers() {
	logging.SetupLogging()
	logging.AdjustLogLevel(logLevel)
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("error during execution")
	}
}
