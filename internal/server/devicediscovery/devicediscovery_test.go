/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package devicediscovery

import (
	"testing"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"
	"github.com/stretchr/testify/assert"
)

func TestSafeSerializeFilter(t *testing.T) {
	t.Run("Should serialize filter", func(t *testing.T) {
		filter := generated.ActiveFilter{
			Key:      "IPRange",
			Operator: generated.ComparisonOperator_EQUAL,
			Value:    &generated.Variant{Value: &generated.Variant_RawData{RawData: []byte("192.168.0.1-192.168.0.3")}},
		}

		outcome := serializeFilterOrOption([]*generated.ActiveFilter{&filter})

		assert.Contains(t, outcome, "IPRange")
		assert.Contains(t, outcome, "192.168.0.1-192.168.0.3")
		assert.Contains(t, outcome, "EQUAL")
	})
}

func TestSafeSerializeOption(t *testing.T) {
	t.Run("Should serialize option", func(t *testing.T) {
		filter := generated.ActiveOption{
			Key:      "OptionKey",
			Operator: generated.ComparisonOperator_EQUAL,
			Value:    &generated.Variant{Value: &generated.Variant_RawData{RawData: []byte("OptionValue")}},
		}

		outcome := serializeFilterOrOption([]*generated.ActiveOption{&filter})

		assert.Contains(t, outcome, "OptionKey")
		assert.Contains(t, outcome, "OptionValue")
		assert.Contains(t, outcome, "EQUAL")
	})
}
