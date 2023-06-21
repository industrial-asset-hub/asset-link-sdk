/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package status

import (
  generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/generated/status"
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/metadata"
  "github.com/rs/zerolog/log"
  "golang.org/x/net/context"
)

type StatusServerEntity struct {
  generated.UnimplementedDcdStatusServer
  Version metadata.Version
}

func (o *StatusServerEntity) GetHealth(ctx context.Context, req *generated.HealthRequest) (*generated.HealthReply, error) {
  log.Info().Msg("GetHealth called")

  return &generated.HealthReply{Health: generated.HealthReply_HEALTHY}, nil
}
func (o *StatusServerEntity) GetVersion(ctx context.Context, req *generated.VersionRequest) (*generated.VersionReply, error) {
  log.Info().Msg("GetVersion called")
  log.Info().Interface("s", o.Version).Msg("--")
  return &generated.VersionReply{Version: o.Version.Version, Commit: o.Version.Commit, Date: o.Version.Date}, nil
}
