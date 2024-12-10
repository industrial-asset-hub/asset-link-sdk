/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package server

import (
	"testing"
	"time"

	pb "github.com/industrial-asset-hub/asset-link-sdk/v3/generated/conn_suite_registry"
	"github.com/stretchr/testify/assert"
)

func TestRemoveExpiredService(t *testing.T) {
	t.Run("serviceIsExpired", func(t *testing.T) {
		// Adding -10 to be on the safe side
		assert.True(t, checkExpired(serviceEntry{
			driver:    &pb.ServiceInfo{AppInstanceId: "driver-1"},
			timeAdded: time.Now().Add(-time.Duration(serviceExpireTime+10) * time.Second)}))
	})

	t.Run("serviceIsNotExpired", func(t *testing.T) {
		assert.False(t, checkExpired(serviceEntry{
			driver:    &pb.ServiceInfo{AppInstanceId: "driver-1"},
			timeAdded: time.Now().Add(time.Duration(200) * time.Second)}))
	})
}

func TestRemoveDuplicates(t *testing.T) {
	t.Run("noDuplicates", func(t *testing.T) {
		elements := []string{"a", "b", "c", "d"}
		expected := []string{"a", "b", "c", "d"}

		result := removeDuplicates(elements)

		assert.Equal(t, expected, result)
	})

	t.Run("withDuplicates", func(t *testing.T) {
		elements := []string{"a", "b", "c", "b", "d", "a"}
		expected := []string{"a", "b", "c", "d"}

		result := removeDuplicates(elements)

		assert.Equal(t, expected, result)
	})

	t.Run("emptySlice", func(t *testing.T) {
		elements := []string{}
		expected := []string{}

		result := removeDuplicates(elements)

		assert.Equal(t, expected, result)
	})
}
