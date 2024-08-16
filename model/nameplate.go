/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

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
	UriOfTheProduct string,
	productArticleNumberOfManufacturer string,
	manufacturerProductDesignation string,
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
		Id:             UriOfTheProduct,
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

// AddSoftware Add software information to an asset
func (d *DeviceInfo) AddSoftware(name string, version string) {
	softwareIdentifier := SoftwareIdentifier{
		Name:    &name,
		Version: &version,
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
