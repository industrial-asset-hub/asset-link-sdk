/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package config

import (
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
)

type DiscoveryConfig interface {
	GetAllFilters() []*generated.ActiveFilter
	GetFilters(filterKey string) []*generated.ActiveFilter

	GetFilterSettingString(filterKey string, defaultValue string) (string, error)
	GetFilterSettingUint64(filterKey string, defaultValue uint64) (uint64, error)
	GetFilterSettingInt64(filterKey string, defaultValue int64) (int64, error)
	GetFilterSettingFloat64(filterKey string, defaultValue float64) (float64, error)

	GetAllOptions() []*generated.ActiveOption
	GetOptions(optionKey string) []*generated.ActiveOption

	GetOptionSettingString(filterKey string, defaultValue string) (string, error)
	GetOptionSettingUint64(filterKey string, defaultValue uint64) (uint64, error)
	GetOptionSettingInt64(filterKey string, defaultValue int64) (int64, error)
	GetOptionSettingFloat64(filterKey string, defaultValue float64) (float64, error)

	// GetTarget() []*generated.Destination

	GetDiscoveryRequest() *generated.DiscoverRequest

	String() string
	JSON() (string, error)
}
