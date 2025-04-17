/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameplate(t *testing.T) {
	t.Run("AddNameplate", func(t *testing.T) {
		m := NewDevice("asset", "device")

		m.AddNameplate(
			"ManufacturerCompany",
			"GuidOfTheProduct",
			"MyOrderNumber",
			"ProductFamily",
			"0.1.2",
			"s-n-1.2.3")

		// ManufacturerProductDesignation
		assert.Equal(t, "ManufacturerCompany", *m.ProductInstanceIdentifier.ManufacturerProduct.Manufacturer.Name)
		assert.Equal(t, "GuidOfTheProduct", m.ProductInstanceIdentifier.ManufacturerProduct.Id)
		assert.Equal(t, "ProductFamily", *m.ProductInstanceIdentifier.ManufacturerProduct.Name)
		assert.Equal(t, "0.1.2", *m.ProductInstanceIdentifier.ManufacturerProduct.ProductVersion)
		assert.Equal(t, "MyOrderNumber", *m.ProductInstanceIdentifier.ManufacturerProduct.ProductId)
		assert.Equal(t, "s-n-1.2.3", *m.ProductInstanceIdentifier.SerialNumber)

		idLinks := m.getIdLink()
		if len(idLinks) != 1 {
			fmt.Printf("Expected 1 id link, got %d\n", len(idLinks))
			t.Fail()
		}
		found := 0
		for _, v := range idLinks {
			found++
			assert.Equal(t, *v.IdLink, "GuidOfTheProduct")
		}
		assert.Equal(t, 1, found)

		// test legacy support for IAH
		legacyDescriptionOk := false
		for _, v := range m.InstanceAnnotations {
			if v.Key != nil && *v.Key == "description" {
				if v.Value != nil && *v.Value == "ProductFamily" {
					legacyDescriptionOk = true
				}
			}
		}
		assert.True(t, legacyDescriptionOk)
	})
}

func TestSoftwareNameplate(t *testing.T) {
	t.Run("AddFirmwareAndOtherSoftware", func(t *testing.T) {
		m := NewDevice("", "")

		firmwareName := "Firmware"
		firmwareVersion := "1.2.3"

		sw1Name := "SoftwareName1"
		sw1Version := "1.0.0"

		sw2Name := "SoftwareName2"
		sw2Version := "2.0.0"

		m.AddSoftware(firmwareName, firmwareVersion, true)
		m.AddSoftware(sw1Name, sw1Version, false)
		m.AddSoftware(sw2Name, sw2Version, false)

		softwareArtifacts := m.getSoftwareArtifacts()

		if len(softwareArtifacts) != 3 {
			fmt.Printf("Expected 3 software entries, got %d\n", len(softwareArtifacts))
			t.Fail()
		}

		fwFound := false
		sw1Found := false
		sw2Found := false
		for _, v := range softwareArtifacts {
			switch *v.SoftwareIdentifier.Name {
			case firmwareName:
				assert.Equal(t, firmwareVersion, *v.SoftwareIdentifier.Version)
				assert.True(t, *v.IsFirmware)
				fwFound = true
			case sw1Name:
				assert.Equal(t, sw1Version, *v.SoftwareIdentifier.Version)
				assert.False(t, *v.IsFirmware)
				sw1Found = true
			case sw2Name:
				assert.Equal(t, sw2Version, *v.SoftwareIdentifier.Version)
				assert.False(t, *v.IsFirmware)
				sw2Found = true
			}

			assert.NotEmpty(t, v.Id)
			stateValue := ManagementStateValuesRegarded
			assert.Equal(t, &stateValue, v.ManagementState.StateValue)
			assert.NotEmpty(t, v.ManagementState.StateTimestamp)
		}

		assert.True(t, fwFound)
		assert.True(t, sw1Found)
		assert.True(t, sw2Found)

		// test legacy support for IAH
		legacyFirmwareVersionOk := false
		for _, v := range m.InstanceAnnotations {
			if v.Key != nil && *v.Key == "firmware_version" {
				if v.Value != nil && *v.Value == firmwareVersion {
					legacyFirmwareVersionOk = true
				}
			}
		}
		assert.True(t, legacyFirmwareVersionOk)
	})
}

func (d *DeviceInfo) getSoftwareArtifacts() []SoftwareArtifact {
	r := []SoftwareArtifact{}
	for _, v := range d.SoftwareComponents {
		if reflect.TypeOf(v) == reflect.TypeOf(SoftwareArtifact{}) {
			r = append(r, v.(SoftwareArtifact))
		}
	}
	return r
}

// Extract IdLink Addresses from model
func (d *DeviceInfo) getIdLink() []IdLink {
	r := []IdLink{}
	for _, v := range d.AssetIdentifiers {
		if reflect.TypeOf(v) == reflect.TypeOf(IdLink{}) {
			r = append(r, v.(IdLink))
		}
	}
	return r
}
