/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	linkmlValidateImage = "linkml/linkml"
)

func RunContainer(service string) error {
	serviceDef := GetServiceDefinition(service, schemaPath, assetPath, targetClass)
	cmdArgs := []string{"run", "-d"}
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
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running docker run for %s: %w, output: %s", service, err, output.String())
	}

	fmt.Printf("Container for %s started successfully.\n", service)
	return nil
}

func GetServiceDefinition(serviceName string, schemaPath string, assetPath string, targetClass string) Service {
	currentDir, _ := os.Getwd()
	baseSchemaVersion, err := getBaseSchemaVersion()
	baseSchemaVersion = baseSchemaVersion
	baseSchemaVersion += ".yaml"
	if err != nil {
		fmt.Println(err)
		return Service{}
	}
	Services := map[string]Service{
		"linkml-validator": {
			Image: linkmlValidateImage,
			Volumes: []string{
				filepath.Join(currentDir, schemaPath) + ":/app/src/cdm/schema.yaml",
				filepath.Join(currentDir, assetPath) + ":/app/src/cdm/asset.json",
				filepath.Join(currentDir, baseSchemaPath) + fmt.Sprintf(":/app/src/cdm/%s", baseSchemaVersion),
			},
			Entrypoint: []string{
				"linkml-validate",
				"--include-range-class-descendants",
				fmt.Sprintf("--target-class=%s", targetClass),
				"-s",
				"/app/src/cdm/schema.yaml",
				"/app/src/cdm/asset.json",
			},
		},
	}
	return Services[serviceName]
}
