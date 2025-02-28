/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"context"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/conn_suite_registry"
	"github.com/rs/zerolog/log"
)

func GetRegisteredServices(endpoint string) []*generated.ServiceInfo {
	conn := shared.GrpcConnection(endpoint)
	defer conn.Close()

	client := generated.NewRegistryApiClient(conn)

	request := &generated.QueryRegisteredServicesRequest{}
	response, err := client.QueryRegisteredServices(context.Background(), request)

	if err != nil {
		log.Err(err).Msg("registry request returned an error")
		return nil
	}
	return response.GetInfos()
}

func validateRegistryParams(service *generated.ServiceInfo, registryParams RegistryParams) []string {
	var unmatchedRegistrationParams []string
	if service.GetAppInstanceId() != registryParams.AppInstanceId {
		unmatchedRegistrationParams = append(unmatchedRegistrationParams, "app_instance_id")
	}
	if !compareRegistryInfoValues(service.GetAppTypes(), registryParams.AppTypes) {
		unmatchedRegistrationParams = append(unmatchedRegistrationParams, "app_types")
	}
	if !compareRegistryInfoValues(service.GetDriverSchemaUris(), registryParams.DeviceSchemaUri) {
		unmatchedRegistrationParams = append(unmatchedRegistrationParams, "device_schema_uri")
	}
	return unmatchedRegistrationParams
}

func compareRegistryInfoValues(registryInfo, registryParam []string) bool {
	if len(registryInfo) != len(registryParam) {
		return false
	}
	for i := range registryInfo {
		if registryInfo[i] != registryParam[i] {
			return false
		}
	}
	return true
}
