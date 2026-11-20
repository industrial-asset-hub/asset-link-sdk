// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT

package reference

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/oci"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cdm-al-reference/simdevices"
	"github.com/rs/zerolog/log"
)

const simulatedAddressSchema = "urn:siemens:cdm:al:reference:simulated-address:v1"

// Match each phase against the same job, device and immutable OCI reference.
type updateTarget struct {
	job       string
	address   simdevices.SimulatedDeviceAddress
	reference oci.Reference
}

type preparedUpdate struct {
	target   updateTarget
	firmware simdevices.PreparedFirmware
}

func parseUpdateTarget(metadata *artifact.ArtifactMetaData) (updateTarget, error) {
	if err := artifact.ValidateMetaData(metadata); err != nil {
		return updateTarget{}, err
	}
	parameters := metadata.GetConnectionParameterSet()
	if parameters.GetSchemaUri() != simulatedAddressSchema {
		return updateTarget{}, errors.New("unexpected simulated address schema")
	}
	var address simdevices.SimulatedDeviceAddress
	if err := json.Unmarshal([]byte(parameters.GetParameterJson()), &address); err != nil {
		return updateTarget{}, err
	}
	if address.AssetLinkNIC == "" || address.DeviceIP == "" || address.SubDeviceID < -1 {
		return updateTarget{}, errors.New("incomplete simulated device address")
	}
	ref := metadata.GetOciReference()
	return updateTarget{
		job: metadata.GetJobId(), address: address,
		reference: oci.Reference{RegistryURL: ref.GetRegistryUrl(), Repository: ref.GetRepository(), ManifestDigest: ref.GetManifestDigest()},
	}, nil
}

// HandlePrepareUpdate validates and stages firmware without changing the device.
func (m *ReferenceAssetLink) HandlePrepareUpdate(metadata *artifact.ArtifactMetaData, reporter artifact.ProgressReporter) (*artifact.Result, error) {
	target, err := parseUpdateTarget(metadata)
	if err != nil {
		return artifact.Failed(artifact.InvalidRequest, "Invalid update metadata"), err
	}
	if !m.updateLock.TryLock() {
		return artifact.Failed(artifact.OperationFailed, "Another simulated update operation is running"), nil
	}
	defer m.updateLock.Unlock()

	if m.prepared == nil {
		return m.prepareUpdate(target, reporter)
	}
	if m.prepared.target == target {
		return artifact.Succeeded(), nil
	}
	return artifact.Failed(artifact.InvalidRequest, "Activate or cancel the existing preparation first"), nil
}

func (m *ReferenceAssetLink) prepareUpdate(target updateTarget, reporter artifact.ProgressReporter) (*artifact.Result, error) {
	device, failure, err := connectUpdateDevice(target.address)
	if failure != nil {
		return failure, err
	}
	firmware, failure, err := downloadFirmware(target.reference, device, reporter)
	if failure != nil {
		return failure, err
	}
	if err := reporter.Context().Err(); err != nil {
		return artifact.Failed(artifact.OperationFailed, "Preparation was interrupted"), err
	}
	// The temporary download has been removed. Only the parsed simulator
	// firmware stays in memory; real devices may need persistent staging.
	m.prepared = &preparedUpdate{target: target, firmware: firmware}
	return artifact.Succeeded(), nil
}

func downloadFirmware(ref oci.Reference, device simdevices.SimulatedDevice, reporter artifact.ProgressReporter) (simdevices.PreparedFirmware, *artifact.Result, error) {
	if err := reportUpdateStep(reporter, "download"); err != nil {
		return simdevices.PreparedFirmware{}, artifact.Failed(artifact.OperationFailed, "Preparation was interrupted"), err
	}
	dir, err := os.MkdirTemp("", "reference-firmware-")
	if err != nil {
		return simdevices.PreparedFirmware{}, artifact.Failed(artifact.OperationFailed, "Could not create a staging directory"), err
	}
	defer removeDownload(dir)

	client := oci.Client{HTTP: http.DefaultClient}
	files, err := client.Download(reporter.Context(), ref, dir)
	if err != nil {
		return simdevices.PreparedFirmware{}, artifact.Failed(artifact.ArtifactUnavailable, "Could not download verified update files"), err
	}
	if len(files) != 1 {
		return simdevices.PreparedFirmware{}, artifact.Failed(artifact.ArtifactInvalid, "Expected one simulated firmware file"), nil
	}
	firmware, err := device.PrepareFirmware(filepath.Join(dir, files[0]))
	if err != nil {
		return simdevices.PreparedFirmware{}, artifact.Failed(artifact.ArtifactInvalid, "Firmware is not valid for this simulated device"), err
	}
	return firmware, nil, nil
}

