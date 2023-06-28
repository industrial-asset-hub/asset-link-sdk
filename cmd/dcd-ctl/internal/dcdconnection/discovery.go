/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package dcdconnection

import (
  generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/device_discovery"
  "github.com/rs/zerolog/log"
  "golang.org/x/net/context"
  "io"
)

func StartDiscovery(endpoint string) {
  log.Trace().Str("Endpoint", endpoint).Msg("Starting discovery job")

  conn := grpcConnection(endpoint)
  defer conn.Close()

  client := generated.NewDeviceDiscoveryApiClient(conn)

  resp, err := client.StartDeviceDiscovery(context.Background(), &generated.DiscoveryRequest{
    Filters: []*generated.DiscoveryFilter{{Key: ""}}})

  if err != nil {
    log.Err(err).Msg("StartDeviceDiscovery request returned an error")
    return
  }

  log.Info().Str("response", resp.String()).Msg("Recevied response")

}

func Subscribe(endpoint string, jobid uint32) []map[string]interface{} {
  log.Trace().Str("Endpoint", endpoint).Msg("Subscribing to discovery job")

  conn := grpcConnection(endpoint)
  defer conn.Close()

  client := generated.NewDeviceDiscoveryApiClient(conn)

  stream, err := client.SubscribeDiscoveryResults(context.Background(), &generated.DiscoveryResultsRequest{DiscoveryId: jobid})
  if err != nil {
    log.Err(err).Msg("open stream error")
    return nil
  }

  devices := make([]map[string]interface{}, 0)
  for {
    resp, err := stream.Recv()

    if err == io.EOF {
      break
    }

    if err != nil {
      log.Err(err).Msg("SubscribeDiscovery request returned an error")
      return nil
    }
    log.Info().Interface("response", resp.Devices).Msg("Received Response")

    for _, d := range resp.Devices {
      log.Trace().Interface("Devices", d).Msg("")
      devices = append(devices, d.Parameters.AsMap())
    }
  }

  return devices
}

func StopDiscovery(endpoint string, jobid uint32) {
  log.Trace().Str("Endpoint", endpoint).Msg("Stop discovery job")

  conn := grpcConnection(endpoint)
  defer conn.Close()

  client := generated.NewDeviceDiscoveryApiClient(conn)

  resp, err := client.StopDeviceDiscovery(context.Background(), &generated.StopDiscoveryRequest{
    DiscoveryId: jobid})

  if err != nil {
    log.Err(err).Msg("StopDeviceDiscovery request returned an error")
    return
  }
  log.Info().Str("response", resp.String()).Msg("Received response")

}
