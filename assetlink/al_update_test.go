/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package assetlink

import (
	"testing"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artefact"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/metadata"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const artefactUpdateApiName = "siemens.industrialassethub.artefact_update.v1.ArtefactUpdateApi"

type updateImplementation struct{}

func (*updateImplementation) HandlePrepareUpdate(artefact.ArtefactMetaData, artefact.ArtefactReceiver) error {
	return nil
}
func (*updateImplementation) HandleActivateUpdate(artefact.ArtefactMetaData, artefact.ArtefactReceiver) error {
	return nil
}
func (*updateImplementation) HandleCancelUpdate(artefact.ArtefactMetaData, artefact.StatusTransmitter) error {
	return nil
}

func TestRegisterUpdateServerWithoutImplementation(t *testing.T) {
	al := New(metadata.Metadata{}).Build()
	al.grpcServer = grpc.NewServer()
	t.Cleanup(al.grpcServer.Stop)

	al.registerUpdateServer()

	services := al.grpcServer.GetServiceInfo()
	require.NotContains(t, services, artefactUpdateApiName)
}

func TestRegisterUpdateServerWithImplementation(t *testing.T) {
	al := New(metadata.Metadata{}).
		Update(&updateImplementation{}).
		Build()
	al.grpcServer = grpc.NewServer()
	t.Cleanup(al.grpcServer.Stop)

	al.registerUpdateServer()

	services := al.grpcServer.GetServiceInfo()
	require.Contains(t, services, artefactUpdateApiName)
	service := services[artefactUpdateApiName]
	require.Len(t, service.Methods, 3)
}

func TestRegisterUpdateServerWithCustomServer(t *testing.T) {
	builder := New(metadata.Metadata{})
	builder.ArtefactUpdateApiServer = &generated.UnimplementedArtefactUpdateApiServer{}
	al := builder.Build()
	al.grpcServer = grpc.NewServer()
	t.Cleanup(al.grpcServer.Stop)

	al.registerUpdateServer()

	services := al.grpcServer.GetServiceInfo()
	require.Contains(t, services, artefactUpdateApiName)
	service := services[artefactUpdateApiName]
	require.Len(t, service.Methods, 3)
}

func TestRegisterUpdateServerWithBothImplementationAndCustomServer(t *testing.T) {
	builder := New(metadata.Metadata{}).Update(&updateImplementation{})
	builder.ArtefactUpdateApiServer = &generated.UnimplementedArtefactUpdateApiServer{}
	al := builder.Build()
	al.grpcServer = grpc.NewServer()
	t.Cleanup(al.grpcServer.Stop)

	al.registerUpdateServer()

	services := al.grpcServer.GetServiceInfo()
	require.Contains(t, services, artefactUpdateApiName)
	service := services[artefactUpdateApiName]
	require.Len(t, service.Methods, 3)
}
