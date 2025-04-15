/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/metadata"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/assetlink"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cdm-al-reference/reference"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/logging"

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
		Msg("Starting cdm-al-reference implementation")

	// Setup log of log infrastructure
	var logLevel string
	var grpcServerAddress, grpcServerEndpointAddress, httpServerAddress, registryAddress, visuServerAddress string
	flag.StringVar(&logLevel, "log-level", "info", fmt.Sprintf("set log level. one of: %s,%s,%s,%s,%s,%s,%s",
		zerolog.TraceLevel.String(),
		zerolog.DebugLevel.String(),
		zerolog.InfoLevel.String(),
		zerolog.WarnLevel.String(),
		zerolog.ErrorLevel.String(),
		zerolog.FatalLevel.String(),
		zerolog.PanicLevel.String()))
	flag.StringVar(&grpcServerAddress, "grpc-server-address", "localhost:8081", "gRPC server endpoint")
	flag.StringVar(&grpcServerEndpointAddress, "grpc-server-endpoint-address", "localhost", "Address which is registered")
	flag.StringVar(&httpServerAddress, "http-address", "localhost:8082", "HTTP server endpoint")
	flag.StringVar(&visuServerAddress, "visu-address", "localhost:8083", "Device visualization server endpoint")
	flag.StringVar(&registryAddress, "grpc-registry-address", "grpc-server-registry:50051", "gRPC registry address")

	// Parse the CLI flags
	flag.Parse()

	// Set log level
	logging.AdjustLogLevel(logLevel)
	// Register al implementation
	myAssetLinkImplementation := new(reference.ReferenceAssetLink)
	alImpl := assetlink.New(metadata.Metadata{
		Version: metadata.Version{Version: version, Commit: commit, Date: date},
		AlId:    "siemens.cdm.al.reference",
		AlName:  "Reference Asset Link",
		Vendor:  "Siemens AG",
	}).
		Discovery(myAssetLinkImplementation).
		Build()

	// Signal handler for a proper shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(d *assetlink.AssetLink) {
		<-c
		log.Info().Msg("Received SIGTERM. Exiting.")
		d.Stop()
		os.Exit(1)
	}(alImpl)

	// Start simulated devices
	simdevices.StartSimulatedDevices(visuServerAddress)

	// Start asset link
	if err := alImpl.Start(grpcServerAddress, grpcServerEndpointAddress, registryAddress, httpServerAddress); err != nil {
		log.Fatal().Err(err).Msg("Could not start asset link instance")
	}
}
