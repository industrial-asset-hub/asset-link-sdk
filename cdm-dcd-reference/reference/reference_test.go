/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package reference

import (
	"testing"

	"github.com/industrial-asset-hub/asset-link-sdk/v2/publish"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDiscovery(t *testing.T) {
	t.Run("discoverySucceeds", func(t *testing.T) {

		devicePublisher := &publish.DevicePublisherMock{}

		var filters map[string]string

		driver := &ReferenceClassDriver{}

		assert.NoError(t, driver.Discover(filters, devicePublisher))
		assert.NotEmpty(t, devicePublisher.GetDevices())
	})

	t.Run("discoveryFails", func(t *testing.T) {
		devicePublisher := &publish.DevicePublisherMock{}

		err := status.Errorf(codes.Canceled, "Discovery was Canceled")
		devicePublisher.SetError(err)

		var filters map[string]string

		driver := &ReferenceClassDriver{}

		assert.Error(t, err, driver.Discover(filters, devicePublisher))
		assert.Empty(t, devicePublisher.GetDevices())
	})
}

func TestOptions(t *testing.T) {
	t.Run("requestFilterOptions", func(t *testing.T) {
		driver := &ReferenceClassDriver{}

		driver.FilterOptions()
	})

	t.Run("requestOptionTypes", func(t *testing.T) {
		driver := &ReferenceClassDriver{}

		driver.FilterTypes()
	})
}
