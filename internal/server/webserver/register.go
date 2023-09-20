//go:build webserver

/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package webserver

import (
  "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/internal/features"
)

func init() {
  features.ObservabilityFeatures().HttpObservabilityServer = true
}
