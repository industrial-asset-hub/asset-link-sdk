/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"fmt"

	apimock "code.siemens.com/common-device-management/shared/cdm-dcd-sdk/v2/cmd/dcd-ctl/internal/api-mock-test"
)

func runTests(address, filters, options string) {
	fmt.Println("Running tests")
	apimock.RunApiMockTests(address, filters, options)
}
