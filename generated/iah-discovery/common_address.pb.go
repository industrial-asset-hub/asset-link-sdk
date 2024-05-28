// ------------------------------------------------------------------
// Common Definition of Device Address and Node Address
// ------------------------------------------------------------------

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: common_address.proto

package iah_discovery

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ConnectionCredential struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the URI of a Connector schema the credential belongs to
	SchemaUri string `protobuf:"bytes,1,opt,name=schema_uri,json=schemaUri,proto3" json:"schema_uri,omitempty"`
	// Connector specific credentials to establish the connection
	// Could be Token, Username + Password, Certificates, ...
	Credentials string `protobuf:"bytes,2,opt,name=credentials,proto3" json:"credentials,omitempty"`
}

func (x *ConnectionCredential) Reset() {
	*x = ConnectionCredential{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_address_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectionCredential) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectionCredential) ProtoMessage() {}

func (x *ConnectionCredential) ProtoReflect() protoreflect.Message {
	mi := &file_common_address_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectionCredential.ProtoReflect.Descriptor instead.
func (*ConnectionCredential) Descriptor() ([]byte, []int) {
	return file_common_address_proto_rawDescGZIP(), []int{0}
}

func (x *ConnectionCredential) GetSchemaUri() string {
	if x != nil {
		return x.SchemaUri
	}
	return ""
}

func (x *ConnectionCredential) GetCredentials() string {
	if x != nil {
		return x.Credentials
	}
	return ""
}

type ConnectionParameterSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the URI of a specific Connector schema
	SchemaUri string `protobuf:"bytes,1,opt,name=schema_uri,json=schemaUri,proto3" json:"schema_uri,omitempty"`
	// A json string containing the connection parameters according to the
	// path $defs/config_connection/parameters in the base schema
	// 'https://siemens.com/connectivity_suite/schemas/base/1.0.0/config.json'
	// specialized for the specific schema given above with 'schema_uri'
	// providing properties like ip address, ...
	ParameterJson string `protobuf:"bytes,2,opt,name=parameter_json,json=parameterJson,proto3" json:"parameter_json,omitempty"`
	// A json string containing the connection subdriver configuration according to the
	// path $defs/config_connection/subdriver in the base schema
	// 'https://siemens.com/connectivity_suite/schemas/base/1.0.0/config.json'
	SubdriverJson string                  `protobuf:"bytes,3,opt,name=subdriver_json,json=subdriverJson,proto3" json:"subdriver_json,omitempty"`
	Credentials   []*ConnectionCredential `protobuf:"bytes,4,rep,name=credentials,proto3" json:"credentials,omitempty"`
}

func (x *ConnectionParameterSet) Reset() {
	*x = ConnectionParameterSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_address_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectionParameterSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectionParameterSet) ProtoMessage() {}

