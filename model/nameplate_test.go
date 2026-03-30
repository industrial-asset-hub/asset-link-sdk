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

const testIDLink = "https://i.siemens.com/1P6ES7131-6BF00-0CA0+SC-P4TM3526"

func TestNameplate(t *testing.T) {
	t.Run("AddNameplate", func(t *testing.T) {
		m, err := NewDevice("asset", "device")
		assert.NoError(t, err)

		err = m.AddNameplate(
			"ManufacturerCompany",
			testIDLink,
			"MyOrderNumber",
			"ProductFamily",
			"0.1.2",
			"s-n-1.2.3")
		assert.NoError(t, err)

		// ManufacturerProductDesignation
		assert.Equal(t, "ManufacturerCompany", *m.ProductInstanceIdentifier.ManufacturerProduct.Manufacturer.Name)
		assert.Equal(t, testIDLink, m.ProductInstanceIdentifier.ManufacturerProduct.Id)
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
			assert.Equal(t, testIDLink, *v.IdLink)
			assert.Equal(t, IdLinkAssetIdentifierTypeIdLink, *v.AssetIdentifierType)
		}
		assert.Equal(t, 1, found)
	})

	t.Run("AddNameplate_InvalidURI_DoesNotAppendIdLink", func(t *testing.T) {
		m, err := NewDevice("asset", "device")
		assert.NoError(t, err)

		err = m.AddNameplate(
			"ManufacturerCompany",
			"https://example.com/not-a-siemens-id-link",
			"MyOrderNumber",
			"ProductFamily",
			"0.1.2",
			"s-n-1.2.3")
		assert.Error(t, err)

		var validationErr *ValidationError
		assert.ErrorAs(t, err, &validationErr)
		assert.Equal(t, "URIOfTheProduct", validationErr.Field)
		assert.Equal(t, "Idlink must be a valid URI", validationErr.Message)
		assert.Empty(t, m.getIdLink())
		assert.Len(t, m.AssetIdentifiers, 0)
	})
}

func TestSoftwareNameplate(t *testing.T) {
	t.Run("AddFirmwareAndOtherSoftware", func(t *testing.T) {
		m, err := NewDevice("asset", "device")
		assert.NoError(t, err)

		firmwareName := "Firmware"
		firmwareVersion := "1.2.3"

		sw1Name := "SoftwareName1"
		sw1Version := "1.0.0"

		sw2Name := "SoftwareName2"
		sw2Version := "2.0.0"

		err = m.AddSoftware(firmwareName, firmwareVersion, true)
		assert.NoError(t, err)
		err = m.AddSoftware(sw1Name, sw1Version, false)
		assert.NoError(t, err)
		err = m.AddSoftware(sw2Name, sw2Version, false)
		assert.NoError(t, err)

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
	})

	t.Run("AddSoftware_EmptyNameOrVersion", func(t *testing.T) {
		m, err := NewDevice("asset", "device")
		assert.NoError(t, err)
		err = m.AddSoftware("", "1.0.0", false)
		assert.Error(t, err)
		err = m.AddSoftware("SoftwareName", "", false)
		assert.Error(t, err)
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
