/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package features

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
)

// This packages provides the interfaces which are needed for a custom asset link

// Interface Discovery provides the methods used the discovery feature
type Discovery interface {
	Discover(discoveryConfig config.DiscoveryConfig, devicePublisher publish.DevicePublisher) error
	GetSupportedFilters() []*generated.SupportedFilter
	GetSupportedOptions() []*generated.SupportedOption
}
