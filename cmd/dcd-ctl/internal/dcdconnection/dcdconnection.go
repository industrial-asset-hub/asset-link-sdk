package dcdconnection

import (
  generated "code.siemens.com/common-device-management/utils/go-modules/discovery.git/pkg/device"
  "github.com/rs/zerolog/log"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "io"
)

func grpcConnection(endpoint string) *grpc.ClientConn {
  conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
  if err != nil {
    log.Err(err).Msg("can not connect with server")
    return nil
  }

  return conn
}

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

func Subscribe(endpoint string) []map[string]interface{} {
  log.Trace().Str("Endpoint", endpoint).Msg("Subscribing to discovery job")

  conn := grpcConnection(endpoint)
  defer conn.Close()

  client := generated.NewDeviceDiscoveryApiClient(conn)

  stream, err := client.SubscribeDiscoveryResults(context.Background(), &generated.DiscoveryResultsRequest{DiscoveryId: 1})
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

func StopDiscovery(endpoint string) {
  log.Trace().Str("Endpoint", endpoint).Msg("Stop discovery job")

  conn := grpcConnection(endpoint)
  defer conn.Close()

  client := generated.NewDeviceDiscoveryApiClient(conn)

  resp, err := client.StopDeviceDiscovery(context.Background(), &generated.StopDiscoveryRequest{
    DiscoveryId: 0})

  if err != nil {
    log.Err(err).Msg("StopDeviceDiscovery request returned an error")
    return
  }
  log.Info().Str("response", resp.String()).Msg("Received response")

}
