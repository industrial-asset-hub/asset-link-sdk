/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

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
	if isNonEmptyValues(uriOfTheProduct) {
		idLinkIdentifier := IdLinkIdentifier{
			AssetIdentifierType:   IdLinkIdentifierAssetIdentifierTypeIdLinkIdentifier,
			IdLink:                uriOfTheProduct,
			IdentifierType:        nil,
			IdentifierUncertainty: nil,
		}
		d.AssetIdentifiers = append(d.AssetIdentifiers, idLinkIdentifier)
	}
}