func removeDownload(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Error().Err(err).Msg("Could not remove reference firmware download")
	}
}

// HandleActivateUpdate installs the prepared firmware and reboots the simulator.
func (m *ReferenceAssetLink) HandleActivateUpdate(metadata *artifact.ArtifactMetaData, reporter artifact.ProgressReporter) (*artifact.Result, error) {
	if !m.updateLock.TryLock() {
		return artifact.Failed(artifact.OperationFailed, "Another simulated update operation is running"), nil
	}
	defer m.updateLock.Unlock()
	if failure, err := m.checkPreparation(metadata); failure != nil {
		return failure, err
	}
	device, failure, err := connectUpdateDevice(m.prepared.target.address)
	if failure != nil {
		return failure, err
	}
	result, err := activateFirmware(device, m.prepared.firmware, reporter)
	if err == nil && result.GetFailure() == nil {
		m.prepared = nil
	}
	return result, err
}

func activateFirmware(device simdevices.SimulatedDevice, firmware simdevices.PreparedFirmware, reporter artifact.ProgressReporter) (*artifact.Result, error) {
	if err := reportUpdateStep(reporter, "install"); err != nil {
		return artifact.Failed(artifact.OperationFailed, "Activation was interrupted before installation"), err
	}
	if err := installFirmware(reporter.Context(), device, firmware); err != nil {
		return installationFailure(err), err
	}
	if err := reportUpdateStep(reporter, "reboot"); err != nil {
		return artifact.Unknown("Firmware is installed; reboot was not started"), err
	}
	if err := device.RebootDeviceContext(reporter.Context()); err != nil {
		return artifact.Unknown("Could not confirm the simulated device reboot"), err
	}
	if device.GetActiveFirmwareVersion() != firmware.Version() {
		return artifact.Failed(artifact.OperationFailed, "The expected firmware is not active after reboot"), nil
	}
	return artifact.Succeeded(), nil
}

func reportUpdateStep(reporter artifact.ProgressReporter, step string) error {
	if err := reporter.Context().Err(); err != nil {
		return err
	}
	return reporter.ReportIndeterminate(step)
}

func installationFailure(err error) *artifact.Result {
	if errors.Is(err, simdevices.ErrFirmwareRejected) {
		return artifact.Failed(artifact.DeviceRejected, "The simulated device rejected the firmware")
	}
	return artifact.Unknown("Could not confirm firmware installation; check the device before retrying")
}

func installFirmware(ctx context.Context, device simdevices.SimulatedDevice, firmware simdevices.PreparedFirmware) error {
	// Resume only if this preparation was installed. An equal version could
	// have been installed by an earlier job and says nothing about this one.
	if device.HasInstalledPreparation(firmware) {
		return nil
	}
	return device.InstallFirmware(ctx, firmware)
}

// HandleCancelUpdate discards a preparation. It never rolls back firmware.
func (m *ReferenceAssetLink) HandleCancelUpdate(metadata *artifact.ArtifactMetaData, _ artifact.ProgressReporter) (*artifact.Result, error) {
	if !m.updateLock.TryLock() {
		return artifact.Failed(artifact.OperationFailed, "Another simulated update operation is running"), nil
	}
	defer m.updateLock.Unlock()
	if failure, err := m.checkPreparation(metadata); failure != nil {
		return failure, err
	}
	m.prepared = nil
	return artifact.Succeeded(), nil
}

func (m *ReferenceAssetLink) checkPreparation(metadata *artifact.ArtifactMetaData) (*artifact.Result, error) {
	target, err := parseUpdateTarget(metadata)
	if err != nil {
		return artifact.Failed(artifact.InvalidRequest, "Invalid update metadata"), err
	}
	if m.prepared == nil || m.prepared.target != target {
		return artifact.Failed(artifact.InvalidRequest, "No matching prepared update"), nil
	}
	return nil, nil
}

func connectUpdateDevice(address simdevices.SimulatedDeviceAddress) (simdevices.SimulatedDevice, *artifact.Result, error) {
	// The current gateway update flow does not provide device credentials.
	device, err := simdevices.ConnectToDevice(address, nil)
	switch {
	case errors.Is(err, simdevices.ErrUnauthenticated):
		return nil, artifact.Failed(artifact.DeviceRejected, "The simulated device requires credentials"), err
	case err != nil:
		return nil, artifact.Failed(artifact.TargetNotFound, "Could not find the simulated device"), err
	case !device.IsUpdateSupported():
		return nil, artifact.Failed(artifact.NotSupported, "The simulated device does not support firmware updates"), nil
	default:
		return device, nil, nil
	}
}
