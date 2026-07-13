/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package driverinfo

import (
	"context"
	"testing"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_drv_info"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVersionInfo(t *testing.T) {
	sut := &DriverInfoServerEntity{
		Metadata: metadata.Metadata{
			AlName:      "Asset Link SDK",
			Vendor:      "Siemens",
			Description: "driver info server",
			DocUrl:      "https://example.com/docs",
			FeedbackUrl: "https://example.com/feedback",
			Version: metadata.Version{
				Version: "1.2.3-beta.1",
				Commit:  "abc123",
				Date:    "2026-07-10",
			},
		},
	}

	resp, err := sut.GetVersionInfo(context.Background(), &generated.GetVersionInfoRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetVersion())

	assert.Equal(t, uint32(1), resp.GetVersion().GetMajor())
	assert.Equal(t, uint32(2), resp.GetVersion().GetMinor())
	assert.Equal(t, uint32(3), resp.GetVersion().GetPatch())
	assert.Equal(t, "beta.1", resp.GetVersion().GetSuffix())
	assert.Equal(t, "Siemens", resp.GetVersion().GetVendorName())
	assert.Equal(t, "Asset Link SDK", resp.GetVersion().GetProductName())
	assert.Equal(t, "driver info server", resp.GetVersion().GetProductDescription())
	assert.Equal(t, "https://example.com/docs", resp.GetVersion().GetDocuUrl())
	assert.Equal(t, "https://example.com/feedback", resp.GetVersion().GetFeedbackUrl())
}

func TestGetVersionInfoUnknownVersion(t *testing.T) {
	sut := &DriverInfoServerEntity{
		Metadata: metadata.Metadata{
			Version: metadata.Version{Version: "unknown"},
		},
	}

	resp, err := sut.GetVersionInfo(context.Background(), &generated.GetVersionInfoRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetVersion())

	assert.Equal(t, uint32(0), resp.GetVersion().GetMajor())
	assert.Equal(t, uint32(0), resp.GetVersion().GetMinor())
	assert.Equal(t, uint32(0), resp.GetVersion().GetPatch())
	assert.Equal(t, "unknown", resp.GetVersion().GetSuffix())
}

func TestGetVersionInfoInvalidVersionFallsBack(t *testing.T) {
	sut := &DriverInfoServerEntity{
		Metadata: metadata.Metadata{
			AlName:  "Asset Link SDK",
			Version: metadata.Version{Version: "not-a-semver"},
		},
	}

	resp, err := sut.GetVersionInfo(context.Background(), &generated.GetVersionInfoRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetVersion())

	assert.Equal(t, uint32(0), resp.GetVersion().GetMajor())
	assert.Equal(t, uint32(0), resp.GetVersion().GetMinor())
	assert.Equal(t, uint32(0), resp.GetVersion().GetPatch())
	assert.Equal(t, "unknown", resp.GetVersion().GetSuffix())
	assert.Equal(t, "Asset Link SDK", resp.GetVersion().GetProductName())
}
