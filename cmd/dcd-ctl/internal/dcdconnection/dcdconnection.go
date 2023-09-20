/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package dcdconnection

import (
  "github.com/rs/zerolog/log"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
)

func grpcConnection(endpoint string) *grpc.ClientConn {
  conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
  if err != nil {
    log.Err(err).Msg("can not connect with server")
    return nil
  }

  return conn
}
