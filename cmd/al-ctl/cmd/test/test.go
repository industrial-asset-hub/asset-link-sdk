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
	"os"
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

var (
	baseSchemaPath          string
	extendedSchemaPath      string
	targetClass             string
	discoveryFile           string
	assetJsonPath           string
	assetValidationRequired bool
)

func init() {
	TestCmd.AddCommand(assetsCmd)
	TestCmd.AddCommand(apiCmd)

	assetsCmd.Flags().StringVarP(&baseSchemaPath, "base-schema-path", "b", "", "Path to the base schema YAML file")
	assetsCmd.Flags().StringVarP(&extendedSchemaPath, "extended-schema-path", "s", "", "Path to the extended schema YAML file")
	assetsCmd.Flags().StringVarP(&assetJsonPath, "asset-path", "a", "", "Path to the asset JSON file")
	assetsCmd.Flags().StringVarP(&targetClass, "target-class", "t", "", "Target class for validation of asset")
	apiCmd.Flags().StringVarP(&discoveryFile, "discovery-file", "d", "", shared.DiscoveryFileDesc)
	apiCmd.Flags().BoolVarP(&assetValidationRequired, "validate-asset-against-schema", "v", false,
		"should be true if discovered asset is to be validated against schema")
	apiCmd.Flags().StringVarP(&baseSchemaPath, "base-schema-path", "b", "", "Path to the base schema YAML file")
	apiCmd.Flags().StringVarP(&extendedSchemaPath, "extended-schema-path", "s", "", "Path to the extended schema YAML file")
	apiCmd.Flags().StringVarP(&targetClass, "target-class", "t", "", "Target class for validation")
}

func runAssetsTests(cmd *cobra.Command, args []string) {
	err := test.ValidateAsset(baseSchemaPath, extendedSchemaPath, assetJsonPath, targetClass)
	if err != nil {
		log.Err(err).Msg("failed to validate asset against schema")
		os.Exit(1)
	}
}

func runApiTests(cmd *cobra.Command, args []string) {
	test.RunApiMockTests(shared.AssetLinkEndpoint, discoveryFile, assetValidationRequired, baseSchemaPath, extendedSchemaPath, targetClass)
}
