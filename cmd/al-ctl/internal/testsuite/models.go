/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package testsuite

type Service struct {
	Image      string   `yaml:"image"`
	Command    string   `yaml:"command,omitempty"`
	Ports      []string `yaml:"ports,omitempty"`
	Volumes    []string `yaml:"volumes,omitempty"`
	Build      string   `yaml:"build,omitempty"`
	DependsOn  []string `yaml:"depends_on,omitempty"`
	Entrypoint []string `yaml:"entrypoint,omitempty"`
}

type RegistryFileParams struct {
	GrpcAddress     string   `json:"grpc_address"`
	AppInstanceId   string   `json:"app_instance_id"`
	AppTypes        []string `json:"app_types"`
	DeviceSchemaUri []string `json:"driver_schema_uris"`
}

type AssetValidationParams struct {
	LinkmlSupported    bool
	ExtendedSchemaPath string
	BaseSchemaPath     string
	AssetJsonPath      string
	TargetClass        string
}

type ApiValidationParams struct {
	DiscoveryFile                  string
	CancellationValidationRequired bool
	AssetValidationRequired        bool
	AssetLinkEndpoint              string
	AssetValidationParams
	ServiceName      string
	TimeoutInSeconds uint
}

type RegistryValidationParams struct {
	GrpcEndpoint      string
	RegistryJsonPath  string
	AssetLinkEndpoint string
}
