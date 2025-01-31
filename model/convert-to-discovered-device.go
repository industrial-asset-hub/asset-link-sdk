/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
)

func ConvertFromDerivedSchemaToDiscoveredDevice[T interface{}](d *T, schemaUri string, deviceClass string) *generated.DiscoveredDevice {
	device := generated.DiscoveredDevice{
		Identifiers:            convertDeviceInfoToDeviceIdentifiers(d, schemaUri, deviceClass),
		ConnectionParameterSet: nil,
		Timestamp:              19347439483904,
	}
	return &device
}

func (d *DeviceInfo) ConvertToDiscoveredDevice() *generated.DiscoveredDevice {
	device := generated.DiscoveredDevice{
		Identifiers:            convertDeviceInfoToDeviceIdentifiers(d, baseSchemaPrefix, "Asset"),
		ConnectionParameterSet: nil,
		Timestamp:              19347439483904,
	}
	return &device
}

func convertDeviceInfoToDeviceIdentifiers[T interface{}](d *T, uri string, deviceClass string) []*generated.DeviceIdentifier {
	valueOfDeviceInfo := reflect.ValueOf(d)
	assetIdentifierUri := fmt.Sprintf("%s/%s#", uri, deviceClass)
	return convertStructTypeToDeviceIdentifiers(valueOfDeviceInfo.Elem(), assetIdentifierUri, 0)
}

func convertToDeviceIdentifiers(valueToConvert reflect.Value, prefixUri string, level int) []*generated.DeviceIdentifier {
	identifiers := []*generated.DeviceIdentifier{}
	switch valueToConvert.Kind() {
	case reflect.Ptr:
		if valueToConvert.IsNil() {
			productIdentifier := convertToDeviceIdentifier(valueToConvert.Interface(), prefixUri)
			identifiers = append(identifiers, productIdentifier)
			return identifiers
		}
		structIdentifiers := convertToDeviceIdentifiers(valueToConvert.Elem(), prefixUri, level)
		identifiers = appendDeviceIdentifiers(identifiers, structIdentifiers)
	case reflect.Struct:
		structIdentifiers := convertStructTypeToDeviceIdentifiers(valueToConvert, prefixUri, level+1)
		identifiers = appendDeviceIdentifiers(identifiers, structIdentifiers)
	case reflect.Slice:
		for index := 0; index < valueToConvert.Len(); index++ {
			sliceIdentifier := convertSliceElementToDeviceIdentifier(valueToConvert.Index(index), prefixUri, level+1)
			identifiers = appendDeviceIdentifiers(identifiers, []*generated.DeviceIdentifier{sliceIdentifier})
		}
	case reflect.String:
		identifier := convertToDeviceIdentifier(valueToConvert.String(), prefixUri)
		identifiers = appendDeviceIdentifiers(identifiers, []*generated.DeviceIdentifier{identifier})
	case reflect.Int, reflect.Float64:
		identifier := convertToDeviceIdentifier(valueToConvert.Interface(), prefixUri)
		identifiers = appendDeviceIdentifiers(identifiers, []*generated.DeviceIdentifier{identifier})
	case reflect.Interface:
		if valueToConvert.IsNil() {
			return identifiers
		}
		interfaceValue := valueToConvert.Elem()
		interfaceIdentifiers := convertToDeviceIdentifiers(interfaceValue, prefixUri, level)
		identifiers = appendDeviceIdentifiers(identifiers, interfaceIdentifiers)
	case reflect.Bool:
		identifier := convertToDeviceIdentifier(valueToConvert.Bool(), prefixUri)
		identifiers = appendDeviceIdentifiers(identifiers, []*generated.DeviceIdentifier{identifier})
	default:
		log.Warn().Msgf("Could not process value of kind %v and type %s", valueToConvert.Kind(), valueToConvert.Type())
	}
	return identifiers
}

