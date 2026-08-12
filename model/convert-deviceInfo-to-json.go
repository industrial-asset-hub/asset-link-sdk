/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (d *DeviceInfo) ConvertToJson() (map[string]interface{}, error) {
	if d == nil {
		return nil, fmt.Errorf("DeviceInfo is nil")
	}
	type alias DeviceInfo // bypass MarshalJSON to avoid infinite recursion
	data, err := json.Marshal((*alias)(d))
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var m map[string]interface{}
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	delete(m, "id")
	return m, nil
}

// MarshalJSON converts DeviceInfo into nested JSON using the existing reflection-based mapper.
func (d *DeviceInfo) MarshalJSON() ([]byte, error) {
	deviceInfoMap, err := d.ConvertToJson()
	if err != nil {
		return nil, err
	}
	return json.Marshal(deviceInfoMap)
}
