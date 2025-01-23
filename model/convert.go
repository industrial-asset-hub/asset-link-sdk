/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */
package model

import (
	"fmt"
	"github.com/industrial-asset-hub/asset-link-sdk/v3/cmd/al-ctl/cmd/test"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
)

const (
	parentKeyIndex = 1
	prefix         = "https://schema.industrial-assets.io/"
)

type classifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Children struct {
	Value []identifier `json:"value"`
}

type identifier struct {
	Text       *string      `json:"text,omitempty"`
	Value      string       `json:"value"`
	Classifier []classifier `json:"classifiers"`
	Children   *Children    `json:"children,omitempty"`
}

type ArrayContext struct {
	IsInArray  bool
	ArrayIndex int
}

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

func filter(identifier *generated.DeviceIdentifier, expectedType string, prefix string) *generated.DeviceIdentifier {
	if hasClassifier(identifier, expectedType, prefix) {
		removeUnsupportedClassifier(identifier, expectedType, prefix)
	} else {
		return nil
	}
	return identifier
}

func hasClassifier(id *generated.DeviceIdentifier, expectedType string, prefix string) bool {
	for _, class := range id.Classifiers {
		if match(class, expectedType, prefix) {
			return true
		}
	}
	return false
}

func removeUnsupportedClassifier(id *generated.DeviceIdentifier, expectedType string, prefix string) {
	found := false
	var index int
	for i, class := range id.Classifiers {
		if match(class, expectedType, prefix) {
			found = true
			index = i
			break
		}
	}
	if found {
		id.Classifiers = id.Classifiers[index : index+1]
	}
}

func match(class *generated.SemanticClassifier, expectedType string, prefix string) bool {
	return (strings.HasPrefix(class.Value, prefix) && class.Type == expectedType)
}

func transformKeys(keys []string, text any, transformedDevice map[string]interface{}) {
	if len(keys) == 0 {
		return
	}

	if len(keys) == 1 {
		transformedDevice[keys[0]] = text
		return
	}
	property := make(map[string]interface{})
	existingValue, ok := transformedDevice[keys[0]]
	if ok {
		property = existingValue.(map[string]interface{})
	}
	transformKeys(keys[1:], text, property)
	transformedDevice[keys[0]] = property
}

func convertTimestampToRFC339(timestamp int64) string {
	timestamp /= 1e7
	timestamp -= 11644473600
	t := time.Unix(timestamp, 0)
	formattedTime := t.Format(time.RFC3339)
	return formattedTime
}

func retrieveAssetTypeFromDiscoveredDevice(device *generated.DiscoveredDevice) string {
	// assumption: all classifiers within a devices identifier are for the same asset type/asset class
	// for current URI-type implementations of this schema conversion the format is as follows:
	// https://schema.industrial-assets.io/<concrete-schema-and-version-info>/<asset-type>#<property-path>
	if len(device.Identifiers) == 0 {
		// This one has no identifiers
		// default to "Asset" as type
		return "Asset"
	}

	firstIdentifier := device.Identifiers[0]
	firstIdentifier = filter(firstIdentifier, "URI", prefix)
	if firstIdentifier == nil || len(firstIdentifier.Classifiers) == 0 {
		// This one has no supported classifiers
		// default to "Asset" as type
		return "Asset"
	}
	firstClassifierThatIsSupported := firstIdentifier.Classifiers[0]
	indexOfPropertyPathIndicator := strings.LastIndex(firstClassifierThatIsSupported.Value, "#")
	if indexOfPropertyPathIndicator == -1 {
		// This one has no supported classifiers
		// default to "Asset" as type
		return "Asset"
	}
	classifierWithoutPropertyPath := firstClassifierThatIsSupported.Value[:indexOfPropertyPathIndicator]
	lastIndexOfSlash := strings.LastIndex(classifierWithoutPropertyPath, "/")
	if lastIndexOfSlash == -1 {
		// This one has no supported classifiers
		// default to "Asset" as type
		return "Asset"
	}
	assetType := classifierWithoutPropertyPath[lastIndexOfSlash+1:]
	return assetType
}

