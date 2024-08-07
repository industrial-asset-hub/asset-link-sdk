/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package dcd

import (
	"fmt"
	"net"
	"os"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/internal/server/webserver"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/metadata"

	generatedDriverInfoServer "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/conn_suite_drv_info"
	generatedDiscoveryServer "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/iah-discovery"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/internal/features"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/internal/registryclient"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/internal/server/devicediscovery"
	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/internal/server/driverinfo"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Asset Link feature builder, according to the GoF build pattern
// The pattern provides methods to register new features in an easy
type dcdFeatureBuilder struct {
	metadata  metadata.Metadata
	discovery features.Discovery
}

// Methods to register new features
func (cb *dcdFeatureBuilder) Discovery(f features.Discovery) *dcdFeatureBuilder {
	cb.discovery = f
	return cb
}

// Builder
func New(metadata metadata.Metadata) *dcdFeatureBuilder {
	return &dcdFeatureBuilder{metadata: metadata}
}

func (cb *dcdFeatureBuilder) Build() *DCD {
	return &DCD{
		discoveryImpl: cb.discovery,
		metadata:      cb.metadata,
	}
}

// Structure of the features
type DCD struct {
	metadata         metadata.Metadata
	discoveryImpl    features.Discovery
	grpcServer       *grpc.Server
	registryClient   *registryclient.GrpcServerRegistry
	driverInfoServer *driverinfo.DriverInfoServerEntity
}

// Method to start the asset link
func (d *DCD) Start(grpcServerAddress, registrationAddress, grpcRegistryAddress, httpServerAddress string) error {
	log.Info().
		Str("Name", d.metadata.DcdId).
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

	d.registryClient = registryclient.New(grpcRegistryAddress, d.metadata.DcdId, fmt.Sprintf("%s:%s", registrationAddress, portNumberString))
	d.registryClient.Register()

	// Start GRPC server
	d.grpcServer = grpc.NewServer()
	// CS Suite Drv Info
	log.Info().Msg("Registered Driver Info endpoint")
	d.driverInfoServer = &driverinfo.DriverInfoServerEntity{
		Metadata: d.metadata}
	generatedDriverInfoServer.RegisterDriverInfoApiServer(d.grpcServer, d.driverInfoServer)

	// Select according to selected features
	if d.discoveryImpl == nil {
		log.Info().
			Msg("Discovery feature implementation not found")
	} else {
		log.Info().
			Msg("Registered Discovery feature implementation")
		discoveryServer := &devicediscovery.DiscoverServerEntity{
			UnimplementedDeviceDiscoverApiServer: generatedDiscoveryServer.UnimplementedDeviceDiscoverApiServer{},
			Discovery:                            d.discoveryImpl,
		}
		generatedDiscoveryServer.RegisterDeviceDiscoverApiServer(d.grpcServer, discoveryServer)
	}

	log.Info().
		Str("address", grpcServerAddress).
		Msg("Serving gPRC Server")

	if err := d.grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("Could not bind server address")
	}
	return nil
}

func (d *DCD) Stop() {
	log.Info().Msg("Stop asset link")

	d.registryClient.Stop()
	d.grpcServer.Stop()

	log.Info().Msg("Asset Link stopped.")
	os.Exit(0)
}
