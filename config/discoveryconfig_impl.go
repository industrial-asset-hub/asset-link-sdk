/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	generated "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/iah-discovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

type discoveryConfigImplementation struct {
	discoveryRequest *generated.DiscoverRequest
	filterMap        map[string][]*generated.ActiveFilter
	optionMap        map[string][]*generated.ActiveOption
}

func (d *discoveryConfigImplementation) GetAllFilters() []*generated.ActiveFilter {
	return d.discoveryRequest.GetFilters()
}

func (d *discoveryConfigImplementation) GetFilters(filterKey string) []*generated.ActiveFilter {
	filters, ok := d.filterMap[filterKey]
	if ok {
		return filters
	}
	return []*generated.ActiveFilter{}
}

func (d *discoveryConfigImplementation) internalGetFilterSetting(filterKey string, defaultValue *generated.Variant) (*generated.Variant, error) {
	filters := d.GetFilters(filterKey)
	if len(filters) == 0 {
		return defaultValue, nil
	} else if len(filters) > 1 {
		return defaultValue, errors.New("More than one filter for setting " + filterKey)
	}

	filter := filters[0]
	if filter.GetOperator().Type() != generated.ComparisonOperator_EQUAL.Type() {
		return defaultValue, errors.New("Operation for filter setting " + filterKey + " is not EQUAL")
	}

	if reflect.TypeOf(filter.GetValue().Value) != reflect.TypeOf(defaultValue.Value) {
		return defaultValue, errors.New("Type for filter setting " + filterKey + " does not match")
	}

	return filter.GetValue(), nil
}

func (d *discoveryConfigImplementation) GetFilterSettingString(filterKey string, defaultValue string) (string, error) {
	variant, err := d.internalGetFilterSetting(filterKey, &generated.Variant{Value: &generated.Variant_Text{Text: defaultValue}})
	return variant.GetText(), err
}

func (d *discoveryConfigImplementation) GetFilterSettingUint64(filterKey string, defaultValue uint64) (uint64, error) {
	variant, err := d.internalGetFilterSetting(filterKey, &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: defaultValue}})
	return variant.GetUint64Value(), err
}

func (d *discoveryConfigImplementation) GetFilterSettingInt64(filterKey string, defaultValue int64) (int64, error) {
	variant, err := d.internalGetFilterSetting(filterKey, &generated.Variant{Value: &generated.Variant_Int64Value{Int64Value: defaultValue}})
	return variant.GetInt64Value(), err
}

func (d *discoveryConfigImplementation) GetFilterSettingFloat64(filterKey string, defaultValue float64) (float64, error) {
	variant, err := d.internalGetFilterSetting(filterKey, &generated.Variant{Value: &generated.Variant_Float64Value{Float64Value: defaultValue}})
	return variant.GetFloat64Value(), err
}

func (d *discoveryConfigImplementation) GetAllOptions() []*generated.ActiveOption {
	return d.discoveryRequest.GetOptions()
}

func (d *discoveryConfigImplementation) GetOptions(optionKey string) []*generated.ActiveOption {
	options, ok := d.optionMap[optionKey]
	if ok {
		return options
	}
	return []*generated.ActiveOption{}
}

func (d *discoveryConfigImplementation) internalGetOptionSetting(optionKey string, defaultValue *generated.Variant) (*generated.Variant, error) {
	options := d.GetOptions(optionKey)
	if len(options) == 0 {
		return defaultValue, nil
	} else if len(options) > 1 {
		return defaultValue, errors.New("More than one option for setting " + optionKey)
	}

	option := options[0]
	if option.GetOperator().Type() != generated.ComparisonOperator_EQUAL.Type() {
		return defaultValue, errors.New("Operation for option setting " + optionKey + " is not EQUAL")
	}

	if reflect.TypeOf(option.GetValue().Value) != reflect.TypeOf(defaultValue.Value) {
		return defaultValue, errors.New("Type for option setting " + optionKey + " does not match")
	}

	return option.GetValue(), nil
}

func (d *discoveryConfigImplementation) GetOptionSettingString(optionKey string, defaultValue string) (string, error) {
	variant, err := d.internalGetOptionSetting(optionKey, &generated.Variant{Value: &generated.Variant_Text{Text: defaultValue}})
	return variant.GetText(), err
}

