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
	v := reflect.ValueOf(d)
	v = deref(v)
	deviceInfoJson := map[string]interface{}{}

	switch v.Kind() {
	case reflect.Struct:
		m, err := structToMap(v)
		if err != nil {
			return nil, err
		}

		// Delete the "id" key from the map
		delete(m, "id")

		return m, nil
	default:
		converted, err := toInterface(v)
		if err != nil {
			return nil, err
		}
		deviceInfoJson[""] = converted
		return deviceInfoJson, nil
	}
}

func structToMap(v reflect.Value) (map[string]interface{}, error) {
	if isTime(v.Type()) {
		return map[string]interface{}{"": v.Interface()}, nil
	}

	out := make(map[string]interface{})
	t := v.Type()

	fields := reflect.VisibleFields(t)
	for _, f := range fields {
		name, opts, include := jsonNameAndOpts(f)
		if !include {
			continue
		}
		fv := v.FieldByIndex(f.Index)

		if opts["omitempty"] && isZero(fv) {
			continue
		}

		val, err := toInterface(fv)
		if err != nil {
			return nil, err
		}

		if name == "Asset" {
			if assetMap, ok := val.(map[string]interface{}); ok {
				for k, v := range assetMap {
					out[k] = v
				}
				continue
			}
		}

		out[name] = val
	}
	return out, nil
}

func toInterface(v reflect.Value) (interface{}, error) {
	v = deref(v)
	if !v.IsValid() {
		return nil, nil
	}

	if v.Kind() == reflect.Struct && isTime(v.Type()) {
		return v.Interface(), nil
	}

	switch v.Kind() {
	case reflect.Struct:
		return structToMap(v)
	case reflect.Slice, reflect.Array:
		n := v.Len()
		out := make([]interface{}, n)
		for i := 0; i < n; i++ {
			x, err := toInterface(v.Index(i))
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
			val, err := toInterface(v.MapIndex(k))
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

func deref(v reflect.Value) reflect.Value {
	for v.IsValid() && (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

func isZero(v reflect.Value) bool {
	return !v.IsValid() || v.IsZero()
}

func isTime(t reflect.Type) bool {
	return t.PkgPath() == "time" && t.Name() == "Time"
}

func jsonNameAndOpts(sf reflect.StructField) (name string, opts map[string]bool, include bool) {
	tag := sf.Tag.Get("json")
	opts = map[string]bool{}

	if tag == "-" {
		return "", opts, false
	}
	if tag == "" {
		return sf.Name, opts, true
	}
	parts := strings.Split(tag, ",")
	if len(parts) == 0 || parts[0] == "" {
		return sf.Name, opts, true
	}
	name = parts[0]
	if name == "" {
		name = sf.Name
	}
	for _, p := range parts[1:] {
		if p != "" {
			opts[p] = true
		}
	}
	return name, opts, true
}
