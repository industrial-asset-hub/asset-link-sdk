/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNameplate(t *testing.T) {
	t.Run("AddNameplate", func(t *testing.T) {
		m := NewDevice("asset", "device")

		m.AddNameplate(
			"ManufacturerCompany",
			"MyOrderNumber",
			"ProductFamily",
			"0.1.2",
			"s-n-1.2.3")

		//ManufacturerProductDesignation
		assert.Equal(t, "ManufacturerCompany", *m.ProductInstanceIdentifier.ManufacturerProduct.Manufacturer.Name)
		assert.Equal(t, "ProductFamily", *m.ProductInstanceIdentifier.ManufacturerProduct.Name)
		assert.Equal(t, "0.1.2", *m.ProductInstanceIdentifier.ManufacturerProduct.ProductVersion)
		assert.Equal(t, "MyOrderNumber", *m.ProductInstanceIdentifier.ManufacturerProduct.ProductId)
		assert.Equal(t, "s-n-1.2.3", *m.ProductInstanceIdentifier.SerialNumber)
	})
}
