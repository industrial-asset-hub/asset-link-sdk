/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package apimock

import (
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/al"
)

func TestStartDiscovery(address, discoveryFile string) bool {
	fmt.Println("Running test for StartDiscovery")
	data := al.Discover(address, discoveryFile)
	return data != nil
}

func TestGetFilterTypes(address, discoveryFile string) bool {
	fmt.Println("Running test for GetFilterTypes")
	data := al.GetFilterTypes(address)
	return data != nil
}

func TestGetFilterOptions(address, discoveryFile string) bool {
	fmt.Println("Running test for GetFilterOptions")
	data := al.GetFilterOptions(address)
	return data != nil
}
