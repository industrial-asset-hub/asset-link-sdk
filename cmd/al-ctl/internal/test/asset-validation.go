/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package test

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	linkmlValidateImage           = "linkml/linkml"
	extendedSchemaFileName        = "schema.yaml"
	assetFileName                 = "asset.json"
	defaultValueForExtendedSchema = "path/to/schema"
)

func ValidateAsset(baseSchemaPath string, extendedSchemaPath string, assetJsonPath string, targetClass string) error {
	serviceDef, err := GetServiceDefinition(extendedSchemaPath, assetJsonPath, targetClass, baseSchemaPath)
	if err != nil {
		return err
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
	fmt.Println("Running command:", cmdArgs)
	cmd := exec.Command("docker", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}

func GetServiceDefinition(extendedSchemaPath string, assetPath string, targetClass string, baseSchemaPath string) (service *Service, err error) {
	currentDir, _ := os.Getwd()
	Service := Service{
		Image: linkmlValidateImage,
		Volumes: []string{
			filepath.Join(currentDir, assetPath) + ":/app/src/cdm/asset.json",
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
	volume = filepath.Join(currentDir, pathInHost) + fmt.Sprintf(":/app/src/cdm/%s", volume)
	service.Volumes = append(service.Volumes, volume)
}

func addSchemaEntrypointInService(service *Service, schemaFileName string) {
	service.Entrypoint = append(service.Entrypoint, "-s", fmt.Sprintf("/app/src/cdm/%s", schemaFileName))
}

func addAssetEntrypointInService(service *Service, assetFileName string) {
	service.Entrypoint = append(service.Entrypoint, "/app/src/cdm/"+assetFileName)
}
