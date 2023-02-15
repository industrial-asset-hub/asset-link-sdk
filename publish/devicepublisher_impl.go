/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package publish

import (
	"sync"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"
)

type DevicePublisherImplementation struct {
	Stream     generated.DeviceDiscoverApi_DiscoverDevicesServer
	streamLock sync.Mutex
}

func (d *DevicePublisherImplementation) PublishDevice(device *generated.DiscoveredDevice) error {
	devices := make([]*generated.DiscoveredDevice, 0)
	devices = append(devices, device)
	return d.PublishDevices(devices)
}

func (d *DevicePublisherImplementation) PublishDevices(devices []*generated.DiscoveredDevice) error {
	response := new(generated.DiscoverResponse)
	response.Devices = devices
	d.streamLock.Lock()
	defer d.streamLock.Unlock()
	return d.Stream.SendMsg(response)
}
