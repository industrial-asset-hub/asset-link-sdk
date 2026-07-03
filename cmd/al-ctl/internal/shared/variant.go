/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package shared

import (
	"fmt"

	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
)

// VariantToInterface converts a protobuf Variant value into native Go data
// structures for JSON marshalling and further schema validation.
func VariantToInterface(value *generated.Variant, path string) (interface{}, error) {
	if value == nil {
		return nil, fmt.Errorf("%s is nil", path)
	}

	switch typed := value.GetValue().(type) {
	case *generated.Variant_BoolValue:
		return typed.BoolValue, nil
	case *generated.Variant_Int64Value:
		return typed.Int64Value, nil
	case *generated.Variant_Uint64Value:
		return typed.Uint64Value, nil
	case *generated.Variant_Float64Value:
		return typed.Float64Value, nil
	case *generated.Variant_Text:
		return typed.Text, nil
	case *generated.Variant_RawData:
		return typed.RawData, nil
	case *generated.Variant_File:
		if typed.File == nil {
			return nil, fmt.Errorf("%s file is nil", path)
		}
		return map[string]interface{}{
			"data":         typed.File.GetData(),
			"content_type": typed.File.GetContentType(),
			"file_name":    typed.File.GetFileName(),
		}, nil
	case *generated.Variant_StructValue:
		if typed.StructValue == nil {
			return nil, fmt.Errorf("%s struct_value is nil", path)
		}
		out := make(map[string]interface{}, len(typed.StructValue.GetFields()))
		for key, child := range typed.StructValue.GetFields() {
			childValue, err := VariantToInterface(child, fmt.Sprintf("%s.%s", path, key))
			if err != nil {
				return nil, err
			}
			out[key] = childValue
		}
		return out, nil
	case *generated.Variant_ArrayValue:
		if typed.ArrayValue == nil {
			return nil, fmt.Errorf("%s array_value is nil", path)
		}
		out := make([]interface{}, len(typed.ArrayValue.GetValues()))
		for i, child := range typed.ArrayValue.GetValues() {
			childValue, err := VariantToInterface(child, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			out[i] = childValue
		}
		return out, nil
	case *generated.Variant_NullValue:
		return nil, nil
	default:
		return nil, fmt.Errorf("%s has unsupported variant type %T", path, typed)
	}
}

// BuildAssetFromPropertyResults converts successful property results into a flat
// asset map keyed by property key. Errors are returned separately for callers
// that want to log or inspect dropped property entries.
func BuildAssetFromPropertyResults(
	propertyResults []*generatedDeviceInfo.PropertyValueResult,
	pathPrefix string,
) (map[string]interface{}, int, []*generatedDeviceInfo.GetPropertyValueError, error) {
	asset := map[string]interface{}{}
	propertyCount := 0
	errorEntries := make([]*generatedDeviceInfo.GetPropertyValueError, 0)

	for index, result := range propertyResults {
		if result.GetError() != nil {
			errorEntries = append(errorEntries, result.GetError())
			continue
		}

		property := result.GetProperty()
		if property == nil {
			continue
		}

		value, err := VariantToInterface(property.GetValue(), fmt.Sprintf("%s[%d].%s", pathPrefix, index, property.GetKey()))
		if err != nil {
			return nil, 0, nil, err
		}

		asset[property.GetKey()] = value
		propertyCount++
	}

	return asset, propertyCount, errorEntries, nil
}
