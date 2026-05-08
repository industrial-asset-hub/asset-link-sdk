/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	linkmlValidateImage           = "linkml/linkml"
	extendedSchemaFileName        = "schema.yaml"
	assetFileName                 = "asset.json"
	linkmlConfigFileName          = "linkml-validation-config.yaml"
	linkmlConfigContainerFilePath = "/app/src/cdm/" + linkmlConfigFileName
	linkmlValidationConfigContent = "include_range_class_descendants: true\n"
	defaultValueForExtendedSchema = "path/to/schema"
)

func ValidateAsset(assetValidationParams AssetValidationParams, linkmlSupported bool) error {
	var cmd *exec.Cmd
	linkmlConfigPath, err := createLinkMLValidationConfigFile()
	if err != nil {
		return err
	}
	defer os.Remove(linkmlConfigPath)

	if linkmlSupported {
		if assetValidationParams.ExtendedSchemaPath == "" {
			assetValidationParams.ExtendedSchemaPath = assetValidationParams.BaseSchemaPath
		}
		cmd = exec.Command(
			"linkml-validate",
			"--config", linkmlConfigPath,
			assetValidationParams.AssetJsonPath,
			"--target-class="+assetValidationParams.TargetClass,
			"-s", assetValidationParams.ExtendedSchemaPath,
		)
	} else {
		cmd, err = getDockerCommand(assetValidationParams, linkmlConfigPath)
		if err != nil {
			log.Err(err).Msg("failed to get docker command")
			return err
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Info().Msgf("Running command: %s", cmd.Args)
	err = cmd.Run()
	return err
}

func createLinkMLValidationConfigFile() (string, error) {
	configFile, err := os.CreateTemp("", "linkml-validation-config-*.yaml")
	if err != nil {
		return "", err
	}
	defer configFile.Close()

	if _, err := configFile.WriteString(linkmlValidationConfigContent); err != nil {
		return "", err
	}

	return configFile.Name(), nil
}

func GetServiceDefinition(assetValidationParams AssetValidationParams, linkmlConfigPath string) (service *Service, err error) {
	currentDir, _ := os.Getwd()
	Service := Service{
		Image: linkmlValidateImage,
		Volumes: []string{
			getHostPath(currentDir, assetValidationParams.AssetJsonPath) + ":/app/src/cdm/asset.json",
			getHostPath(currentDir, linkmlConfigPath) + ":" + linkmlConfigContainerFilePath,
		},
		Entrypoint: []string{
			"linkml-validate",
			"-D",
			"--config",
			linkmlConfigContainerFilePath,
			fmt.Sprintf("--target-class=%s", assetValidationParams.TargetClass),
		},
	}
	var baseSchemaFileName string
	switch assetValidationParams.ExtendedSchemaPath {
	case "", defaultValueForExtendedSchema:
		baseSchemaFileName = filepath.Base(assetValidationParams.BaseSchemaPath)
		addVolumeInService(&Service, currentDir, assetValidationParams.BaseSchemaPath, baseSchemaFileName)
		addSchemaEntrypointInService(&Service, baseSchemaFileName)
	default:
		baseSchemaFileName, err := getBaseSchemaVersionFromExtendedSchema(assetValidationParams.ExtendedSchemaPath)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		baseSchemaFileName += ".yaml"
		addVolumeInService(&Service, currentDir, assetValidationParams.BaseSchemaPath, baseSchemaFileName)
		addVolumeInService(&Service, currentDir, assetValidationParams.ExtendedSchemaPath, extendedSchemaFileName)
		addSchemaEntrypointInService(&Service, extendedSchemaFileName)
	}
	addAssetEntrypointInService(&Service, assetFileName)
	return &Service, nil
}

func getBaseSchemaVersionFromExtendedSchema(extendedSchemaPath string) (string, error) {
	file, err := os.Open(extendedSchemaPath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var data map[string]interface{}

	if err := decoder.Decode(&data); err != nil {
		fmt.Println(err)
		return "", err
	}
	imports, ok := data["imports"].([]interface{})
	if !ok {
		return "", fmt.Errorf("imports not found in extended schema")
	}
	if len(imports) < 2 {
		return "", fmt.Errorf("reference to base schema not found in extended schema")
	}
	baseSchemaVersion := imports[1].(string)
	return baseSchemaVersion, nil
}

func addVolumeInService(service *Service, currentDir string, pathInHost string, volume string) {
	volume = getHostPath(currentDir, pathInHost) + fmt.Sprintf(":/app/src/cdm/%s", volume)
	service.Volumes = append(service.Volumes, volume)
}

func getHostPath(currentDir string, pathInHost string) string {
	if filepath.IsAbs(pathInHost) {
		return pathInHost
	}
	return filepath.Join(currentDir, pathInHost)
}

func addSchemaEntrypointInService(service *Service, schemaFileName string) {
	service.Entrypoint = append(service.Entrypoint, "-s", fmt.Sprintf("/app/src/cdm/%s", schemaFileName))
}

func addAssetEntrypointInService(service *Service, assetFileName string) {
	service.Entrypoint = append(service.Entrypoint, "/app/src/cdm/"+assetFileName)
}

func getDockerCommand(assetValidationParams AssetValidationParams, linkmlConfigPath string) (cmd *exec.Cmd, err error) {
	serviceDef, err := GetServiceDefinition(assetValidationParams, linkmlConfigPath)
	if err != nil {
		return nil, err
	}
	cmdArgs := []string{"run", "-i"}
	for _, port := range serviceDef.Ports {
		cmdArgs = append(cmdArgs, "-p", port)
	}
	for _, volume := range serviceDef.Volumes {
		cmdArgs = append(cmdArgs, "-v", volume)
	}
	if len(serviceDef.Entrypoint) > 0 {
		cmdArgs = append(cmdArgs, "--entrypoint", serviceDef.Entrypoint[0])
	}
	cmdArgs = append(cmdArgs, serviceDef.Image)
	if len(serviceDef.Entrypoint) > 1 {
		cmdArgs = append(cmdArgs, serviceDef.Entrypoint[1:]...)
	}
	return exec.Command("docker", cmdArgs...), err
}
