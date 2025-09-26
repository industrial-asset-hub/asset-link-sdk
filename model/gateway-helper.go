/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"time"

	"github.com/rs/zerolog/log"
)

// Add TrustEstablishmentState to the gateway. Allowed values are
// TrustEstablishedStateValuesTrusted ("trusted")
// TrustEstablishedStateValuesPending ("pending")
// TrustEstablishedStateValuesFailed ("failed")
func (d *GatewayInfo) AddTrustEstablishmentState(stateValue TrustEstablishedStateValues) {
	if !isNonEmptyValues(string(stateValue)) {
		log.Warn().Msg("Trust establishment state value is empty")
		return
	}
	if stateValue != TrustEstablishedStateValuesTrusted && stateValue != TrustEstablishedStateValuesPending && stateValue != TrustEstablishedStateValuesFailed {
		log.Warn().Msgf("Trust establishment state value %s is not valid", stateValue)
		return
	}

	currentTime := time.Now().UTC()
	trustEstState := TrustEstablishedState{
		StateTimestamp: &currentTime,
		StateValue:     &stateValue,
	}

	d.TrustEstablishedState = &trustEstState
}

// Add ProductInstanceIdentifier to the gateway
func (d *GatewayInfo) AddProductInstanceIdentifier(productId, productVersion, productName,
	manufacturerName, serialNumber string) {

	if !isNonEmptyValues(productId, productVersion, productName, manufacturerName, serialNumber) {
		log.Warn().Msg("One or more ProductInstanceIdentifier values are empty")
		return
	}

	d.ProductInstanceIdentifier = &ProductSerialIdentifier{
		SerialNumber: &serialNumber,
		ManufacturerProduct: &Product{
			Name:      &productName,
			ProductId: &productId,
			Manufacturer: &Organization{
				Name: &manufacturerName,
			},
			ProductVersion: &productVersion,
		},
	}
}

// Add Reachability state to the gateway. Allowed values are
// ReachabilityStateValuesReached ("reached")
// ReachabilityStateValuesFailed ("failed")
// ReachabilityStateValuesUnknown ("unknown")
func (d *GatewayInfo) AddReachabilityState(stateValue ReachabilityStateValues) {
	if !isNonEmptyValues(string(stateValue)) {
		log.Warn().Msg("Reachability state value is empty")
		return
	}
	if stateValue != ReachabilityStateValuesReached && stateValue != ReachabilityStateValuesFailed && stateValue != ReachabilityStateValuesUnknown {
		log.Warn().Msgf("Reachability state value %s is not valid", stateValue)
		return
	}

	timestamp := time.Now().UTC()
	d.ReachabilityState = &ReachabilityState{
		StateTimestamp: &timestamp,
		StateValue:     &stateValue,
	}
}

// Add RunningSoftwareType to the gateway. Allowed values are
// RunningSoftwareValuesCdmGateway ("cdm_gateway")
// RunningSoftwareValuesIahGateway ("iah_gateway")
// RunningSoftwareValuesOther ("other")
func (d *GatewayInfo) AddRunningSoftwareType(softwareType RunningSoftwareValues) {
	if !isNonEmptyValues(string(softwareType)) {
		log.Warn().Msg("Running software type is empty")
		return
	}
	if softwareType != RunningSoftwareValuesCdmGateway && softwareType != RunningSoftwareValuesIahGateway && softwareType != RunningSoftwareValuesOther {
		log.Warn().Msgf("Running software type %s is not valid", softwareType)
		return
	}

	d.RunningSoftwareType = &softwareType
}
