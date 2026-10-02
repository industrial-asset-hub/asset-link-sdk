// SPDX-FileCopyrightText: 2026 Siemens AG
// SPDX-License-Identifier: MIT
package artifact

import (
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	common "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
	"testing"
)

func TestTypedConnectionTargetAndOCIReference(t *testing.T) {
	m := &ArtifactMetaData{JobId: "job", OciReference: &generated.OciArtifactReference{RegistryUrl: "http://oci-registry:5000", Repository: "firmware", ManifestDigest: "sha256:" + strings.Repeat("a", 64)}, ConnectionParameterSet: &common.ConnectionParameterSet{SchemaUri: "urn:sat", ParameterJson: `{"key":"device"}`}}
	require.NoError(t, ValidateMetaData(m))
	m.ConnectionParameterSet.SubdriverJson = "invalid"
	require.Error(t, ValidateMetaData(m))
	m.ConnectionParameterSet.SubdriverJson = ""
	m.ConnectionParameterSet.ParameterJson = "invalid"
	require.Error(t, ValidateMetaData(m))
	require.Error(t, ValidateMetaData(nil))
}

func TestUpdateContractHasNoFileTransfer(t *testing.T) {
	file := generated.File_artifact_update_proto
	require.Equal(t, protoreflect.FullName("siemens.industrialassethub.artifact_update.v1"), file.Package())
	for _, name := range []string{"ArtifactChunk", "ArtifactMessage", "ArtifactOperationRequest", "DeviceIdentifier", "JobIdentifier", "ArtifactIdentifier", "DeviceCredentials", "ArtifactCredentials"} {
		require.Nil(t, file.Messages().ByName(protoreflect.Name(name)), name)
	}
	metadata := file.Messages().ByName("ArtifactMetaData")
	require.Equal(t, 3, metadata.Fields().Len())
	require.Nil(t, metadata.Fields().ByName("artifact_type"))
	require.Nil(t, file.Enums().ByName("ArtifactType"))
	require.Equal(t, 0, metadata.ReservedRanges().Len())
	require.Equal(t, 0, metadata.ReservedNames().Len())
	for i, name := range []protoreflect.Name{"job_id", "connection_parameter_set", "oci_reference"} {
		require.Equal(t, protoreflect.FieldNumber(i+1), metadata.Fields().ByName(name).Number())
	}
	require.Equal(t, protoreflect.FullName("siemens.industrialassethub.artifact_update.v1.OciArtifactReference"), metadata.Fields().ByName("oci_reference").Message().FullName())
	require.Equal(t, protoreflect.FullName("siemens.common.address.v1.ConnectionParameterSet"), metadata.Fields().ByName("connection_parameter_set").Message().FullName())
	methods := file.Services().ByName("ArtifactUpdateApi").Methods()
	require.Equal(t, 3, methods.Len())
	for i := 0; i < methods.Len(); i++ {
		m := methods.Get(i)
		require.False(t, m.IsStreamingClient())
		require.True(t, m.IsStreamingServer())
		require.Equal(t, metadata.FullName(), m.Input().FullName())
		require.Equal(t, protoreflect.FullName("siemens.industrialassethub.artifact_update.v1.UpdateEvent"), m.Output().FullName())
	}
}
