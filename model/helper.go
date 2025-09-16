/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import "time"

func (d *DeviceInfo) addIdentifier(mac string) {

	if isNonEmptyValues(mac) {
		identifierUncertainty := 1
		identifier := MacIdentifier{
			IdentifierType:        nil,
			IdentifierUncertainty: &identifierUncertainty,
			MacAddress:            &mac,
		}
		d.MacIdentifiers = append(d.MacIdentifiers, identifier)
	}
}

// Add reachability state to the asset
func (d *DeviceInfo) addReachabilityState() {
	timestamp := d.getAssetCreationTimestamp()
	state := ReachabilityStateValuesReached

	reachabilityState := ReachabilityState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}

	d.ReachabilityState = &reachabilityState
}

func (d *DeviceInfo) getAssetCreationTimestamp() time.Time {
	if d.ManagementState.StateTimestamp != nil {
		return *d.ManagementState.StateTimestamp
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
