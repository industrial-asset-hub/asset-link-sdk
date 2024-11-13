//go:build webserver

/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package webserver

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v2/internal/features"
)

func init() {
	features.ObservabilityFeatures().HttpObservabilityServer = true
}
