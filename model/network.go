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

	if isNonEmptyValues(macAddress) {
		nicId = uuid.New().String()

		connectionPointType := EthernetPortConnectionPointTypeEthernetPort
		nic := EthernetPort{
			ConnectionPointType:     &connectionPointType,
			Id:                      nicId,
			MacAddress:              &macAddress,
			RelatedConnectionPoints: nil,
		}

		nameKey := "name"
		if isNonEmptyValues(name) {
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
	return ""
}

// AddIPv4 Add an IPv4 address to a network card
//
// The given network mask should consist of 4 octets (aaa.bbb.ccc.ddd)
//
// No validation of is currently done
func (d *DeviceInfo) AddIPv4(nicId string, ipv4Address string, networkMask string, routerAddress string) (id string) {

	if isNonEmptyValues(ipv4Address, networkMask, routerAddress) {
		id = uuid.New().String()

		connectionPointType := Ipv4ConnectivityConnectionPointTypeIpv4Connectivity
		customRelName := "Relies on"
		relationship := RelatedConnectionPoint{
			ConnectionPoint:    &nicId,
			CustomRelationship: &customRelName,
		}
		ipv4 := Ipv4Connectivity{
			ConnectionPointType:     &connectionPointType,
			Id:                      id,
			InstanceAnnotations:     nil,
			Ipv4Address:             &ipv4Address,
			RelatedConnectionPoints: []RelatedConnectionPoint{relationship},
		}
		if isNonEmptyValues(networkMask, routerAddress) {
			ipv4.NetworkMask = &networkMask
			ipv4.RouterIpv4Address = &routerAddress
		}
		d.ConnectionPoints = append(d.ConnectionPoints, ipv4)

		return id
	}
	return ""
}

// AddIPv6 Add an IPv6 address to a network card
//
// No validation of is currently done
func (d *DeviceInfo) AddIPv6(nicId string, ipv6Address string, networkPrefix string, routerAddress string) (id string) {

	if isNonEmptyValues(ipv6Address) {
		id = uuid.New().String()

		connectionPointType := Ipv6ConnectivityConnectionPointTypeIpv6Connectivity
		customRelName := "Relies on"
		relationship := RelatedConnectionPoint{
			ConnectionPoint:    &nicId,
			CustomRelationship: &customRelName,
		}
		ipv6 := Ipv6Connectivity{
			ConnectionPointType:     &connectionPointType,
			Id:                      id,
			InstanceAnnotations:     nil,
			Ipv6Address:             &ipv6Address,
			RelatedConnectionPoints: []RelatedConnectionPoint{relationship},
		}
		if isNonEmptyValues(networkPrefix, routerAddress) {
			ipv6.Ipv6NetworkPrefix = &networkPrefix
			ipv6.RouterIpv6Address = &routerAddress
		}
		d.ConnectionPoints = append(d.ConnectionPoints, ipv6)

		return id
	}
	return ""
}
