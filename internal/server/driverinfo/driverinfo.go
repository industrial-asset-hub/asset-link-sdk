/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package driverinfo

import (
  generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/conn_suite_drv_info"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"
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

  var major uint32 = 1
  var minor uint32 = 2
  var patch uint32 = 3
  var suffix = "unknown"
  if o.Metadata.Version.Version != "unknown" {
    // @TODO: Split semver from version
  }

  var product = o.Metadata.DcdName
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
