/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package testsuite

import (
	"encoding/json"
	"fmt"
	"os"
)

func ValidateAssetLinkRegistration(params RegistryValidationParams) error {
	var registryParams RegistryFileParams
	jsonData, err := os.ReadFile(params.RegistryJsonPath)
	if err != nil {
		return fmt.Errorf("error opening file at %s", params.RegistryJsonPath)
	}

	err = json.Unmarshal(jsonData, &registryParams)
	if err != nil {
		return fmt.Errorf("error unmarshalling registration params")
	}
	registeredServiceInfos := GetRegisteredServices(params.GrpcEndpoint)
	if len(registeredServiceInfos) == 0 {
		return fmt.Errorf("no registered services found")
	}
	AssetLinkFound := false
	for _, serviceInfo := range registeredServiceInfos {
		if params.AssetLinkEndpoint == getServiceAddress(serviceInfo) {
			AssetLinkFound = true
			unmatchedRegistrationParams := validateRegistryParams(serviceInfo, registryParams)
			if len(unmatchedRegistrationParams) > 0 {
				return fmt.Errorf("Unmatched registration parameters %s for service-info %s", unmatchedRegistrationParams, serviceInfo)
			}
		}
	}
	if !AssetLinkFound {
		return fmt.Errorf("assetLink was not found in registered services")
	}
	return nil
}
