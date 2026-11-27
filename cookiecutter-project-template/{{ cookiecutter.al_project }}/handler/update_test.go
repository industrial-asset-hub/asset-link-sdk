/*
 * SPDX-FileCopyrightText: {{cookiecutter.year}} {{cookiecutter.company}}
 *
 * SPDX-License-Identifier: MIT
 *
 */

package handler

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/contracttest"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/oci"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	discovery "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/stretchr/testify/assert"
)

const (
	updateSucceeded = generated.UpdateOutcome_UPDATE_OUTCOME_SUCCEEDED
	updateFailed    = generated.UpdateOutcome_UPDATE_OUTCOME_FAILED
)

func TestUpdate(t *testing.T) {
	t.Run("prepareAndActivate", func(t *testing.T) {
		ctx, driver := setupUpdateTest(t)
		request := updateRequest(t, "update-1", dummyDeviceMAC)

		contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateSucceeded)
		contracttest.Check(t, ctx, request, driver.HandleActivateUpdate, updateSucceeded)
		assert.Equal(t, request.GetOciReference().GetManifestDigest(), dummyDevice.activeDigest)

		// The preparation is removed after a successful activation.
		result := contracttest.Check(t, ctx, request, driver.HandleActivateUpdate, updateFailed)
		assert.Equal(t, artifact.InvalidRequest, result.GetFailure().GetReason())
	})

	t.Run("prepareAndCancel", func(t *testing.T) {
		ctx, driver := setupUpdateTest(t)
		request := updateRequest(t, "update-2", dummyDeviceMAC)

		contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateSucceeded)
		contracttest.Check(t, ctx, request, driver.HandleCancelUpdate, updateSucceeded)
		contracttest.Check(t, ctx, request, driver.HandleActivateUpdate, updateFailed)
		assert.Empty(t, dummyDevice.activeDigest)
	})

	t.Run("unknownDevice", func(t *testing.T) {
		ctx, driver := setupUpdateTest(t)
		request := updateRequest(t, "update-3", "02:00:00:00:00:99")

		result := contracttest.Check(t, ctx, request, driver.HandlePrepareUpdate, updateFailed)
		assert.Equal(t, artifact.TargetNotFound, result.GetFailure().GetReason())
	})
}

func setupUpdateTest(t *testing.T) (context.Context, *AssetLinkImplementation) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	t.Cleanup(func() {
		for jobID := range preparedUpdates.jobs {
			_ = releasePrepared(jobID, true)
		}
		dummyDevice.activeJobID, dummyDevice.activeDigest = "", ""
	})
	return ctx, &AssetLinkImplementation{}
}

// updateRequest builds the request the Gateway or al-ctl sends, using the
// connection parameters published by discovery and a test OCI registry.
func updateRequest(t *testing.T, jobID, macAddress string) *artifact.ArtifactMetaData {
	return &artifact.ArtifactMetaData{
		JobId: jobID,
		ConnectionParameterSet: &discovery.ConnectionParameterSet{
			SchemaUri:     connectionSchema,
			ParameterJson: fmt.Sprintf(`{"macAddress":%q}`, macAddress),
		},
		OciReference: serveUpdateArtifact(t, []byte("dummy update")),
	}
}

// serveUpdateArtifact serves one file as an OCI artifact from a local test server.
func serveUpdateArtifact(t *testing.T, content []byte) *generated.OciArtifactReference {
	t.Helper()
	config := []byte("{}")
	layer := oci.Descriptor{
		Digest:      fmt.Sprintf("sha256:%x", sha256.Sum256(content)),
		Size:        int64(len(content)),
		Annotations: map[string]string{"org.opencontainers.image.title": "update.bin"},
	}
	configDescriptor := oci.Descriptor{Digest: fmt.Sprintf("sha256:%x", sha256.Sum256(config)), Size: int64(len(config))}
	manifest, err := json.Marshal(oci.Manifest{
		SchemaVersion: 2, MediaType: oci.ManifestMediaType, Config: configDescriptor, Layers: []oci.Descriptor{layer},
	})
	assert.NoError(t, err)
	digest := fmt.Sprintf("sha256:%x", sha256.Sum256(manifest))

	mux := http.NewServeMux()
	mux.HandleFunc("/v2/my-update/manifests/"+digest, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", oci.ManifestMediaType)
		_, _ = w.Write(manifest)
	})
	mux.HandleFunc("/v2/my-update/blobs/"+layer.Digest, func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(content) })
	mux.HandleFunc("/v2/my-update/blobs/"+configDescriptor.Digest, func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(config) })
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return &generated.OciArtifactReference{RegistryUrl: server.URL, Repository: "my-update", ManifestDigest: digest}
}
