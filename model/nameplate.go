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
// productFamily: family of products the device belongs to
// hardwareVersion: version of the hardware supplied with the device
// serialNumber: unique combination of numbers and letters used to identify
// the device once it has been manufactured
func (d *DeviceInfo) AddNameplate(manufacturerName string,
	uriOfTheProduct string,
	productArticleNumberOfManufacturer string,
	productFamily string,
	hardwareVersion string,
	serialNumber string,
) error {
	if !isNonEmptyValues(manufacturerName, uriOfTheProduct, productArticleNumberOfManufacturer,
		productFamily, hardwareVersion, serialNumber) {
		err := &EmptyError{
			Field:   "ProductInstanceInformation",
			Message: "All fields for ProductInstanceInformation are empty",
		}
		return err
	}

	if isNonEmptyValues(uriOfTheProduct) {
		if !validateByPattern(uriOfTheProduct, IdLinkPattern) {
			return &ValidationError{
				Field:   "ProductLink",
				Message: "Product link format is invalid. Please refer to the base schema for the supported pattern.",
				Value:   uriOfTheProduct,
				Details: uriOfTheProduct,
			}
		}
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
		ProductFamily:  &productFamily,
	}

	d.ProductInstanceInformation = &ProductInstanceInformation{
		ManufacturerProduct: &mp,
		SerialNumber:        &serialNumber,
	}

	return nil
}

// AddSoftwareArtifactComponent Add software artifact component information to an asset
func (d *DeviceInfo) AddSoftwareArtifactComponent(name string, version string, isFirmware bool) error {
	if err := validateField(name, "SoftwareName", "Software name is empty", "", ""); err != nil {
		return err
	}
	if err := validateField(version, "SoftwareVersion", "Software version is empty", "", ""); err != nil {
		return err
	}
	if err := validateField(
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

	softwareArtifactComponent := SoftwareArtifactComponent{}
	softwareArtifactComponent.Artifact = &softwareArtifact
	softwareArtifactComponent.SoftwareComponentType = SoftwareArtifactComponentSoftwareComponentTypeSoftwareArtifactComponent

	d.SoftwareComponents = append(d.SoftwareComponents, softwareArtifactComponent)
	return nil
}

// AddRunningSoftwareComponent Add running software component information to an asset
func (d *DeviceInfo) AddRunningSoftwareComponent(name string, version string, isFirmware bool, runningSoftwareId string) error {
	if err := validateField(name, "SoftwareName", "Software name is empty", "", ""); err != nil {
		return err
	}
	if err := validateField(version, "SoftwareVersion", "Software version is empty", "", ""); err != nil {
		return err
	}
	if err := validateField(runningSoftwareId, "RunningSoftwareId", "Running software ID is empty", "", ""); err != nil {
		return err
	}
	if err := validateField(
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

	runningSoftwareComponent := RunningSoftwareComponent{}
	runningSoftwareComponent.Artifact = &softwareArtifact
	runningSoftwareComponent.RunningSoftwareId = &runningSoftwareId
	runningSoftwareComponent.SoftwareComponentType = RunningSoftwareComponentSoftwareComponentTypeRunningSoftwareComponent

	d.SoftwareComponents = append(d.SoftwareComponents, runningSoftwareComponent)
	return nil
}
