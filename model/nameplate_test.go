/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
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

		productInfo, ok := m.ProductInstanceInformation.(*ProductInstanceInformation)
		if assert.True(t, ok) {
			manufacturerProduct, ok := productInfo.ManufacturerProduct.(*Product)
			if assert.True(t, ok) {
				manufacturer, ok := manufacturerProduct.Manufacturer.(*Organization)
				if assert.True(t, ok) {
					if assert.NotNil(t, manufacturer.Name) {
						assert.Equal(t, "ManufacturerCompany", *manufacturer.Name)
					}
				}

				if assert.NotNil(t, manufacturerProduct.ProductLink) {
					assert.Equal(t, testIDLink, *manufacturerProduct.ProductLink)
				}
				if assert.NotNil(t, manufacturerProduct.ProductVersion) {
					assert.Equal(t, "0.1.2", *manufacturerProduct.ProductVersion)
				}
				if assert.NotNil(t, manufacturerProduct.ProductId) {
					assert.Equal(t, "MyOrderNumber", *manufacturerProduct.ProductId)
				}
			}

			if assert.NotNil(t, productInfo.SerialNumber) {
				assert.Equal(t, "s-n-1.2.3", *productInfo.SerialNumber)
			}
		}

		idLinks := m.getIdLink()
		if !assert.Len(t, idLinks, 1) {
			return
		}
		found := 0
		for _, v := range idLinks {
			found++
			assert.Equal(t, testIDLink, v.IdLink)
			assert.Equal(t, IdLinkIdentifierAssetIdentifierTypeIdLinkIdentifier, v.AssetIdentifierType)
		}
		assert.Equal(t, 1, found)
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

		if !assert.Len(t, softwareArtifacts, 3) {
			return
		}

		fwFound := false
		sw1Found := false
		sw2Found := false
		for _, v := range softwareArtifacts {
			if assert.Len(t, v.AssetIdentifiers, 1) {
				softwareIdentifier, ok := v.AssetIdentifiers[0].(SoftwareIdentifier)
				if !assert.True(t, ok) {
					continue
				}

				switch softwareIdentifier.Name {
				case firmwareName:
					assert.Equal(t, firmwareVersion, softwareIdentifier.Version)
					assert.True(t, *v.IsFirmware)
					fwFound = true
				case sw1Name:
					assert.Equal(t, sw1Version, softwareIdentifier.Version)
					assert.False(t, *v.IsFirmware)
					sw1Found = true
				case sw2Name:
					assert.Equal(t, sw2Version, softwareIdentifier.Version)
					assert.False(t, *v.IsFirmware)
					sw2Found = true
				}
			}

			assert.Equal(t, FunctionalObjectSchemaUrl, v.FunctionalObjectSchemaUrl)
			assert.NotNil(t, v.IsFirmware)
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
		artifact, ok := v.(SoftwareArtifact)
		if ok {
			r = append(r, artifact)
		}
	}
	return r
}

// Extract IdLink Addresses from model
func (d *DeviceInfo) getIdLink() []IdLinkIdentifier {
	r := []IdLinkIdentifier{}
	for _, v := range d.AssetIdentifiers {
		identifier, ok := v.(IdLinkIdentifier)
		if ok {
			r = append(r, identifier)
		}
	}
	return r
}
