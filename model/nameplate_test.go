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
	})
}

func TestSoftwareNameplate(t *testing.T) {
	t.Run("AddFirmware", func(t *testing.T) {
		m := NewDevice("", "")

		m.AddSoftware("ArtifactName", "0.1.2")
		m.AddSoftware("ArtifactName1", "2.1.3")

		firmware := m.getFirmware()
		if len(firmware) != 2 {
			fmt.Printf("Expected 1 added firmware, got %d\n", len(firmware))
			t.Fail()
		}
		found := 0
		for _, v := range firmware {
			if *v.Artifact.SoftwareIdentifier.Name == "ArtifactName" {
				found++
				assert.Equal(t, "ArtifactName", *v.Artifact.SoftwareIdentifier.Name)
				assert.Equal(t, "0.1.2", *v.Artifact.SoftwareIdentifier.Version)
				break
			}
		}

		assert.Equal(t, 1, found)
	})
}

func (d *DeviceInfo) getFirmware() []RunningSoftware {

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
