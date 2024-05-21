//
// Copyright (c) Siemens AG 2022 ALL RIGHTS RESERVED.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: device-discovery.proto

package device_discovery

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ComparisonOperator int32

const (
	ComparisonOperator_EQUAL                    ComparisonOperator = 0 // ==
	ComparisonOperator_NOT_EQUAL                ComparisonOperator = 1 // !=
	ComparisonOperator_GREATER_THAN             ComparisonOperator = 2 // >
	ComparisonOperator_GREATER_THAN_OR_EQUAL_TO ComparisonOperator = 3 // >=
	ComparisonOperator_LESS_THAN                ComparisonOperator = 4 // <
	ComparisonOperator_LESS_THAN_OR_EQUAL_TO    ComparisonOperator = 5 // <=
)

// Enum value maps for ComparisonOperator.
var (
	ComparisonOperator_name = map[int32]string{
		0: "EQUAL",
		1: "NOT_EQUAL",
		2: "GREATER_THAN",
		3: "GREATER_THAN_OR_EQUAL_TO",
		4: "LESS_THAN",
		5: "LESS_THAN_OR_EQUAL_TO",
	}
	ComparisonOperator_value = map[string]int32{
		"EQUAL":                    0,
		"NOT_EQUAL":                1,
		"GREATER_THAN":             2,
		"GREATER_THAN_OR_EQUAL_TO": 3,
		"LESS_THAN":                4,
		"LESS_THAN_OR_EQUAL_TO":    5,
	}
)

func (x ComparisonOperator) Enum() *ComparisonOperator {
	p := new(ComparisonOperator)
	*p = x
	return p
}

func (x ComparisonOperator) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ComparisonOperator) Descriptor() protoreflect.EnumDescriptor {
	return file_device_discovery_proto_enumTypes[0].Descriptor()
}

func (ComparisonOperator) Type() protoreflect.EnumType {
	return &file_device_discovery_proto_enumTypes[0]
}

func (x ComparisonOperator) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ComparisonOperator.Descriptor instead.
func (ComparisonOperator) EnumDescriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{0}
}

type DiscoveryFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string             `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value    string             `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Operator ComparisonOperator `protobuf:"varint,3,opt,name=operator,proto3,enum=siemens.commondevicemanagement.devicediscovery.v1.ComparisonOperator" json:"operator,omitempty"`
}

func (x *DiscoveryFilter) Reset() {
	*x = DiscoveryFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_discovery_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoveryFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryFilter) ProtoMessage() {}

func (x *DiscoveryFilter) ProtoReflect() protoreflect.Message {
	mi := &file_device_discovery_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryFilter.ProtoReflect.Descriptor instead.
func (*DiscoveryFilter) Descriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{0}
}

func (x *DiscoveryFilter) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *DiscoveryFilter) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *DiscoveryFilter) GetOperator() ComparisonOperator {
	if x != nil {
		return x.Operator
	}
	return ComparisonOperator_EQUAL
}

type DiscoveryOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string             `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value    string             `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Operator ComparisonOperator `protobuf:"varint,3,opt,name=operator,proto3,enum=siemens.commondevicemanagement.devicediscovery.v1.ComparisonOperator" json:"operator,omitempty"`
}

func (x *DiscoveryOption) Reset() {
	*x = DiscoveryOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_discovery_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoveryOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryOption) ProtoMessage() {}

func (x *DiscoveryOption) ProtoReflect() protoreflect.Message {
	mi := &file_device_discovery_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryOption.ProtoReflect.Descriptor instead.
func (*DiscoveryOption) Descriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{1}
}

func (x *DiscoveryOption) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *DiscoveryOption) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *DiscoveryOption) GetOperator() ComparisonOperator {
	if x != nil {
		return x.Operator
	}
	return ComparisonOperator_EQUAL
}

type DiscoveryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// values of identical keys are logically "OR" combined
	// values of different keys are logically "AND" combined
	Filters []*DiscoveryFilter `protobuf:"bytes,1,rep,name=filters,proto3" json:"filters,omitempty"`
	// values of identical keys are logically "OR" combined
	// values of different keys are logically "AND" combined
	Options []*DiscoveryOption `protobuf:"bytes,2,rep,name=options,proto3" json:"options,omitempty"`
}

func (x *DiscoveryRequest) Reset() {
	*x = DiscoveryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_discovery_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoveryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryRequest) ProtoMessage() {}

