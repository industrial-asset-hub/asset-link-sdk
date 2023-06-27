/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package dcdconnection

import (
  generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/status"
  "github.com/rs/zerolog/log"
  "golang.org/x/net/context"
)

func GetHealth(endpoint string) string {
  log.Trace().Str("Endpoint", endpoint).Msg("Fetching status")

  conn := grpcConnection(endpoint)
  defer conn.Close()

  client := generated.NewDcdStatusClient(conn)

  resp, err := client.GetHealth(context.Background(), &generated.HealthRequest{})
  if err != nil {
    log.Err(err).Msg("Status request returned an error")
    return ""
  }
  var health = resp.GetHealth().String()
  log.Info().Str("Health", health).Msg("DCD health")
  return health
}

func GetVersion(endpoint string) (string, string, string) {
  log.Trace().Str("Endpoint", endpoint).Msg("Fetching health")

  conn := grpcConnection(endpoint)
  defer conn.Close()

  client := generated.NewDcdStatusClient(conn)

  resp, err := client.GetVersion(context.Background(), &generated.VersionRequest{})

  if err != nil {
    log.Err(err).Msg("Health request returned an error")
    return "", "", ""
  }
  var version = resp.GetVersion()
  var commit = resp.GetCommit()
  var date = resp.GetDate()

  log.Info().Str("Version", version).Str("Commit", commit).Str("Date", date).Msg("DCD version")
  return version, commit, date
}
