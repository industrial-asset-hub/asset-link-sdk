/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artefact"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandlePrepareUpdate prepares firmware for later activation.
func (m *AssetLinkImplementation) HandlePrepareUpdate(metadata *artefact.ArtefactMetaData, transmitter artefact.StatusTransmitter) error {
	// Locate/authenticate the device using metadata.GetConnectionParameterSet().
	// Download from metadata.GetOciReference() over HTTP if needed.
	// Prepare without installing firmware or disrupting the device.
	// Report progress using transmitter.UpdateStatus() and observe transmitter.Context().
	// Retain any job state needed for activation or cancellation.
	return status.Error(codes.Unimplemented, "firmware update preparation is not implemented")
}

// HandleActivateUpdate activates the prepared firmware on the device.
func (m *AssetLinkImplementation) HandleActivateUpdate(metadata *artefact.ArtefactMetaData, transmitter artefact.StatusTransmitter) error {
	// Find the prepared update using metadata.GetJobId() and the device identifier.
	// Download firmware from the local OCI registry here if it was deferred.
	// Activate the firmware and report progress using transmitter.UpdateStatus().
	return status.Error(codes.Unimplemented, "firmware update activation is not implemented")
}

// HandleCancelUpdate cleans up a prepared firmware update.
func (m *AssetLinkImplementation) HandleCancelUpdate(metadata *artefact.ArtefactMetaData, transmitter artefact.StatusTransmitter) error {
	// Find the prepared update using metadata.GetJobId() and the device identifier.
	// Clean up the device and job state, reporting progress with transmitter.UpdateStatus().
	return status.Error(codes.Unimplemented, "firmware update cancellation is not implemented")
}
