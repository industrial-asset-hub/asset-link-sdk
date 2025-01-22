/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package shared

var (
	RegistryEndpoint        string
	AssetLinkEndpoint       string
	AssetJsonPath           string
	AssetValidationRequired bool
)

const (
	DiscoveryFileDesc string = "discovery file allows the configuration of discovery filters and options (see discovery.json for an example)"
)
