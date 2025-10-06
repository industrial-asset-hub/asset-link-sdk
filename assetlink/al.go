/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package assetlink

import (
	"fmt"
	"net"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/server/webserver"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/metadata"

	generatedDriverInfoServer "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/conn_suite_drv_info"
	generatedDiscoveryServer "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/features"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/registryclient"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/server/devicediscovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/internal/server/driverinfo"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Asset Link feature builder, according to the GoF build pattern
// The pattern provides methods to register new features in an easy
type alFeatureBuilder struct {
	metadata metadata.Metadata

	discovery features.Discovery
	generatedDiscoveryServer.DeviceDiscoverApiServer
}

// Methods to register new features
func (cb *alFeatureBuilder) Discovery(f features.Discovery) *alFeatureBuilder {
	cb.discovery = f
	return cb
}

// Builder
func New(metadata metadata.Metadata) *alFeatureBuilder {
	return &alFeatureBuilder{metadata: metadata}
}

func (cb *alFeatureBuilder) Build() *AssetLink {
	return &AssetLink{
		discoveryImpl:         cb.discovery,
		customDiscoveryServer: cb.DeviceDiscoverApiServer,
		metadata:              cb.metadata,
	}
}

// Structure of the features
type AssetLink struct {
	metadata              metadata.Metadata
	discoveryImpl         features.Discovery
	customDiscoveryServer generatedDiscoveryServer.DeviceDiscoverApiServer
	grpcServer            *grpc.Server
	registryClient        *registryclient.GrpcServerRegistry
	driverInfoServer      *driverinfo.DriverInfoServerEntity
}

// Method to start the asset link
func (d *AssetLink) Start(grpcServerAddress, registrationAddress, grpcRegistryAddress, httpServerAddress string) error {
	log.Info().
		Str("ID", d.metadata.AlId).
		Str("gRPC Address", grpcServerAddress).
		Str("grpcRegistryAddress", grpcRegistryAddress).
		Str("registrationName", registrationAddress).
		Msg("Starting asset link")

	// Webserver for observerability purposes
	if features.ObservabilityFeatures().HttpObservabilityServer {
		log.Info().Str("HTTP address", httpServerAddress).Msg("Starting RestAPI Observability Endpoint")

		s := webserver.NewServerWithParameters(httpServerAddress,
			metadata.Version{
				Version: d.metadata.Version.Version,
				Commit:  d.metadata.Version.Commit,
				Date:    d.metadata.Version.Date,
			})

		go s.Run()

	}
	// GRPC Server
	listener, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not bind server address")
	}

	// Register at the grpc server registry
	// Split into host and port. The registered endpoint is assembled by an dedicated flag and the grpc server endpoint.
	// Since, a grpc endpoint can also be ":8081" which listens on all ports, the endpoint needs to be explicitly set.
	_, portNumberString, err := net.SplitHostPort(grpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not determine port of gRPC server address")
	}

	d.registryClient = registryclient.New(grpcRegistryAddress, d.metadata.AlId, fmt.Sprintf("%s:%s", registrationAddress, portNumberString))
	d.registryClient.Register()

	// Start GRPC server
	d.grpcServer = grpc.NewServer()
	// CS Suite Drv Info
	log.Info().Msg("Registered Driver Info endpoint")
	d.driverInfoServer = &driverinfo.DriverInfoServerEntity{
		Metadata: d.metadata}
	generatedDriverInfoServer.RegisterDriverInfoApiServer(d.grpcServer, d.driverInfoServer)

	switch {
	// if a custom discovery server is provided, register it
	case d.customDiscoveryServer != nil:
		log.Info().Msg("Registered existing discovery server")
		generatedDiscoveryServer.RegisterDeviceDiscoverApiServer(d.grpcServer, d.customDiscoveryServer)

	// if a discovery implementation is provided, register it
	case d.discoveryImpl != nil:
		log.Info().
			Msg("Registered Discovery feature implementation")
		discoveryServer := &devicediscovery.DiscoverServerEntity{
			UnimplementedDeviceDiscoverApiServer: generatedDiscoveryServer.UnimplementedDeviceDiscoverApiServer{},
			Discovery:                            d.discoveryImpl,
		}
		generatedDiscoveryServer.RegisterDeviceDiscoverApiServer(d.grpcServer, discoveryServer)

	// if no discovery implementation is provided, log it
	default:
		log.Info().
			Msg("Discovery feature implementation not found")
	}

	log.Info().
		Str("address", grpcServerAddress).
		Msg("Serving gPRC Server")

	if err := d.grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("Could not bind server address")
	}
	return nil
}

func (d *AssetLink) Stop() {
	log.Info().Msg("Stop asset link")

	d.registryClient.Stop()
	d.grpcServer.Stop()

	log.Info().Msg("Asset Link stopped.")
	os.Exit(0)
}
