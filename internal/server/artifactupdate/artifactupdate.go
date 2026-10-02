/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artifactupdate

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/internal/features"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ArtifactUpdateServerEntity struct {
	generated.UnimplementedArtifactUpdateApiServer
	features.Update
}

func (d *ArtifactUpdateServerEntity) PrepareUpdate(metadata *generated.ArtifactMetaData, stream generated.ArtifactUpdateApi_PrepareUpdateServer) error {
	if d.Update == nil {
		return status.Error(codes.Unimplemented, "no Update implementation found")
	}
	return handleUpdate(metadata, stream, d.HandlePrepareUpdate)
}

func (d *ArtifactUpdateServerEntity) ActivateUpdate(metadata *generated.ArtifactMetaData, stream generated.ArtifactUpdateApi_ActivateUpdateServer) error {
	if d.Update == nil {
		return status.Error(codes.Unimplemented, "no Update implementation found")
	}
	return handleUpdate(metadata, stream, d.HandleActivateUpdate)
}

func handleUpdate(metadata *generated.ArtifactMetaData, stream grpc.ServerStreamingServer[generated.UpdateEvent], handler func(*artifact.ArtifactMetaData, artifact.ProgressReporter) (*artifact.Result, error)) error {
	if err := artifact.ValidateMetaData(metadata); err != nil {
		return err
	}
	return artifact.RunUpdate(metadata, stream, handler)
}

func (d *ArtifactUpdateServerEntity) CancelUpdate(metadata *generated.ArtifactMetaData, stream generated.ArtifactUpdateApi_CancelUpdateServer) error {
	if d.Update == nil {
		return status.Error(codes.Unimplemented, "no Update implementation found")
	}
	return handleUpdate(metadata, stream, d.HandleCancelUpdate)
}
