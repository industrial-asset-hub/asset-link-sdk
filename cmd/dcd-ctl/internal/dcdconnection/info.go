/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package dcdconnection

import (
	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/conn_suite_drv_info"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

func GetInfo(endpoint string) string {
	log.Trace().Str("Endpoint", endpoint).Msg("Fetching health")

	conn := grpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewDriverInfoApiClient(conn)

	resp, err := client.GetVersionInfo(context.Background(), &generated.GetVersionInfoRequest{})

	if err != nil {
		log.Err(err).Msg("version request returned an error")
		return ""
	}
	var version = resp.GetVersion().String()

	log.Info().Str("Version", version).Msg("DCD version")
	return version
}
