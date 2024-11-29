/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package apimock

import (
	"fmt"

	"github.com/industrial-asset-hub/asset-link-sdk/v2/cmd/dcd-ctl/internal/dcd"
)

func TestStartDiscovery(address, filters, options string) bool {
	fmt.Println("Running test for StartDiscovery")
	data := dcd.StartDiscovery(address, options, filters)
	return data != nil
}

func TestGetFilterTypes(address, filters, options string) bool {
	fmt.Println("Running test for GetFilterTypes")
	data := dcd.GetFilterTypes(address)
	return data != nil
}

func TestGetFilterOptions(address, filters, options string) bool {
	fmt.Println("Running test for GetFilterOptions")
	data := dcd.GetFilterOptions(address)
	return data != nil
}
