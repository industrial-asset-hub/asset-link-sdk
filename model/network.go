/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package model

import (
	"github.com/google/uuid"
)

// Add an Network card
//
// No validation of is currently done
func (d *DeviceInfo) AddNic(name string, macAddress string) (id string) {
	id = uuid.New().String()

	t := "EthernetPort"
	nameKey := "name"
	nic := EthernetPort{
		ConnectionPointType: &t,
		Id:                  id,
		InstanceAnnotations: []InstanceAnnotation{InstanceAnnotation{
			Key:   &nameKey,
			Value: &name,
		}},
		MacAddress:              &macAddress,
		RelatedConnectionPoints: nil,
	}
	d.ConnectionPoints = append(d.ConnectionPoints, nic)
	return id
}

// Add an IPv4 address to a network card
//
// networkMask should be a consists of 4 octets (aaa.bbb.ccc.ddd)
//
// No validation of is currently done
func (d *DeviceInfo) AddIPv4(nicId string, address string, networkMask string, router string) (id string) {
	id = uuid.New().String()

	t := "Ipv4Connectivity"
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

// Add an IPv6 address to a network card
//
// networkMask is currently missing
//
// No validation of is currently done
func (d *DeviceInfo) AddIPv6(nicId string, address string, networkMask string, router string) (id string) {
	id = uuid.New().String()

	t := "Ipv6Connectivity"
	customRelName := "Relies on"
	relationship := RelatedConnectionPoint{
		ConnectionPoint:    &nicId,
		CustomRelationship: &customRelName,
	}
	ipv6 := Ipv6Connectivity{
		ConnectionPointType:     &t,
		Id:                      id,
		InstanceAnnotations:     nil,
		Ipv6Address:             []string{address},
		RelatedConnectionPoints: []RelatedConnectionPoint{relationship},
		RouterIpv6Address:       []string{router},
	}
	d.ConnectionPoints = append(d.ConnectionPoints, ipv6)

	return id
}
