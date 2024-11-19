/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package features

import (
	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v2/publish"
)

// This packages provides the interfaces which are needed for a custom asset link

// Interface Discovery provides the methods used the discovery feature
type Discovery interface {
	Discover(filters map[string]string, devicePublisher publish.DevicePublisher) error //TODO: why should we provide a string map here instead of the filters/options themselfes (like everywhere else)
	FilterTypes() []*generated.SupportedFilter
	FilterOptions() []*generated.SupportedOption
}
