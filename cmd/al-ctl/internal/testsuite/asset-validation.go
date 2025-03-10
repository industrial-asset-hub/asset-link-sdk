/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package testsuite

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	linkmlValidateImage           = "linkml/linkml"
	extendedSchemaFileName        = "schema.yaml"
	dockerAssetFileMapping        = "asset.json"
	defaultValueForExtendedSchema = "path/to/schema"
)

func ValidateAsset(params AssetValidationParams) error {
	var cmd *exec.Cmd
	if params.LinkmlSupported {
		if params.ExtendedSchemaPath == "" {
			params.ExtendedSchemaPath = params.BaseSchemaPath
		}
		cmd = exec.Command("linkml-validate", params.AssetJsonPath, "--include-range-class-descendants",
			"--target-class="+params.TargetClass, "-s", params.ExtendedSchemaPath)
	} else {
		var err error
		cmd, err = getDockerCommand(params.AssetJsonPath, params.TargetClass, params.ExtendedSchemaPath, params.BaseSchemaPath)
		if err != nil {
			log.Err(err).Msg("failed to get docker command")
			return err
		}
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Info().Msgf("Running command: %s", cmd.Args)
	err := cmd.Run()
	return err
}

func GetServiceDefinition(assetJsonPath string, targetClass string, extendedSchemaPath string, baseSchemaPath string) (service *Service, err error) {
	currentDir, _ := os.Getwd()
	Service := Service{
		Image: linkmlValidateImage,
		Volumes: []string{
			filepath.Join(currentDir, assetJsonPath) + ":/app/src/cdm/asset.json",
		},
		Entrypoint: []string{
			"linkml-validate",
			"--include-range-class-descendants",
			"-D",
			fmt.Sprintf("--target-class=%s", targetClass),
		},
	}
	var baseSchemaFileName string
	switch extendedSchemaPath {
	case "", defaultValueForExtendedSchema:
		baseSchemaFileName = filepath.Base(baseSchemaPath)
		addVolumeInService(&Service, currentDir, baseSchemaPath, baseSchemaFileName)
		addSchemaEntrypointInService(&Service, baseSchemaFileName)
	default:
		baseSchemaFileName, err := getBaseSchemaVersionFromExtendedSchema(extendedSchemaPath)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		baseSchemaFileName += ".yaml"
		addVolumeInService(&Service, currentDir, baseSchemaPath, baseSchemaFileName)
		addVolumeInService(&Service, currentDir, extendedSchemaPath, extendedSchemaFileName)
		addSchemaEntrypointInService(&Service, extendedSchemaFileName)
	}
	addAssetEntrypointInService(&Service, dockerAssetFileMapping)
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
	volume = filepath.Join(currentDir, pathInHost) + fmt.Sprintf(":/app/src/cdm/%s", volume)
	service.Volumes = append(service.Volumes, volume)
}

func addSchemaEntrypointInService(service *Service, schemaFileName string) {
	service.Entrypoint = append(service.Entrypoint, "-s", fmt.Sprintf("/app/src/cdm/%s", schemaFileName))
}

func addAssetEntrypointInService(service *Service, assetFileName string) {
	service.Entrypoint = append(service.Entrypoint, "/app/src/cdm/"+assetFileName)
}

func getDockerCommand(assetJsonPath string, targetClass string, extendedSchemaPath string, baseSchemaPath string) (cmd *exec.Cmd, err error) {
	serviceDef, err := GetServiceDefinition(assetJsonPath, targetClass, extendedSchemaPath, baseSchemaPath)
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
