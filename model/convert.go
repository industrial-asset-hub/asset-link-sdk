package model

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"strings"

	generated "code.siemens.com/common-device-management/device-class-drivers/cdm-dcd-sdk/v2/generated/iah-discovery"
)

func (d *DeviceInfo) ConvertToDiscoveredDevice() *generated.DiscoveredDevice {
	device := generated.DiscoveredDevice{
		Identifiers:            convertDeviceInfoToDeviceIdentifiers(d),
		ConnectionParameterSet: nil,
		Timestamp:              19347439483904,
	}
	return &device
}

func convertDeviceInfoToDeviceIdentifiers(d *DeviceInfo) []*generated.DeviceIdentifier {
	valueOfDeviceInfo := reflect.ValueOf(d)
	assetIdentifierUri := fmt.Sprintf("%s/Asset#", baseSchemaPrefix)
	return convertStructTypeToDeviceIdentifiers(valueOfDeviceInfo.Elem(), assetIdentifierUri, 0)
}

func convertToDeviceIdentifiers(valueToConvert reflect.Value, prefixUri string, level int) []*generated.DeviceIdentifier {
	identifiers := []*generated.DeviceIdentifier{}
	switch valueToConvert.Kind() {
	case reflect.Ptr:
		if valueToConvert.IsNil() {
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
	case reflect.Int, reflect.Float64, reflect.Bool:
		identifier := convertToDeviceIdentifier(valueToConvert.Interface(), prefixUri)
		identifiers = appendDeviceIdentifiers(identifiers, []*generated.DeviceIdentifier{identifier})
	default:
		log.Warn().Msgf(fmt.Sprintf("Coudn't process value of kind %v and type %s", valueToConvert.Kind(), valueToConvert.Type()))
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
	switch value.(type) {
	case string:
		if value.(string) != "" {
			identifier.Value = &generated.DeviceIdentifier_Text{Text: value.(string)}
		}
	case *string:
		if *value.(*string) != "" {
			identifier.Value = &generated.DeviceIdentifier_Text{Text: *value.(*string)}
		}
	case int64:
		identifier.Value = &generated.DeviceIdentifier_Int64Value{Int64Value: value.(int64)}
	case uint64:
		identifier.Value = &generated.DeviceIdentifier_Uint64Value{Uint64Value: value.(uint64)}
	case []byte:
		identifier.Value = &generated.DeviceIdentifier_RawData{RawData: value.([]byte)}
	}
	if identifier.Value == nil {
		return nil
	}
	classifier := &generated.SemanticClassifier{
		Type:  "URI",
		Value: identifierUri,
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