func TransformDevice(device *generated.DiscoveredDevice, expectedType string) map[string]interface{} {
	DeviceInIahSchema := make(map[string]interface{})
	DeviceInIahSchema["@type"] = retrieveAssetTypeFromDiscoveredDevice(device)
	baseSchemaVersion, err := test.GetBaseSchemaVersionFromExtendedSchema()
	if err != nil {
		log.Warn().Msg("base schema version not found in extended schema")
		baseSchemaVersion = "v0.9.0"
	}
	DeviceInIahSchema["@context"] = map[string]interface{}{
		"base":      fmt.Sprintf("https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/%s/", baseSchemaVersion),
		"linkml":    "https://w3id.org/linkml/",
		"lis":       "http://rds.posccaesar.org/ontology/lis14/rdl/",
		"schemaorg": "https://schema.org/",
		"skos":      "http://www.w3.org/2004/02/skos/core#",
		"@vocab":    "https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/v0.7.5/",
	}
	timestamp := device.Timestamp
	formattedTimestamp := convertTimestampToRFC339(int64(timestamp))
	DeviceInIahSchema["management_state"] = map[string]interface{}{
		"state_value":     "unknown",
		"state_timestamp": formattedTimestamp,
	}
	mapPropertiesIntoDevice(device.Identifiers, DeviceInIahSchema, parentKeyIndex, expectedType, prefix)
	return DeviceInIahSchema
}

func mapPropertiesIntoDevice(identifiers []*generated.DeviceIdentifier,
	device map[string]interface{}, parentKeyIndex int,
	expectedType string, prefix string) {
	for _, identifier := range identifiers {
		identifier := filter(identifier, expectedType, prefix)
		if identifier == nil {
			continue
		}
		classifierValue := adjustClassifierFormatForSchemaTransformation(identifier.Classifiers[0].GetValue())

		keys := strings.Split(classifierValue, "/")[parentKeyIndex:]
		switch identifier.GetValue().(type) {
		case *generated.DeviceIdentifier_Text:
			{
				// special case mapping for profinet devices to map their semantic name to our asset name
				semanticProfinetNameMapping(keys, identifier.GetText(), device)
				semanticProfinetFirmwareVersionMapping(keys, identifier.GetText(), device)
				transformKeys(keys, identifier.GetText(), device)
			}
		case *generated.DeviceIdentifier_Children:
			{
				createChildProperty(identifier, device, keys, parentKeyIndex, expectedType, prefix)
			}
		case *generated.DeviceIdentifier_Int64Value:
			{
				transformKeys(keys, identifier.GetInt64Value(), device)
			}
		case *generated.DeviceIdentifier_Float64Value:
			{
				transformKeys(keys, identifier.GetFloat64Value(), device)
			}
		case *generated.DeviceIdentifier_RawData:
			{
				transformKeys(keys, identifier.GetRawData(), device)
			}
		default:
			_ = fmt.Errorf("unknown datatype %s with classifier value %s", identifier.GetValue(), classifierValue)
		}
	}
}

func mapPropertiesIntoChildObjects(identifiers []*generated.DeviceIdentifier,
	childrenContainer []map[string]interface{}, parentKeyIndex int,
	expectedType string, prefix string) []map[string]interface{} {
	arrayContext := ArrayContext{IsInArray: false, ArrayIndex: -1}

	for identifierIndex, identifier := range identifiers {
		identifier := filter(identifier, expectedType, prefix)
		if identifier == nil {
			continue
		}
		classifierValue := adjustClassifierFormatForSchemaTransformation(identifier.Classifiers[0].GetValue())
		if strings.Contains(classifierValue, "[") {
			arrayContext.IsInArray = true
			arrayIndex := getArrayIndexFromClassifier(classifierValue)
			if arrayIndex == -1 {
				continue
			}
			// Use arrayIndex as needed
			arrayContext.ArrayIndex = arrayIndex
			// create array element, if not exists
			for len(childrenContainer) <= arrayContext.ArrayIndex {
				child := make(map[string]interface{})
				childrenContainer = append(childrenContainer, child)
			}
		} else if identifierIndex == 0 {
			// still in array context, but no index in property path, so we assume just one element is present
			// Sooo turns out this branch is actually in use
			// this part works by specifying the same children container once for every child
			// therefor each container specification should only result in one child
			// we work around this, by saying that only the first child property created its own container object
			child := make(map[string]interface{})
			childrenContainer = append(childrenContainer, child)
			arrayContext.ArrayIndex = len(childrenContainer) - 1
			arrayContext.IsInArray = true
		}
		// hand down array element to insert properties into
		device := childrenContainer[arrayContext.ArrayIndex]

		keys := strings.Split(classifierValue, "/")[parentKeyIndex:]
		switch identifier.GetValue().(type) {
		case *generated.DeviceIdentifier_Text:
			{
				// special case mapping for profinet devices to map their semantic name to our asset name
				semanticProfinetNameMapping(keys, identifier.GetText(), device)
				semanticProfinetFirmwareVersionMapping(keys, identifier.GetText(), device)
				transformKeys(keys, identifier.GetText(), device)
			}
		case *generated.DeviceIdentifier_Children:
			{
				createChildProperty(identifier, device, keys, parentKeyIndex, expectedType, prefix)
			}
		case *generated.DeviceIdentifier_Int64Value:
			{
				transformKeys(keys, identifier.GetInt64Value(), device)
			}
		case *generated.DeviceIdentifier_Float64Value:
			{
				transformKeys(keys, identifier.GetFloat64Value(), device)
			}
		case *generated.DeviceIdentifier_RawData:
			{
				transformKeys(keys, identifier.GetRawData(), device)
			}
		default:
			_ = fmt.Errorf("unknown datatype %s with classifier value %s", identifier.GetValue(), classifierValue)
		}
	}
	return childrenContainer
}

