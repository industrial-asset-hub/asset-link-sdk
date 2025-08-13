/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package test

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/test"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test suite for asset-link",
	Long: `Run tests for asset-link.
You can run tests for validation of Assets and API.
`,
}

var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Run tests for asset validation",
	Run:   runAssetsTests,
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run tests for API validation",
	Run:   runApiTests,
}

var registerCmd = &cobra.Command{
	Use:   "registration",
	Short: "Validate Registration of Asset Link",
	Run:   runRegistrationTests,
}

var (
	baseSchemaPath          string
	extendedSchemaPath      string
	targetClass             string
	discoveryFile           string
	assetJsonPath           string
	credentialPath          string
	assetValidationRequired bool

	serviceName              string
	registryJsonPath         string
	linkmlSupported          bool
	cancelValidationRequired bool
)

func init() {
	TestCmd.AddCommand(assetsCmd)
	TestCmd.AddCommand(apiCmd)
	TestCmd.AddCommand(registerCmd)
	assetsCmd.Flags().StringVarP(&baseSchemaPath, "base-schema-path", "b", "", "Path to the base schema YAML file")
	assetsCmd.Flags().StringVarP(&extendedSchemaPath, "extended-schema-path", "s", "", "Path to the extended schema YAML file")
	assetsCmd.Flags().StringVarP(&assetJsonPath, "asset-path", "a", "", "Path to the asset JSON file")
	assetsCmd.Flags().StringVarP(&targetClass, "target-class", "t", "", "Target class for validation of asset")
	assetsCmd.Flags().BoolVarP(&linkmlSupported, "linkml-is-supported", "l", false,
		"should be true if linkml is already supported in the test environment")
	apiCmd.Flags().StringVarP(&discoveryFile, "discovery-file", "d", "", shared.DiscoveryFileDesc)
	apiCmd.Flags().BoolVarP(&assetValidationRequired, "validate-asset-against-schema", "v", false,
		"should be true if discovered asset is to be validated against schema")
	apiCmd.Flags().BoolVarP(&linkmlSupported, "linkml-is-supported", "l", false,
		"should be true if linkml is already supported in the test environment")
	apiCmd.Flags().StringVarP(&baseSchemaPath, "base-schema-path", "b", "", "Path to the base schema YAML file")
	apiCmd.Flags().StringVarP(&extendedSchemaPath, "extended-schema-path", "s", "", "Path to the extended schema YAML file")
	apiCmd.Flags().StringVarP(&credentialPath, "credential-path", "p", "", "Path to the credential file")
	apiCmd.Flags().StringVarP(&targetClass, "target-class", "t", "", "Target class for validation")
	apiCmd.Flags().StringVarP(&serviceName, "service-name", "u", "",
		"Service to be validated (supported services: discovery, identifiers)")
	apiCmd.Flags().BoolVarP(&cancelValidationRequired, "cancel", "c", false, "Check cancellation of the service request")
	registerCmd.Flags().StringVarP(&registryJsonPath, "registry-json-path", "f", "", "Registration param file path")
}

func runAssetsTests(cmd *cobra.Command, args []string) {
	assetValidationParams := test.AssetValidationParams{
		BaseSchemaPath:     baseSchemaPath,
		ExtendedSchemaPath: extendedSchemaPath,
		TargetClass:        targetClass,
		AssetJsonPath:      assetJsonPath,
	}
	err := test.ValidateAsset(assetValidationParams, linkmlSupported)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to validate asset against schema")
	}
}

func runApiTests(cmd *cobra.Command, args []string) {
	assetValidationParams := test.AssetValidationParams{
		BaseSchemaPath:     baseSchemaPath,
		ExtendedSchemaPath: extendedSchemaPath,
		TargetClass:        targetClass,
	}

	testConfig := test.TestConfig{
		DiscoveryFile:           discoveryFile,
		Credential:              credentialPath,
		AssetValidationRequired: assetValidationRequired,
		LinkMLSupported:         linkmlSupported,
		AssetValidationParams:   assetValidationParams,
	}

	test.RunApiTests(serviceName, cancelValidationRequired, testConfig)
}

func runRegistrationTests(cmd *cobra.Command, args []string) {
	test.RunRegistrationTests(shared.RegistryEndpoint, registryJsonPath)
}
