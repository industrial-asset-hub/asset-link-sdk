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

func TestStartDiscovery(address, discoveryFile string) interface{} {
	fmt.Println("Running Test for StartDiscovery")
	data := al.Discover(address, discoveryFile)
	return data
}

func TestGetFilterTypes(address, discoveryFile string) interface{} {
	fmt.Println("Running Test for GetFilterTypes")
	data := al.GetFilterTypes(address)
	return data
}

func TestGetFilterOptions(address, discoveryFile string) interface{} {
	fmt.Println("Running Test for GetFilterOptions")
	data := al.GetFilterOptions(address)
	return data
}
