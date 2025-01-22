/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
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

func RunContainer(service string) error {
	serviceDef, err := GetServiceDefinition(schemaPath, shared.AssetJsonPath, targetClass)
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

func GetServiceDefinition(schemaPath string, assetPath string, targetClass string) (service *Service, err error) {
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
	switch schemaPath {
	case "", defaultValueForExtendedSchema:
		baseSchemaFileName = filepath.Base(baseSchemaPath)
		addVolumeInService(&Service, currentDir, baseSchemaPath, baseSchemaFileName)
		addSchemaEntrypointInService(&Service, baseSchemaFileName)
	default:
		baseSchemaFileName, err := getBaseSchemaVersionFromExtendedSchema()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		baseSchemaFileName += ".yaml"
		addVolumeInService(&Service, currentDir, baseSchemaPath, baseSchemaFileName)
		addVolumeInService(&Service, currentDir, schemaPath, extendedSchemaFileName)
		addSchemaEntrypointInService(&Service, extendedSchemaFileName)
	}
	addAssetEntrypointInService(&Service, assetFileName)
	return &Service, nil
}
