/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

// AddCapabilities
func (d *DeviceInfo) AddCapabilities(name string, enabled bool) {
	opperation := AssetOperation{
		ActivationFlag: &enabled,
		OperationName:  &name,
	}

	d.AssetOperations = append(d.AssetOperations, opperation)
}
