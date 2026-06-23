/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/al"
	"github.com/industrial-asset-hub/asset-link-sdk/v4/cmd/al-ctl/internal/shared"
	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	generatedVariant "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

const propertyTestResultFile = "test_result.json"
const supportedPropertyTestResultFile = "supported_properties_test_result.json"

func TestGetDeviceInfoProperties(testConfig TestConfig) bool {
	log.Info().Msg("Running Test for GetPropertyValues")
	properties, err := al.GetDeviceInfoProperties(shared.AssetLinkEndpoint, testConfig.Credential)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GetPropertyValuesRequest from file")
		return false
	}
	if properties == nil {
		log.Error().Msg("get-property-values test failed")
		return false
	}

	err = al.WritePropertyResponsesFile(propertyTestResultFile, properties)
	if err != nil {
		log.Err(err).Msg("Error writing property responses to file")
		return false
	}

	err = validatePropertyResponsesFile(propertyTestResultFile, testConfig)
	if err != nil {
		log.Err(err).Str("result-file", propertyTestResultFile).Msg("Property response file schema validation failed")
		return false
	}

	return true
}

func TestGetSupportedDeviceInfoProperties(testConfig TestConfig) bool {
	log.Info().Msg("Running Test for GetSupportedProperties")

	supportedProperties, err := al.GetSupportedDeviceInfoProperties(shared.AssetLinkEndpoint, testConfig.Credential)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create GetSupportedPropertiesRequest from file")
		return false
	}
	if supportedProperties == nil {
		log.Error().Msg("get-supported-properties test failed")
		return false
	}

	err = al.WriteSupportedPropertyResponsesFile(supportedPropertyTestResultFile, supportedProperties)
	if err != nil {
		log.Err(err).Msg("Error writing supported property responses to file")
		return false
	}

	err = validateSupportedPropertyResponsesFile(supportedPropertyTestResultFile)
	if err != nil {
		log.Err(err).Str("result-file", supportedPropertyTestResultFile).Msg("Supported property response file schema validation failed")
		return false
	}

	return true
}

func validatePropertyResponsesFile(filePath string, testConfig TestConfig) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read property response file: %w", err)
	}

	var response generatedDeviceInfo.GetPropertyValuesResponse
	if err := protojson.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("unmarshal property response file to schema: %w", err)
	}

	if len(response.GetPropertyResults()) == 0 {
		return fmt.Errorf("property response file contains no property_results")
	}

	return validatePropertyResults(&response, testConfig)
}

func validateSupportedPropertyResponsesFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read supported property response file: %w", err)
	}

	var response generatedDeviceInfo.GetSupportedPropertiesResponse
	if err := protojson.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("unmarshal supported property response file to schema: %w", err)
	}

	if len(response.GetProperties()) == 0 {
		return fmt.Errorf("supported property response file contains no properties")
	}

	for index, property := range response.GetProperties() {
		if property == nil {
			return fmt.Errorf("properties[%d] is nil", index)
		}
		if property.GetKey() == "" {
			return fmt.Errorf("properties[%d] key is empty", index)
		}

		typeSet := 0
		if property.GetDatatype() != 0 {
			typeSet++
		}
		if property.GetSchemaUri() != "" {
			typeSet++
		}
		if property.GetContentType() != "" {
			typeSet++
		}

		if typeSet > 1 {
			return fmt.Errorf("properties[%d] has multiple type fields set", index)
		}
	}

	return nil
}

func validatePropertyResults(response *generatedDeviceInfo.GetPropertyValuesResponse, testConfig TestConfig) error {
	if response == nil {
		return fmt.Errorf("property response is nil")
	}

	propertyCount := 0
	errorCount := 0
	deviceModel := map[string]interface{}{}

	for index, result := range response.GetPropertyResults() {
		if err := validatePropertyResult(result, index); err != nil {
			return err
		}

		if property := result.GetProperty(); property != nil {
			propertyCount++

			if testConfig.AssetValidationRequired {
				value, err := shared.VariantToInterface(property.GetValue(), fmt.Sprintf("property[%d].%s", index, property.GetKey()))
				if err != nil {
					return err
				}
				deviceModel[property.GetKey()] = value
			}
		}

		if result.GetError() != nil {
			errorCount++
		}
	}

	if propertyCount == 0 {
		return fmt.Errorf("no successful properties found")
	}

	if errorCount > 0 {
		log.Warn().Msgf("GetPropertyValues returned %d property error entries", errorCount)
	}

	if testConfig.AssetValidationRequired {
		if err := validateAgainstAssetSchema(deviceModel, propertyCount, testConfig); err != nil {
			return err
		}
	}

	return nil
}

