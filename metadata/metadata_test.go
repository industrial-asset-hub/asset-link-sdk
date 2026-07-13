/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package metadata

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionJSONRoundTrip(t *testing.T) {
	original := Version{
		Version: "1.2.3",
		Commit:  "abc123",
		Date:    "2026-07-10",
	}

	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	assert.JSONEq(t, `{"version":"1.2.3","commit":"abc123","date":"2026-07-10"}`, string(encoded))

	var decoded Version
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	assert.Equal(t, original, decoded)
}

func TestMetadataFieldsArePreserved(t *testing.T) {
	original := Metadata{
		AlId:        "asset-01",
		AlName:      "Asset Link",
		Version:     Version{Version: "1.2.3", Commit: "abc123", Date: "2026-07-10"},
		Vendor:      "Siemens",
		Description: "test metadata",
		DocUrl:      "https://example.com/docs",
		FeedbackUrl: "https://example.com/feedback",
	}

	assert.Equal(t, "asset-01", original.AlId)
	assert.Equal(t, "Asset Link", original.AlName)
	assert.Equal(t, Version{Version: "1.2.3", Commit: "abc123", Date: "2026-07-10"}, original.Version)
	assert.Equal(t, "Siemens", original.Vendor)
	assert.Equal(t, "test metadata", original.Description)
	assert.Equal(t, "https://example.com/docs", original.DocUrl)
	assert.Equal(t, "https://example.com/feedback", original.FeedbackUrl)
}
