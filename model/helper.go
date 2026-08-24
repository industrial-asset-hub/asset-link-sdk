/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import "strings"

// AddMacIdentifier appends a MacIdentifier only when provided value is non-empty.
func (d *DeviceInfo) AddMacIdentifier(macAddress string) {

	if isNonEmptyValues(macAddress) {
		identifierUncertainty := 1
		identifier := MacIdentifier{
			AssetIdentifierType:   MacIdentifierAssetIdentifierTypeMacIdentifier,
			IdentifierType:        nil,
			IdentifierUncertainty: &identifierUncertainty,
			MacAddress:            macAddress,
		}
		d.AssetIdentifiers = append(d.AssetIdentifiers, identifier)
	}
}

func (d *DeviceInfo) AddIdLinkIdentifier(uriOfTheProduct string) {
	idLinkIdentifier := IdLinkIdentifier{
		AssetIdentifierType:   IdLinkIdentifierAssetIdentifierTypeIdLinkIdentifier,
		IdLink:                uriOfTheProduct,
		IdentifierType:        nil,
		IdentifierUncertainty: nil,
	}
	d.AssetIdentifiers = append(d.AssetIdentifiers, idLinkIdentifier)
}

// AddCustomIdentifier appends a CustomIdentifier only when provided values are valid.
func (d *DeviceInfo) AddCustomIdentifier(name, value string) {
	if strings.TrimSpace(name) != "" && strings.TrimSpace(value) != "" && validateByPattern(value, CustomIdentifierValuePattern) {
		identifier := CustomIdentifier{
			AssetIdentifierType:   CustomIdentifierAssetIdentifierTypeCustomIdentifier,
			IdentifierType:        nil,
			IdentifierUncertainty: nil,
			Name:                  name,
			Value:                 value,
		}
		d.AssetIdentifiers = append(d.AssetIdentifiers, identifier)
	}
}

// AddCertificateIdentifier appends a CertificateIdentifier only when provided value is non-empty.
func (d *DeviceInfo) AddCertificateIdentifier(certificateID string) {
	if isNonEmptyValues(certificateID) {
		identifier := CertificateIdentifier{
			AssetIdentifierType:   CertificateIdentifierAssetIdentifierTypeCertificateIdentifier,
			CertificateId:         certificateID,
			IdentifierType:        nil,
			IdentifierUncertainty: nil,
		}
		d.AssetIdentifiers = append(d.AssetIdentifiers, identifier)
	}
}

// AddHostBasedSoftwareIdentifier appends a HostBasedSoftwareIdentifier only when all required fields are non-empty.
func (d *DeviceInfo) AddHostBasedSoftwareIdentifier(name, version, hostIdentifier string) {
	if strings.TrimSpace(name) != "" && strings.TrimSpace(version) != "" && strings.TrimSpace(hostIdentifier) != "" {
		identifier := HostBasedSoftwareIdentifier{
			AssetIdentifierType: HostBasedSoftwareIdentifierAssetIdentifierTypeHostBasedSoftwareIdentifier,
			Name:                name,
			Version:             version,
			HostIdentifier:      hostIdentifier,
		}
		d.AssetIdentifiers = append(d.AssetIdentifiers, identifier)
	}
}

// AddProductInstanceIdentifier appends a ProductInstanceIdentifier; skipped only when all three fields are empty.
func (d *DeviceInfo) AddProductInstanceIdentifier(vendor, articleNumber, serialNumber string) {
	if !isNonEmptyValues(vendor, articleNumber, serialNumber) {
		return
	}
	identifier := ProductInstanceIdentifier{
		AssetIdentifierType: ProductInstanceIdentifierAssetIdentifierTypeProductInstanceIdentifier,
		Vendor:              vendor,
		ArticleNumber:       articleNumber,
		SerialNumber:        serialNumber,
	}
	d.AssetIdentifiers = append(d.AssetIdentifiers, identifier)
}

// AddParentRelativeIdentifier appends a ParentRelativeIdentifier; parentIdentifier must be a MacIdentifier or IdLinkIdentifier.
// Nil slot or subslot defaults to 0.
func (d *DeviceInfo) AddParentRelativeIdentifier(parentIdentifier interface{}, slot, subslot *int) {
	if parentIdentifier == nil {
		return
	}
	slotVal, subslotVal := 0, 0
	if slot != nil {
		slotVal = *slot
	}
	if subslot != nil {
		subslotVal = *subslot
	}
	identifier := ParentRelativeIdentifier{
		AssetIdentifierType: ParentRelativeIdentifierAssetIdentifierTypeParentRelativeIdentifier,
		ParentIdentifier:    parentIdentifier,
		Slot:                slotVal,
		Subslot:             subslotVal,
	}
	d.AssetIdentifiers = append(d.AssetIdentifiers, identifier)
}

// AddAssetRelation appends an AssetRelation when provided values are valid.
// Set isBidirectional to true if the relation can be interpreted in both directions.
func (d *DeviceInfo) AddAssetRelation(predicate string, relatedAsset RelatedAsset, role RelationalRoleOfRelatedAssetValues, isBidirectional bool) error {
	err := validateField(predicate, "Predicate", "Predicate is empty", PredicatePattern,
		"Predicate format is invalid. Please refer to the base schema for the supported pattern.")
	if err != nil {
		return err
	}

	if err := validateIdentifiers(relatedAsset.AssetIdentifiers); err != nil {
		return err
	}

	var bidirectional *bool
	if isBidirectional {
		bidirectional = &isBidirectional
	}
	d.AssetRelations = append(d.AssetRelations, AssetRelation{
		Predicate:                    predicate,
		RelatedAsset:                 relatedAsset,
		RelationalRoleOfRelatedAsset: role,
		IsBidirectional:              bidirectional,
	})
	return nil
}

// validateIdentifiers validates all identifiers in the RelatedAsset.
func validateIdentifiers(identifiers []interface{}) error {
	if len(identifiers) == 0 {
		return &EmptyError{
			Field:   "AssetIdentifiers",
			Message: "Related asset identifiers are empty",
		}
	}

	for _, identifier := range identifiers {
		switch id := identifier.(type) {
		case MacIdentifier:
			if err := validateField(id.MacAddress, "MacAddress", "MAC address is empty", MacAddressPattern,
				"MAC address format is invalid. Please refer to the base schema for the supported pattern."); err != nil {
				return err
			}
		case CustomIdentifier:
			if err := validateField(id.Value, "CustomIdentifier.Value", "Custom identifier value is empty", CustomIdentifierValuePattern,
				"Custom identifier value format is invalid. Please refer to the base schema for the supported pattern."); err != nil {
				return err
			}
		case IdLinkIdentifier:
			if err := validateField(id.IdLink, "IdLinkIdentifier.IdLink", "IdLink identifier is empty", IdLinkPattern,
				"IdLink format is invalid. Please refer to the base schema for the supported pattern."); err != nil {
				return err
			}
		}
	}
	return nil
}
