/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import "strings"

func (d *DeviceInfo) addMacIdentifier(macAddress string) {

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

func (d *DeviceInfo) addIdLinkIdentifier(uriOfTheProduct string) {
	if isNonEmptyValues(uriOfTheProduct) && ValidateByPattern(uriOfTheProduct, IdLinkPattern) {
		idLinkIdentifier := IdLinkIdentifier{
			AssetIdentifierType:   IdLinkIdentifierAssetIdentifierTypeIdLinkIdentifier,
			IdLink:                uriOfTheProduct,
			IdentifierType:        nil,
			IdentifierUncertainty: nil,
		}
		d.AssetIdentifiers = append(d.AssetIdentifiers, idLinkIdentifier)
	}
}

// addCustomIdentifier appends a CustomIdentifier only when provided values are valid.
func (d *DeviceInfo) addCustomIdentifier(name, value string) {
	if strings.TrimSpace(name) != "" && strings.TrimSpace(value) != "" && ValidateByPattern(value, CustomIdentifierValuePattern) {
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

// addCertificateIdentifier appends a CertificateIdentifier only when provided value is non-empty.
func (d *DeviceInfo) addCertificateIdentifier(certificateID string) {
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
