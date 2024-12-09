/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package driverinfo

import (
	semver "github.com/Masterminds/semver/v3"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/conn_suite_drv_info"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/metadata"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type DriverInfoServerEntity struct {
	generated.UnimplementedDriverInfoApiServer
	Metadata metadata.Metadata
}

func (o *DriverInfoServerEntity) GetVersionInfo(c context.Context, request *generated.GetVersionInfoRequest) (*generated.GetVersionInfoResponse, error) {
	log.Info().Msg("GetVersionInfo called")
	log.Debug().Interface("Metadata", o.Metadata).Msg("Metadata")

	var major uint32 = 0
	var minor uint32 = 0
	var patch uint32 = 0
	var suffix = "unknown"
	if o.Metadata.Version.Version != "unknown" {
		parsedVersion, err := semver.NewVersion(o.Metadata.Version.Version)
		if err != nil {
			log.Err(err).Str("version string", o.Metadata.Version.Version).Msg("Parsing of Semantic Version")
		} else {
			major = uint32(parsedVersion.Major())
			minor = uint32(parsedVersion.Minor())
			patch = uint32(parsedVersion.Patch())
			suffix = parsedVersion.Prerelease()
		}
	}

	var product = o.Metadata.AlName
	var vendor = o.Metadata.Vendor
	// Currently not used.
	var docu = ""

	return &generated.GetVersionInfoResponse{Version: &generated.VersionInfo{
		Major:       major,
		Minor:       minor,
		Patch:       patch,
		Suffix:      suffix,
		VendorName:  vendor,
		ProductName: product,
		DocuUrl:     docu,
	}}, nil
}