func validateAgainstAssetSchema(deviceModel map[string]interface{}, propertyCount int, testConfig TestConfig) error {
	if propertyCount == 0 {
		return fmt.Errorf("no successful properties available for base schema validation")
	}

	payload, err := json.Marshal(deviceModel)
	if err != nil {
		return fmt.Errorf("marshal properties for base schema validation: %w", err)
	}

	tempAssetFile, err := os.CreateTemp("", "deviceinfo-asset-*.json")
	if err != nil {
		return fmt.Errorf("create temp file for schema validation: %w", err)
	}
	defer os.Remove(tempAssetFile.Name())

	if _, err := tempAssetFile.Write(payload); err != nil {
		tempAssetFile.Close()
		return fmt.Errorf("write temp file for schema validation: %w", err)
	}
	if err := tempAssetFile.Close(); err != nil {
		return fmt.Errorf("close temp file for schema validation: %w", err)
	}

	assetValidationParams := testConfig.AssetValidationParams
	assetValidationParams.AssetJsonPath = tempAssetFile.Name()

	if assetValidationParams.TargetClass == "" {
		functionalObjectType, _ := deviceModel["functional_object_type"].(string)
		if functionalObjectType == "" {
			return fmt.Errorf("asset schema validation failed: missing target class and functional_object_type")
		}
		assetValidationParams.TargetClass = functionalObjectType
	}

	if err := ValidateAsset(assetValidationParams, testConfig.LinkMLSupported); err != nil {
		return fmt.Errorf("asset schema validation failed: %w", err)
	}

	return nil
}

func validatePropertyResult(result *generatedDeviceInfo.PropertyValueResult, index int) error {
	if result == nil {
		return fmt.Errorf("result[%d] is nil", index)
	}

	property := result.GetProperty()
	propertyErr := result.GetError()

	if (property == nil && propertyErr == nil) || (property != nil && propertyErr != nil) {
		return fmt.Errorf("result[%d] must contain exactly one of property or error", index)
	}

	if propertyErr != nil {
		if propertyErr.GetKey() == "" {
			return fmt.Errorf("result[%d] error key is empty", index)
		}
		if propertyErr.GetDescription() == "" {
			return fmt.Errorf("result[%d] error description is empty", index)
		}
		return nil
	}

	if property.GetKey() == "" {
		return fmt.Errorf("result[%d] property key is empty", index)
	}
	if property.GetValue() == nil {
		return fmt.Errorf("result[%d] property value is nil for key %q", index, property.GetKey())
	}

	return validateVariantValue(property.GetValue(), fmt.Sprintf("result[%d].property[%s]", index, property.GetKey()))
}

func validateVariantValue(value *generatedVariant.Variant, path string) error {
	if value == nil {
		return fmt.Errorf("%s is nil", path)
	}

	switch typed := value.GetValue().(type) {
	case nil:
		return fmt.Errorf("%s has no value type set", path)
	case *generatedVariant.Variant_StructValue:
		if typed.StructValue == nil {
			return fmt.Errorf("%s struct_value is nil", path)
		}
		for key, child := range typed.StructValue.GetFields() {
			if key == "" {
				return fmt.Errorf("%s has empty struct field key", path)
			}
			if err := validateVariantValue(child, fmt.Sprintf("%s.%s", path, key)); err != nil {
				return err
			}
		}
	case *generatedVariant.Variant_ArrayValue:
		if typed.ArrayValue == nil {
			return fmt.Errorf("%s array_value is nil", path)
		}
		for i, child := range typed.ArrayValue.GetValues() {
			if err := validateVariantValue(child, fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	case *generatedVariant.Variant_File:
		if typed.File == nil {
			return fmt.Errorf("%s file is nil", path)
		}
	}

	return nil
}
