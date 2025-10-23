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

	"github.com/google/uuid"
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
) {
	if isNonEmptyValues(manufacturerName, uriOfTheProduct, productArticleNumberOfManufacturer, manufacturerProductDesignation, hardwareVersion, serialNumber) {

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

		// Duplicate IDLink field to explict field
		t := IdLinkAssetIdentifierTypeIdLink
		idLink := IdLink{
			AssetIdentifierType:   &t,
			IdLink:                &uriOfTheProduct,
			IdentifierType:        nil,
			IdentifierUncertainty: nil,
		}
		d.AssetIdentifiers = append(d.AssetIdentifiers, idLink)

		d.addInstanceAnnotation("description", manufacturerProductDesignation) // legacy support for IAH (should be removed eventually)
	}
}

// AddSoftware Add software information to an asset
func (d *DeviceInfo) AddSoftware(name string, version string, isFirmware bool) {
	if isNonEmptyValues(name, version) {
		softwareIdentifier := SoftwareIdentifier{}
		softwareIdentifier.Name = &name
		softwareIdentifier.Version = &version

		softwareArtifactId := uuid.New().String()

		stateValue := ManagementStateValuesRegarded
		stateTimestamp := getAssetCreationTimestamp(d.ManagementState.StateTimestamp)

		softwareArtifact := SoftwareArtifact{
			Id:                  softwareArtifactId,
			AssetOperations:     nil,
			ChecksumIdentifier:  nil,
			ConnectionPoints:    nil,
			CustomUiProperties:  nil,
			FunctionalParts:     nil,
			InstanceAnnotations: nil,
			ManagementState: ManagementState{
				StateTimestamp: &stateTimestamp,
				StateValue:     &stateValue,
			},
			Name:                      nil,
			OtherStates:               nil,
			ProductInstanceIdentifier: nil,
			ReachabilityState:         nil,
			SoftwareComponents:        nil,
			SoftwareIdentifier:        &softwareIdentifier,
			IsFirmware:                &isFirmware,
		}

		d.SoftwareComponents = append(d.SoftwareComponents, softwareArtifact)

		if isFirmware {
			// legacy support for IAH (should be removed eventually)
			d.addInstanceAnnotation("firmware_version", version)
		}
	}
}

// instance annotations should only be used for backwards-compatibilty
func (d *DeviceInfo) addInstanceAnnotation(key, value string) {
	for index := range d.InstanceAnnotations {
		annotation := &d.InstanceAnnotations[index]
		if annotation.Key != nil && *annotation.Key == key {
			annotation.Value = &value
			return
		}
	}

	newAnnotation := InstanceAnnotation{
		Key:   &key,
		Value: &value,
	}

	d.InstanceAnnotations = append(d.InstanceAnnotations, newAnnotation)
}
