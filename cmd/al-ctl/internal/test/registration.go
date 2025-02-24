/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"encoding/json"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/generated/conn_suite_registry"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func RunRegistrationTests(grpcEndpoint, registryJsonPath string) {
	result := TestGetRegisteredServices(grpcEndpoint, registryJsonPath)
	if !result {
		log.Fatal().Msg("registration test failed")
	}
	log.Info().Msg("registration test passed")
}

func TestGetRegisteredServices(grpcEndpoint, registryJsonPath string) bool {
	log.Info().Msg("Running test for GetRegisteredServices")
	var registryParams RegistryParams
	jsonData, err := os.ReadFile(registryJsonPath)
	if err != nil {
		log.Error().Err(err).Msgf("Error opening file at %s", registryJsonPath)
		return false
	}

	err = json.Unmarshal(jsonData, &registryParams)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling registration params")
		return false
	}
	registeredServiceInfos := GetRegisteredServices(grpcEndpoint)
	if len(registeredServiceInfos) == 0 {
		log.Error().Msg("No registered services found")
		return false
	}
	AssetLinkFound := false
	testPassed := true
	for _, serviceInfo := range registeredServiceInfos {
		if shared.AssetLinkEndpoint == getServiceAddress(serviceInfo) {
			AssetLinkFound = true
			unmatchedRegistrationParams := validateRegistryParams(serviceInfo, registryParams)
			if len(unmatchedRegistrationParams) > 0 {
				testPassed = false
				log.Error().Msgf("Unmatched registration parameters %s for service-info %s", unmatchedRegistrationParams, serviceInfo)
			}
		}
	}
	if !AssetLinkFound {
		log.Error().Msgf("assetLink was not found in registered services")
		testPassed = false
	}
	return testPassed
}

func getServiceAddress(serviceInfo *conn_suite_registry.ServiceInfo) string {
	var ip string
	if serviceInfo.GetDnsDomainname() != "" {
		ip = serviceInfo.GetDnsDomainname()
	} else if serviceInfo.GetIpv4Address() != "" {
		ip = serviceInfo.GetIpv4Address()
	}
	grpcAddressOfService := ip + ":" + strconv.Itoa(int(serviceInfo.GetGrpcIpPortNumber()))
	return grpcAddressOfService
}