func (x *DiscoveryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_device_discovery_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryRequest.ProtoReflect.Descriptor instead.
func (*DiscoveryRequest) Descriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{2}
}

func (x *DiscoveryRequest) GetFilters() []*DiscoveryFilter {
	if x != nil {
		return x.Filters
	}
	return nil
}

func (x *DiscoveryRequest) GetOptions() []*DiscoveryOption {
	if x != nil {
		return x.Options
	}
	return nil
}

type DiscoveryReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiscoveryId uint32 `protobuf:"varint,1,opt,name=discovery_id,json=discoveryId,proto3" json:"discovery_id,omitempty"`
}

func (x *DiscoveryReply) Reset() {
	*x = DiscoveryReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_discovery_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoveryReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryReply) ProtoMessage() {}

func (x *DiscoveryReply) ProtoReflect() protoreflect.Message {
	mi := &file_device_discovery_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryReply.ProtoReflect.Descriptor instead.
func (*DiscoveryReply) Descriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{3}
}

func (x *DiscoveryReply) GetDiscoveryId() uint32 {
	if x != nil {
		return x.DiscoveryId
	}
	return 0
}

type DiscoveryResultsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiscoveryId uint32 `protobuf:"varint,1,opt,name=discovery_id,json=discoveryId,proto3" json:"discovery_id,omitempty"`
}

func (x *DiscoveryResultsRequest) Reset() {
	*x = DiscoveryResultsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_discovery_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoveryResultsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryResultsRequest) ProtoMessage() {}

func (x *DiscoveryResultsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_device_discovery_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryResultsRequest.ProtoReflect.Descriptor instead.
func (*DiscoveryResultsRequest) Descriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{4}
}

func (x *DiscoveryResultsRequest) GetDiscoveryId() uint32 {
	if x != nil {
		return x.DiscoveryId
	}
	return 0
}

type DiscoveryDevice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the semantics expected by the CDM discovery client is defined in the device-parameter-schema.json
	Parameters *structpb.Struct `protobuf:"bytes,1,opt,name=parameters,proto3" json:"parameters,omitempty"`
}

func (x *DiscoveryDevice) Reset() {
	*x = DiscoveryDevice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_discovery_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoveryDevice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryDevice) ProtoMessage() {}

func (x *DiscoveryDevice) ProtoReflect() protoreflect.Message {
	mi := &file_device_discovery_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryDevice.ProtoReflect.Descriptor instead.
func (*DiscoveryDevice) Descriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{5}
}

func (x *DiscoveryDevice) GetParameters() *structpb.Struct {
	if x != nil {
		return x.Parameters
	}
	return nil
}

type DiscoveryResultsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Devices []*DiscoveryDevice `protobuf:"bytes,1,rep,name=devices,proto3" json:"devices,omitempty"`
}

func (x *DiscoveryResultsReply) Reset() {
	*x = DiscoveryResultsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_discovery_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoveryResultsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveryResultsReply) ProtoMessage() {}

func (x *DiscoveryResultsReply) ProtoReflect() protoreflect.Message {
	mi := &file_device_discovery_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveryResultsReply.ProtoReflect.Descriptor instead.
func (*DiscoveryResultsReply) Descriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{6}
}

func (x *DiscoveryResultsReply) GetDevices() []*DiscoveryDevice {
	if x != nil {
		return x.Devices
	}
	return nil
}

type StopDiscoveryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiscoveryId uint32 `protobuf:"varint,1,opt,name=discovery_id,json=discoveryId,proto3" json:"discovery_id,omitempty"`
}

func (x *StopDiscoveryRequest) Reset() {
	*x = StopDiscoveryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_discovery_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopDiscoveryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopDiscoveryRequest) ProtoMessage() {}

func (x *StopDiscoveryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_device_discovery_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopDiscoveryRequest.ProtoReflect.Descriptor instead.
func (*StopDiscoveryRequest) Descriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{7}
}

func (x *StopDiscoveryRequest) GetDiscoveryId() uint32 {
	if x != nil {
		return x.DiscoveryId
	}
	return 0
}

type StopDiscoveryReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StopDiscoveryReply) Reset() {
	*x = StopDiscoveryReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_discovery_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopDiscoveryReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopDiscoveryReply) ProtoMessage() {}

