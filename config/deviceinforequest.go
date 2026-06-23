/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package config

import generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"

type DeviceInfoRequest interface {
	GetParameterJson() string
	GetCredentials() []*generated.ConnectionCredential
}
