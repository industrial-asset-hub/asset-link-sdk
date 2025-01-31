/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package test

import (
	"fmt"
	apimock "github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/api-mock-test"
)

func runTests(address, discoveryFile, assetJsonPath string, assetValidationRequired bool) (countOfAssetDiscovered int) {
	fmt.Println("Running tests")
	return apimock.RunApiMockTests(address, discoveryFile, assetJsonPath, assetValidationRequired)
}
