/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork(t *testing.T) {
	t.Run("AddNic", func(t *testing.T) {
		m := NewDevice("asset", "MyDevice")

		assert.NotEmpty(t, m.AddNic("nic2", "AA:AA:AA:AA:AA:AA"))
		m.AddNic("nic0", "AA:BB:CC:DD:EE:FF")

		nics := m.getNics()
		if len(nics) != 2 {
			fmt.Printf("Expected 2 address, got %d\n", len(nics))
			t.Fail()
		}
		found := 0
		for _, v := range nics {
			for _, ik := range v.InstanceAnnotations {
				if *ik.Value == "nic0" {
					found++
					assert.Equal(t, "name", *v.InstanceAnnotations[0].Key)
					assert.Equal(t, "nic0", *v.InstanceAnnotations[0].Value)
					assert.Equal(t, "AA:BB:CC:DD:EE:FF", *v.MacAddress)

					uncertainity := 1
					id := MacIdentifier{MacAddress: v.MacAddress, IdentifierUncertainty: &uncertainity}
					assert.Contains(t, m.MacIdentifiers, id)
					break
				}
			}
		}

		assert.Equal(t, 1, found)

	})

	t.Run("AddIpv4", func(t *testing.T) {
		m := NewDevice("asset", "device")
		assert.NotEmpty(t, m.AddIPv4("nic0",
			"10.0.0.1", "255.0.0.0",
			"10.0.0.254"))
		assert.NotEmpty(t, m.AddIPv4("nic99",
			"172.16.0.1", "255.255.255.0",
			""))

		addresses := m.getIPv4()
		if len(addresses) != 2 {
			fmt.Printf("Expected 1 address, got %d\n", len(addresses))
			t.Fail()
		}
		found := 0
		for _, v := range addresses {
			for _, ik := range v.RelatedConnectionPoints {
				if *ik.ConnectionPoint == "nic0" {
					found++
					assert.Equal(t, "10.0.0.1", *v.Ipv4Address)
					assert.Equal(t, "255.0.0.0", *v.NetworkMask)
					assert.Equal(t, "10.0.0.254", *v.RouterIpv4Address)
					break
				}
			}
		}
		assert.Equal(t, 1, found)
	})

	t.Run("AddIpv6", func(t *testing.T) {
		m := NewDevice("asset", "device")
		m.AddIPv6("nic0",
			"fd00::42", "/64",
			"fd00::1")
		m.AddIPv6("nic2",
			"fd06:1:2:3::1", "",
			"")

		addresses := m.getIPv6()
		if len(addresses) != 2 {
			fmt.Printf("Expected 1 address, got %d\n", len(addresses))
			t.Fail()
		}
		found := 0
		for _, v := range addresses {
			for _, ik := range v.RelatedConnectionPoints {
				if *ik.ConnectionPoint == "nic0" {
					found++
					assert.Equal(t, "fd00::42", *v.Ipv6Address)
					assert.Equal(t, "/64", *v.Ipv6NetworkPrefix)
					assert.Equal(t, "fd00::1", *v.RouterIpv6Address)
					break
				}
			}
		}
		assert.Equal(t, 1, found)
	})

	t.Run("When Name is not present instance annotations not be added", func(t *testing.T) {
		m := NewDevice("asset", "MyDevice")
		m.AddNic("", "AA:BB:CC:DD:EE:FF")
		nics := m.getNics()
		assert.Nil(t, nics[0].InstanceAnnotations)
	})

	t.Run("When Name is present instance annotations should be added", func(t *testing.T) {
		m := NewDevice("asset", "MyDevice")
		m.AddNic("Test-Nic", "AA:BB:CC:DD:EE:FF")
		nics := m.getNics()
		assert.NotNil(t, nics[0].InstanceAnnotations)
	})

}

// TODO: Use templating
// Extract ethernet ports from model
func (d *DeviceInfo) getNics() []EthernetPort {
	r := []EthernetPort{}
	for _, v := range d.ConnectionPoints {
		if reflect.TypeOf(v) == reflect.TypeOf(EthernetPort{}) {
			r = append(r, v.(EthernetPort))
		}
	}
	return r
}

// Extract IPv4 Addresses from model
func (d *DeviceInfo) getIPv4() []Ipv4Connectivity {
	r := []Ipv4Connectivity{}
	for _, v := range d.ConnectionPoints {
		if reflect.TypeOf(v) == reflect.TypeOf(Ipv4Connectivity{}) {
			r = append(r, v.(Ipv4Connectivity))
		}
	}
	return r
}

// Extract IPv4 addresses from model
func (d *DeviceInfo) getIPv6() []Ipv6Connectivity {
	r := []Ipv6Connectivity{}
	for _, v := range d.ConnectionPoints {
		if reflect.TypeOf(v) == reflect.TypeOf(Ipv6Connectivity{}) {
			r = append(r, v.(Ipv6Connectivity))
		}
	}
	return r
}
