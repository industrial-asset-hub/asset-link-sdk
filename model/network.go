/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"github.com/google/uuid"
)

// AddNic Add a network interface card
//
// The function returns an UUID, which can be used to add an IP address to the NIC
// No validation of is currently done
func (d *DeviceInfo) AddNic(name string, macAddress string) (nicId string) {
	nicId = uuid.New().String()

	t := EthernetPortConnectionPointTypeEthernetPort
	nic := EthernetPort{
		ConnectionPointType:     &t,
		Id:                      nicId,
		MacAddress:              &macAddress,
		RelatedConnectionPoints: nil,
	}

	nameKey := "name"
	if name != "" {
		nic.InstanceAnnotations = []InstanceAnnotation{{
			Key:   &nameKey,
			Value: &name,
		}}
	}

	d.ConnectionPoints = append(d.ConnectionPoints, nic)

	// automatically an MAC identifier, as it is required currently.
	d.addIdentifier(macAddress)

	return nicId
}

// AddIPv4 Add an IPv4 address to a network card
//
// The given network mask should consist of 4 octets (aaa.bbb.ccc.ddd)
//
// No validation of is currently done
func (d *DeviceInfo) AddIPv4(nicId string, address string, networkMask string, router string) (id string) {
	id = uuid.New().String()

	t := Ipv4ConnectivityConnectionPointTypeIpv4Connectivity
	customRelName := "Relies on"
	relationship := RelatedConnectionPoint{
		ConnectionPoint:    &nicId,
		CustomRelationship: &customRelName,
	}
	ipv4 := Ipv4Connectivity{
		ConnectionPointType:     &t,
		Id:                      id,
		InstanceAnnotations:     nil,
		Ipv4Address:             &address,
		NetworkMask:             &networkMask,
		RelatedConnectionPoints: []RelatedConnectionPoint{relationship},
		RouterIpv4Address:       &router,
	}
	d.ConnectionPoints = append(d.ConnectionPoints, ipv4)

	return id
}

// AddIPv6 Add an IPv6 address to a network card
//
// No validation of is currently done
func (d *DeviceInfo) AddIPv6(nicId string, address string, networkMask string, router string) (id string) {
	id = uuid.New().String()

	t := Ipv6ConnectivityConnectionPointTypeIpv6Connectivity
	customRelName := "Relies on"
	relationship := RelatedConnectionPoint{
		ConnectionPoint:    &nicId,
		CustomRelationship: &customRelName,
	}
	ipv6 := Ipv6Connectivity{
		ConnectionPointType:     &t,
		Id:                      id,
		InstanceAnnotations:     nil,
		Ipv6Address:             &address,
		Ipv6NetworkPrefix:       &networkMask,
		RelatedConnectionPoints: []RelatedConnectionPoint{relationship},
		RouterIpv6Address:       &router,
	}
	d.ConnectionPoints = append(d.ConnectionPoints, ipv6)

	return id
}
