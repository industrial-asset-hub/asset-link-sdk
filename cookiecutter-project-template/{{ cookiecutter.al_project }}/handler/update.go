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
func (m *AssetLinkImplementation) HandlePrepareUpdate(metadata artefact.ArtefactMetaData, receiver artefact.ArtefactReceiver) error {
	// Locate the device using metadata.GetDeviceIdentifierBlob() and authenticate
	// using metadata.GetDeviceCredentials() if required.
	// Request the firmware using receiver.ReceiveArtefactToWriter() or
	// receiver.ReceiveArtefactToFile(), then prepare it on the device.
	// Report progress using receiver.UpdateStatus() and observe receiver.Context().
	// Retain any job state needed for activation or cancellation.
	return status.Error(codes.Unimplemented, "firmware update preparation is not implemented")
}

// HandleActivateUpdate activates the prepared firmware on the device.
func (m *AssetLinkImplementation) HandleActivateUpdate(metadata artefact.ArtefactMetaData, receiver artefact.ArtefactReceiver) error {
	// Find the prepared update using metadata.GetJobId() and the device identifier.
	// Devices without two-stage updates can request firmware here instead.
	// Activate the firmware and report progress using receiver.UpdateStatus().
	return status.Error(codes.Unimplemented, "firmware update activation is not implemented")
}

// HandleCancelUpdate cleans up a prepared firmware update.
func (m *AssetLinkImplementation) HandleCancelUpdate(metadata artefact.ArtefactMetaData, transmitter artefact.StatusTransmitter) error {
	// Find the prepared update using metadata.GetJobId() and the device identifier.
	// Clean up the device and job state, reporting progress with transmitter.UpdateStatus().
	return status.Error(codes.Unimplemented, "firmware update cancellation is not implemented")
}
