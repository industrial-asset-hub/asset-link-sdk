/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

// NewDeviceRelationship Generates a new device relationship skeleton
func NewDeviceRelationship(subject *DeviceInfo, predicate PredicateValues, object *DeviceInfo) *DeviceRelationshipInfo {
	d := DeviceRelationshipInfo{}
	d.AssetRelationship = AssetRelationship{
		Subject:   &subject.Id,
		Predicate: &predicate,
		Object:    &object.Id,
	}

	return &d
}

type DeviceRelationshipInfo struct {
	AssetRelationship AssetRelationship `json:"asset_relationship"`
}
