/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package shared

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	RegistryEndpoint        string
	AssetLinkEndpoint       string
	AssetJsonPath           string
	AssetValidationRequired bool
)

const (
	DiscoveryFileDesc string = "discovery file allows the configuration of discovery filters and options (see discovery.json for an example)"
)

func GrpcConnection(endpoint string) *grpc.ClientConn {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Err(err).Msg("can not connect with server")
		return nil
	}

	return conn
}
