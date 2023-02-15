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

func TestStartDiscovery(address, discoveryFile string) bool {
	fmt.Println("Running test for StartDiscovery")
	data := dcd.Discover(address, discoveryFile)
	return data != nil
}

func TestGetFilterTypes(address, discoveryFile string) bool {
	fmt.Println("Running test for GetFilterTypes")
	data := dcd.GetFilterTypes(address)
	return data != nil
}

func TestGetFilterOptions(address, discoveryFile string) bool {
	fmt.Println("Running test for GetFilterOptions")
	data := dcd.GetFilterOptions(address)
	return data != nil
}
