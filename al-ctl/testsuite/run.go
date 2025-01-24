/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package testsuite

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/al-ctl/shared"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var TestCmd = &cobra.Command{
	Use:   "testsuite",
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

func init() {
	TestCmd.AddCommand(assetsCmd)
	TestCmd.AddCommand(apiCmd)

	assetsCmd.Flags().StringVarP(&BaseSchemaPath, "base-schema-path", "b", "", "Path to the base schema YAML file")
	assetsCmd.Flags().StringVarP(&SchemaPath, "extended-schema-path", "s", "", "Path to the extended schema YAML file")
	assetsCmd.Flags().StringVarP(&assetPath, "asset-path", "a", "", "Path to the asset JSON file")
	assetsCmd.Flags().BoolVarP(&semanticIdentifierInputType, "semantic-identifier-input-type", "i", false,
		"should be true if asset input is of type semantic identifiers")
	assetsCmd.Flags().StringVarP(&TargetClass, "target-class", "t", "", "Target class for validation of asset")
	apiCmd.Flags().StringVarP(&discoveryFile, "discovery-file", "d", "", shared.DiscoveryFileDesc)
	apiCmd.Flags().BoolVarP(&shared.AssetValidationRequired, "validate-asset-against-schema", "v", false,
		"should be true if discovered asset is to be validated against schema")
	apiCmd.Flags().StringVarP(&BaseSchemaPath, "base-schema-path", "b", "", "Path to the base schema YAML file")
	apiCmd.Flags().StringVarP(&SchemaPath, "extended-schema-path", "s", "", "Path to the extended schema YAML file")
	apiCmd.Flags().StringVarP(&TargetClass, "target-class", "t", "", "Target class for validation")
}

func runAssetsTests(cmd *cobra.Command, args []string) {
	if semanticIdentifierInputType {
		err := transformSemanticIdentifierToAsset()
		if err != nil {
			log.Err(err).Msg("failed to transform semantic identifier to asset")
			return
		}
	}

	err := RunContainer("linkml-validator")
	if err != nil {
		log.Err(err).Msg("failed to validate asset against schema")
	}
}

func runApiTests(cmd *cobra.Command, args []string) {
	RunApiMockTests(shared.AssetLinkEndpoint, discoveryFile)
	if shared.AssetValidationRequired {
		runAssetsTests(cmd, args)
	}

}
