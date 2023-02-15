/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package publish

import (
	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"
)

type DevicePublisher interface {
	PublishDevice(device *generated.DiscoveredDevice) error
	PublishDevices(devices []*generated.DiscoveredDevice) error
}
