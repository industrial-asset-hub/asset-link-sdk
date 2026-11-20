// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package reference

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/contracttest"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/oci"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cdm-al-reference/simdevices"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/config"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	discovery "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/publish"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

var updateDeviceAddress = simdevices.SimulatedDeviceAddress{AssetLinkNIC: "eth0", DeviceIP: "192.168.0.10", SubDeviceID: -1}

const (
	updateSucceeded = generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED
	updateFailed    = generated.UpdateOutcome_UPDATE_OUTCOME_FAILED
)

func TestReferenceArtifactUpdate(t *testing.T) {
	simdevices.StartSimulatedDevices("")
	driver := new(ReferenceAssetLink)
	device, err := simdevices.ConnectToDevice(updateDeviceAddress, nil)
	require.NoError(t, err)
	request := updateRequest(t)
	request.ConnectionParameterSet = discoverUpdateTarget(t, driver)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateSucceeded)
	contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateSucceeded)
	require.Equal(t, "1.0.0", device.GetInstalledFirmwareVersion())
	require.Equal(t, "1.0.0", device.GetActiveFirmwareVersion())

	conflict := proto.CloneOf(request)
	conflict.OciReference.Repository = "different-artifact"
	contracttest.Check(t, ctx, conflict, driver.HandlePrepareUpdate, updateFailed)
	contracttest.Check(t, ctx, conflict, driver.HandleActivateUpdate, updateFailed)
	contracttest.Check(t, ctx, conflict, driver.HandleCancelUpdate, updateFailed)

	contracttest.Check(t, ctx, request, driver.HandleCancelUpdate, updateSucceeded)
	contracttest.Check(t, ctx, request, driver.HandleActivateUpdate, updateFailed)
	require.Equal(t, "1.0.0", device.GetInstalledFirmwareVersion())
	require.Equal(t, "1.0.0", device.GetActiveFirmwareVersion())

	contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateSucceeded)
	contracttest.Check(t, ctx, request, driver.HandleActivateUpdate, updateSucceeded)
	require.Equal(t, "2.0.0", device.GetInstalledFirmwareVersion())
	require.Equal(t, "2.0.0", device.GetActiveFirmwareVersion())
	contracttest.Check(t, ctx, request, driver.HandleActivateUpdate, updateFailed)
}

func TestReferenceUpdateRejectsWrongTarget(t *testing.T) {
	simdevices.StartSimulatedDevices("")
	driver := new(ReferenceAssetLink)
	base := updateRequest(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	for _, tc := range []struct {
		name   string
		params string
		reason artifact.FailureReason
	}{
		{"wrong JSON shape", `[]`, artifact.InvalidRequest},
		{"missing address", `{}`, artifact.InvalidRequest},
		{"unknown device", `{"alNic":"eth0","ipAddress":"192.168.0.99","subDeviceID":-1}`, artifact.TargetNotFound},
		{"protected device", `{"alNic":"eth0","ipAddress":"192.168.0.11","subDeviceID":-1}`, artifact.DeviceRejected},
		{"unsupported subdevice", `{"alNic":"eth1","ipAddress":"192.168.1.10","subDeviceID":0}`, artifact.NotSupported},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := proto.CloneOf(base)
			request.ConnectionParameterSet.ParameterJson = tc.params
			result := contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateFailed)
			require.Equal(t, tc.reason, result.GetFailure().GetReason())
		})
	}
}

func TestReferenceUpdateInterruptedBeforeReboot(t *testing.T) {
	simdevices.StartSimulatedDevices("")
	driver := new(ReferenceAssetLink)
	device, err := simdevices.ConnectToDevice(updateDeviceAddress, nil)
	require.NoError(t, err)
	request := updateRequest(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateSucceeded)

	interrupted, interrupt := context.WithCancel(ctx)
	defer interrupt()
	result, err := driver.HandleActivateUpdate(request, rebootInterrupt{ctx: interrupted, cancel: interrupt})
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN, result.GetOutcome())
	require.Equal(t, "2.0.0", device.GetInstalledFirmwareVersion())
	require.Equal(t, "1.0.0", device.GetActiveFirmwareVersion())

	contracttest.Check(t, ctx, request, driver.HandleActivateUpdate, updateSucceeded)
	require.Equal(t, "2.0.0", device.GetActiveFirmwareVersion())
}

