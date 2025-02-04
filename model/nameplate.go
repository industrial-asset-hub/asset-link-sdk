/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
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
func (d *DeviceInfo) AddNameplate(manufacturerName string, uriOfTheProduct string,
	productArticleNumberOfManufacturer string, manufacturerProductDesignation string, hardwareVersion string, serialNumber string) error {

	// URI of the product is a required property
	if !checkIfAnyValueIsNonEmpty(uriOfTheProduct) {
		return errors.New("URI of the product should not be empty")
	}

	// this check ensures if any field is non-empty then its value is set
	if checkIfAnyValueIsNonEmpty(manufacturerName, productArticleNumberOfManufacturer, manufacturerProductDesignation, hardwareVersion, serialNumber) {

		// We hash the manufacturer to get a unique identifier
		h := sha1.New()
		h.Write([]byte(manufacturerName))
		manufacturerId := hex.EncodeToString(h.Sum(nil))

		organisation := Organization{
			Address:        nil,
			AlternateNames: nil,
			ContactPoint:   nil,
			Id:             manufacturerId,
			Name:           &manufacturerName,
		}

		mp := Product{
			Id:             uriOfTheProduct,
			Manufacturer:   &organisation,
			Name:           &manufacturerProductDesignation,
			ProductId:      &productArticleNumberOfManufacturer,
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
	return nil
}

// AddSoftware Add software information to an asset
func (d *DeviceInfo) AddSoftware(name string, version string) {
	softwareIdentifier := SoftwareIdentifier{}

	if checkIfAnyValueIsNonEmpty(name, version) {
		softwareIdentifier.Name = &name
		softwareIdentifier.Version = &version
	}

	softwareArtifact := SoftwareArtifact{
		AssetOperations:           nil,
		ChecksumIdentifier:        nil,
		ConnectionPoints:          nil,
		CustomUiProperties:        nil,
		FunctionalParts:           nil,
		Id:                        "",
		InstanceAnnotations:       nil,
		ManagementState:           ManagementState{},
		Name:                      nil,
		OtherStates:               nil,
		ProductInstanceIdentifier: nil,
		ReachabilityState:         nil,
		SoftwareComponents:        nil,
		SoftwareIdentifier:        &softwareIdentifier,
	}

	runningSoftware := RunningSoftware{
		Artifact:                  &softwareArtifact,
		AssetOperations:           nil,
		ConnectionPoints:          nil,
		CustomRunningSoftwareType: nil,
		CustomUiProperties:        nil,
		FunctionalParts:           nil,
		Id:                        "",
		InstanceAnnotations:       nil,
		ManagementState:           ManagementState{},
		Name:                      nil,
		OtherStates:               nil,
		ProductInstanceIdentifier: nil,
		ReachabilityState:         nil,
		RunningSoftwareType:       nil,
		RunningSwId:               nil,
		SoftwareComponents:        nil,
	}

	d.SoftwareComponents = append(d.SoftwareComponents, runningSoftware)
}
