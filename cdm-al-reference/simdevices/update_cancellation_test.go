// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package simdevices

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCancelSimulatedFirmwareInstall(t *testing.T) {
	device := newSimulatedDevice("eth0", "test", "", "192.168.0.10", "test", nil, nil)
	firmware, err := device.PrepareFirmware(getFilenamePath("misc/simulated_device_firmware_2.0.0.fwu"))
	require.NoError(t, err)

	checkFirmwareStepCancellation(t, device, StateUpdating, func(ctx context.Context) error {
		return device.installFirmware(ctx, firmware, time.Hour)
	})
	require.Equal(t, "1.0.0", device.GetInstalledFirmwareVersion())
	require.Equal(t, "1.0.0", device.GetActiveFirmwareVersion())
}

func TestCancelSimulatedFirmwareReboot(t *testing.T) {
	device := newSimulatedDevice("eth0", "test", "", "192.168.0.10", "test", nil, nil)
	require.NoError(t, device.UpdateFirmware(getFilenamePath("misc/simulated_device_firmware_2.0.0.fwu")))

	checkFirmwareStepCancellation(t, device, StateBooting, func(ctx context.Context) error {
		return device.rebootDevice(ctx, time.Hour)
	})
	require.Equal(t, "2.0.0", device.GetInstalledFirmwareVersion(), "cancellation does not undo installation")
	require.Equal(t, "1.0.0", device.GetActiveFirmwareVersion())
}

func checkFirmwareStepCancellation(t *testing.T, device *simulatedDeviceInfo, state SimulatedDeviceState, step func(context.Context) error) {
	t.Helper()
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		done := make(chan error, 1)
		go func() { done <- step(ctx) }()

		synctest.Wait()
		require.Equal(t, state, device.DeviceState, "operation must already be in progress")
		startedCancelling := time.Now()
		cancel()
		require.ErrorIs(t, <-done, context.Canceled)
		require.Zero(t, time.Since(startedCancelling), "must stop without waiting for the step timer")
		require.Equal(t, StateActive, device.DeviceState)
	})
}
