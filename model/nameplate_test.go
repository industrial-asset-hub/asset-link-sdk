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
	t.Run("AddFirmware", func(t *testing.T) {
		m := NewDevice("", "")

		firmwareVersion := "1.2.3"

		a1Version := "1.0.0"
		a2Version := "2.0.0"

		m.AddSoftware("firmware", firmwareVersion)

		m.AddSoftware("ArtifactName1", a1Version)
		m.AddSoftware("ArtifactName2", a2Version)

		runningSoftware := m.getRunningSoftware()

		if len(runningSoftware) != 3 {
			fmt.Printf("Expected 3 software entries, got %d\n", len(runningSoftware))
			t.Fail()
		}

		fwFound := false
		a1Found := false
		a2Found := false
		for _, v := range runningSoftware {
			switch *v.Artifact.SoftwareIdentifier.Name {
			case "firmware":
				assert.Equal(t, firmwareVersion, *v.Artifact.SoftwareIdentifier.Version)
				fwFound = true
			case "ArtifactName1":
				assert.Equal(t, a1Version, *v.Artifact.SoftwareIdentifier.Version)
				a1Found = true
			case "ArtifactName2":
				assert.Equal(t, a2Version, *v.Artifact.SoftwareIdentifier.Version)
				a2Found = true
			}
		}
		assert.True(t, fwFound)
		assert.True(t, a1Found)
		assert.True(t, a2Found)

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

func (d *DeviceInfo) getRunningSoftware() []RunningSoftware {

	r := []RunningSoftware{}
	for _, v := range d.SoftwareComponents {
		if reflect.TypeOf(v) == reflect.TypeOf(RunningSoftware{}) {
			r = append(r, v.(RunningSoftware))
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
