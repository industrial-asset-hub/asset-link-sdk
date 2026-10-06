/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artefactupdate

import (
	"context"

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

// GetUpdateTargetRequirements delegates only when the update implementation
// explicitly provides this optional capability.
func (d *ArtefactUpdateServerEntity) GetUpdateTargetRequirements(ctx context.Context, request *generated.GetUpdateTargetRequirementsRequest) (*generated.GetUpdateTargetRequirementsResponse, error) {
	provider, ok := d.Update.(artefact.UpdateTargetRequirementsProvider)
	if !ok {
		return nil, status.Error(codes.Unimplemented, "update target requirements are not implemented")
	}
	if request == nil || request.ArtefactIdentifier == nil {
		return nil, status.Error(codes.InvalidArgument, "artefact identifier is required")
	}
	if _, ok := generated.ArtefactType_name[int32(request.ArtefactIdentifier.Type)]; !ok {
		return nil, status.Error(codes.InvalidArgument, "unknown artefact type")
	}
	return provider.GetUpdateTargetRequirements(ctx, request)
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
