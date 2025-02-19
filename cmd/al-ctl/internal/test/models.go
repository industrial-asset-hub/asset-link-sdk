/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

type Service struct {
	Image      string   `yaml:"image"`
	Command    string   `yaml:"command,omitempty"`
	Ports      []string `yaml:"ports,omitempty"`
	Volumes    []string `yaml:"volumes,omitempty"`
	Build      string   `yaml:"build,omitempty"`
	DependsOn  []string `yaml:"depends_on,omitempty"`
	Entrypoint []string `yaml:"entrypoint,omitempty"`
}

type AssetValidationParams struct {
	AssetJsonPath      string
	BaseSchemaPath     string
	ExtendedSchemaPath string
	TargetClass        string
}

type RegistryParams struct {
	AppInstanceId   string   `json:"app_instance_id"`
	AppTypes        []string `json:"app_types"`
	DeviceSchemaUri []string `json:"driver_schema_uris"`
}
