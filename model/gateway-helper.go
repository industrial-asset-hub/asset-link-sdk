/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"time"
)

// Add TrustEstablishmentState to the gateway. Allowed values are
// TrustEstablishedStateValuesTrusted ("trusted")
// TrustEstablishedStateValuesPending ("pending")
// TrustEstablishedStateValuesFailed ("failed")
func (d *GatewayInfo) AddTrustEstablishmentState(stateValue TrustEstablishedStateValues) error {

	if !isNonEmptyValues(string(stateValue)) {
		err := &EmptyError{
			Field:   "TrustEstablishmentState",
			Message: "Trust establishment state value is empty",
			Value:   stateValue,
		}
		return err
	}

	if stateValue != TrustEstablishedStateValuesTrusted && stateValue != TrustEstablishedStateValuesPending && stateValue != TrustEstablishedStateValuesFailed {
		err := &PermissibleValuesError{
			Field:   "TrustEstablishmentState",
			Value:   stateValue,
			Allowed: []interface{}{TrustEstablishedStateValuesTrusted, TrustEstablishedStateValuesPending, TrustEstablishedStateValuesFailed},
		}
		return err
	}

	currentTime := time.Now().UTC()
	trustEstState := TrustEstablishedState{
		StateTimestamp: &currentTime,
		StateValue:     &stateValue,
	}

	d.TrustEstablishedState = &trustEstState
	return nil
}

// Add ProductInstanceIdentifier to the gateway
func (d *GatewayInfo) AddProductInstanceIdentifier(productId, productVersion, productName,
	manufacturerName, serialNumber string) error {

	if !isNonEmptyValues(productId, productVersion, productName, manufacturerName, serialNumber) {
		err := &EmptyError{
			Field:   "ProductInstanceIdentifier",
			Message: "One or more required fields for ProductInstanceIdentifier are empty",
			Value: map[string]string{
				"productId":        productId,
				"productVersion":   productVersion,
				"productName":      productName,
				"manufacturerName": manufacturerName,
				"serialNumber":     serialNumber,
			},
		}
		return err
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
	return nil
}

// Add Reachability state to the gateway. Allowed values are
// ReachabilityStateValuesReached ("reached")
// ReachabilityStateValuesFailed ("failed")
// ReachabilityStateValuesUnknown ("unknown")
func (d *GatewayInfo) AddReachabilityState(stateValue ReachabilityStateValues) error {
	if !isNonEmptyValues(string(stateValue)) {
		err := &EmptyError{
			Field:   "ReachabilityState",
			Message: "Reachability state value is empty",
			Value:   stateValue,
		}
		return err
	}
	if stateValue != ReachabilityStateValuesReached && stateValue != ReachabilityStateValuesFailed && stateValue != ReachabilityStateValuesUnknown {
		err := &PermissibleValuesError{
			Field:   "ReachabilityState",
			Value:   stateValue,
			Allowed: []interface{}{ReachabilityStateValuesReached, ReachabilityStateValuesFailed, ReachabilityStateValuesUnknown},
		}
		return err
	}

	timestamp := time.Now().UTC()
	d.ReachabilityState = &ReachabilityState{
		StateTimestamp: &timestamp,
		StateValue:     &stateValue,
	}
	return nil
}

// Add RunningSoftwareType to the gateway. Allowed values are
// RunningSoftwareValuesCdmGateway ("cdm_gateway")
// RunningSoftwareValuesIahGateway ("iah_gateway")
// RunningSoftwareValuesOther ("other")
func (d *GatewayInfo) AddRunningSoftwareType(softwareType RunningSoftwareValues) error {
	if !isNonEmptyValues(string(softwareType)) {
		err := &EmptyError{
			Field:   "RunningSoftwareType",
			Message: "Running software type value is empty",
			Value:   softwareType,
		}
		return err
	}
	if softwareType != RunningSoftwareValuesCdmGateway && softwareType != RunningSoftwareValuesIahGateway && softwareType != RunningSoftwareValuesOther {
		err := &PermissibleValuesError{
			Field:   "RunningSoftwareType",
			Value:   softwareType,
			Allowed: []interface{}{RunningSoftwareValuesCdmGateway, RunningSoftwareValuesIahGateway, RunningSoftwareValuesOther},
		}
		return err
	}

	d.RunningSoftwareType = &softwareType
	return nil
}