func createChildProperty(identifier *generated.DeviceIdentifier, device map[string]interface{},
	keys []string, parentKeyIndex int,
	expectedType string, prefix string) {
	// before creating the children container (array property), we need to check, if it already exists
	var arrayProperty []map[string]interface{}
	existingArrayProperty, ok := device[keys[len(keys)-1]].([]map[string]interface{})
	if ok {
		arrayProperty = existingArrayProperty
	} else {
		arrayProperty = make([]map[string]interface{}, 0)
	}
	// create a container for the child objects
	// every child is an object in an array for children identifiers
	// child objects  will be created in a lower level, since only then,
	// we know, the array index of the child object
	// so just provide the container for the child objects

	arrayProperty = mapPropertiesIntoChildObjects(identifier.GetChildren().GetValue(),
		arrayProperty, parentKeyIndex+len(keys), expectedType, prefix)

	// this should rather work like the recursive property level creation in <transformKeys>
	transformKeys(keys, arrayProperty, device)
}

func semanticProfinetNameMapping(keys []string, identifierTextValue string, device map[string]interface{}) {
	for _, v := range keys {
		if v == "profinet_name" {
			device["name"] = identifierTextValue
		}
	}

}

func semanticProfinetFirmwareVersionMapping(keys []string, identifierTextValue string, device map[string]interface{}) {
	// firmware%2Fsoftware_identifier%2Fversion -> firmware/software_identifier/version
	// to instance_annotations: [{"key": "firmware_version", "value": "1.0.0"}]
	if (len(keys) == 3 && keys[0] == "firmware") && (keys[1] == "software_identifier") && (keys[2] == "version") {
		// check if instance_annotations already exists
		instanceAnnotations, ok := device["instance_annotations"].([]interface{})
		if ok {
			instanceAnnotations = append(instanceAnnotations, map[string]interface{}{"key": "firmware_version", "value": identifierTextValue})
			device["instance_annotations"] = instanceAnnotations
			return
		}

		// if not create them
		newAnnotations := []interface{}{}
		newAnnotations = append(newAnnotations, map[string]interface{}{"key": "firmware_version", "value": identifierTextValue})
		device["instance_annotations"] = newAnnotations
	}
}

func adjustClassifierFormatForSchemaTransformation(classifierValue string) string {
	classifierValue, err := url.PathUnescape(classifierValue)
	if err != nil {
		_ = fmt.Errorf("error decoding the string %v", err)
	}
	_, after, success := strings.Cut(classifierValue, "#")
	if success {
		classifierValue = "/" + after
	}
	return classifierValue
}

func getArrayIndexFromClassifier(classifierValue string) int {
	start := strings.LastIndex(classifierValue, "[")
	end := strings.LastIndex(classifierValue, "]")
	if start != -1 && end != -1 && start < end {
		arrayIndex, err := strconv.Atoi(classifierValue[start+1 : end])
		if err != nil {
			_ = fmt.Errorf("error parsing array index %v", err)
			return -1
		}
		return arrayIndex
	}
	return -1
}