func (x *StopDiscoveryReply) ProtoReflect() protoreflect.Message {
	mi := &file_device_discovery_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopDiscoveryReply.ProtoReflect.Descriptor instead.
func (*StopDiscoveryReply) Descriptor() ([]byte, []int) {
	return file_device_discovery_proto_rawDescGZIP(), []int{8}
}

var File_device_discovery_proto protoreflect.FileDescriptor

var file_device_discovery_proto_rawDesc = []byte{
	0x0a, 0x16, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2d, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x31, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x01, 0x0a, 0x0f, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x61, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x45, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70,
	0x61, 0x72, 0x69, 0x73, 0x6f, 0x6e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x08,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x9c, 0x01, 0x0a, 0x0f, 0x44, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x61, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x45, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x61,
	0x72, 0x69, 0x73, 0x6f, 0x6e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x08, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x22, 0xce, 0x01, 0x0a, 0x10, 0x44, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x5c, 0x0a, 0x07,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x42, 0x2e,
	0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x46, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x12, 0x5c, 0x0a, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x42, 0x2e, 0x73, 0x69,
	0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x33, 0x0a, 0x0e, 0x44, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0b, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x49, 0x64, 0x22, 0x3c, 0x0a,
	0x17, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x49, 0x64, 0x22, 0x4a, 0x0a, 0x0f, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37,
	0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x0a, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x22, 0x75, 0x0a, 0x15, 0x44, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x5c, 0x0a, 0x07, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x42, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x07, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x22, 0x39,
	0x0a, 0x14, 0x53, 0x74, 0x6f, 0x70, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x49, 0x64, 0x22, 0x14, 0x0a, 0x12, 0x53, 0x74, 0x6f,
	0x70, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2a,
	0x88, 0x01, 0x0a, 0x12, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x69, 0x73, 0x6f, 0x6e, 0x4f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x10,
	0x00, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x4f, 0x54, 0x5f, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x10, 0x01,
	0x12, 0x10, 0x0a, 0x0c, 0x47, 0x52, 0x45, 0x41, 0x54, 0x45, 0x52, 0x5f, 0x54, 0x48, 0x41, 0x4e,
	0x10, 0x02, 0x12, 0x1c, 0x0a, 0x18, 0x47, 0x52, 0x45, 0x41, 0x54, 0x45, 0x52, 0x5f, 0x54, 0x48,
	0x41, 0x4e, 0x5f, 0x4f, 0x52, 0x5f, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x5f, 0x54, 0x4f, 0x10, 0x03,
	0x12, 0x0d, 0x0a, 0x09, 0x4c, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x48, 0x41, 0x4e, 0x10, 0x04, 0x12,
	0x19, 0x0a, 0x15, 0x4c, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x48, 0x41, 0x4e, 0x5f, 0x4f, 0x52, 0x5f,
	0x45, 0x51, 0x55, 0x41, 0x4c, 0x5f, 0x54, 0x4f, 0x10, 0x05, 0x32, 0x99, 0x04, 0x0a, 0x12, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x41, 0x70,
	0x69, 0x12, 0xa0, 0x01, 0x0a, 0x14, 0x53, 0x74, 0x61, 0x72, 0x74, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x12, 0x43, 0x2e, 0x73, 0x69, 0x65,
	0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x41, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x00, 0x12, 0xb5, 0x01, 0x0a, 0x19, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x73, 0x12, 0x4a, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x48,
	0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e,
	0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x30, 0x01, 0x12, 0xa7, 0x01, 0x0a,
	0x13, 0x53, 0x74, 0x6f, 0x70, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x12, 0x47, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x44, 0x69, 0x73,
	0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x45, 0x2e,
	0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_device_discovery_proto_rawDescOnce sync.Once
	file_device_discovery_proto_rawDescData = file_device_discovery_proto_rawDesc
)

func file_device_discovery_proto_rawDescGZIP() []byte {
	file_device_discovery_proto_rawDescOnce.Do(func() {
		file_device_discovery_proto_rawDescData = protoimpl.X.CompressGZIP(file_device_discovery_proto_rawDescData)
	})
	return file_device_discovery_proto_rawDescData
}

var file_device_discovery_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_device_discovery_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_device_discovery_proto_goTypes = []interface{}{
	(ComparisonOperator)(0),         // 0: siemens.commondevicemanagement.devicediscovery.v1.ComparisonOperator
	(*DiscoveryFilter)(nil),         // 1: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryFilter
	(*DiscoveryOption)(nil),         // 2: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryOption
	(*DiscoveryRequest)(nil),        // 3: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryRequest
	(*DiscoveryReply)(nil),          // 4: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryReply
	(*DiscoveryResultsRequest)(nil), // 5: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryResultsRequest
	(*DiscoveryDevice)(nil),         // 6: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryDevice
	(*DiscoveryResultsReply)(nil),   // 7: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryResultsReply
	(*StopDiscoveryRequest)(nil),    // 8: siemens.commondevicemanagement.devicediscovery.v1.StopDiscoveryRequest
	(*StopDiscoveryReply)(nil),      // 9: siemens.commondevicemanagement.devicediscovery.v1.StopDiscoveryReply
	(*structpb.Struct)(nil),         // 10: google.protobuf.Struct
}
var file_device_discovery_proto_depIdxs = []int32{
	0,  // 0: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryFilter.operator:type_name -> siemens.commondevicemanagement.devicediscovery.v1.ComparisonOperator
	0,  // 1: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryOption.operator:type_name -> siemens.commondevicemanagement.devicediscovery.v1.ComparisonOperator
	1,  // 2: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryRequest.filters:type_name -> siemens.commondevicemanagement.devicediscovery.v1.DiscoveryFilter
	2,  // 3: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryRequest.options:type_name -> siemens.commondevicemanagement.devicediscovery.v1.DiscoveryOption
	10, // 4: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryDevice.parameters:type_name -> google.protobuf.Struct
	6,  // 5: siemens.commondevicemanagement.devicediscovery.v1.DiscoveryResultsReply.devices:type_name -> siemens.commondevicemanagement.devicediscovery.v1.DiscoveryDevice
	3,  // 6: siemens.commondevicemanagement.devicediscovery.v1.DeviceDiscoveryApi.StartDeviceDiscovery:input_type -> siemens.commondevicemanagement.devicediscovery.v1.DiscoveryRequest
	5,  // 7: siemens.commondevicemanagement.devicediscovery.v1.DeviceDiscoveryApi.SubscribeDiscoveryResults:input_type -> siemens.commondevicemanagement.devicediscovery.v1.DiscoveryResultsRequest
	8,  // 8: siemens.commondevicemanagement.devicediscovery.v1.DeviceDiscoveryApi.StopDeviceDiscovery:input_type -> siemens.commondevicemanagement.devicediscovery.v1.StopDiscoveryRequest
	4,  // 9: siemens.commondevicemanagement.devicediscovery.v1.DeviceDiscoveryApi.StartDeviceDiscovery:output_type -> siemens.commondevicemanagement.devicediscovery.v1.DiscoveryReply
	7,  // 10: siemens.commondevicemanagement.devicediscovery.v1.DeviceDiscoveryApi.SubscribeDiscoveryResults:output_type -> siemens.commondevicemanagement.devicediscovery.v1.DiscoveryResultsReply
	9,  // 11: siemens.commondevicemanagement.devicediscovery.v1.DeviceDiscoveryApi.StopDeviceDiscovery:output_type -> siemens.commondevicemanagement.devicediscovery.v1.StopDiscoveryReply
	9,  // [9:12] is the sub-list for method output_type
	6,  // [6:9] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_device_discovery_proto_init() }
func file_device_discovery_proto_init() {
	if File_device_discovery_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_device_discovery_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoveryFilter); i {
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
		file_device_discovery_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoveryOption); i {
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
		file_device_discovery_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoveryRequest); i {
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
		file_device_discovery_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoveryReply); i {
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
		file_device_discovery_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoveryResultsRequest); i {
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
		file_device_discovery_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoveryDevice); i {
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
		file_device_discovery_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoveryResultsReply); i {
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
		file_device_discovery_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopDiscoveryRequest); i {
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
		file_device_discovery_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopDiscoveryReply); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_device_discovery_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_device_discovery_proto_goTypes,
		DependencyIndexes: file_device_discovery_proto_depIdxs,
		EnumInfos:         file_device_discovery_proto_enumTypes,
		MessageInfos:      file_device_discovery_proto_msgTypes,
	}.Build()
	File_device_discovery_proto = out.File
	file_device_discovery_proto_rawDesc = nil
	file_device_discovery_proto_goTypes = nil
	file_device_discovery_proto_depIdxs = nil
}