/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

// AddCapabilities Add device capabilities to an asset
//
// Name can be for example firmware_update, backup or restore.
func (d *DeviceInfo) AddCapabilities(name string, enabled bool) {
	opperation := AssetOperation{
		ActivationFlag: &enabled,
		OperationName:  &name,
	}

	d.AssetOperations = append(d.AssetOperations, opperation)
}
