/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
  generated "code.siemens.com/common-device-management/utils/go-modules/discovery.git/pkg/device"
  "context"
  "fmt"
  "github.com/rs/zerolog/log"
  "github.com/spf13/cobra"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "io"
)

// discoveryCmd represents the discovery command
var discoveryCmd = &cobra.Command{
  Use:   "discovery",
  Short: "Use discovery feature of an DCD",
  Long: `This command allows to start/stop and receive the results of an
discovery job.`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("discovery called")
    do()
  },
}

func do() {
  grpcEndpoint := "localhost:8081"
  conn, err := grpc.Dial(grpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
  if err != nil {
    log.Fatal().Err(err).Msg("can not connect with server")
  }
  defer conn.Close()

  client := generated.NewDeviceDiscoveryApiClient(conn)

  _, err = client.StartDeviceDiscovery(context.Background(), &generated.DiscoveryRequest{
    Filters: []*generated.DiscoveryFilter{{Key: ""}}})
  const test uint32 = 12
  if err != nil {
    log.Fatal().Err(err).Msg("StartDeviceDiscovery request returned an error")
  }
  defer conn.Close()

  stream, err := client.SubscribeDiscoveryResults(context.Background(), &generated.DiscoveryResultsRequest{DiscoveryId: 1})
  if err != nil {
    log.Err(err).Msg("open stream error")
  }

  // Inline function to progress
  var progressResponse = func(devices []*generated.DiscoveryDevice) {
    log.Trace().Interface("response", devices).Msg("Progress Devices")
    for _, d := range devices {
      log.Trace().Interface("Devices", d).Msg("")

      for _, k := range []string{"ManagementState", "test"} {

        if d.Parameters != nil {
          _, ok := d.Parameters.Fields[k]
          if !ok {
            log.Warn().Str("Key", k).Msg("Key not found")
          } else {
            log.Trace().Str("Key", k).Msg("Key found")
          }
        } else {
          log.Info().Msg("received asset description empty")

        }
      }
    }
  }

  for {
    resp, err := stream.Recv()
    if err == io.EOF {
      return
    }
    if err != nil {
      log.Err(err).Msg("StartDeviceDiscovery request returned an error")
    }
    progressResponse(resp.Devices)
    log.Trace().Interface("response", resp.Devices).Msg("Received Response")
  }
}

func init() {
  rootCmd.AddCommand(discoveryCmd)
}
