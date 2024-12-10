/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/logging"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/registry/internal/server"

	pb "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/conn_suite_registry"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
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
		Msg("Starting grpc-server-registry")

	pflag.String("log-level", "info", fmt.Sprintf("set log level. one of: %s,%s,%s,%s,%s,%s,%s",
		zerolog.TraceLevel.String(),
		zerolog.DebugLevel.String(),
		zerolog.InfoLevel.String(),
		zerolog.WarnLevel.String(),
		zerolog.ErrorLevel.String(),
		zerolog.FatalLevel.String(),
		zerolog.PanicLevel.String()))
	pflag.String("server-address", ":50051", "gRPC server address")
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen.")
	}

	viper.SetEnvPrefix("registry")
	viper.AutomaticEnv()

	// parse the CLI flags
	pflag.Parse()

	// Replace - with _ inside environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Read Flags/environment vars
	logLevel := viper.GetString("log-level")
	serverAddress := viper.GetString("server-address")

	// Set log level
	logging.AdjustLogLevel(logLevel)

	log.Debug().Msg("Debug level enabled.")

	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen.")
	}

	// start GRPC Server
	grpcServer := grpc.NewServer()

	pb.RegisterRegistryApiServer(grpcServer, server.NewServer())
	log.Info().
		Str("address", serverAddress).
		Msg("Serving gRPC Server")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().
			Err(err).
			Str("address", serverAddress).
			Msg("Failed to serve gPRC server!")
	}
}
