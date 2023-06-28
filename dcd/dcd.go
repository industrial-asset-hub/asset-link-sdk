/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package dcd

import (
  generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/status"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/server/webserver"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"
  "net"
  "os"

  generatedDiscoveryServer "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/device_discovery"
  generatedFirmwareUpdate "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/firmware_update"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/features"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/registryclient"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/server/devicediscovery"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/server/firmwareupdate"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/server/status"

  "github.com/rs/zerolog/log"
  "google.golang.org/grpc"
)

var shouldStartHttpObservabilityServer = true

const httpServeraddress = ":8091"

var shouldRegisterGrpcObservabilityServer = true

// Device class driver feature builder, according to the GoF build pattern
// The pattern provides methods to register new features in an easy
type dcdFeatureBuilder struct {
  driverName     string
  version        metadata.Version
  discovery      features.Discovery
  softwareUpdate features.SoftwareUpdate
}

// Methods to register new features
func (cb *dcdFeatureBuilder) Discovery(f features.Discovery) *dcdFeatureBuilder {
  cb.discovery = f
  return cb
}

func (cb *dcdFeatureBuilder) SoftwareUpdate(f features.SoftwareUpdate) *dcdFeatureBuilder {
  cb.softwareUpdate = f
  return cb
}

// Builder
func New(driverName string, version metadata.Version) *dcdFeatureBuilder {
  return &dcdFeatureBuilder{driverName: driverName, version: version}
}

func (cb *dcdFeatureBuilder) Build() *DCD {
  return &DCD{
    discoveryImpl:      cb.discovery,
    softwareUpdateImpl: cb.softwareUpdate,
    name:               cb.driverName,
    version:            cb.version,
  }
}

// Structure of the features
type DCD struct {
  name               string
  version            metadata.Version
  discoveryImpl      features.Discovery
  softwareUpdateImpl features.SoftwareUpdate
  grpcServer         *grpc.Server
  registryClient     *registryclient.GrpcServerRegistry
  statusServer       *status.StatusServerEntity
}

// Method to start the device class driver
func (d *DCD) Start(grpcServerAddress string, grpcRegistryAddress string) error {
  log.Info().
    Str("Name", d.name).
    Str("gRPC Address", grpcServerAddress).
    Msg("Starting device class driver")

  // Webserver for observerability purposes
  if shouldStartHttpObservabilityServer {
    log.Info().Str("HTTP address", httpServeraddress).Msg("Starting RestAPI Observability Endpoint")

    s := webserver.NewServerWithParameters(httpServeraddress,
      metadata.Version{
        Version: d.version.Version,
        Commit:  d.version.Commit,
        Date:    d.version.Date,
      })

    go s.Run()

  }
  // GRPC Server
  listener, err := net.Listen("tcp", grpcServerAddress)
  if err != nil {
    log.Fatal().Err(err).Msg("Could not bind server address")
  }

  // Register at the grpc server registry
  d.registryClient = registryclient.New(grpcRegistryAddress, d.name, grpcServerAddress)
  d.registryClient.Register()

  // Start GRPC server
  d.grpcServer = grpc.NewServer()
  // ObservabilityServer
  if shouldRegisterGrpcObservabilityServer {
    log.Info().Msg("Registered gRPC Observability Endpoint")
    d.statusServer = &status.StatusServerEntity{Version: metadata.Version{
      Version: d.version.Version,
      Commit:  d.version.Commit,
      Date:    d.version.Date,
    }}
    generated.RegisterDcdStatusServer(d.grpcServer, d.statusServer)
  }

  // Select according to selected features
  if d.discoveryImpl == nil {
    log.Info().
      Msg("Discovery feature implementation not found")
  } else {
    log.Info().
      Msg("Registered Discovery feature implementation")
    discoveryServer := &devicediscovery.DiscoveryServerEntity{Discovery: d.discoveryImpl}
    generatedDiscoveryServer.RegisterDeviceDiscoveryApiServer(d.grpcServer, discoveryServer)
  }

  if d.softwareUpdateImpl == nil {
    log.Info().
      Msg("Software Update feature implementation not found")
  } else {
    log.Info().
      Msg("Registered feature Software Update feature implementation")
    firmwareUpdateServer := &firmwareupdate.FirmwareUpdateServerEntity{SoftwareUpdate: d.softwareUpdateImpl}
    generatedFirmwareUpdate.RegisterFirmwareupdateApiServer(d.grpcServer, firmwareUpdateServer)
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
  log.Info().Msg("Stop device class driver")

  d.registryClient.Stop()
  d.grpcServer.Stop()

  log.Info().Msg("Device class driver stopped.")
  os.Exit(0)
}
