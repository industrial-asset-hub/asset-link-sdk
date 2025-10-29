/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import "time"

func (d *DeviceInfo) addIdentifier(macAddress string) {

	if isNonEmptyValues(macAddress) {
		identifierUncertainty := 1
		identifier := MacIdentifier{
			IdentifierType:        nil,
			IdentifierUncertainty: &identifierUncertainty,
			MacAddress:            &macAddress,
		}
		d.MacIdentifiers = append(d.MacIdentifiers, identifier)
	}
}

// Add reachability state to the asset
func (d *DeviceInfo) addReachabilityState() {
	timestamp := getAssetCreationTimestamp(d.ManagementState.StateTimestamp)
	state := ReachabilityStateValuesReached

	reachabilityState := ReachabilityState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}

	d.ReachabilityState = &reachabilityState
}

func getAssetCreationTimestamp(timestamp *time.Time) time.Time {
	if timestamp != nil {
		return *timestamp
	}
	return time.Now().UTC()
}

func getAssetContext() *AssetContext {
	return &AssetContext{
		Lis:       "http://rds.posccaesar.org/ontology/lis14/rdl/",
		Base:      baseSchemaInContext,
		Skos:      "http://www.w3.org/2004/02/skos/core#",
		Vocab:     baseSchemaInContext,
		Linkml:    "https://w3id.org/linkml/",
		SchemaOrg: "https://schema.org/",
	}
}
