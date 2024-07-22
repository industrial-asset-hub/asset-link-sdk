/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

// Add nameplate to an discovered asset
func (d *DeviceInfo) AddNameplate(manufacturerName string,
	productArticleNumberOfManufacturer string,
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

// Add Firmware information to an discovered asset
func (d *DeviceInfo) AddFirmware(name string, version string) {
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
