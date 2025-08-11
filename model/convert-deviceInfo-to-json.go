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
	"strings"
)

func (d *DeviceInfo) convertToJson() (map[string]interface{}, error) {
	if d == nil {
		return nil, fmt.Errorf("DeviceInfo is nil")
	}
	deviceInfoValue := reflect.ValueOf(d)
	deviceInfoValue = dereferenceDeviceInfoValue(deviceInfoValue)

	switch deviceInfoValue.Kind() {
	case reflect.Struct:
		deviceInfoMap, err := convertDeviceInfoToMap(deviceInfoValue)
		if err != nil {
			return nil, err
		}
		// Delete the "id" key from the map
		delete(deviceInfoMap, "id")
		return deviceInfoMap, nil
	default:
		return map[string]interface{}{}, nil
	}
}

func convertDeviceInfoToMap(deviceInfoValue reflect.Value) (map[string]interface{}, error) {
	if isTime(deviceInfoValue.Type()) {
		return map[string]interface{}{"": deviceInfoValue.Interface()}, nil
	}

	deviceInfoMap := make(map[string]interface{})
	t := deviceInfoValue.Type()

	deviceInfoFields := reflect.VisibleFields(t)
	for _, deviceInfoField := range deviceInfoFields {
		name, opts, include := extractJsonNameAndOptions(deviceInfoField)
		if !include {
			continue
		}
		fieldValue := deviceInfoValue.FieldByIndex(deviceInfoField.Index)

		if opts["omitempty"] && isFieldValueZero(fieldValue) {
			continue
		}

		propertyValue, err := changePropertyValueToInterface(fieldValue)
		if err != nil {
			return nil, err
		}

		// Handle special case for "Asset" property in device-info model
		if name == "Asset" {
			if assetMap, ok := propertyValue.(map[string]interface{}); ok {
				for key, value := range assetMap {
					deviceInfoMap[key] = value
				}
				continue
			}
		}

		deviceInfoMap[name] = propertyValue
	}
	return deviceInfoMap, nil
}

func changePropertyValueToInterface(v reflect.Value) (interface{}, error) {
	v = dereferenceDeviceInfoValue(v)
	if !v.IsValid() {
		return nil, nil
	}

	if v.Kind() == reflect.Struct && isTime(v.Type()) {
		return v.Interface(), nil
	}

	switch v.Kind() {
	case reflect.Struct:
		return convertDeviceInfoToMap(v)
	case reflect.Slice, reflect.Array:
		n := v.Len()
		out := make([]interface{}, n)
		for i := 0; i < n; i++ {
			x, err := changePropertyValueToInterface(v.Index(i))
			if err != nil {
				return nil, err
			}
			out[i] = x
		}
		return out, nil
	case reflect.Map:
		out := make(map[string]interface{}, v.Len())
		for _, k := range v.MapKeys() {
			key := k.Interface()
			var ks string
			if k.Kind() == reflect.String {
				ks = key.(string)
			} else {
				ks = fmt.Sprint(key)
			}
			val, err := changePropertyValueToInterface(v.MapIndex(k))
			if err != nil {
				return nil, err
			}
			out[ks] = val
		}
		return out, nil
	case reflect.Bool,
		reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		if v.CanInterface() {
			return v.Interface(), nil
		}
		return nil, nil
	case reflect.Invalid:
		return nil, nil
	default:
		if v.CanInterface() {
			return v.Interface(), nil
		}
		return nil, nil
	}
}

func dereferenceDeviceInfoValue(deviceInfoValue reflect.Value) reflect.Value {
	for deviceInfoValue.IsValid() && (deviceInfoValue.Kind() == reflect.Pointer || deviceInfoValue.Kind() == reflect.Interface) {
		if deviceInfoValue.IsNil() {
			return reflect.Value{}
		}
		deviceInfoValue = deviceInfoValue.Elem()
	}
	return deviceInfoValue
}

func isFieldValueZero(v reflect.Value) bool {
	return !v.IsValid() || v.IsZero()
}

func isTime(t reflect.Type) bool {
	return t.PkgPath() == "time" && t.Name() == "Time"
}

func extractJsonNameAndOptions(structField reflect.StructField) (name string, opts map[string]bool, include bool) {
	jsonTag := structField.Tag.Get("json")
	opts = map[string]bool{}

	if jsonTag == "-" {
		return "", opts, false
	}
	if jsonTag == "" {
		return structField.Name, opts, true
	}
	parts := strings.Split(jsonTag, ",")
	if len(parts) == 0 || parts[0] == "" {
		return structField.Name, opts, true
	}
	name = parts[0]
	if name == "" {
		name = structField.Name
	}
	for _, part := range parts[1:] {
		if part != "" {
			opts[part] = true
		}
	}
	return name, opts, true
}
