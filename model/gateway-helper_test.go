/*
 * SPDX-FileCopyrightText: 2025 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGatewayConvenience(t *testing.T) {
	t.Run("AddTrustEstablishmentState - Valid Trusted State", func(t *testing.T) {
		g := &GatewayInfo{}
		beforeTime := time.Now().UTC()

		g.AddTrustEstablishmentState(TrustEstablishedStateValuesTrusted)

		afterTime := time.Now().UTC()

		assert.NotNil(t, g.TrustEstablishedState)
		assert.Equal(t, TrustEstablishedStateValuesTrusted, *g.TrustEstablishedState.StateValue)
		assert.NotNil(t, g.TrustEstablishedState.StateTimestamp)
		assert.True(t, g.TrustEstablishedState.StateTimestamp.After(beforeTime) || g.TrustEstablishedState.StateTimestamp.Equal(beforeTime))
		assert.True(t, g.TrustEstablishedState.StateTimestamp.Before(afterTime) || g.TrustEstablishedState.StateTimestamp.Equal(afterTime))
	})

	t.Run("AddTrustEstablishmentState - Valid Pending State", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddTrustEstablishmentState(TrustEstablishedStateValuesPending)

		assert.NotNil(t, g.TrustEstablishedState)
		assert.Equal(t, TrustEstablishedStateValuesPending, *g.TrustEstablishedState.StateValue)
		assert.NotNil(t, g.TrustEstablishedState.StateTimestamp)
	})

	t.Run("AddTrustEstablishmentState - Valid Failed State", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddTrustEstablishmentState(TrustEstablishedStateValuesFailed)

		assert.NotNil(t, g.TrustEstablishedState)
		assert.Equal(t, TrustEstablishedStateValuesFailed, *g.TrustEstablishedState.StateValue)
		assert.NotNil(t, g.TrustEstablishedState.StateTimestamp)
	})

	t.Run("AddTrustEstablishmentState - Empty State Value", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddTrustEstablishmentState(TrustEstablishedStateValues(""))

		assert.Nil(t, g.TrustEstablishedState)
	})

	t.Run("AddTrustEstablishmentState - Invalid State Value", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddTrustEstablishmentState(TrustEstablishedStateValues("invalid"))

		assert.Nil(t, g.TrustEstablishedState)
	})

	t.Run("AddProductInstanceIdentifier - Valid Parameters", func(t *testing.T) {
		g := &GatewayInfo{}
		productId := "PROD123"
		productVersion := "v1.0.0"
		productName := "Test Gateway"
		manufacturerName := "Test Manufacturer"
		serialNumber := "SN123456"

		g.AddProductInstanceIdentifier(productId, productVersion, productName, manufacturerName, serialNumber)

		assert.NotNil(t, g.ProductInstanceIdentifier)
		assert.Equal(t, serialNumber, *g.ProductInstanceIdentifier.SerialNumber)
		assert.NotNil(t, g.ProductInstanceIdentifier.ManufacturerProduct)
		assert.Equal(t, productName, *g.ProductInstanceIdentifier.ManufacturerProduct.Name)
		assert.Equal(t, productId, *g.ProductInstanceIdentifier.ManufacturerProduct.ProductId)
		assert.Equal(t, productVersion, *g.ProductInstanceIdentifier.ManufacturerProduct.ProductVersion)
		assert.NotNil(t, g.ProductInstanceIdentifier.ManufacturerProduct.Manufacturer)
		assert.Equal(t, manufacturerName, *g.ProductInstanceIdentifier.ManufacturerProduct.Manufacturer.Name)
	})

	t.Run("AddProductInstanceIdentifier - Empty Parameters", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddProductInstanceIdentifier("", "", "", "", "")

		assert.Nil(t, g.ProductInstanceIdentifier)
	})

	t.Run("AddReachabilityState - Valid Reached State", func(t *testing.T) {
		g := &GatewayInfo{}
		beforeTime := time.Now().UTC()

		g.AddReachabilityState(ReachabilityStateValuesReached)

		afterTime := time.Now().UTC()

		assert.NotNil(t, g.ReachabilityState)
		assert.Equal(t, ReachabilityStateValuesReached, *g.ReachabilityState.StateValue)
		assert.NotNil(t, g.ReachabilityState.StateTimestamp)
		assert.True(t, g.ReachabilityState.StateTimestamp.After(beforeTime) || g.ReachabilityState.StateTimestamp.Equal(beforeTime))
		assert.True(t, g.ReachabilityState.StateTimestamp.Before(afterTime) || g.ReachabilityState.StateTimestamp.Equal(afterTime))
	})

	t.Run("AddReachabilityState - Valid Failed State", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddReachabilityState(ReachabilityStateValuesFailed)

		assert.NotNil(t, g.ReachabilityState)
		assert.Equal(t, ReachabilityStateValuesFailed, *g.ReachabilityState.StateValue)
		assert.NotNil(t, g.ReachabilityState.StateTimestamp)
	})

	t.Run("AddReachabilityState - Valid Unknown State", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddReachabilityState(ReachabilityStateValuesUnknown)

		assert.NotNil(t, g.ReachabilityState)
		assert.Equal(t, ReachabilityStateValuesUnknown, *g.ReachabilityState.StateValue)
		assert.NotNil(t, g.ReachabilityState.StateTimestamp)
	})

	t.Run("AddReachabilityState - Empty State Value", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddReachabilityState(ReachabilityStateValues(""))

		assert.Nil(t, g.ReachabilityState)
	})

	t.Run("AddReachabilityState - Invalid State Value", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddReachabilityState(ReachabilityStateValues("invalid"))

		assert.Nil(t, g.ReachabilityState)
	})

	t.Run("AddTrustEstablishmentState - Override Previous State", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddTrustEstablishmentState(TrustEstablishedStateValuesTrusted)
		firstTimestamp := *g.TrustEstablishedState.StateTimestamp

		time.Sleep(1 * time.Millisecond) // Ensure different timestamps

		g.AddTrustEstablishmentState(TrustEstablishedStateValuesFailed)

		assert.NotNil(t, g.TrustEstablishedState)
		assert.Equal(t, TrustEstablishedStateValuesFailed, *g.TrustEstablishedState.StateValue)
		assert.True(t, g.TrustEstablishedState.StateTimestamp.After(firstTimestamp))
	})

	t.Run("AddReachabilityState - Override Previous State", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddReachabilityState(ReachabilityStateValuesReached)
		firstTimestamp := *g.ReachabilityState.StateTimestamp

		time.Sleep(1 * time.Millisecond) // Ensure different timestamps

		g.AddReachabilityState(ReachabilityStateValuesFailed)

		assert.NotNil(t, g.ReachabilityState)
		assert.Equal(t, ReachabilityStateValuesFailed, *g.ReachabilityState.StateValue)
		assert.True(t, g.ReachabilityState.StateTimestamp.After(firstTimestamp))
	})

	t.Run("AddRunningSoftwareType - Valid CdmGateway Type", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddRunningSoftwareType(RunningSoftwareValuesCdmGateway)

		assert.NotNil(t, g.RunningSoftwareType)
		assert.Equal(t, RunningSoftwareValuesCdmGateway, *g.RunningSoftwareType)
	})

	t.Run("AddRunningSoftwareType - Valid IahGateway Type", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddRunningSoftwareType(RunningSoftwareValuesIahGateway)

		assert.NotNil(t, g.RunningSoftwareType)
		assert.Equal(t, RunningSoftwareValuesIahGateway, *g.RunningSoftwareType)
	})

	t.Run("AddRunningSoftwareType - Valid Other Type", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddRunningSoftwareType(RunningSoftwareValuesOther)

		assert.NotNil(t, g.RunningSoftwareType)
		assert.Equal(t, RunningSoftwareValuesOther, *g.RunningSoftwareType)
	})

	t.Run("AddRunningSoftwareType - Empty Software Type", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddRunningSoftwareType(RunningSoftwareValues(""))

		assert.Nil(t, g.RunningSoftwareType)
	})

	t.Run("AddRunningSoftwareType - Invalid Software Type", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddRunningSoftwareType(RunningSoftwareValues("invalid_type"))

		assert.Nil(t, g.RunningSoftwareType)
	})

	t.Run("AddRunningSoftwareType - Override Previous Type", func(t *testing.T) {
		g := &GatewayInfo{}

		g.AddRunningSoftwareType(RunningSoftwareValuesCdmGateway)
		assert.NotNil(t, g.RunningSoftwareType)
		assert.Equal(t, RunningSoftwareValuesCdmGateway, *g.RunningSoftwareType)

		g.AddRunningSoftwareType(RunningSoftwareValuesOther)
		assert.NotNil(t, g.RunningSoftwareType)
		assert.Equal(t, RunningSoftwareValuesOther, *g.RunningSoftwareType)
	})
}
