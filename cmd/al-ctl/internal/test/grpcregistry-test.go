/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type TestRegistryFunction func(string, string) bool

type TestRegistry struct {
	name     string
	function TestRegistryFunction
}

func RunRegistrationTests(grpcEndpoint, registryJsonPath string) {
	test := TestRegistry{"TestGetRegisteredServices", TestGetRegisteredServices}
	totalTests := 1
	testPassed := 0
	result := test.function(grpcEndpoint, registryJsonPath)
	if !result {
		log.Error().Str("test-name", test.name).Msg("test failed")
	} else {
		testPassed++
	}
	log.Info().Msgf("Total tests passed: %d/%d, failed: %d\n", testPassed, totalTests, totalTests-testPassed)
}

func TestGetRegisteredServices(grpcEndpoint, registryJsonPath string) bool {
	log.Info().Msg("Running test for GetRegisteredServices")
	var registryParams RegistryParams
	jsonData, err := os.ReadFile(registryJsonPath)
	if err != nil {
		log.Err(err).Msg("Error opening registry.json file: ")
		return false
	}

	err = json.Unmarshal(jsonData, &registryParams)
	if err != nil {
		log.Err(err).Msg("Error unmarshalling registry.json data: ")
		return false
	}
	data := GetRegisteredServices(grpcEndpoint)
	for _, service := range data {
		return validateRegistryParams(service, registryParams)
	}
	return false
}
