/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package model

import (
	"bytes"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type StringInterface interface {
	String() string
}

func createType(typeName string) StringInterface {
	switch typeName {
	case "AssetIdentifier":
		return &AssetIdentifier{}
	case "IdLink":
		return &IdLink{}
	case "EthernetPort":
		return &EthernetPort{}
	case "Ipv4Connectivity":
		return &Ipv4Connectivity{}
	case "asset":
		return &SoftwareAsset{}
	case "artifact":
		return &SoftwareArtifact{}
	default:
		return nil
	}
}

func printPolyNestedMapArray(objectArray []interface{}, buffer *bytes.Buffer) {
	for _, object := range objectArray {
		if nestedMap, ok := object.(map[string]interface{}); ok {
			printPolyNestedMap(nestedMap, buffer)
		}
	}
}

func printPolyNestedMap(nestedMap map[string]interface{}, buffer *bytes.Buffer) {
	for key, value := range nestedMap {
		if valueMap, ok := value.(map[string]interface{}); ok {
			printPolyMapInternal(valueMap, key, buffer)
		}
	}
}

func printPolyMapArray(objectArray []interface{}, typeKey string, buffer *bytes.Buffer) {
	for _, object := range objectArray {
		if objectMap, ok := object.(map[string]interface{}); ok {
			printPolyMap(objectMap, typeKey, buffer)
		}
	}
}

func printPolyMap(objectMap map[string]interface{}, typeKey string, buffer *bytes.Buffer) {
	if typePtr := objectMap[typeKey]; typePtr != nil {
		if typeName, ok := typePtr.(string); ok {
			printPolyMapInternal(objectMap, typeName, buffer)
		}
	}
}

func printPolyMapInternal(objectMap map[string]interface{}, typeName string, buffer *bytes.Buffer) {
	if object := createType(typeName); object != nil {
		if err := mapstructure.Decode(objectMap, &object); err == nil {
			buffer.WriteString(object.String())
		}
	}
}

func rds(fieldPointer *string) string {
	if fieldPointer == nil {
		return ""
	}
	return *fieldPointer
}

func rdb(fieldPointer *bool) string {
	if fieldPointer == nil {
		return ""
	}
	if *fieldPointer {
		return "true"
	}
	return "false"
}

func (a *Asset) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Asset %s\n", rds(a.Name)))
	buffer.WriteString(a.ProductInstanceIdentifier.String())
	printPolyMapArray(a.ConnectionPoints, "connection_point_type", &buffer)
	printPolyMapArray(a.AssetIdentifiers, "asset_identifier_type", &buffer)
	printPolyNestedMapArray(a.SoftwareComponents, &buffer)

	for _, op := range a.AssetOperations {
		buffer.WriteString(op.String())
	}

	return buffer.String()
}

func (o *ConnectionPoint) String() string {
	return fmt.Sprintf("  ConnectionPoint\n    type:%s\n", *o.ConnectionPointType)
}

func (o *EthernetPort) String() string {
	return fmt.Sprintf("  EthernetPort\n    mac: %s\n", rds(o.MacAddress))
}

func (o *Ipv4Connectivity) String() string {
	return fmt.Sprintf("  IPv4Connectivity\n    ip: %s\n    netmask: %s\n    router: %s\n",
		rds(o.Ipv4Address), rds(o.NetworkMask), rds(o.RouterIpv4Address))
}

func (o *AssetIdentifier) String() string {
	return fmt.Sprintf("  AssetIdentifier\n    type: %s\n", rds((*string)(o.AssetIdentifierType)))
}

func (o *IdLink) String() string {
	return fmt.Sprintf("  IdLink\n    uri: %s\n", rds(o.IdLink))
}

func (o *ProductSerialIdentifier) String() string {
	return fmt.Sprintf("  ProductSerialIdentifier\n    manufacturer: %s\n    product: %s\n    serial_number: %s\n    product_id: %s\n    hardware_version: %s\n",
		rds(o.ManufacturerProduct.Manufacturer.Name), rds(o.ManufacturerProduct.Name), rds(o.SerialNumber), rds(o.ManufacturerProduct.ProductId), rds(o.ManufacturerProduct.ProductVersion))
}

func (o *SoftwareAsset) String() string {
	return fmt.Sprintf("  SoftwareAsset\n    name=%s\n", rds(o.Name))
}

func (o *SoftwareArtifact) String() string {
	return fmt.Sprintf("  SoftwareArtifact\n    name: %s\n    version: %s\n", rds(o.SoftwareIdentifier.Name), rds(o.SoftwareIdentifier.Version))
}

func (o *AssetOperation) String() string {
	return fmt.Sprintf("  AssetOperation\n    name: %s\n    activation_flag: %s\n", rds(o.OperationName), rdb(o.ActivationFlag))
}
