/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

// Add nameplate to an discovered asset
//
// Name = ManufacturerProductDesignation
func (d *DeviceInfo) AddNameplate(manufacturerName string,
	ProductArticleNumberOfManufacturer string,
	manufacturerProductFamily string,
	hardwareVersion string,
	serialNumber string,
) {

	organisation := Organization{
		Address:        nil,
		AlternateNames: nil,
		ContactPoint:   nil,
		Id:             "",
		Name:           &manufacturerName,
	}

	mp := Product{
		Id:             "",
		Manufacturer:   &organisation,
		Name:           &manufacturerProductFamily,
		ProductId:      &ProductArticleNumberOfManufacturer,
		ProductVersion: &hardwareVersion,
	}
	pi := ProductSerialIdentifier{
		IdentifierType:        nil,
		IdentifierUncertainty: nil,
		ManufacturerProduct:   &mp,
		SerialNumber:          &serialNumber,
	}

	d.ProductInstanceIdentifier = &pi

}
