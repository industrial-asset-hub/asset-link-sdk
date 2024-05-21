/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */
package cmd

import (
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/cmd/info"
	"fmt"
	"os"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/cmd/discovery"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cmd/dcd-ctl/internal/shared"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/logging"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	logLevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dcd-ctl",
	Short: "command line interface to interact with device-class-drivers",
	Long: `This command line interfaces allows to interact with the so called
DCDs (Device Class Drivers).

This can be useful for validation purposes inside CI/CD pipelines or just
to ease development efforts.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initHandlers)
	rootCmd.PersistentFlags().StringVarP(&shared.DcdEndpoint, "endpoint", "e", "localhost:8081", "gRPC Server Address of the DCD")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "info",
		fmt.Sprintf("set log level. one of: %s,%s,%s,%s,%s,%s,%s",
			zerolog.TraceLevel.String(),
			zerolog.DebugLevel.String(),
			zerolog.InfoLevel.String(),
			zerolog.WarnLevel.String(),
			zerolog.ErrorLevel.String(),
			zerolog.FatalLevel.String(),
			zerolog.PanicLevel.String()))

	rootCmd.AddCommand(discovery.DiscoveryCmd)
	rootCmd.AddCommand(info.InfoCmd)

}
func initHandlers() {
	logging.SetupLogging()
	logging.AdjustLogLevel(logLevel)
}