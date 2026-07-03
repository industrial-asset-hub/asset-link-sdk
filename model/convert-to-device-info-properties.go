/*
 * SPDX-FileCopyrightText: 2026 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"fmt"
	"sort"

	generatedDeviceInfo "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/conn_suite_device_info"
	generated "github.com/industrial-asset-hub/asset-link-sdk/v4/generated/iah-discovery"
)

func (d *DeviceInfo) ConvertToPropertyValueResults() ([]*generatedDeviceInfo.PropertyValueResult, error) {
	properties, err := d.ConvertToJson()
	if err != nil {
		return nil, err
	}

	keys := sortedKeysOf(properties)

	results := make([]*generatedDeviceInfo.PropertyValueResult, 0, len(keys))
	for _, key := range keys {
		value := properties[key]
		results = append(results, &generatedDeviceInfo.PropertyValueResult{
			Result: &generatedDeviceInfo.PropertyValueResult_Property{
				Property: &generatedDeviceInfo.PropertyKeyValuePair{
					Key:   key,
					Value: interfaceToPropertyVariant(value),
				},
			},
		})
	}

	return results, nil
}

// sortedKeysOf returns the keys of m in sorted order.
func sortedKeysOf(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func interfaceToPropertyVariant(value interface{}) *generated.Variant {
	switch typed := value.(type) {
	case nil:
		return &generated.Variant{Value: &generated.Variant_NullValue{NullValue: generated.NullValue_NULL_VALUE}}
	case bool:
		return &generated.Variant{Value: &generated.Variant_BoolValue{BoolValue: typed}}
	case string:
		return &generated.Variant{Value: &generated.Variant_Text{Text: typed}}
	case map[string]interface{}:
		fields := make(map[string]*generated.Variant, len(typed))
		for key, item := range typed {
			fields[key] = interfaceToPropertyVariant(item)
		}
		return &generated.Variant{Value: &generated.Variant_StructValue{StructValue: &generated.VariantStruct{Fields: fields}}}
	case []interface{}:
		values := make([]*generated.Variant, 0, len(typed))
		for _, item := range typed {
			values = append(values, interfaceToPropertyVariant(item))
		}
		return &generated.Variant{Value: &generated.Variant_ArrayValue{ArrayValue: &generated.VariantArray{Values: values}}}
	case int:
		return &generated.Variant{Value: &generated.Variant_Int64Value{Int64Value: int64(typed)}}
	case int8:
		return &generated.Variant{Value: &generated.Variant_Int64Value{Int64Value: int64(typed)}}
	case int16:
		return &generated.Variant{Value: &generated.Variant_Int64Value{Int64Value: int64(typed)}}
	case int32:
		return &generated.Variant{Value: &generated.Variant_Int64Value{Int64Value: int64(typed)}}
	case int64:
		return &generated.Variant{Value: &generated.Variant_Int64Value{Int64Value: typed}}
	case uint:
		return &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: uint64(typed)}}
	case uint8:
		return &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: uint64(typed)}}
	case uint16:
		return &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: uint64(typed)}}
	case uint32:
		return &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: uint64(typed)}}
	case uint64:
		return &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: typed}}
	case float32:
		return &generated.Variant{Value: &generated.Variant_Float64Value{Float64Value: float64(typed)}}
	case float64:
		return &generated.Variant{Value: &generated.Variant_Float64Value{Float64Value: typed}}
	case []byte:
		return &generated.Variant{Value: &generated.Variant_RawData{RawData: typed}}
	default:
		return &generated.Variant{Value: &generated.Variant_Text{Text: fmt.Sprint(typed)}}
	}
}
