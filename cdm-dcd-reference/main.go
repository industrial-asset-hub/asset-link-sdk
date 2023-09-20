/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cdm-dcd-reference/reference"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/dcd"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/logging"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// values provided by linker
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	logging.SetupLogging()
	log.Info().
		Str("Version", version).
		Str("commit", commit).
		Str("date", date).
		Msg("Starting cdm-dcd-reference implementation")

	// Setup log of log infrastructure
	var logLevel string
	var grpcServerAddress string
	var httpServerAddress string
	var registryAddress string
	flag.StringVar(&logLevel, "log-level", "info", fmt.Sprintf("set log level. one of: %s,%s,%s,%s,%s,%s,%s",
		zerolog.TraceLevel.String(),
		zerolog.DebugLevel.String(),
		zerolog.InfoLevel.String(),
		zerolog.WarnLevel.String(),
		zerolog.ErrorLevel.String(),
		zerolog.FatalLevel.String(),
		zerolog.PanicLevel.String()))
	flag.StringVar(&grpcServerAddress, "grpc-address", "localhost:8081", "gRPC server endpoint")
	flag.StringVar(&httpServerAddress, "http-address", "localhost:8082", "HTTP server endpoint")
	flag.StringVar(&registryAddress, "grpc-registry-address", "grpc-server-registry:50051", "gRPC registry address")

	// Parse the CLI flags
	flag.Parse()

	// Set log level
	logging.AdjustLogLevel(logLevel)
	// Register dcd implementation
	myDCDImplementation := new(reference.ReferenceClassDriver)
	dcdImpl := dcd.New(metadata.Metadata{
		Version: metadata.Version{Version: version, Commit: commit, Date: date},
		DcdName: "cdm-dcd-reference",
		Vendor:  "Siemens AG",
	}).
		Discovery(myDCDImplementation).
		SoftwareUpdate(myDCDImplementation).
		Build()

	// Signal handler for a proper shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(d *dcd.DCD) {
		<-c
		log.Info().Msg("Received SIGTERM. Exiting.")
		d.Stop()
		os.Exit(1)
	}(dcdImpl)

	// Start device class driver
	if err := dcdImpl.Start(grpcServerAddress, registryAddress, httpServerAddress); err != nil {
		log.Fatal().Err(err).Msg("Could not start device class driver instance")
	}
}
