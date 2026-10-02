/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
)

// HandlePrepareUpdate prepares an artifact update for later activation.
func (m *AssetLinkImplementation) HandlePrepareUpdate(metadata *artifact.ArtifactMetaData, transmitter artifact.ProgressReporter) (*artifact.Result, error) {
	// Locate/authenticate the device using metadata.GetConnectionParameterSet().
	// Download from metadata.GetOciReference() over HTTP if needed.
	// Prepare without installing the update or disrupting the device.
	// Report progress using transmitter.ReportProgress() and observe transmitter.Context().
	// Retain any information needed later for activation or cancellation.
	return artifact.Failed(artifact.NotSupported, "Update preparation is not supported"), nil
}

// HandleActivateUpdate activates the prepared update on the device.
func (m *AssetLinkImplementation) HandleActivateUpdate(metadata *artifact.ArtifactMetaData, transmitter artifact.ProgressReporter) (*artifact.Result, error) {
	// Identify the device using metadata.GetConnectionParameterSet().
	// Retrieve any preparation state needed for this operation.
	// Download the artifact from the local OCI registry here if it was deferred.
	// Activate the update and report progress using transmitter.ReportProgress().
	return artifact.Failed(artifact.NotSupported, "Update activation is not supported"), nil
}

// HandleCancelUpdate cleans up a prepared update.
func (m *AssetLinkImplementation) HandleCancelUpdate(metadata *artifact.ArtifactMetaData, transmitter artifact.ProgressReporter) (*artifact.Result, error) {
	// Identify the device using metadata.GetConnectionParameterSet().
	// Retrieve any preparation state needed for this operation.
	// Clean up the device and stored state, reporting progress with transmitter.ReportProgress().
	return artifact.Failed(artifact.NotSupported, "Update cancellation is not supported"), nil
}
