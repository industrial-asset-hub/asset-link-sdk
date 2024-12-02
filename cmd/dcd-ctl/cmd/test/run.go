/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"fmt"

	apimock "github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/dcd-ctl/internal/api-mock-test"
)

func runTests(address, discoveryFile string) {
	fmt.Println("Running tests")
	apimock.RunApiMockTests(address, discoveryFile)
}
