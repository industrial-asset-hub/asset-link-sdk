/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: {{cookiecutter.company}}
 *
 * Author: {{cookiecutter.author_name}} <{{cookiecutter.author_email}}>
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"

	"{{cookiecutter.dcd_name}}/handler"

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
	dcdName = "{{cookiecutter.dcd_name}}"
	vendor  = "{{cookiecutter.company}}"
)

func main() {
	logging.SetupLogging()
	log.Info().
		Str("Version", version).
		Str("commit", commit).
		Str("date", date).
		Msg("Starting " + dcdName + " driver")

	// Setup log of log infrastructure
	var logLevel string
	var grpcServerAddress string
	var registryAddress string
	var httpServerAddress string
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
	dcdImpl := new(handler.DCDImplementation)
	dcdInstance := dcd.New(metadata.Metadata{
		Version: metadata.Version{Version: version, Commit: commit, Date: date},
		DcdName: dcdName,
		Vendor:  vendor,
	}).
		Discovery(dcdImpl).
		//SoftwareUpdate(dcdImpl).
		Build()

	// Signal handler for a proper shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(d *dcd.DCD) {
		<-c
		log.Info().Msg("Received SIGTERM. Exiting.")
		d.Stop()
		os.Exit(1)
	}(dcdInstance)

	// Start device class driver
	if err := dcdInstance.Start(grpcServerAddress, registryAddress, httpServerAddress); err != nil {
		log.Fatal().Err(err).Msg("Could not start device class driver instance")
	}

}
