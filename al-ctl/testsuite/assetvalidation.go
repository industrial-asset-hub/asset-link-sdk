/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package testsuite

import (
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/al-ctl/shared"
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

type Service struct {
	Image      string   `yaml:"image"`
	Command    string   `yaml:"command,omitempty"`
	Ports      []string `yaml:"ports,omitempty"`
	Volumes    []string `yaml:"volumes,omitempty"`
	Build      string   `yaml:"build,omitempty"`
	DependsOn  []string `yaml:"depends_on,omitempty"`
	Entrypoint []string `yaml:"entrypoint,omitempty"`
}

func RunContainer(service string) error {
	serviceDef, err := GetServiceDefinition(SchemaPath, shared.AssetJsonPath, TargetClass)
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
		baseSchemaFileName = filepath.Base(BaseSchemaPath)
		AddVolumeInService(&Service, currentDir, BaseSchemaPath, baseSchemaFileName)
		AddSchemaEntrypointInService(&Service, baseSchemaFileName)
	default:
		baseSchemaFileName, err := GetBaseSchemaVersionFromExtendedSchema()
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		baseSchemaFileName += ".yaml"
		AddVolumeInService(&Service, currentDir, BaseSchemaPath, baseSchemaFileName)
		AddVolumeInService(&Service, currentDir, schemaPath, extendedSchemaFileName)
		AddSchemaEntrypointInService(&Service, extendedSchemaFileName)
	}
	AddAssetEntrypointInService(&Service, assetFileName)
	return &Service, nil
}

func AddVolumeInService(service *Service, currentDir string, pathInHost string, volume string) {
	volume = filepath.Join(currentDir, pathInHost) + fmt.Sprintf(":/app/src/cdm/%s", volume)
	service.Volumes = append(service.Volumes, volume)
}

func AddSchemaEntrypointInService(service *Service, schemaFileName string) {
	service.Entrypoint = append(service.Entrypoint, "-s", fmt.Sprintf("/app/src/cdm/%s", schemaFileName))
}

func AddAssetEntrypointInService(service *Service, assetFileName string) {
	service.Entrypoint = append(service.Entrypoint, "/app/src/cdm/"+assetFileName)
}
