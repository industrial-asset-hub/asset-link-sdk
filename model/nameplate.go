/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import (
	"github.com/google/uuid"
)

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
		Id:                        uuid.New().String(),
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
		Id:                        uuid.New().String(),
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
