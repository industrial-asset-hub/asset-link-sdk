/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package config

import (
	"testing"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v2/generated/iah-discovery"

	"github.com/stretchr/testify/assert"
)

func createDummyConfig() DiscoveryConfig {
	dr := &generated.DiscoverRequest{}

	dr.Filters = []*generated.ActiveFilter{}
	dr.Options = []*generated.ActiveOption{}

	dr.Filters = append(dr.Filters, &generated.ActiveFilter{
		Key:      "filterkey",
		Operator: *generated.ComparisonOperator_EQUAL.Enum(),
		Value:    &generated.Variant{Value: &generated.Variant_Text{Text: "filtervalue"}},
	})

	dr.Filters = append(dr.Filters, &generated.ActiveFilter{
		Key:      "filterkey_multiple",
		Operator: *generated.ComparisonOperator_GREATER_THAN_OR_EQUAL_TO.Enum().Enum(),
		Value:    &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: 1}},
	})

	dr.Filters = append(dr.Filters, &generated.ActiveFilter{
		Key:      "filterkey_multiple",
		Operator: *generated.ComparisonOperator_GREATER_THAN_OR_EQUAL_TO.Enum().Enum(),
		Value:    &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: 5}},
	})

	dr.Options = append(dr.Options, &generated.ActiveOption{
		Key:      "optionkey",
		Operator: *generated.ComparisonOperator_EQUAL.Enum(),
		Value:    &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: 1102}},
	})

	dr.Options = append(dr.Options, &generated.ActiveOption{
		Key:      "optionkey_multiple",
		Operator: *generated.ComparisonOperator_GREATER_THAN_OR_EQUAL_TO.Enum().Enum(),
		Value:    &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: 1}},
	})

	dr.Options = append(dr.Options, &generated.ActiveOption{
		Key:      "optionkey_multiple",
		Operator: *generated.ComparisonOperator_LESS_THAN_OR_EQUAL_TO.Enum(),
		Value:    &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: 5}},
	})

	return NewDiscoveryConfigFromDiscoveryRequest(dr)
}

func TestDiscoveryConfig(t *testing.T) {
	t.Run("testLookupFiltersSucceeds", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		filters := discoveryConfig.GetFilters("filterkey")
		assert.NotEmpty(t, filters)
		assert.Len(t, filters, 1)

		filter := filters[0]
		assert.Equal(t, filter.GetKey(), "filterkey")
		assert.Equal(t, filter.GetValue().GetText(), "filtervalue")
		assert.Equal(t, filter.GetOperator().Type(), generated.ComparisonOperator_EQUAL.Type())
	})

	t.Run("testLookupFiltersFails", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		filters := discoveryConfig.GetFilters("filterkey_missing")
		assert.Empty(t, filters)
	})

	t.Run("testLookupAllFilters", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		filters := discoveryConfig.GetAllFilters()
		assert.NotEmpty(t, filters)
		assert.Len(t, filters, 3)
	})

	t.Run("testGetFilterSettingSucceeds", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		value, err := discoveryConfig.GetFilterSettingString("filterkey", "default")
		assert.NoError(t, err)
		assert.Equal(t, "filtervalue", value)
	})

	t.Run("testGetFilterSettingSucceedsDefault", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		value, err := discoveryConfig.GetFilterSettingString("filterkey_missing", "default")
		assert.NoError(t, err)
		assert.Equal(t, "default", value)
	})

	t.Run("testGetFilterSettingFailsWrongType", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		value, err := discoveryConfig.GetFilterSettingUint64("filterkey", 1234)
		assert.Error(t, err)
		assert.Equal(t, uint64(1234), value)
	})

	t.Run("testLookupOptionsSucceeds", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		options := discoveryConfig.GetOptions("optionkey")
		assert.NotEmpty(t, options)
		assert.Len(t, options, 1)

		option := options[0]
		assert.Equal(t, option.GetKey(), "optionkey")
		assert.Equal(t, option.GetValue().GetUint64Value(), uint64(1102))
		assert.Equal(t, option.GetOperator().Type(), generated.ComparisonOperator_NOT_EQUAL.Type())
	})

	t.Run("testLookupOptionsFails", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		options := discoveryConfig.GetOptions("missingoptionkey")
		assert.Empty(t, options)
	})

	t.Run("testLookupAllOptions", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		options := discoveryConfig.GetAllOptions()
		assert.NotEmpty(t, options)
		assert.Len(t, options, 3)
	})

	t.Run("testGetOptionSettingSucceeds", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		value, err := discoveryConfig.GetOptionSettingUint64("optionkey", 1234)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1102), value)

	})

	t.Run("testGetOptionSettingSucceedsDefault", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		value, err := discoveryConfig.GetOptionSettingUint64("missingoptionkey", 1234)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1234), value)
	})

	t.Run("testGetOptionSettingFailsWrongType", func(t *testing.T) {
		discoveryConfig := createDummyConfig()

		value, err := discoveryConfig.GetOptionSettingString("optionkey", "default")
		assert.Error(t, err)
		assert.Equal(t, "default", value)
	})
}