func (d *discoveryConfigImplementation) GetOptionSettingUint64(optionKey string, defaultValue uint64) (uint64, error) {
	variant, err := d.internalGetOptionSetting(optionKey, &generated.Variant{Value: &generated.Variant_Uint64Value{Uint64Value: defaultValue}})
	return variant.GetUint64Value(), err
}

func (d *discoveryConfigImplementation) GetOptionSettingInt64(optionKey string, defaultValue int64) (int64, error) {
	variant, err := d.internalGetOptionSetting(optionKey, &generated.Variant{Value: &generated.Variant_Int64Value{Int64Value: defaultValue}})
	return variant.GetInt64Value(), err
}

func (d *discoveryConfigImplementation) GetOptionSettingFloat64(optionKey string, defaultValue float64) (float64, error) {
	variant, err := d.internalGetOptionSetting(optionKey, &generated.Variant{Value: &generated.Variant_Float64Value{Float64Value: defaultValue}})
	return variant.GetFloat64Value(), err
}

// func (d *DiscoveryConfigImplementation) GetTarget() []*generated.Destination {
// 	return d.discoveryRequest.GetTarget()
// }

func (d *discoveryConfigImplementation) initLookup() {
	d.filterMap = map[string][]*generated.ActiveFilter{}
	d.optionMap = map[string][]*generated.ActiveOption{}

	for _, filter := range d.discoveryRequest.GetFilters() {
		list, contains := d.filterMap[filter.Key]
		if !contains {
			d.filterMap[filter.Key] = []*generated.ActiveFilter{}
		}
		d.filterMap[filter.Key] = append(list, filter)
	}

	for _, option := range d.discoveryRequest.GetOptions() {
		list, contains := d.optionMap[option.Key]
		if !contains {
			d.optionMap[option.Key] = []*generated.ActiveOption{}
		}
		d.optionMap[option.Key] = append(list, option)
	}
}

func (d *discoveryConfigImplementation) GetDiscoveryRequest() *generated.DiscoverRequest {
	return d.discoveryRequest
}

func (d *discoveryConfigImplementation) String() string {
	return fmt.Sprintf("%+v\n", d.discoveryRequest)
}

func (d *discoveryConfigImplementation) JSON() (string, error) {
	result, mErr := protojson.Marshal(d.discoveryRequest)
	if mErr != nil {
		log.Err(mErr).Msg("Marshal Error")
		return "", mErr
	}

	stringResult := string(result[:])
	log.Info().Str("Marshal Result", stringResult).Msg("Marshal Result")
	return stringResult, nil
}

func NewDiscoveryConfigFromDiscoveryRequest(discoveryRequest *generated.DiscoverRequest) *discoveryConfigImplementation {
	dr := &discoveryConfigImplementation{discoveryRequest: discoveryRequest}
	dr.initLookup()

	return dr
}

func NewDiscoveryConfigFromFile(discoveryFile string) (*discoveryConfigImplementation, error) {
	discoveryRequest := &generated.DiscoverRequest{
		Options: []*generated.ActiveOption{},
		Filters: []*generated.ActiveFilter{},
		Target:  nil,
	}

	if discoveryFile != "" {
		file, openErr := os.Open(discoveryFile)

		if openErr != nil {
			return nil, openErr
		}

		defer file.Close()

		configReader := bufio.NewReader(file)

		configBuffer, readErr := io.ReadAll(configReader)
		if readErr != nil {
			return nil, readErr
		}

		unmarshalErr := protojson.Unmarshal(configBuffer, discoveryRequest)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
	}

	discoveryConfig := NewDiscoveryConfigFromDiscoveryRequest(discoveryRequest)
	discoveryConfig.initLookup()

	return discoveryConfig, nil
}

func NewDiscoveryConfigWithDefaults() *discoveryConfigImplementation {
	discoveryRequest := &generated.DiscoverRequest{
		Options: []*generated.ActiveOption{},
		Filters: []*generated.ActiveFilter{},
		Target:  nil,
	}

	discoveryConfig := NewDiscoveryConfigFromDiscoveryRequest(discoveryRequest)
	discoveryConfig.initLookup()

	return discoveryConfig
}
