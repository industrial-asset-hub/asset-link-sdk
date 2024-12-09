/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
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

var jsonSchemaCmd = &cobra.Command{
	Use:   "json-schema",
	Short: "Validate JSON schema",
	Run:   runJsonSchemaValidation,
}

var (
	baseSchemaPath string
	schemaPath     string
	assetPath      string
	targetClass    string
	discoveryFile  string
)

func init() {
	TestCmd.AddCommand(assetsCmd)
	TestCmd.AddCommand(apiCmd)
	TestCmd.AddCommand(jsonSchemaCmd)

	assetsCmd.Flags().StringVarP(&baseSchemaPath, "base-schema-path", "b", "path/to/base/schema", "Path to the base schema YAML file")
	assetsCmd.Flags().StringVarP(&schemaPath, "schema-path", "s", "path/to/schema", "Path to the schema file")
	assetsCmd.Flags().StringVarP(&assetPath, "asset-path", "a", "path/to/asset", "Path to the asset JSON file")
	assetsCmd.Flags().StringVarP(&targetClass, "target-class", "t", "targetClass", "Target class for validation")
	jsonSchemaCmd.Flags().StringVarP(&schemaPath, "schema-path", "s", "path/to/schema", "Path to the schema file")
	jsonSchemaCmd.Flags().StringVarP(&assetPath, "asset-path", "a", "path/to/asset", "Path to the asset JSON file")
	apiCmd.Flags().StringVarP(&discoveryFile, "discovery-file", "d", "", shared.DiscoveryFileDesc)
}

func runAssetsTests(cmd *cobra.Command, args []string) {
	err := RunContainer("linkml-validator")
	if err != nil {
		log.Err(err).Msg("failed to validate asset against schema")
	}
}

func runApiTests(cmd *cobra.Command, args []string) {
	runTests(shared.AssetLinkEndpoint, discoveryFile)
}

func runJsonSchemaValidation(cmd *cobra.Command, args []string) {
	err := ValidateJsonSchema(schemaPath, assetPath)
	if err != nil {
		log.Err(err).Msg("Failed to validate JSON schema")
	}
}
