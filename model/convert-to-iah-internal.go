/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"fmt"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	DeviceInIahSchema["@context"] = map[string]interface{}{
		"base":      "https://common-device-management.code.siemens.io/documentation/asset-modeling/base-schema/v0.7.5/",
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
