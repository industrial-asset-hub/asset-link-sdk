/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

// Pattern constants for validation
const (
	MacAddressPattern         = "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
	IPv4AddressPattern        = "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	IPv6AddressPattern        = "^([0-9a-fA-F]{1,4}:){1,7}[0-9a-fA-F]{1,4}$"
	IPv6NetworkPrefixPattern  = "^[a-fA-F0-9]{1,4}:(?:[a-fA-F0-9]{0,4}:){1,7}/([0-9]{1,2}|1[01][0-9]|12[0-8])$"
	NetworkMaskPattern        = "^(255)\\.(0|128|192|224|240|248|252|254|255)\\.(0|128|192|224|240|248|252|254|255)\\.(0|128|192|224|240|248|252|254|255)$"
	RouterIPv4AddressPattern  = "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	RouterIPv6AddressPattern  = "^([0-9a-fA-F]{1,4}:){1,7}[0-9a-fA-F]{1,4}$" // Updated IPv6 patterns to be more permissive and match real-world values
	FunctionalObjectSchemaUrl = "https://industrial-assets.io/schemas/iah/base-schema/released/v1/iah-base.json"
)
