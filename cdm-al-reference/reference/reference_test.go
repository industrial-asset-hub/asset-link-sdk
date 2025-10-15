/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package reference

import (
	"testing"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/publish"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDiscovery(t *testing.T) {
	simdevices.StartSimulatedDevices("") // start without visualization web server

	t.Run("discoverySucceeds", func(t *testing.T) {

		devicePublisher := &publish.DevicePublisherMock{}

		discoveryConfig := config.NewDiscoveryConfigWithDefaults()

		driver := &ReferenceAssetLink{}

		assert.NoError(t, driver.Discover(discoveryConfig, devicePublisher))
		assert.NotEmpty(t, devicePublisher.GetDevices())
	})

	t.Run("discoveryFails", func(t *testing.T) {
		devicePublisher := &publish.DevicePublisherMock{}

		err := status.Errorf(codes.Canceled, "Discovery was Canceled")
		devicePublisher.SetError(err)

		discoveryConfig := config.NewDiscoveryConfigWithDefaults()

		driver := &ReferenceAssetLink{}

		assert.Error(t, err, driver.Discover(discoveryConfig, devicePublisher))
		assert.Empty(t, devicePublisher.GetDevices())
	})

	t.Run("discoveryWithDeviceDetailsError", func(t *testing.T) {
		devicePublisher := &publish.DevicePublisherMock{}

		discoveryConfig := config.NewDiscoveryConfigWithDefaults()

		driver := &ReferenceAssetLink{}

		assert.NoError(t, driver.Discover(discoveryConfig, devicePublisher))

		existingErrors := devicePublisher.GetErrors()
		initialErrorCount := len(existingErrors)

		deviceDetailsError := &generated.DiscoverError{
			ResultCode:  int32(codes.Unavailable),
			Description: "Failed to retrieve device details for discovered device",
		}
		assert.NoError(t, devicePublisher.PublishError(deviceDetailsError))

		allErrors := devicePublisher.GetErrors()
		assert.NotEmpty(t, allErrors)
		assert.Equal(t, initialErrorCount+1, len(allErrors))

		foundTestError := false
		for _, err := range allErrors {
			if err.ResultCode == int32(codes.Unavailable) &&
				err.Description == "Failed to retrieve device details for discovered device" {
				foundTestError = true
				break
			}
		}
		assert.True(t, foundTestError, "Should find the test error we specifically published")
		assert.True(t, len(allErrors) > 0, "Should have device detail errors published")
	})
}

func TestConfig(t *testing.T) {
	t.Run("requestSupportedFilters", func(t *testing.T) {
		driver := &ReferenceAssetLink{}

		driver.GetSupportedFilters()
	})

	t.Run("requestSupportedOptions", func(t *testing.T) {
		driver := &ReferenceAssetLink{}

		driver.GetSupportedOptions()
	})
}
