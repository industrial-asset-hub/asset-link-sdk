/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

// This file contains a working dummy implementation of artifact updates for
// the dummy device published by discovery. Replace the parts marked with TODO
// with the logic for your devices.
// Guide: https://github.com/industrial-asset-hub/asset-link-sdk/blob/main/docs/artifact-management.md

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/oci"
	"github.com/rs/zerolog/log"
)

// TODO: Use the SchemaUri and ParameterJson format your discovery publishes.
const connectionSchema = "urn:example:asset-link:mac-address:v1"

type connectionParameters struct {
	MacAddress string `json:"macAddress"`
}

// dummyDeviceMAC identifies the dummy device published by discovery.
const dummyDeviceMAC = "00:16:3e:01:02:03"

// dummyDevice simulates the update state of the dummy device.
// TODO: Replace it with communication with your devices.
var dummyDevice = struct {
	sync.Mutex
	activeJobID  string
	activeDigest string
}{}

// preparedUpdate is what Activate and Cancel need from Prepare.
type preparedUpdate struct {
	device         connectionParameters
	manifestDigest string
	stagingDir     string
	files          []string
	busy           bool // Activate or Cancel is running for this job
	activated      bool // Activate has started and may have changed the device
}

// preparedUpdates is kept in memory, so a restart discards preparations.
// Preparations that are never activated or cancelled keep their files until
// the process exits. Remove stale ones if that matters for your device.
var preparedUpdates = struct {
	sync.Mutex
	jobs map[string]preparedUpdate
}{jobs: map[string]preparedUpdate{}}

// HandlePrepareUpdate downloads and stages an artifact without applying it.
func (m *AssetLinkImplementation) HandlePrepareUpdate(metadata *artifact.ArtifactMetaData, transmitter artifact.ProgressReporter) (*artifact.Result, error) {
	device, failed := parseConnection(metadata)
	if failed != nil {
		return failed, nil
	}
	// TODO: Check that your device is reachable.
	if device.MacAddress != dummyDeviceMAC {
		return artifact.Failed(artifact.TargetNotFound, "Device is not reachable"), nil
	}

	if err := transmitter.ReportIndeterminate("download"); err != nil {
		return artifact.Failed(artifact.OperationFailed, "Preparation was interrupted"), err
	}
	stagingDir, err := os.MkdirTemp("", "asset-link-update-")
	if err != nil {
		return artifact.Failed(artifact.OperationFailed, "Could not create a staging directory"), err
	}
	ref := metadata.GetOciReference()
	client := oci.Client{HTTP: http.DefaultClient}
	files, err := client.Download(transmitter.Context(), oci.Reference{
		RegistryURL:    ref.GetRegistryUrl(),
		Repository:     ref.GetRepository(),
		ManifestDigest: ref.GetManifestDigest(),
	}, stagingDir)
	if err != nil {
		return artifact.Failed(artifact.ArtifactUnavailable, "Update files could not be downloaded"),
			errors.Join(err, os.RemoveAll(stagingDir))
	}

	// TODO: Validate the files (filepath.Join(stagingDir, name)) for your
	// device, and optionally transfer them to the device without applying them.
	if len(files) == 0 {
		return artifact.Failed(artifact.ArtifactInvalid, "The artifact contains no update files"), os.RemoveAll(stagingDir)
	}

	return storePrepared(metadata.GetJobId(), preparedUpdate{
		device: device, manifestDigest: ref.GetManifestDigest(), stagingDir: stagingDir, files: files,
	})
}

// HandleActivateUpdate applies the prepared artifact on the device.
func (m *AssetLinkImplementation) HandleActivateUpdate(metadata *artifact.ArtifactMetaData, transmitter artifact.ProgressReporter) (*artifact.Result, error) {
	device, failed := parseConnection(metadata)
	if failed != nil {
		return failed, nil
	}
	prepared, failed := claimPrepared(metadata, device, true)
	if failed != nil {
		return failed, nil
	}
	// Keep the preparation unless activation succeeds, so the same job can be
	// activated again after a failure.
	succeeded := false
	defer func() {
		if err := releasePrepared(metadata.GetJobId(), succeeded); err != nil {
			log.Warn().Err(err).Str("jobId", metadata.GetJobId()).Msg("Could not remove prepared update files")
		}
	}()

	dummyDevice.Lock()
	defer dummyDevice.Unlock()
	// A previous attempt for this job may have ended with Unknown. If the
	// device already runs this update, do not install it again.
	// TODO: Read this state from your device.
	if dummyDevice.activeJobID == metadata.GetJobId() {
		succeeded = true
		return artifact.Succeeded(), nil
	}

	if err := transmitter.ReportIndeterminate("install"); err != nil {
		return artifact.Failed(artifact.OperationFailed, "Activation was interrupted before installation"), err
	}
	// TODO: Apply prepared.files on your device and verify the result. Pass
	// transmitter.Context() to device calls so they stop when the Gateway
	// aborts the call; stopping does not undo earlier device work.
	// - Device rejected the update: return artifact.Failed(artifact.DeviceRejected, "..."), nil
	// - Connection lost or outcome unclear: return artifact.Unknown("Check the device before retrying"), err
	dummyDevice.activeJobID = metadata.GetJobId()
	dummyDevice.activeDigest = prepared.manifestDigest
	log.Info().Str("device", device.MacAddress).Strs("files", prepared.files).
		Str("digest", prepared.manifestDigest).Msg("Dummy device updated")

	succeeded = true
	return artifact.Succeeded(), nil
}

