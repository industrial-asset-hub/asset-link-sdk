/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package features

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v4/config"
	deviceinfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/publish"
)

// This packages provides the interfaces which are needed for a custom asset link

// Interface Discovery provides the methods used the discovery feature
type Discovery interface {
	Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error
	GetSupportedFilters() []*generated.SupportedFilter
	GetSupportedOptions() []*generated.SupportedOption
}

// Interface DeviceInfo provides the methods used by the DeviceInfo feature
type DeviceInfo interface {
	GetPropertyValues(request *deviceinfo.GetPropertyValuesRequest) (*deviceinfo.GetPropertyValuesResponse, error)
	GetSupportedProperties(request *deviceinfo.GetSupportedPropertiesRequest) (*deviceinfo.GetSupportedPropertiesResponse, error)
}
