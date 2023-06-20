/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package main

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"
  "flag"
  "fmt"
  "os"
  "os/signal"
  "syscall"

  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/cdm-dcd-reference/reference"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/dcd"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/logging"

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
  var serverAddress string
  var registryAddress string
  flag.StringVar(&logLevel, "log-level", "info", fmt.Sprintf("set log level. one of: %s,%s,%s,%s,%s,%s,%s",
    zerolog.TraceLevel.String(),
    zerolog.DebugLevel.String(),
    zerolog.InfoLevel.String(),
    zerolog.WarnLevel.String(),
    zerolog.ErrorLevel.String(),
    zerolog.FatalLevel.String(),
    zerolog.PanicLevel.String()))
  flag.StringVar(&serverAddress, "grpc-address", "mydcd:8081", "gRPC server address")
  flag.StringVar(&registryAddress, "grpc-registry-address", "grpc-server-registry:50051", "gRPC server address")

  // Parse the CLI flags
  flag.Parse()

  // Set log level
  logging.AdjustLogLevel(logLevel)
  // Register dcd implementation
  myDCDImplementation := new(reference.ReferenceClassDriver)
  dcdImpl := dcd.New("cdm-dcd-reference", metadata.Version{version, commit, date}).
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
  if err := dcdImpl.Start(serverAddress, registryAddress); err != nil {
    log.Fatal().Err(err).Msg("Could not start device class driver instance")
  }
}
