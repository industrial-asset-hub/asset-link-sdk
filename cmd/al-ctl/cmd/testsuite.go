/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package cmd

import (
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/testsuite"
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
	Run: func(cmd *cobra.Command, args []string) {
		params := testsuite.AssetValidationParams{
			LinkmlSupported:    linkmlSupported,
			ExtendedSchemaPath: extendedSchemaPath,
			BaseSchemaPath:     baseSchemaPath,
			AssetJsonPath:      assetPath,
			TargetClass:        targetClass,
		}
		testsuite.RunAssetValidation(params)
	},
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run tests for API validation",
	Run: func(cmd *cobra.Command, args []string) {
		params := testsuite.ApiValidationParams{
			DiscoveryFile:                  discoveryFile,
			AssetLinkEndpoint:              assetLinkEndpoint,
			CancellationValidationRequired: validateCancellation,
			ServiceName:                    serviceName,
			TimeoutInSeconds:               timeoutInSeconds,
			AssetValidationRequired:        validateAsset,
			AssetValidationParams: testsuite.AssetValidationParams{
				LinkmlSupported:    linkmlSupported,
				ExtendedSchemaPath: extendedSchemaPath,
				BaseSchemaPath:     baseSchemaPath,
				AssetJsonPath:      assetPath,
				TargetClass:        targetClass,
			},
		}
		testsuite.RunAPIValidation(params)
	},
}

var registerCmd = &cobra.Command{
	Use:   "registration",
	Short: "Validate Registration of Asset Link",
	Run: func(cmd *cobra.Command, args []string) {
		params := testsuite.RegistryValidationParams{
			GrpcEndpoint:      registryEndpoint,
			RegistryJsonPath:  registryParamsPath,
			AssetLinkEndpoint: assetLinkEndpoint,
		}
		testsuite.RunRegistrationValidation(params)
	},
}

func init() {
	apiCmd.Flags().BoolVarP(&linkmlSupported, "linkml-is-supported", "l", false,
		"should be true if linkml is already supported in the test environment")
	// api validation flags
	apiCmd.Flags().BoolVarP(&validateAsset, "validate-asset-against-schema", "v", false,
		"should be true if discovered asset is to be validated against schema")
	apiCmd.Flags().StringVarP(&serviceName, "service-name", "u", "", "Service to be valdiated (supported services: discovery)")
	apiCmd.Flags().BoolVarP(&validateCancellation, "cancel", "c", false, "Check cancellation of the service request")
	apiCmd.Flags().StringVarP(&baseSchemaPath, "base-schema-path", "b", "", "Path to the base schema YAML file")
	apiCmd.Flags().StringVarP(&extendedSchemaPath, "extended-schema-path", "s", "", "Path to the extended schema YAML file")
	apiCmd.Flags().StringVarP(&targetClass, "target-class", "t", "", "Target class for validation of asset")
	// asset validation flags
	assetsCmd.Flags().BoolVarP(&linkmlSupported, "linkml-is-supported", "l", false,
		"should be true if linkml is already supported in the test environment")
	assetsCmd.Flags().StringVarP(&assetPath, "asset-path", "a", "", "Path to the asset JSON file")
	assetsCmd.Flags().StringVarP(&baseSchemaPath, "base-schema-path", "b", "", "Path to the base schema YAML file")
	assetsCmd.Flags().StringVarP(&extendedSchemaPath, "extended-schema-path", "s", "", "Path to the extended schema YAML file")
	assetsCmd.Flags().StringVarP(&targetClass, "target-class", "t", "", "Target class for validation of asset")
	// registration validation flags
	registerCmd.Flags().StringVarP(&registryParamsPath, "registry-json-path", "f", "", "Registration param file path")
	TestCmd.AddCommand(assetsCmd)
	TestCmd.AddCommand(apiCmd)
	TestCmd.AddCommand(registerCmd)
}
