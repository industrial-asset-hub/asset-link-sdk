/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork(t *testing.T) {
	t.Run("AddNic", func(t *testing.T) {
		m, err := NewDevice("Asset", "MyDevice")
		assert.NoError(t, err)

		nic2Id, err := m.AddNic("nic2", "AA:AA:AA:AA:AA:AA")
		assert.NoError(t, err)
		assert.NotEmpty(t, nic2Id)
		_, err = m.AddNic("nic0", "AA:BB:CC:DD:EE:FF")
		assert.NoError(t, err)

		nics := m.getNics()
		if len(nics) != 2 {
			fmt.Printf("Expected 2 address, got %d\n", len(nics))
			t.Fail()
		}
		found := 0
		for _, v := range nics {
			if v.Name != nil && *v.Name == "nic0" {
				found++
				assert.Equal(t, "AA:BB:CC:DD:EE:FF", v.MacAddress)

				uncertainity := 1
				id := MacIdentifier{
					AssetIdentifierType:   MacIdentifierAssetIdentifierTypeMacIdentifier,
					MacAddress:            v.MacAddress,
					IdentifierUncertainty: &uncertainity,
				}
				assert.Contains(t, m.AssetIdentifiers, id)
				break
			}
		}

		assert.Equal(t, 1, found)

	})

	t.Run("AddIpv4", func(t *testing.T) {
		m, err := NewDevice("Asset", "device")
		assert.NoError(t, err)

		id, err := m.AddIPv4("nic0", "10.0.0.1", "255.0.0.0", "10.0.0.254")
		assert.NotEmpty(t, id)
		assert.NoError(t, err)

		id, err = m.AddIPv4("nic1", "172.16.0.1", "", "")
		assert.NotEmpty(t, id)
		assert.NoError(t, err)
		cp := m.findIPv4(id)
		assert.NotNil(t, cp)
		assert.Equal(t, "172.16.0.1", cp.Ipv4Address)
		assert.Nil(t, cp.NetworkMask)
		assert.Nil(t, cp.RouterIpv4Address)

		id, err = m.AddIPv4("nic101", "", "", "")
		assert.Empty(t, id)
		assert.Error(t, err)
		var ee *EmptyError
		if errors.As(err, &ee) {
			assert.Equal(t, "Ipv4Connectivity", ee.Field)
			assert.Equal(t, "All fields for Ipv4Connectivity are empty", ee.Message)
		}

		id, err = m.AddIPv4("nic102", "999.999.999.999", "255.255.255.0", "10.0.0.254")
		assert.NotEmpty(t, id)
		assert.NoError(t, err)
		cp = m.findIPv4(id)
		assert.NotNil(t, cp)
		assert.Equal(t, "", cp.Ipv4Address)
		assert.Equal(t, "255.255.255.0", *cp.NetworkMask)
		assert.Equal(t, "10.0.0.254", *cp.RouterIpv4Address)

		found := 0
		for _, v := range m.getIPv4() {
			for _, ik := range v.RelatedConnectionPoints {
				if ik.ConnectionPointId == "nic0" {
					found++
					assert.Equal(t, "10.0.0.1", v.Ipv4Address)
					assert.Equal(t, "255.0.0.0", *v.NetworkMask)
					assert.Equal(t, "10.0.0.254", *v.RouterIpv4Address)
					break
				}
			}
		}
		assert.Equal(t, 1, found)
	})

	t.Run("AddIpv6", func(t *testing.T) {
		m, err := NewDevice("Asset", "device")
		assert.NoError(t, err)

		id1, err1 := m.AddIPv6("nic0", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:db8:/64", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		assert.NoError(t, err1)
		assert.NotEmpty(t, id1)

		id2, err2 := m.AddIPv6("nic2", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "", "")
		assert.NoError(t, err2)
		assert.NotEmpty(t, id2)
		cp := m.findIPv6(id2)
		assert.NotNil(t, cp)
		assert.Equal(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", cp.Ipv6Address)
		assert.Nil(t, cp.Ipv6NetworkPrefix)
		assert.Nil(t, cp.RouterIpv6Address)

		id3, err3 := m.AddIPv6("nic101", "", "", "")
		assert.Empty(t, id3)
		assert.Error(t, err3)
		var ee *EmptyError
		if errors.As(err3, &ee) {
			assert.Equal(t, "Ipv6Connectivity", ee.Field)
			assert.Equal(t, "All fields for Ipv6Connectivity are empty", ee.Message)
		}

		found := 0
		for _, v := range m.getIPv6() {
			for _, ik := range v.RelatedConnectionPoints {
				if ik.ConnectionPointId == "nic0" {
					found++
					assert.Equal(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", v.Ipv6Address)
					assert.Equal(t, "2001:db8:/64", *v.Ipv6NetworkPrefix)
					assert.Equal(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", *v.RouterIpv6Address)
					break
				}
			}
		}
		assert.Equal(t, 1, found)
	})

	t.Run("AddNic do not return error when the NIC name is empty", func(t *testing.T) {
		m, err := NewDevice("Asset", "MyDevice")
		assert.NoError(t, err)
		_, err = m.AddNic("", "AA:BB:CC:DD:EE:FF")
		assert.NoError(t, err)
		nics := m.getNics()
		assert.NotNil(t, nics[0].Name)
		assert.Equal(t, "", *nics[0].Name)
	})

	t.Run("AddNic stores the provided NIC name on EthernetPort.Name", func(t *testing.T) {
		m, err := NewDevice("Asset", "MyDevice")
		assert.NoError(t, err)
		_, err = m.AddNic("Test-Nic", "AA:BB:CC:DD:EE:FF")
		assert.NoError(t, err)
		nics := m.getNics()
		assert.NotNil(t, nics[0].Name)
		assert.Equal(t, "Test-Nic", *nics[0].Name)
	})
}

func TestAddNic_MacAddressValidation(t *testing.T) {
	m, err := NewDevice("Asset", "TestDevice")
	assert.NoError(t, err)

	t.Run("Empty MAC address should return EmptyError", func(t *testing.T) {
		_, err := m.AddNic("nic1", "")
		assert.Error(t, err)
		var ee *EmptyError
		if errors.As(err, &ee) {
			assert.Equal(t, "MacAddress", ee.Field)
			assert.Equal(t, "MAC address is empty", ee.Message)
			assert.Equal(t, "", ee.Value)
		}
	})

	t.Run("Invalid MAC address format should return ValidationError", func(t *testing.T) {
		_, err := m.AddNic("nic2", "invalid-mac")
		assert.Error(t, err)
		var ve *ValidationError
		if errors.As(err, &ve) {
			assert.Equal(t, "MacAddress", ve.Field)
			assert.Equal(t, "MAC address format is invalid. Please refer to the base schema for the supported pattern.", ve.Message)
			assert.Equal(t, "invalid-mac", ve.Value)
		}
	})

	t.Run("Valid MAC address should succeed", func(t *testing.T) {
		nicId, err := m.AddNic("nic3", "AA:BB:CC:DD:EE:FF")
		assert.NoError(t, err)
		assert.NotEmpty(t, nicId)
	})
}

func TestAddIPv6_Validation(t *testing.T) {
	m, err := NewDevice("Asset", "TestDevice")
	assert.NoError(t, err)

	nicId, err := m.AddNic("nic6", "AA:BB:CC:DD:EE:FF")
	assert.NoError(t, err)
	assert.NotEmpty(t, nicId)

	t.Run("Valid IPv6 address, prefix, and router", func(t *testing.T) {
		id, err := m.AddIPv6(nicId,
			"2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:db8:/64",
			"2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("Empty IPv6 address is ignored with warning", func(t *testing.T) {
		id, err := m.AddIPv6(nicId, "", "2001:db8:/64", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		assert.NotEmpty(t, id)
		assert.NoError(t, err)
		cp := m.findIPv6(id)
		assert.NotNil(t, cp)
		assert.Equal(t, "", cp.Ipv6Address)
		assert.Equal(t, "2001:db8:/64", *cp.Ipv6NetworkPrefix)
		assert.Equal(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", *cp.RouterIpv6Address)
	})

	t.Run("Invalid IPv6 address format is ignored with warning", func(t *testing.T) {
		id, err := m.AddIPv6(nicId, "invalid-ipv6", "2001:db8:/64", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		assert.NotEmpty(t, id)
		assert.NoError(t, err)
		cp := m.findIPv6(id)
		assert.NotNil(t, cp)
		assert.Equal(t, "", cp.Ipv6Address)
	})

	t.Run("Invalid IPv6 network prefix is ignored with warning", func(t *testing.T) {
		id, err := m.AddIPv6(nicId, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "invalid-prefix", "2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		assert.NotEmpty(t, id)
		assert.NoError(t, err)
		cp := m.findIPv6(id)
		assert.NotNil(t, cp)
		assert.Nil(t, cp.Ipv6NetworkPrefix)
	})

	t.Run("Invalid router IPv6 address should be ignored with warning", func(t *testing.T) {
		id, err := m.AddIPv6(nicId, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:db8:/64", "invalid-router")
		assert.NotEmpty(t, id)
		assert.NoError(t, err)
		cp := m.findIPv6(id)
		assert.NotNil(t, cp)
		assert.Nil(t, cp.RouterIpv6Address)
	})
}

func TestAddIPv4_ZeroAddressValues(t *testing.T) {
	m, err := NewDevice("Asset", "TestDevice")
	assert.NoError(t, err)

	t.Run("IPv4 address 0.0.0.0 should succeed", func(t *testing.T) {
		id, err := m.AddIPv4("nic0", "0.0.0.0", "255.255.255.0", "10.0.0.254")
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("Network mask 0.0.0.0 is ignored with warning", func(t *testing.T) {
		id, err := m.AddIPv4("nic0", "10.0.0.1", "0.0.0.0", "10.0.0.254")
		assert.NotEmpty(t, id)
		assert.NoError(t, err)
		cp := m.findIPv4(id)
		assert.NotNil(t, cp)
		assert.Equal(t, "10.0.0.1", cp.Ipv4Address)
		assert.Nil(t, cp.NetworkMask)
		assert.Equal(t, "10.0.0.254", *cp.RouterIpv4Address)
	})

	t.Run("Gateway router 0.0.0.0 should succeed", func(t *testing.T) {
		id, err := m.AddIPv4("nic0", "10.0.0.1", "255.255.255.0", "0.0.0.0")
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
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

func (d *DeviceInfo) findIPv4(id string) *Ipv4Connectivity {
	for _, v := range d.getIPv4() {
		if v.Id != nil && *v.Id == id {
			c := v
			return &c
		}
	}
	return nil
}

func (d *DeviceInfo) findIPv6(id string) *Ipv6Connectivity {
	for _, v := range d.getIPv6() {
		if v.Id != nil && *v.Id == id {
			c := v
			return &c
		}
	}
	return nil
}