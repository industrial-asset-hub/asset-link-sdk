/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

const (
	baseSchemaVersion = "v1"
	baseSchemaPrefix  = "https://industrial-assets.io/schemas/iah/base-schema/released/" + baseSchemaVersion
)

// NewDevice Generates a new asset skeleton
func NewDevice(functionalObjectType string, assetName string) (*DeviceInfo, error) {

	d := DeviceInfo{}
	if !isNonEmptyValues(functionalObjectType) {
		err := &EmptyError{
			Field:   "FunctionalObjectType",
			Message: "Functional object type is required and cannot be empty",
			Value:   functionalObjectType,
		}
		return &d, err
	}

	if err := ValidateField(
		FunctionalObjectSchemaUrl,
		"FunctionalObjectSchemaUrl",
		"Functional object schema URL is empty",
		FunctionalObjectSchemaUrlPattern,
		"Functional object schema URL format is invalid. Please refer to the base schema for the supported pattern.",
	); err != nil {
		return &d, err
	}

	d.FunctionalObjectType = functionalObjectType
	d.FunctionalObjectSchemaUrl = FunctionalObjectSchemaUrl
	d.Name = &assetName

	return &d, nil
}

type DeviceInfo struct {
	FunctionalObjectType      any `json:"functional_object_type"`
	FunctionalObjectSchemaUrl any `json:"functional_object_schema_url"`
	// Override connection point, since generated base schema does not provide derived types
	ConnectionPoints []any `json:"connection_points,omitempty"`
	Asset
	AssetIdentifiers []any `json:"asset_identifiers"`
	// To Be clarified
	SoftwareComponents []any `json:"software_components,omitempty"`
}

func (d *DeviceInfo) AddDescription(description string) error {
	if !isNonEmptyValues(description) {
		err := &EmptyError{
			Field:   "Description",
			Message: "Description is empty",
		}
		return err
	}
	d.Description = &description
	return nil
}
