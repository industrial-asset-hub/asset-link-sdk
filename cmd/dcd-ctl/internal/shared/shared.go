/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package shared

var (
	RegistryEndpoint  string
	AssetLinkEndpoint string
)

const (
	DiscoveryOptionsDesc string = `Discovery options

Key/Value: TODO Description

Operator:
    EQUAL = 0
    NOT_EQUAL = 1
    GREATER_THAN = 2
    GREATER_THAN_OR_EQUAL_TO = 3
    LESS_THAN = 4
    LESS_THAN_OR_EQUAL_TO = 5

Please be aware to use quotes on our commandline

Example options:
  - [{"key": "test", "value": "value", "operator": 1}]`
	DiscoveryFiltersDesc string = `Discovery filters

Key/Value: TODO Description

Operator:
    EQUAL = 0
    NOT_EQUAL = 1
    GREATER_THAN = 2
    GREATER_THAN_OR_EQUAL_TO = 3
    LESS_THAN = 4
    LESS_THAN_OR_EQUAL_TO = 5

Please be aware to use quotes on our commandline

Example filters:
  - [{"key": "test", "value": "value", "operator": 1}]`
)
