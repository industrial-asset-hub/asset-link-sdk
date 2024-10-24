/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package registry

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/cmd/dcd-ctl/internal/shared"
	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/conn_suite_registry"
	"github.com/rs/zerolog/log"
)

func GetList(endpoint string) string {

	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewRegistryApiClient(conn)

	request := &generated.QueryRegisteredServicesRequest{}

	response, err := client.QueryRegisteredServices(context.Background(), request)
	if err != nil {
		log.Err(err).Msg("registry request returned an error")
		return ""
	} else {
		var stringBuilder strings.Builder

		for i, info := range response.GetInfos() {
			if i != 0 {
				stringBuilder.WriteString("---------------------------\n")
			}

			putStringValue(&stringBuilder, "App Instance ID", info.GetAppInstanceId())

			putStringValueArray(&stringBuilder, "App Types", info.GetAppTypes())
			putStringValueArray(&stringBuilder, "Driver Schema URIs", info.GetDriverSchemaUris())
			putStringValueArray(&stringBuilder, "Interfaces", info.GetInterfaces())

			putStringValue(&stringBuilder, "gRPC IPv4 Address", info.GetIpv4Address())
			putStringValue(&stringBuilder, "gRPC DNS Domain Name", info.GetDnsDomainname())
			putStringValue(&stringBuilder, "gRPC Port Number", strconv.Itoa(int(info.GetGrpcIpPortNumber())))
		}

		listText := stringBuilder.String()

		// log.Info().Str("List", listText).Msg("List of registered applications")
		fmt.Print(listText)

		return listText
	}

}

func putStringValueArray(stringBuilder *strings.Builder, label string, stringValueArray []string) {
	if len(stringValueArray) > 0 {
		stringBuilder.WriteString(label)
		stringBuilder.WriteString(":\n")
		for _, aType := range stringValueArray {
			stringBuilder.WriteString("  ")
			stringBuilder.WriteString(aType)
			stringBuilder.WriteString("\n")
		}
	}
}

func putStringValue(stringBuilder *strings.Builder, label string, stringValue string) {
	if len(stringValue) > 0 {
		stringBuilder.WriteString(label)
		stringBuilder.WriteString(":\n")
		stringBuilder.WriteString("  ")
		stringBuilder.WriteString(stringValue)
		stringBuilder.WriteString("\n")
	}
}
