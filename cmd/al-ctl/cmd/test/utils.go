/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package test

import (
	"encoding/json"
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/internal/shared"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func getBaseSchemaVersionFromExtendedSchema() (string, error) {
	file, err := os.Open(schemaPath)
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

func transformSemanticIdentifierToAsset() error {
	testDevice := map[string]interface{}{}
	var discoveryResponse generated.DiscoverResponse
	// Read the asset JSON file
	file, err := os.Open(shared.AssetJsonPath)
	if err != nil {
		log.Err(err).Msg("failed to read asset json file")
		return err
	}
	// unmarshal the asset JSON file
	fileInfo, err := file.Stat()
	if err != nil {
		log.Err(err).Msg("failed to read asset json file structure")
		return err
	}
	byteBuffer := make([]byte, fileInfo.Size())
	_, _ = file.Read(byteBuffer)
	unmarshalOptions := protojson.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}
	if err := unmarshalOptions.Unmarshal(byteBuffer, &discoveryResponse); err != nil {
		log.Err(err).Msg("failed to unmarshal asset")
		return err
	}
	// Transform the semantic-identifiers to asset
	testDevice = shared.TransformDevice(discoveryResponse.Devices[0], "URI")
	file, err = os.Create("test.json")
	if err != nil {
		log.Err(err).Msg("failed to create test json file")
		return err
	}
	// change the path to the asset json file
	arr := strings.Split(shared.AssetJsonPath, "/")
	newArr := arr[:len(arr)-1]
	shared.AssetJsonPath = strings.Join(newArr, "/")
	shared.AssetJsonPath = "/test.json"
	// Write the transformed asset to a file
	jsonWriter := json.NewEncoder(file)
	if err := jsonWriter.Encode(testDevice); err != nil {
		log.Err(err).Msg("failed to write transformed asset to file")
		return err
	}
	return nil
}