func getUri(prefixUri string, fieldName string, level int) string {
	currentPrefixUri := prefixUri
	if fieldName == "" {
		return currentPrefixUri
	}

	if level <= 2 {
		currentPrefixUri = fmt.Sprintf("%s%s", prefixUri, fieldName)
	} else {
		currentPrefixUri = fmt.Sprintf("%s/%s", prefixUri, fieldName)
	}
	return currentPrefixUri
}

func convertSliceElementToDeviceIdentifier(sliceElementValue reflect.Value, identifierUri string, level int) *generated.DeviceIdentifier {
	return &generated.DeviceIdentifier{
		Value: &generated.DeviceIdentifier_Children{
			Children: &generated.DeviceIdentifierValueList{
				Value: convertToDeviceIdentifiers(sliceElementValue, identifierUri, level),
			},
		},
		Classifiers: []*generated.SemanticClassifier{{
			Type:  "URI",
			Value: identifierUri,
		}},
	}
}

func convertToDeviceIdentifier(value interface{}, identifierUri string) *generated.DeviceIdentifier {
	identifier := &generated.DeviceIdentifier{
		Value:       nil,
		Classifiers: nil,
	}
	switch v := value.(type) {
	case string:
		if isNonEmptyValues(value.(string)) {
			identifier.Value = &generated.DeviceIdentifier_Text{Text: v}
		}
	case *string:
		if value.(*string) != nil {
			identifier.Value = &generated.DeviceIdentifier_Text{Text: *v}
		}
	case int64:
		identifier.Value = &generated.DeviceIdentifier_Int64Value{Int64Value: v}
	case uint64:
		identifier.Value = &generated.DeviceIdentifier_Uint64Value{Uint64Value: v}
	case []byte:
		identifier.Value = &generated.DeviceIdentifier_RawData{RawData: v}
	case bool:
		// TODO: Convert to bool if datatype is available
		// For now we convert to an string, and accepting a schema violation
		// Related to, that capabilities may have a separate interface.
		identifier.Value = &generated.DeviceIdentifier_Text{Text: strconv.FormatBool(v)}
	}
	if identifier.Value == nil {
		return nil
	}

	classifier := &generated.SemanticClassifier{
		Type: "URI", Value: identifierUri,
	}
	identifier.Classifiers = []*generated.SemanticClassifier{classifier}
	return identifier
}

func appendDeviceIdentifiers(destinationIdentifiers []*generated.DeviceIdentifier, sourceIdentifiers []*generated.DeviceIdentifier) []*generated.DeviceIdentifier {
	if destinationIdentifiers == nil || sourceIdentifiers == nil {
		return destinationIdentifiers
	}
	for _, sourceIdentifier := range sourceIdentifiers {
		if sourceIdentifier != nil {
			destinationIdentifiers = append(destinationIdentifiers, sourceIdentifier)
		}
	}
	return destinationIdentifiers
}

func convertStructTypeToDeviceIdentifiers(valueStruct reflect.Value, prefixUri string, level int) []*generated.DeviceIdentifier {
	structIdentifiers := []*generated.DeviceIdentifier{}
	for index := 0; index < valueStruct.NumField(); index++ {
		field := valueStruct.Type().Field(index)
		fieldValue := valueStruct.Field(index)
		jsonTag := field.Tag.Get("json")
		parts := strings.Split(jsonTag, ",")
		fieldName := ""
		currentFieldLevel := level
		if len(parts) > 0 {
			fieldName = parts[0]
			if fieldName != "" {
				currentFieldLevel = level + 1
			}
		}
		currentFieldUri := getUri(prefixUri, fieldName, currentFieldLevel)
		elementIdentifiers := convertToDeviceIdentifiers(fieldValue, currentFieldUri, currentFieldLevel)
		structIdentifiers = appendDeviceIdentifiers(structIdentifiers, elementIdentifiers)
	}
	return structIdentifiers
}

func isNonEmptyValues(values ...string) bool {
	for _, value := range values {
		if value != "" {
			return true
		}
	}
	return false
}
