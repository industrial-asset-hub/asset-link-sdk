/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

import (
	"strings"
	"testing"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/config"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/publish"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestDiscovery(t *testing.T) {
	t.Run("discoverySucceeds", func(t *testing.T) {

		devicePublisher := &publish.DevicePublisherMock{}

		discoveryConfig := config.NewDiscoveryConfigWithDefaults()

		driver := &AssetLinkImplementation{}

		assert.NoError(t, driver.Discover(discoveryConfig, devicePublisher))

		devices := devicePublisher.GetDevices()
		assert.Len(t, devices, 1)

		devicePayload, err := protojson.Marshal(devices[0])
		assert.NoError(t, err)
		assert.True(t, strings.Contains(string(devicePayload), "00:16:3e:01:02:03"), "discovered device should contain MAC-based identifier")
	})

	t.Run("discoveryFails", func(t *testing.T) {
		devicePublisher := &publish.DevicePublisherMock{}

		err := status.Errorf(codes.Canceled, "Discovery was canceled")
		devicePublisher.SetError(err)

		discoveryConfig := config.NewDiscoveryConfigWithDefaults()

		driver := &AssetLinkImplementation{}

		assert.Error(t, err, driver.Discover(discoveryConfig, devicePublisher))
		assert.Empty(t, devicePublisher.GetDevices())
	})
}

func TestConfig(t *testing.T) {
	t.Run("requestSupportedFilters", func(t *testing.T) {
		driver := &AssetLinkImplementation{}

		driver.GetSupportedFilters()
	})

	t.Run("requestSupportedOptions", func(t *testing.T) {
		driver := &AssetLinkImplementation{}

		driver.GetSupportedOptions()
	})
}
