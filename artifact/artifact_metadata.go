/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 * SPDX-License-Identifier: MIT
 */

package artifact

import (
	"encoding/json"
	"strings"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/artifact/oci"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/artifact-update"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ArtifactMetaData is the same typed metadata on the wire and in handlers.
// No blob encoding, credential wrappers or conversion layer is needed.
type ArtifactMetaData = generated.ArtifactMetaData

func ValidateMetaData(metadata *ArtifactMetaData) error {
	if metadata == nil || strings.TrimSpace(metadata.JobId) == "" {
		return status.Error(codes.InvalidArgument, "job_id is required")
	}
	if err := validateTarget(metadata); err != nil {
		return err
	}
	ref := metadata.OciReference
	if ref == nil {
		return status.Error(codes.InvalidArgument, "oci_reference is required")
	}
	if err := (oci.Reference{RegistryURL: ref.RegistryUrl, Repository: ref.Repository, ManifestDigest: ref.ManifestDigest}).Validate(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func validateTarget(metadata *ArtifactMetaData) error {
	target := metadata.ConnectionParameterSet
	if target == nil || strings.TrimSpace(target.SchemaUri) == "" || !json.Valid([]byte(target.ParameterJson)) {
		return status.Error(codes.InvalidArgument, "a schema URI and JSON connection parameters are required")
	}
	if target.SubdriverJson != "" && !json.Valid([]byte(target.SubdriverJson)) {
		return status.Error(codes.InvalidArgument, "subdriver configuration must be JSON")
	}
	return nil
}
