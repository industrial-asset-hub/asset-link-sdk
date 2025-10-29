/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package config

import generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"

type IdentifiersRequest interface {
	GetParameterJson() string
	GetCredentials() []*generated.ConnectionCredential
}