func TestReferenceUpdateWithDeviceConnectionLoss(t *testing.T) {
	simdevices.StartSimulatedDevices("")
	driver := new(ReferenceAssetLink)
	device, err := simdevices.ConnectToDevice(updateDeviceAddress, nil)
	require.NoError(t, err)
	request := updateRequest(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateSucceeded)
	firmware := driver.prepared.firmware

	rejected := func(_ *artifact.ArtifactMetaData, reporter artifact.ProgressReporter) (*artifact.Result, error) {
		return activateFirmware(device, simdevices.PreparedFirmware{}, reporter)
	}
	result := contracttest.Check(t, ctx, request, rejected, updateFailed)
	require.Equal(t, artifact.DeviceRejected, result.GetFailure().GetReason())

	disconnected := func(_ *artifact.ArtifactMetaData, reporter artifact.ProgressReporter) (*artifact.Result, error) {
		return activateFirmware(losesInstallReply{device}, firmware, reporter)
	}
	result = contracttest.Check(t, ctx, request, disconnected, generated.UpdateOutcome_UPDATE_OUTCOME_UNKNOWN)
	require.Equal(t, artifact.OutcomeNotConfirmed, result.GetFailure().GetReason())
	require.NoError(t, ctx.Err(), "device loss must not cancel the gateway stream")
	require.Equal(t, "2.0.0", device.GetInstalledFirmwareVersion())
	require.Equal(t, "1.0.0", device.GetActiveFirmwareVersion(), "must not reboot after losing the device connection")

	contracttest.Check(t, ctx, request, driver.HandleActivateUpdate, updateSucceeded)
	require.Equal(t, "2.0.0", device.GetActiveFirmwareVersion())
}

// Let the real simulator install the firmware, then lose the device reply
// before the Asset Link can confirm the result. The gateway stream stays open.
type losesInstallReply struct{ simdevices.SimulatedDevice }

func (d losesInstallReply) InstallFirmware(ctx context.Context, firmware simdevices.PreparedFirmware) error {
	if err := d.SimulatedDevice.InstallFirmware(ctx, firmware); err != nil {
		return err
	}
	return errors.New("simulated device connection lost before install reply")
}

type rebootInterrupt struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (r rebootInterrupt) Context() context.Context            { return r.ctx }
func (r rebootInterrupt) ReportProgress(string, uint32) error { return r.ctx.Err() }

func (r rebootInterrupt) ReportIndeterminate(step string) error {
	if step == "reboot" {
		r.cancel()
	}
	return r.ctx.Err()
}

func TestReferenceUpdateRejectsBadArtifact(t *testing.T) {
	simdevices.StartSimulatedDevices("")
	driver := new(ReferenceAssetLink)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	request := updateRequest(t)
	wrongRepo := proto.CloneOf(request)
	wrongRepo.OciReference.Repository = "missing"
	result := contracttest.Check(t, ctx, wrongRepo, driver.HandlePrepareUpdate, updateFailed)
	require.Equal(t, artifact.ArtifactUnavailable, result.GetFailure().GetReason())

	request.OciReference = serveUpdateArtifact(t, []byte(`{"artefact_type":"configuration"}`))
	result = contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateFailed)
	require.Equal(t, artifact.ArtifactInvalid, result.GetFailure().GetReason())
	require.Nil(t, driver.prepared)
}

func discoverUpdateTarget(t *testing.T, driver *ReferenceAssetLink) *discovery.ConnectionParameterSet {
	t.Helper()
	publisher := new(publish.DevicePublisherMock)
	require.NoError(t, driver.Discover(config.NewDiscoveryConfigWithDefaults(), publisher))
	want := simulatedConnectionParameters(updateDeviceAddress)
	for _, device := range publisher.GetDevices() {
		if proto.Equal(want, device.GetConnectionParameterSet()) {
			return device.GetConnectionParameterSet()
		}
	}
	t.Fatal("target not published by discovery")
	return nil
}

func updateRequest(t *testing.T) *artifact.ArtifactMetaData {
	t.Helper()
	firmware, err := os.ReadFile("../../misc/simulated_device_firmware_2.0.0.fwu")
	require.NoError(t, err)
	return &artifact.ArtifactMetaData{
		JobId: "update-1", ConnectionParameterSet: simulatedConnectionParameters(updateDeviceAddress),
		OciReference: serveUpdateArtifact(t, firmware),
	}
}

func serveUpdateArtifact(t *testing.T, firmware []byte) *generated.OciArtifactReference {
	t.Helper()
	layer := oci.Descriptor{Digest: fmt.Sprintf("sha256:%x", sha256.Sum256(firmware)), Size: int64(len(firmware))}
	manifest, err := json.Marshal(oci.Manifest{
		SchemaVersion: 2, MediaType: oci.ManifestMediaType,
		Config: oci.Descriptor{Digest: fmt.Sprintf("sha256:%x", sha256.Sum256(nil))}, Layers: []oci.Descriptor{layer},
	})
	require.NoError(t, err)
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/firmware/manifests/"+digest, func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(manifest) })
	mux.HandleFunc("/v2/firmware/blobs/"+layer.Digest, func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(firmware) })
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return &generated.OciArtifactReference{RegistryUrl: server.URL, Repository: "firmware", ManifestDigest: digest}
}