func (x *ConnectionParameterSet) ProtoReflect() protoreflect.Message {
	mi := &file_common_address_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectionParameterSet.ProtoReflect.Descriptor instead.
func (*ConnectionParameterSet) Descriptor() ([]byte, []int) {
	return file_common_address_proto_rawDescGZIP(), []int{1}
}

func (x *ConnectionParameterSet) GetSchemaUri() string {
	if x != nil {
		return x.SchemaUri
	}
	return ""
}

func (x *ConnectionParameterSet) GetParameterJson() string {
	if x != nil {
		return x.ParameterJson
	}
	return ""
}

func (x *ConnectionParameterSet) GetSubdriverJson() string {
	if x != nil {
		return x.SubdriverJson
	}
	return ""
}

func (x *ConnectionParameterSet) GetCredentials() []*ConnectionCredential {
	if x != nil {
		return x.Credentials
	}
	return nil
}

type DatapointConfiguration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the name of the datapoint
	DatapointName string `protobuf:"bytes,1,opt,name=datapoint_name,json=datapointName,proto3" json:"datapoint_name,omitempty"`
	// The configurations for the datapoints which has to be added
	DatapointParameterSet *DatapointParameterSet `protobuf:"bytes,2,opt,name=datapoint_parameter_set,json=datapointParameterSet,proto3" json:"datapoint_parameter_set,omitempty"`
	// The access mode for this datapoint "r", "w" or "rw"
	AccessMode string `protobuf:"bytes,3,opt,name=access_mode,json=accessMode,proto3" json:"access_mode,omitempty"`
	// The owner of this datapoint
	Owner string `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
	// A json string containing the datapoint array lower bounds according to the
	// path $defs/connection_datapoint/array_lower_bounds in the base schema
	// 'https://siemens.com/connectivity_suite/schemas/base/1.0.0/config.json'
	// specialized for the specific schema given above with 'schema_uri'
	ArrayLowerBounds []int32 `protobuf:"varint,5,rep,packed,name=array_lower_bounds,json=arrayLowerBounds,proto3" json:"array_lower_bounds,omitempty"`
}

func (x *DatapointConfiguration) Reset() {
	*x = DatapointConfiguration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_address_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatapointConfiguration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatapointConfiguration) ProtoMessage() {}

func (x *DatapointConfiguration) ProtoReflect() protoreflect.Message {
	mi := &file_common_address_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatapointConfiguration.ProtoReflect.Descriptor instead.
func (*DatapointConfiguration) Descriptor() ([]byte, []int) {
	return file_common_address_proto_rawDescGZIP(), []int{2}
}

func (x *DatapointConfiguration) GetDatapointName() string {
	if x != nil {
		return x.DatapointName
	}
	return ""
}

func (x *DatapointConfiguration) GetDatapointParameterSet() *DatapointParameterSet {
	if x != nil {
		return x.DatapointParameterSet
	}
	return nil
}

func (x *DatapointConfiguration) GetAccessMode() string {
	if x != nil {
		return x.AccessMode
	}
	return ""
}

func (x *DatapointConfiguration) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *DatapointConfiguration) GetArrayLowerBounds() []int32 {
	if x != nil {
		return x.ArrayLowerBounds
	}
	return nil
}

// Benötigt bei generic read/write, bestandteil von configuration
type DatapointParameterSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A json string containing the datapoint address according to the
	// path $defs/connection_datapoint/address in the base schema
	// 'https://siemens.com/connectivity_suite/schemas/base/1.0.0/config.json'
	// specialized for the specific schema given by the schema uri of the connection
	AddressJson string `protobuf:"bytes,2,opt,name=address_json,json=addressJson,proto3" json:"address_json,omitempty"`
	// A json string containing the datapoint specific parameters according to the
	// path $defs/connection_datapoint/parameters in the base schema
	// 'https://siemens.com/connectivity_suite/schemas/base/1.0.0/config.json'
	// specialized for the specific schema given above with 'schema_uri'
	ParameterJson string `protobuf:"bytes,3,opt,name=parameter_json,json=parameterJson,proto3" json:"parameter_json,omitempty"`
	// This is the device specific datatype and not the Connectivity Suite datatype.
	// The types are defined by the driver.
	ConnectorSpecificDatatype string `protobuf:"bytes,4,opt,name=connector_specific_datatype,json=connectorSpecificDatatype,proto3" json:"connector_specific_datatype,omitempty"`
	// Size of the array dimension(s)
	// Examples:
	// [] ( empty ) - scalar value
	// [2]          - 1-dim array with size 2
	// [0]          - 1-dim array with dynamic size, i.e. size is part of payload
	// [2,3]        - 2-dim array with size 2 x 3
	// [2,0]        - 2-dim array with dynamic size 2 x n, i.e. size of 2nd dimension is part of payload
	ArrayDimensions []int32 `protobuf:"varint,5,rep,packed,name=array_dimensions,json=arrayDimensions,proto3" json:"array_dimensions,omitempty"`
}

func (x *DatapointParameterSet) Reset() {
	*x = DatapointParameterSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_address_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatapointParameterSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatapointParameterSet) ProtoMessage() {}

func (x *DatapointParameterSet) ProtoReflect() protoreflect.Message {
	mi := &file_common_address_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatapointParameterSet.ProtoReflect.Descriptor instead.
func (*DatapointParameterSet) Descriptor() ([]byte, []int) {
	return file_common_address_proto_rawDescGZIP(), []int{3}
}

func (x *DatapointParameterSet) GetAddressJson() string {
	if x != nil {
		return x.AddressJson
	}
	return ""
}

func (x *DatapointParameterSet) GetParameterJson() string {
	if x != nil {
		return x.ParameterJson
	}
	return ""
}

func (x *DatapointParameterSet) GetConnectorSpecificDatatype() string {
	if x != nil {
		return x.ConnectorSpecificDatatype
	}
	return ""
}

func (x *DatapointParameterSet) GetArrayDimensions() []int32 {
	if x != nil {
		return x.ArrayDimensions
	}
	return nil
}

// this is an array of browsenames which together form one browse path
// to identify one Node during browsing (e.g. starting node)
type BrowsePath struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Names []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
}

func (x *BrowsePath) Reset() {
	*x = BrowsePath{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_address_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BrowsePath) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BrowsePath) ProtoMessage() {}

func (x *BrowsePath) ProtoReflect() protoreflect.Message {
	mi := &file_common_address_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BrowsePath.ProtoReflect.Descriptor instead.
func (*BrowsePath) Descriptor() ([]byte, []int) {
	return file_common_address_proto_rawDescGZIP(), []int{4}
}

func (x *BrowsePath) GetNames() []string {
	if x != nil {
		return x.Names
	}
	return nil
}

// Address of the device / instance you want to communicate with
type Destination struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Target:
	//
	//	*Destination_ConnectionParameterSet
	//	*Destination_ConnectionName
	Target isDestination_Target `protobuf_oneof:"target"`
}

func (x *Destination) Reset() {
	*x = Destination{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_address_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Destination) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Destination) ProtoMessage() {}

func (x *Destination) ProtoReflect() protoreflect.Message {
	mi := &file_common_address_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Destination.ProtoReflect.Descriptor instead.
func (*Destination) Descriptor() ([]byte, []int) {
	return file_common_address_proto_rawDescGZIP(), []int{5}
}

func (m *Destination) GetTarget() isDestination_Target {
	if m != nil {
		return m.Target
	}
	return nil
}

func (x *Destination) GetConnectionParameterSet() *ConnectionParameterSet {
	if x, ok := x.GetTarget().(*Destination_ConnectionParameterSet); ok {
		return x.ConnectionParameterSet
	}
	return nil
}

func (x *Destination) GetConnectionName() string {
	if x, ok := x.GetTarget().(*Destination_ConnectionName); ok {
		return x.ConnectionName
	}
	return ""
}

type isDestination_Target interface {
	isDestination_Target()
}

type Destination_ConnectionParameterSet struct {
	// e.g. PROFINET name, IP-Address, PA-Address, ...
	ConnectionParameterSet *ConnectionParameterSet `protobuf:"bytes,1,opt,name=connection_parameter_set,json=connectionParameterSet,proto3,oneof"`
}

type Destination_ConnectionName struct {
	// alternative, if a connection exist and can be referenced by name
	ConnectionName string `protobuf:"bytes,2,opt,name=connection_name,json=connectionName,proto3,oneof"`
}

func (*Destination_ConnectionParameterSet) isDestination_Target() {}

func (*Destination_ConnectionName) isDestination_Target() {}

// This is a work-around for "repeated oneof" in protobuf. See https://github.com/protocolbuffers/protobuf/issues/2592
type NodeAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*NodeAddress_DatapointJson
	//	*NodeAddress_BrowsePath
	Type isNodeAddress_Type `protobuf_oneof:"type"`
}

func (x *NodeAddress) Reset() {
	*x = NodeAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_address_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeAddress) ProtoMessage() {}

func (x *NodeAddress) ProtoReflect() protoreflect.Message {
	mi := &file_common_address_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeAddress.ProtoReflect.Descriptor instead.
func (*NodeAddress) Descriptor() ([]byte, []int) {
	return file_common_address_proto_rawDescGZIP(), []int{6}
}

func (m *NodeAddress) GetType() isNodeAddress_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *NodeAddress) GetDatapointJson() string {
	if x, ok := x.GetType().(*NodeAddress_DatapointJson); ok {
		return x.DatapointJson
	}
	return ""
}

func (x *NodeAddress) GetBrowsePath() *BrowsePath {
	if x, ok := x.GetType().(*NodeAddress_BrowsePath); ok {
		return x.BrowsePath
	}
	return nil
}

type isNodeAddress_Type interface {
	isNodeAddress_Type()
}

type NodeAddress_DatapointJson struct {
	// A json string containing the datapoint address according to the
	// path $defs/connection_datapoint/address in the base schema
	// 'https://siemens.com/connectivity_suite/schemas/base/1.0.0/config.json'
	// specialized for the specific schema given by the schema uri of the connection
	DatapointJson string `protobuf:"bytes,1,opt,name=datapoint_json,json=datapointJson,proto3,oneof"`
}

type NodeAddress_BrowsePath struct {
	// should work for any browse server
	BrowsePath *BrowsePath `protobuf:"bytes,2,opt,name=browse_path,json=browsePath,proto3,oneof"`
}

func (*NodeAddress_DatapointJson) isNodeAddress_Type() {}

func (*NodeAddress_BrowsePath) isNodeAddress_Type() {}

var File_common_address_proto protoreflect.FileDescriptor

var file_common_address_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x76,
	0x31, 0x22, 0x57, 0x0a, 0x14, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43,
	0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x55, 0x72, 0x69, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63,
	0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x22, 0xd8, 0x01, 0x0a, 0x16, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x53, 0x65, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x5f,
	0x75, 0x72, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x55, 0x72, 0x69, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x4a, 0x73, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x0e, 0x73,
	0x75, 0x62, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x75, 0x62, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x4a, 0x73,
	0x6f, 0x6e, 0x12, 0x51, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x72,
	0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x61, 0x6c, 0x73, 0x22, 0x8e, 0x02, 0x0a, 0x16, 0x44, 0x61, 0x74, 0x61, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x25, 0x0a, 0x0e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x68, 0x0a, 0x17, 0x64, 0x61, 0x74, 0x61, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x5f, 0x73,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65,
	0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x52, 0x15, 0x64, 0x61, 0x74, 0x61,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x53, 0x65,
	0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6d, 0x6f, 0x64, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4d, 0x6f,
	0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x12, 0x61, 0x72, 0x72, 0x61,
	0x79, 0x5f, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x62, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x05, 0x52, 0x10, 0x61, 0x72, 0x72, 0x61, 0x79, 0x4c, 0x6f, 0x77, 0x65, 0x72,
	0x42, 0x6f, 0x75, 0x6e, 0x64, 0x73, 0x22, 0xcc, 0x01, 0x0a, 0x15, 0x44, 0x61, 0x74, 0x61, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x6a, 0x73, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x4a,
	0x73, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x4a, 0x73, 0x6f, 0x6e, 0x12, 0x3e, 0x0a, 0x1b, 0x63, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x19, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66,
	0x69, 0x63, 0x44, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x61, 0x72,
	0x72, 0x61, 0x79, 0x5f, 0x64, 0x69, 0x6d, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x05, 0x52, 0x0f, 0x61, 0x72, 0x72, 0x61, 0x79, 0x44, 0x69, 0x6d, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x22, 0x0a, 0x0a, 0x42, 0x72, 0x6f, 0x77, 0x73, 0x65, 0x50,
	0x61, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x22, 0xb1, 0x01, 0x0a, 0x0b, 0x44, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x6d, 0x0a, 0x18, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x5f, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x73, 0x69,
	0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x48, 0x00,
	0x52, 0x16, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x65, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x12, 0x29, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e,
	0x61, 0x6d, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x22, 0x88, 0x01,
	0x0a, 0x0b, 0x4e, 0x6f, 0x64, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x27, 0x0a,
	0x0e, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0d, 0x64, 0x61, 0x74, 0x61, 0x70, 0x6f, 0x69,
	0x6e, 0x74, 0x4a, 0x73, 0x6f, 0x6e, 0x12, 0x48, 0x0a, 0x0b, 0x62, 0x72, 0x6f, 0x77, 0x73, 0x65,
	0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x73, 0x69,
	0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x72, 0x6f, 0x77, 0x73, 0x65, 0x50, 0x61,
	0x74, 0x68, 0x48, 0x00, 0x52, 0x0a, 0x62, 0x72, 0x6f, 0x77, 0x73, 0x65, 0x50, 0x61, 0x74, 0x68,
	0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_address_proto_rawDescOnce sync.Once
	file_common_address_proto_rawDescData = file_common_address_proto_rawDesc
)

func file_common_address_proto_rawDescGZIP() []byte {
	file_common_address_proto_rawDescOnce.Do(func() {
		file_common_address_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_address_proto_rawDescData)
	})
	return file_common_address_proto_rawDescData
}

var file_common_address_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_common_address_proto_goTypes = []interface{}{
	(*ConnectionCredential)(nil),   // 0: siemens.common.address.v1.ConnectionCredential
	(*ConnectionParameterSet)(nil), // 1: siemens.common.address.v1.ConnectionParameterSet
	(*DatapointConfiguration)(nil), // 2: siemens.common.address.v1.DatapointConfiguration
	(*DatapointParameterSet)(nil),  // 3: siemens.common.address.v1.DatapointParameterSet
	(*BrowsePath)(nil),             // 4: siemens.common.address.v1.BrowsePath
	(*Destination)(nil),            // 5: siemens.common.address.v1.Destination
	(*NodeAddress)(nil),            // 6: siemens.common.address.v1.NodeAddress
}
var file_common_address_proto_depIdxs = []int32{
	0, // 0: siemens.common.address.v1.ConnectionParameterSet.credentials:type_name -> siemens.common.address.v1.ConnectionCredential
	3, // 1: siemens.common.address.v1.DatapointConfiguration.datapoint_parameter_set:type_name -> siemens.common.address.v1.DatapointParameterSet
	1, // 2: siemens.common.address.v1.Destination.connection_parameter_set:type_name -> siemens.common.address.v1.ConnectionParameterSet
	4, // 3: siemens.common.address.v1.NodeAddress.browse_path:type_name -> siemens.common.address.v1.BrowsePath
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_common_address_proto_init() }
func file_common_address_proto_init() {
	if File_common_address_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_address_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectionCredential); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_address_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectionParameterSet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_address_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatapointConfiguration); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_address_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatapointParameterSet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_address_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BrowsePath); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_address_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Destination); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_common_address_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeAddress); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_common_address_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*Destination_ConnectionParameterSet)(nil),
		(*Destination_ConnectionName)(nil),
	}
	file_common_address_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*NodeAddress_DatapointJson)(nil),
		(*NodeAddress_BrowsePath)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_address_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_address_proto_goTypes,
		DependencyIndexes: file_common_address_proto_depIdxs,
		MessageInfos:      file_common_address_proto_msgTypes,
	}.Build()
	File_common_address_proto = out.File
	file_common_address_proto_rawDesc = nil
	file_common_address_proto_goTypes = nil
	file_common_address_proto_depIdxs = nil
}