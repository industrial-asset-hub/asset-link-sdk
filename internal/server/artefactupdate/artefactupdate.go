/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artefactupdate

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artefact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/internal/features"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ArtefactUpdateServerEntity struct {
	generated.UnimplementedArtefactUpdateApiServer
	features.Update
}

func (d *ArtefactUpdateServerEntity) PrepareUpdate(metadata *generated.ArtefactMetaData, stream generated.ArtefactUpdateApi_PrepareUpdateServer) error {
	if d.Update == nil {
		return status.Error(codes.Unimplemented, "no Update implementation found")
	}
	return handleUpdate(metadata, stream, d.HandlePrepareUpdate)
}

func (d *ArtefactUpdateServerEntity) ActivateUpdate(metadata *generated.ArtefactMetaData, stream generated.ArtefactUpdateApi_ActivateUpdateServer) error {
	if d.Update == nil {
		return status.Error(codes.Unimplemented, "no Update implementation found")
	}
	return handleUpdate(metadata, stream, d.HandleActivateUpdate)
}

func handleUpdate(metadata *generated.ArtefactMetaData, stream grpc.ServerStreamingServer[generated.UpdateEvent], handler func(*artefact.ArtefactMetaData, artefact.ProgressReporter) (*artefact.Result, error)) error {
	if err := artefact.ValidateMetaData(metadata); err != nil {
		return err
	}
	return artefact.RunUpdate(metadata, stream, handler)
}

func (d *ArtefactUpdateServerEntity) CancelUpdate(metadata *generated.ArtefactMetaData, stream generated.ArtefactUpdateApi_CancelUpdateServer) error {
	if d.Update == nil {
		return status.Error(codes.Unimplemented, "no Update implementation found")
	}
	if err := artefact.ValidateMetaData(metadata); err != nil {
		return err
	}
	return d.HandleCancelUpdate(metadata, artefact.NewStatusTransmitter(stream))
}
