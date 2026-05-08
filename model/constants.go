/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

// Pattern constants for validation - sourced from cdm_base.schema_v1.9.0.json
const (
	// MacAddressPattern validates MAC addresses (e.g. AC:64:17:01:1E:52 or AC-64-17-01-1E-52).
	MacAddressPattern = "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"

	// IPv4AddressPattern validates IPv4 addresses (e.g. 192.168.0.1).
	IPv4AddressPattern = "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"

	// RouterIPv4AddressPattern validates the IPv4 address of a router/gateway.
	RouterIPv4AddressPattern = "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"

	// NetworkMaskPattern validates IPv4 subnet masks (e.g. 255.255.255.0).
	NetworkMaskPattern = "^(255)\\.(0|128|192|224|240|248|252|254|255)\\.(0|128|192|224|240|248|252|254|255)\\.(0|128|192|224|240|248|252|254|255)$"

	// IPv6AddressPattern validates full (non-compressed) IPv6 addresses with exactly 8 groups.
	IPv6AddressPattern = "^(?:[a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}$"

	// IPv6NetworkPrefixPattern validates IPv6 network prefixes with CIDR notation.
	IPv6NetworkPrefixPattern = "^[a-fA-F0-9]{1,4}:(?:[a-fA-F0-9]{0,4}:){1,7}/([0-9]{1,2}|1[01][0-9]|12[0-8])$"

	// RouterIPv6AddressPattern validates the IPv6 address of a default router/gateway.
	RouterIPv6AddressPattern = "^(?:[a-fA-F0-9]{0,4}:){1,7}[a-fA-F0-9]{1,4}$"

	// FunctionalObjectSchemaUrlPattern validates the schema URL for functional objects. .
	FunctionalObjectSchemaUrlPattern = "^https://industrial-assets\\.io.*$"

	// IdLinkPattern validates ID Link URLs following IEC 61406 / DIN SPEC 91406.
	IdLinkPattern = "^https?://(?:www\\.)?[a-z0-9.-]+(?:\\.[a-z]+)?/[^\\s#]*$"

	// CustomIdentifierValuePattern validates custom identifier values (max 256 chars, URL-safe chars).
	CustomIdentifierValuePattern = "^[A-Za-z0-9._~!$&'()*+,;=:/?@%-]{1,256}$"

	// FunctionalObjectSchemaUrl is the canonical schema URL used for all IAH functional objects.
	FunctionalObjectSchemaUrl = "https://industrial-assets.io/schemas/iah/base-schema/released/v1/iah-base.json"
)
