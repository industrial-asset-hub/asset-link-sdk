/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package features

import generated "code.siemens.com/common-device-management/shared/cdm-dcd-sdk/v2/generated/iah-discovery"

// This packages provides the interfaces which are needed for a custom asset link

// Interface Discovery provides the methods used the discovery feature
type Discovery interface {
	Start(jobId uint32, deviceChannel chan []*generated.DiscoveredDevice, err chan error, filters map[string]string)
	FilterTypes(filterTypesChannel chan []*generated.SupportedFilter)
	FilterOptions(filterOptionsChannel chan []*generated.SupportedOption)
}
