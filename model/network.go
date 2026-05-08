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
func (d *DeviceInfo) AddNic(name, macAddress string) (string, error) {
	err := ValidateField(macAddress, "MacAddress", "MAC address is empty", MacAddressPattern,
		"MAC address format is invalid. Please refer to the base schema for the supported pattern.")
	if err != nil {
		return "", err
	}

	err = ValidateField(name, "Name", "Name is empty", "", "")
	if err != nil {
		return "", err
	}

	nicId := uuid.New().String()
	connectionPointType := EthernetPortConnectionPointTypeEthernetPort
	nic := EthernetPort{
		ConnectionPointType:     connectionPointType,
		Id:                      &nicId,
		Name:                    &name,
		MacAddress:              macAddress,
		RelatedConnectionPoints: nil,
	}

	d.ConnectionPoints = append(d.ConnectionPoints, nic)
	// automatically an MAC identifier, as it is required currently.
	d.addMacIdentifier(macAddress)

	return nicId, nil
}

// AddIPv4 Add an IPv4 address to a network card
//
// The given network mask should consist of 4 octets (aaa.bbb.ccc.ddd)
//
// No validation of is currently done
func (d *DeviceInfo) AddIPv4(nicId string, ipv4Address string, networkMask string, routerAddress string) (id string, err error) {
	err = ValidateField(ipv4Address, "IPv4Address", "IPv4 address is empty", IPv4AddressPattern,
		"IPv4 address format is invalid. Please refer to the base schema for the supported pattern.")
	if err != nil {
		return "", err
	}
	err = ValidateField(networkMask, "NetworkMask", "Network mask is empty", NetworkMaskPattern,
		"Network mask format is invalid. Please refer to the base schema for the supported pattern.")
	if err != nil {
		return "", err
	}
	err = ValidateField(routerAddress, "RouterIPv4Address", "Router IPv4 address is empty", RouterIPv4AddressPattern,
		"Router IPv4 address format is invalid. Please refer to the base schema for the supported pattern.")
	if err != nil {
		return "", err
	}

	id = uuid.New().String()
	connectionPointType := Ipv4ConnectivityConnectionPointTypeIpv4Connectivity
	customRelName := "Relies on"
	relationship := RelatedConnectionPoint{
		ConnectionPointId:  nicId,
		CustomRelationship: customRelName,
	}
	ipv4 := Ipv4Connectivity{
		ConnectionPointType:     connectionPointType,
		Id:                      &id,
		Ipv4Address:             ipv4Address,
		RelatedConnectionPoints: []RelatedConnectionPoint{relationship},
	}
	ipv4.NetworkMask = &networkMask
	ipv4.RouterIpv4Address = &routerAddress
	d.ConnectionPoints = append(d.ConnectionPoints, ipv4)

	return id, nil
}

// AddIPv6 Add an IPv6 address to a network card
//
// No validation of is currently done
func (d *DeviceInfo) AddIPv6(nicId string, ipv6Address string, networkPrefix string, routerAddress string) (id string, err error) {
	err = ValidateField(ipv6Address, "IPv6Address", "IPv6 address is empty", IPv6AddressPattern, "IPv6 address format is invalid. Please refer to the base schema for the supported pattern.")
	if err != nil {
		return "", err
	}
	err = ValidateField(networkPrefix, "IPv6NetworkPrefix", "IPv6 network prefix is empty", IPv6NetworkPrefixPattern, "IPv6 network prefix format is invalid. Please refer to the base schema for the supported pattern.")
	if err != nil {
		return "", err
	}
	err = ValidateField(routerAddress, "RouterIPv6Address", "Router IPv6 address is empty", RouterIPv6AddressPattern, "Router IPv6 address format is invalid. Please refer to the base schema for the supported pattern.")
	if err != nil {
		return "", err
	}

	id = uuid.New().String()
	connectionPointType := Ipv6ConnectivityConnectionPointTypeIpv6Connectivity
	customRelName := "Relies on"
	relationship := RelatedConnectionPoint{
		ConnectionPointId:  nicId,
		CustomRelationship: customRelName,
	}
	ipv6 := Ipv6Connectivity{
		ConnectionPointType:     connectionPointType,
		Id:                      &id,
		Ipv6Address:             ipv6Address,
		RelatedConnectionPoints: []RelatedConnectionPoint{relationship},
	}
	ipv6.Ipv6NetworkPrefix = &networkPrefix
	ipv6.RouterIpv6Address = &routerAddress
	d.ConnectionPoints = append(d.ConnectionPoints, ipv6)

	return id, nil
}
