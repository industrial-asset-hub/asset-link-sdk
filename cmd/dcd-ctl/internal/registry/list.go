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

	"github.com/industrial-asset-hub/asset-link-sdk/v2/cmd/dcd-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/conn_suite_registry"
	"github.com/rs/zerolog/log"
)

func PrintList(endpoint string) {
	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewRegistryApiClient(conn)

	request := &generated.QueryRegisteredServicesRequest{}

	response, err := client.QueryRegisteredServices(context.Background(), request)
	if err != nil {
		log.Err(err).Msg("registry request returned an error")
		return
	} else {
		var stringBuilder strings.Builder

		for _, info := range response.GetInfos() {
			println("Asset Link:")

			printStringValue("App Instance ID", info.GetAppInstanceId())

			printStringValueArray("App Types", info.GetAppTypes())
			printStringValueArray("Driver Schema URIs", info.GetDriverSchemaUris())
			printStringValueArray("Interfaces", info.GetInterfaces())

			printStringValue("gRPC IPv4 Address", info.GetIpv4Address())
			printStringValue("gRPC DNS Domain Name", info.GetDnsDomainname())
			printStringValue("gRPC Port Number", strconv.Itoa(int(info.GetGrpcIpPortNumber())))
		}

		listText := stringBuilder.String()

		fmt.Print(listText)
	}
}

func printStringValueArray(label string, stringValueArray []string) {
	num := len(stringValueArray)
	if num > 0 {
		print("  " + label + ": ")
		num := len(stringValueArray)
		for n, aType := range stringValueArray {
			if n < num-1 {
				print(aType + ", ")
			} else {
				println(aType)
			}
		}
	}
}

func printStringValue(label string, stringValue string) {
	if len(stringValue) > 0 {
		println("  " + label + ": " + stringValue)
	}
}
