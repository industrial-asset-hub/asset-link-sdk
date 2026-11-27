/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package assetlink

import (
	"testing"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/metadata"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const artifactUpdateApiName = "siemens.industrialassethub.artifact_update.v1.ArtifactUpdateApi"

type updateImplementation struct{}

func (*updateImplementation) HandlePrepareUpdate(*artifact.ArtifactMetaData, artifact.ProgressReporter) (*artifact.Result, error) {
	return artifact.Succeeded(), nil
}
func (*updateImplementation) HandleActivateUpdate(*artifact.ArtifactMetaData, artifact.ProgressReporter) (*artifact.Result, error) {
	return artifact.Succeeded(), nil
}
func (*updateImplementation) HandleCancelUpdate(*artifact.ArtifactMetaData, artifact.ProgressReporter) (*artifact.Result, error) {
	return artifact.Succeeded(), nil
}

func TestRegisterUpdateServerWithoutImplementation(t *testing.T) {
	al := New(metadata.Metadata{}).Build()
	al.grpcServer = grpc.NewServer()
	t.Cleanup(al.grpcServer.Stop)

	al.registerUpdateServer()

	services := al.grpcServer.GetServiceInfo()
	require.NotContains(t, services, artifactUpdateApiName)
}

func TestRegisterUpdateServerWithImplementation(t *testing.T) {
	al := New(metadata.Metadata{}).
		Update(&updateImplementation{}).
		Build()
	al.grpcServer = grpc.NewServer()
	t.Cleanup(al.grpcServer.Stop)

	al.registerUpdateServer()

	services := al.grpcServer.GetServiceInfo()
	require.Contains(t, services, artifactUpdateApiName)
	service := services[artifactUpdateApiName]
	require.Len(t, service.Methods, 3)
}

func TestRegisterUpdateServerWithCustomServer(t *testing.T) {
	builder := New(metadata.Metadata{})
	builder.ArtifactUpdateApiServer = &generated.UnimplementedArtifactUpdateApiServer{}
	al := builder.Build()
	al.grpcServer = grpc.NewServer()
	t.Cleanup(al.grpcServer.Stop)

	al.registerUpdateServer()

	services := al.grpcServer.GetServiceInfo()
	require.Contains(t, services, artifactUpdateApiName)
	service := services[artifactUpdateApiName]
	require.Len(t, service.Methods, 3)
}

func TestRegisterUpdateServerWithBothImplementationAndCustomServer(t *testing.T) {
	builder := New(metadata.Metadata{}).Update(&updateImplementation{})
	builder.ArtifactUpdateApiServer = &generated.UnimplementedArtifactUpdateApiServer{}
	al := builder.Build()
	al.grpcServer = grpc.NewServer()
	t.Cleanup(al.grpcServer.Stop)

	al.registerUpdateServer()

	services := al.grpcServer.GetServiceInfo()
	require.Contains(t, services, artifactUpdateApiName)
	service := services[artifactUpdateApiName]
	require.Len(t, service.Methods, 3)
}
