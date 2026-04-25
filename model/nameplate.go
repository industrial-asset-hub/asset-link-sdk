/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"crypto/sha1"
)

// AddNameplate Add a digital nameplate to an asset.
// The nameplate is inspired by IDTA 02006-0-0 Digital Nameplate for industrial equippment
//
// manufacturerName: legally valid designation of the natural or judicial person
// URIOfTheProduct: unique global identification of the product using an universal resource identifier (URI
// productArticleNumberOfManufacturer: unique product identifier of the manufacturer
// manufacturerProductDesignation: short description of the product (short text)
// hardwareVersion: version of the hardware supplied with the device
// serialNumber: unique combination of numbers and letters used to identify
// the device once it has been manufactured
func (d *DeviceInfo) AddNameplate(manufacturerName string,
	uriOfTheProduct string,
	productArticleNumberOfManufacturer string,
	manufacturerProductDesignation string,
	hardwareVersion string,
	serialNumber string,
) error {
	if !isNonEmptyValues(manufacturerName, uriOfTheProduct, productArticleNumberOfManufacturer,
		manufacturerProductDesignation, hardwareVersion, serialNumber) {
		err := &EmptyError{
			Field:   "ProductInstanceInformation",
			Message: "All fields for ProductInstanceInformation are empty",
		}
		return err
	}

	if err := ValidateField(
		uriOfTheProduct,
		"ProductLink",
		"Product link is empty",
		IdLinkPattern,
		"Product link format is invalid. Please refer to the base schema for the supported pattern.",
	); err != nil {
		return err
	}

	// We hash the manufacturer to get a unique identifier
	h := sha1.New()
	h.Write([]byte(manufacturerName))

	organisation := Organization{
		Name: &manufacturerName,
	}
	mp := Product{
		ProductLink:    &uriOfTheProduct,
		Manufacturer:   &organisation,
		ProductId:      &productArticleNumberOfManufacturer,
		ProductVersion: &hardwareVersion,
	}

	d.ProductInstanceInformation = &ProductInstanceInformation{
		ManufacturerProduct: &mp,
		SerialNumber:        &serialNumber,
	}

	d.addIdLinkIdentifier(uriOfTheProduct)
	return nil
}

// AddSoftware Add software information to an asset
func (d *DeviceInfo) AddSoftware(name string, version string, isFirmware bool) error {
	if err := ValidateField(name, "SoftwareName", "Software name is empty", "", ""); err != nil {
		return err
	}
	if err := ValidateField(version, "SoftwareVersion", "Software version is empty", "", ""); err != nil {
		return err
	}
	if err := ValidateField(
		FunctionalObjectSchemaUrl,
		"FunctionalObjectSchemaUrl",
		"Functional object schema URL is empty",
		FunctionalObjectSchemaUrlPattern,
		"Functional object schema URL format is invalid. Please refer to the base schema for the supported pattern.",
	); err != nil {
		return err
	}

	softwareArtifact := SoftwareArtifact{
		AssetOperations:            nil,
		ConnectionPoints:           nil,
		FunctionalObjectType:       SoftwareArtifactFunctionalObjectTypeSoftwareArtifact,
		InstanceAnnotations:        nil,
		Name:                       nil,
		ProductInstanceInformation: nil,
		SoftwareComponents:         nil,
		IsFirmware:                 &isFirmware,
		FunctionalObjectSchemaUrl:  FunctionalObjectSchemaUrl,
	}

	softwareIdentifier := SoftwareIdentifier{}
	softwareIdentifier.AssetIdentifierType = SoftwareIdentifierAssetIdentifierTypeSoftwareIdentifier
	softwareIdentifier.Name = name
	softwareIdentifier.Version = version
	softwareArtifact.AssetIdentifiers = append(softwareArtifact.AssetIdentifiers, softwareIdentifier)

	d.SoftwareComponents = append(d.SoftwareComponents, softwareArtifact)
	return nil
}
