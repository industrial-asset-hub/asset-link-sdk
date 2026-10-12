// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT
package artefact

import (
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artefact-update"
	common "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
	"testing"
)

func TestTypedConnectionTargetAndOCIReference(t *testing.T) {
	m := &ArtefactMetaData{JobId: "job", AssetId: "asset", ViewId: "view", OciReference: &generated.OciArtefactReference{RegistryUrl: "http://zot:5000", Repository: "firmware", ManifestDigest: "sha256:" + strings.Repeat("a", 64)}, ConnectionParameterSet: &common.ConnectionParameterSet{SchemaUri: "urn:sat", ParameterJson: `{"key":"device"}`}}
	require.NoError(t, ValidateMetaData(m))
	require.Equal(t, "view", m.GetViewId())
	require.Equal(t, "asset", m.GetAssetId())
	m.ConnectionParameterSet.SubdriverJson = "invalid"
	require.Error(t, ValidateMetaData(m))
	m.ConnectionParameterSet.SubdriverJson = ""
	m.ConnectionParameterSet.ParameterJson = "invalid"
	require.Error(t, ValidateMetaData(m))
	require.Error(t, ValidateMetaData(nil))
}

func TestUpdateContractHasNoFileTransfer(t *testing.T) {
	file := generated.File_artefact_update_proto
	for _, name := range []string{"ArtefactChunk", "ArtefactMessage", "ArtefactOperationRequest", "DeviceIdentifier", "JobIdentifier", "ArtefactIdentifier", "DeviceCredentials", "ArtefactCredentials"} {
		require.Nil(t, file.Messages().ByName(protoreflect.Name(name)), name)
	}
	metadata := file.Messages().ByName("ArtefactMetaData")
	require.Equal(t, 6, metadata.Fields().Len())
	require.Equal(t, protoreflect.FullName("siemens.common.address.v1.ConnectionParameterSet"), metadata.Fields().ByName("connection_parameter_set").Message().FullName())
	methods := file.Services().ByName("ArtefactUpdateApi").Methods()
	require.Equal(t, 3, methods.Len())
	for i := 0; i < methods.Len(); i++ {
		m := methods.Get(i)
		require.False(t, m.IsStreamingClient())
		require.True(t, m.IsStreamingServer())
		require.Equal(t, metadata.FullName(), m.Input().FullName())
		require.Equal(t, protoreflect.FullName("siemens.industrialassethub.artefact_update.v1.ArtefactOperationStatus"), m.Output().FullName())
	}
}
