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

func (d *ArtefactUpdateServerEntity) PrepareUpdate(stream generated.ArtefactUpdateApi_PrepareUpdateServer) error {
	if d.Update == nil {
		return status.Error(codes.Unimplemented, "no Update implementation found")
	}
	return handleUpdate(stream, d.HandlePrepareUpdate)
}

func (d *ArtefactUpdateServerEntity) ActivateUpdate(stream generated.ArtefactUpdateApi_ActivateUpdateServer) error {
	if d.Update == nil {
		return status.Error(codes.Unimplemented, "no Update implementation found")
	}
	return handleUpdate(stream, d.HandleActivateUpdate)
}

func handleUpdate(stream grpc.BidiStreamingServer[generated.ArtefactChunk, generated.ArtefactMessage], handler func(artefact.ArtefactMetaData, artefact.ArtefactReceiver) error) error {
	receiver := artefact.NewArtefactReceiver(stream)
	metadata, err := receiver.ReceiveArtefactMetaData()
	if err != nil {
		return err
	}
	return handler(metadata, receiver)
}

func (d *ArtefactUpdateServerEntity) CancelUpdate(metadata *generated.ArtefactMetaData, stream generated.ArtefactUpdateApi_CancelUpdateServer) error {
	if d.Update == nil {
		return status.Error(codes.Unimplemented, "no Update implementation found")
	}
	converted, err := artefact.NewArtefactMetaDataFromInternal(metadata)
	if err != nil {
		return err
	}
	return d.HandleCancelUpdate(converted, artefact.NewStatusTransmitter(stream))
}