// HandleCancelUpdate discards a prepared artifact without applying it.
func (m *AssetLinkImplementation) HandleCancelUpdate(metadata *artifact.ArtifactMetaData, transmitter artifact.ProgressReporter) (*artifact.Result, error) {
	device, failed := parseConnection(metadata)
	if failed != nil {
		return failed, nil
	}
	if _, failed := claimPrepared(metadata, device, false); failed != nil {
		return failed, nil
	}
	// TODO: Remove anything Prepare transferred to your device. If that
	// fails, keep the preparation for another attempt:
	//   _ = releasePrepared(metadata.GetJobId(), false)
	//   return artifact.Failed(artifact.OperationFailed, "..."), err
	if err := releasePrepared(metadata.GetJobId(), true); err != nil {
		return artifact.Failed(artifact.OperationFailed, "Could not remove the prepared files"), err
	}
	return artifact.Succeeded(), nil
}

func parseConnection(metadata *artifact.ArtifactMetaData) (connectionParameters, *artifact.Result) {
	var device connectionParameters
	set := metadata.GetConnectionParameterSet()
	if set.GetSchemaUri() != connectionSchema {
		return device, artifact.Failed(artifact.NotSupported, "Connection parameter format is not supported")
	}
	if err := json.Unmarshal([]byte(set.GetParameterJson()), &device); err != nil || device.MacAddress == "" {
		return device, artifact.Failed(artifact.InvalidRequest, "Connection parameters are invalid")
	}
	return device, nil
}

// storePrepared saves a new preparation. A retried Prepare for the same job
// replaces the old preparation unless Activate or Cancel already uses it.
func storePrepared(jobID string, prepared preparedUpdate) (*artifact.Result, error) {
	preparedUpdates.Lock()
	defer preparedUpdates.Unlock()
	if old, ok := preparedUpdates.jobs[jobID]; ok {
		if old.busy || old.activated {
			return artifact.Failed(artifact.OperationFailed, "Activation already started; check the device"),
				os.RemoveAll(prepared.stagingDir)
		}
		if err := os.RemoveAll(old.stagingDir); err != nil {
			log.Warn().Err(err).Str("jobId", jobID).Msg("Could not remove previously prepared update files")
		}
	}
	preparedUpdates.jobs[jobID] = prepared
	return artifact.Succeeded(), nil
}

// claimPrepared returns the preparation matching the request and marks it busy
// so that Activate and Cancel cannot run for the same job at the same time.
// Cancel is rejected once Activate has started, because it cannot undo
// changes Activate may already have made on the device.
func claimPrepared(metadata *artifact.ArtifactMetaData, device connectionParameters, activate bool) (preparedUpdate, *artifact.Result) {
	preparedUpdates.Lock()
	defer preparedUpdates.Unlock()
	prepared, ok := preparedUpdates.jobs[metadata.GetJobId()]
	if !ok {
		return prepared, artifact.Failed(artifact.InvalidRequest, "No preparation found for this job")
	}
	if prepared.device != device || prepared.manifestDigest != metadata.GetOciReference().GetManifestDigest() {
		return prepared, artifact.Failed(artifact.InvalidRequest, "Request does not match the preparation")
	}
	if prepared.busy {
		return prepared, artifact.Failed(artifact.OperationFailed, "Another operation is running for this job")
	}
	if !activate && prepared.activated {
		return prepared, artifact.Failed(artifact.OperationFailed, "Activation already started; check the device")
	}
	prepared.busy = true
	prepared.activated = prepared.activated || activate
	preparedUpdates.jobs[metadata.GetJobId()] = prepared
	return prepared, nil
}

// releasePrepared ends an operation started with claimPrepared. If done is
// true, it removes the preparation and its files; otherwise it keeps them so
// the same job can be activated again.
func releasePrepared(jobID string, done bool) error {
	preparedUpdates.Lock()
	defer preparedUpdates.Unlock()
	prepared, ok := preparedUpdates.jobs[jobID]
	if !ok {
		return nil
	}
	if !done {
		prepared.busy = false
		preparedUpdates.jobs[jobID] = prepared
		return nil
	}
	delete(preparedUpdates.jobs, jobID)
	return os.RemoveAll(prepared.stagingDir)
}
