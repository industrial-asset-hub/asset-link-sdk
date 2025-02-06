/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

func (d *DeviceInfo) addIdentifier(mac string) {

	if checkIfAnyValueIsNonEmpty(mac) {
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
	timestamp := createTimestamp()
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

	timestamp := createTimestamp()
	state := ManagementStateValuesUnknown

	mgmtState := ManagementState{
		StateTimestamp: &timestamp,
		StateValue:     &state,
	}

	d.ManagementState = mgmtState
}
