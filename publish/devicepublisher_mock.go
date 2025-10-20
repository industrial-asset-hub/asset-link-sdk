/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package publish

import (
	"sync"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
)

type DevicePublisherMock struct {
	deviceList []*generated.DiscoveredDevice
	errorList  []*generated.DiscoverError
	dataLock   sync.Mutex
	callError  error
}

func (d *DevicePublisherMock) PublishDevice(device *generated.DiscoveredDevice) error {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	if d.callError != nil {
		return d.callError
	}

	d.deviceList = append(d.deviceList, device)
	return nil
}

func (d *DevicePublisherMock) PublishDevices(devices []*generated.DiscoveredDevice) error {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	if d.callError != nil {
		return d.callError
	}

	d.deviceList = append(d.deviceList, devices...)
	return nil
}

func (d *DevicePublisherMock) PublishDeviceRelationship(realationship *generated.DiscoveredDevice) error {
	return d.PublishDevice(realationship)
}

func (d *DevicePublisherMock) PublishDeviceRelationships(realationships []*generated.DiscoveredDevice) error {
	return d.PublishDevices(realationships)
}

func (d *DevicePublisherMock) GetDevices() []*generated.DiscoveredDevice {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	return d.deviceList
}

func (d *DevicePublisherMock) ClearDevices() {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	d.deviceList = d.deviceList[:0]
}

func (d *DevicePublisherMock) PublishError(err *generated.DiscoverError) error {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	if d.callError != nil {
		return d.callError
	}

	d.errorList = append(d.errorList, err)
	return nil
}

func (d *DevicePublisherMock) PublishErrors(errors []*generated.DiscoverError) error {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	if d.callError != nil {
		return d.callError
	}

	d.errorList = append(d.errorList, errors...)
	return nil
}

func (d *DevicePublisherMock) GetErrors() []*generated.DiscoverError {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	return d.errorList
}

func (d *DevicePublisherMock) ClearErrors() {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	d.errorList = d.errorList[:0]
}

func (d *DevicePublisherMock) SetError(err error) {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	d.callError = err
}

func (d *DevicePublisherMock) GetError() error {
	d.dataLock.Lock()
	defer d.dataLock.Unlock()

	return d.callError
}
