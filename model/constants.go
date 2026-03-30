/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

// Pattern constants for validation
const (
	MacAddressPattern        = "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
	IPv4AddressPattern       = "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	IPv6AddressPattern       = "^([0-9a-fA-F]{1,4}:){1,7}[0-9a-fA-F]{1,4}$"            // Updated IPv6 patterns to be more permissive and match real-world values
	IPv6NetworkPrefixPattern = "^([0-9a-fA-F]{1,4}:){1,7}[0-9a-fA-F]{1,4}/[0-9]{1,2}$" // Updated IPv6 patterns to be more permissive and match real-world values
	NetworkMaskPattern       = "^(255)\\.(0|128|192|224|240|248|252|254|255)\\.(0|128|192|224|240|248|252|254|255)\\.(0|128|192|224|240|248|252|254|255)$"
	RouterIPv4AddressPattern = "^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	RouterIPv6AddressPattern = "^([0-9a-fA-F]{1,4}:){1,7}[0-9a-fA-F]{1,4}$" // Updated IPv6 patterns to be more permissive and match real-world values
	IdLinkPattern            = "^https://i\\.siemens\\.com((/1P(.*?)\\+S(.*?)(\\+23S(.*?))?(\\+30P(.*?))?)|(\\?1P=(.*?)&S=(.*?)(&23S=(.*?))?(&30P=(.*?))?))$"
)
