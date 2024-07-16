/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

func (d *DeviceInfo) addIdentifier(mac string) {

	identifierUncertainty := 1
	identifier := MacIdentifier{
		IdentifierType:        nil,
		IdentifierUncertainty: &identifierUncertainty,
		MacAddress:            &mac,
	}
	d.MacIdentifiers = append(d.MacIdentifiers, identifier)
}

// Add reachability state to the asset
func (d *DeviceInfo) addReachabilityState() {
	timestamp := CreateTimestamp()
	state := ReachabilityStateValuesReached

	reachabilityState := ReachabilityState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}

	d.ReachabilityState = &reachabilityState
}

// Add Management state to the asset
// Only used internal
func (d *DeviceInfo) addManagementState() {

	timestamp := CreateTimestamp()
	state := ManagementStateValuesUnknown

	mgmtState := ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}

	d.ManagementState = mgmtState
}