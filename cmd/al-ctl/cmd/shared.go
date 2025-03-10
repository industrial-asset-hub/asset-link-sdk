/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package cmd

var (
	logLevel             string
	assetPath            string
	registryEndpoint     string
	assetLinkEndpoint    string
	timeoutInSeconds     uint
	linkmlSupported      bool
	baseSchemaPath       string
	extendedSchemaPath   string
	targetClass          string
	outputFile           string
	discoveryFile        string
	registryParamsPath   string
	serviceName          string
	validateCancellation bool
	validateAsset        bool
)

const (
	DiscoveryFileDesc string = "discovery file allows the configuration of discovery filters and options (see discovery.json for an example)"
)
