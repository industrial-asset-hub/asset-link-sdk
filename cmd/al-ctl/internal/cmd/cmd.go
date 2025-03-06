/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/discovery"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/registry"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/testsuite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v3/logging"
	"github.com/spf13/cobra"
)

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

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "al-ctl",
	Short: "command line interface to interact with Asset Links",
	Long: `This command line interfaces allows to interact with the so called Asset Links (AL).

This can be useful for validation purposes inside CI/CD pipelines or just
to ease development efforts.`,
}

// discoverCmd represents the discovery command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Start discovery job",
	Long:  `This command starts an discovery job and prints the result.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := discovery.Discover(assetLinkEndpoint, discoveryFile, timeoutInSeconds)
		if err != nil {
			log.Fatal().Err(err).Msg("error during discovery")
		}
		log.Trace().Str("File", outputFile).Msg("Saving to file")
		f, err := os.Create(outputFile)
		if err != nil {
			log.Fatal().Err(err).Msg("error creating file")
		}
		defer f.Close()

		asJson, _ := json.MarshalIndent(resp, "", "  ")
		_, err = f.Write(asJson)
		if err != nil {
			log.Fatal().Err(err).Msg("error during writing of the json file")
		}

	},
}

// InfoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print asset link information",
	Long:  `This command prints information on the asset link.`,
	Run: func(cmd *cobra.Command, args []string) {
		registry.PrintInfo(assetLinkEndpoint)
	},
}

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered asset links",
	Long:  `This command lists all asset links registered in the registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		registry.PrintList(registryEndpoint)
	},
}

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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("error during execution")
	}
}

func init() {
	cobra.OnInitialize(initHandlers)
	setupRootCommand()
	setupDiscoveryCommands()
	setupTestCommands()
}
func initHandlers() {
	logging.SetupLogging()
	logging.AdjustLogLevel(logLevel)
}

func setupRootCommand() {
	RootCmd.AddCommand(discoverCmd)
	RootCmd.AddCommand(InfoCmd)
	RootCmd.AddCommand(ListCmd)
	RootCmd.AddCommand(TestCmd)
	RootCmd.PersistentFlags().StringVarP(&registryEndpoint, "registry", "r", "localhost:50051", "gRPC Server Address of the Registry")
	RootCmd.PersistentFlags().StringVarP(&assetLinkEndpoint, "endpoint", "e", "localhost:8081", "gRPC Server Address of the AssetLink")
	RootCmd.PersistentFlags().UintVarP(&timeoutInSeconds, "timeout", "n", 0, "timeout in seconds (default none)")
	RootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "info",
		fmt.Sprintf("set log level. one of: %s,%s,%s,%s,%s,%s,%s",
			zerolog.TraceLevel.String(),
			zerolog.DebugLevel.String(),
			zerolog.InfoLevel.String(),
			zerolog.WarnLevel.String(),
			zerolog.ErrorLevel.String(),
			zerolog.FatalLevel.String(),
			zerolog.PanicLevel.String()))
}

func setupDiscoveryCommands() {
	discoverCmd.Flags().StringVarP(&outputFile, "output-file", "o", "result.json", "output file")
	discoverCmd.Flags().StringVarP(&discoveryFile, "discovery-file", "d", "", DiscoveryFileDesc)
}

func setupTestCommands() {
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
