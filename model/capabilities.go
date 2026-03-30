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
func (d *DeviceInfo) AddCapabilities(name string, enabled bool) error {
	if err := ValidateField(name, "CapabilityName", "Capability name is empty", "", ""); err != nil {
		return err
	}
	operation := AssetOperation{}
	operation.OperationName = &name
	operation.ActivationFlag = &enabled
	d.AssetOperations = append(d.AssetOperations, operation)
	return nil
}
