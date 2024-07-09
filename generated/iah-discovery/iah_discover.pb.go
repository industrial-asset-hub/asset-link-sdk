//
// Copyright (c) Siemens AG 2023 ALL RIGHTS RESERVED.

// Device Discover Interface

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.24.4
// source: iah_discover.proto

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

type DiscoverRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// optional: Filters which are used to filter the discover result, e.g. a special device type
	// values of identical keys are logically "OR" combined
	// values of different keys are logically "AND" combined
	Filters []*ActiveFilter `protobuf:"bytes,1,rep,name=filters,proto3" json:"filters,omitempty"`
	// optional: Options which are used to perform the discover, e.g. timeout, hircharchie level, ...
	// values of identical keys are logically "OR" combined
	// values of different keys are logically "AND" combined
	Options []*ActiveOption `protobuf:"bytes,2,rep,name=options,proto3" json:"options,omitempty"`
	// optional: Specify the target where the discovery should be performed
	// If it's not specified, the whole system is scanned
	Target []*Destination `protobuf:"bytes,3,rep,name=target,proto3" json:"target,omitempty"`
}

func (x *DiscoverRequest) Reset() {
	*x = DiscoverRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iah_discover_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoverRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoverRequest) ProtoMessage() {}

func (x *DiscoverRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iah_discover_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoverRequest.ProtoReflect.Descriptor instead.
func (*DiscoverRequest) Descriptor() ([]byte, []int) {
	return file_iah_discover_proto_rawDescGZIP(), []int{0}
}

func (x *DiscoverRequest) GetFilters() []*ActiveFilter {
	if x != nil {
		return x.Filters
	}
	return nil
}

func (x *DiscoverRequest) GetOptions() []*ActiveOption {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *DiscoverRequest) GetTarget() []*Destination {
	if x != nil {
		return x.Target
	}
	return nil
}

type DiscoverResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Holds information on one or more discovered devices
	Devices []*DiscoveredDevice `protobuf:"bytes,1,rep,name=devices,proto3" json:"devices,omitempty"`
}

func (x *DiscoverResponse) Reset() {
	*x = DiscoverResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iah_discover_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoverResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoverResponse) ProtoMessage() {}

func (x *DiscoverResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iah_discover_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoverResponse.ProtoReflect.Descriptor instead.
func (*DiscoverResponse) Descriptor() ([]byte, []int) {
	return file_iah_discover_proto_rawDescGZIP(), []int{1}
}

func (x *DiscoverResponse) GetDevices() []*DiscoveredDevice {
	if x != nil {
		return x.Devices
	}
	return nil
}

type DiscoveredDevice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of default device identifiers of a device supported by the protocol
	Identifiers []*DeviceIdentifier `protobuf:"bytes,1,rep,name=identifiers,proto3" json:"identifiers,omitempty"`
	// The whole connection parameters, to establish the communciation with
	// the discovered device, including the schema and subdriver configuration
	ConnectionParameterSet *ConnectionParameterSet `protobuf:"bytes,2,opt,name=connection_parameter_set,json=connectionParameterSet,proto3" json:"connection_parameter_set,omitempty"`
	// Timestamp when device was last seen
	// 64 bit unsigned integer which represents the number
	// of 100 nano-second [0.1 microsec] intervals since January 1, 1601 (UTC).
	Timestamp uint64 `protobuf:"fixed64,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *DiscoveredDevice) Reset() {
	*x = DiscoveredDevice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iah_discover_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscoveredDevice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscoveredDevice) ProtoMessage() {}

func (x *DiscoveredDevice) ProtoReflect() protoreflect.Message {
	mi := &file_iah_discover_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscoveredDevice.ProtoReflect.Descriptor instead.
func (*DiscoveredDevice) Descriptor() ([]byte, []int) {
	return file_iah_discover_proto_rawDescGZIP(), []int{2}
}

func (x *DiscoveredDevice) GetIdentifiers() []*DeviceIdentifier {
	if x != nil {
		return x.Identifiers
	}
	return nil
}

func (x *DiscoveredDevice) GetConnectionParameterSet() *ConnectionParameterSet {
	if x != nil {
		return x.ConnectionParameterSet
	}
	return nil
}

func (x *DiscoveredDevice) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

var File_iah_discover_proto protoreflect.FileDescriptor

var file_iah_discover_proto_rawDesc = []byte{
	0x0a, 0x12, 0x69, 0x61, 0x68, 0x5f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x26, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x69, 0x6e,
	0x64, 0x75, 0x73, 0x74, 0x72, 0x69, 0x61, 0x6c, 0x61, 0x73, 0x73, 0x65, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x14, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xd7, 0x01, 0x0a, 0x0f, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x41, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x52, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x12, 0x41, 0x0a, 0x07, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x73, 0x69, 0x65,
	0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3e, 0x0a, 0x06,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x73,
	0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x22, 0x66, 0x0a, 0x10,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x52, 0x0a, 0x07, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x38, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x69, 0x6e, 0x64, 0x75,
	0x73, 0x74, 0x72, 0x69, 0x61, 0x6c, 0x61, 0x73, 0x73, 0x65, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x65, 0x64, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x07, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x22, 0xf0, 0x01, 0x0a, 0x10, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x65, 0x64, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x51, 0x0a, 0x0b, 0x69, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f,
	0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52,
	0x0b, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x12, 0x6b, 0x0a, 0x18,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x65, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31,
	0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x53, 0x65,
	0x74, 0x52, 0x16, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x53, 0x65, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x06, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x32, 0x8a, 0x03, 0x0a, 0x11, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x41, 0x70, 0x69, 0x12, 0x71, 0x0a,
	0x0e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x73, 0x12,
	0x2d, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e,
	0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x77, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2f, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x88, 0x01, 0x0a, 0x0f, 0x44, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x37, 0x2e,
	0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x69, 0x6e, 0x64, 0x75, 0x73, 0x74, 0x72, 0x69,
	0x61, 0x6c, 0x61, 0x73, 0x73, 0x65, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x38, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73,
	0x2e, 0x69, 0x6e, 0x64, 0x75, 0x73, 0x74, 0x72, 0x69, 0x61, 0x6c, 0x61, 0x73, 0x73, 0x65, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_iah_discover_proto_rawDescOnce sync.Once
	file_iah_discover_proto_rawDescData = file_iah_discover_proto_rawDesc
)

func file_iah_discover_proto_rawDescGZIP() []byte {
	file_iah_discover_proto_rawDescOnce.Do(func() {
		file_iah_discover_proto_rawDescData = protoimpl.X.CompressGZIP(file_iah_discover_proto_rawDescData)
	})
	return file_iah_discover_proto_rawDescData
}

var file_iah_discover_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_iah_discover_proto_goTypes = []interface{}{
	(*DiscoverRequest)(nil),        // 0: siemens.industrialassethub.discover.v1.DiscoverRequest
	(*DiscoverResponse)(nil),       // 1: siemens.industrialassethub.discover.v1.DiscoverResponse
	(*DiscoveredDevice)(nil),       // 2: siemens.industrialassethub.discover.v1.DiscoveredDevice
	(*ActiveFilter)(nil),           // 3: siemens.common.filters.v1.ActiveFilter
	(*ActiveOption)(nil),           // 4: siemens.common.filters.v1.ActiveOption
	(*Destination)(nil),            // 5: siemens.common.address.v1.Destination
	(*DeviceIdentifier)(nil),       // 6: siemens.common.identifiers.v1.DeviceIdentifier
	(*ConnectionParameterSet)(nil), // 7: siemens.common.address.v1.ConnectionParameterSet
	(*FilterTypesRequest)(nil),     // 8: siemens.common.filters.v1.FilterTypesRequest
	(*FilterOptionsRequest)(nil),   // 9: siemens.common.filters.v1.FilterOptionsRequest
	(*FilterTypesResponse)(nil),    // 10: siemens.common.filters.v1.FilterTypesResponse
	(*FilterOptionsResponse)(nil),  // 11: siemens.common.filters.v1.FilterOptionsResponse
}
var file_iah_discover_proto_depIdxs = []int32{
	3,  // 0: siemens.industrialassethub.discover.v1.DiscoverRequest.filters:type_name -> siemens.common.filters.v1.ActiveFilter
	4,  // 1: siemens.industrialassethub.discover.v1.DiscoverRequest.options:type_name -> siemens.common.filters.v1.ActiveOption
	5,  // 2: siemens.industrialassethub.discover.v1.DiscoverRequest.target:type_name -> siemens.common.address.v1.Destination
	2,  // 3: siemens.industrialassethub.discover.v1.DiscoverResponse.devices:type_name -> siemens.industrialassethub.discover.v1.DiscoveredDevice
	6,  // 4: siemens.industrialassethub.discover.v1.DiscoveredDevice.identifiers:type_name -> siemens.common.identifiers.v1.DeviceIdentifier
	7,  // 5: siemens.industrialassethub.discover.v1.DiscoveredDevice.connection_parameter_set:type_name -> siemens.common.address.v1.ConnectionParameterSet
	8,  // 6: siemens.industrialassethub.discover.v1.DeviceDiscoverApi.GetFilterTypes:input_type -> siemens.common.filters.v1.FilterTypesRequest
	9,  // 7: siemens.industrialassethub.discover.v1.DeviceDiscoverApi.GetFilterOptions:input_type -> siemens.common.filters.v1.FilterOptionsRequest
	0,  // 8: siemens.industrialassethub.discover.v1.DeviceDiscoverApi.DiscoverDevices:input_type -> siemens.industrialassethub.discover.v1.DiscoverRequest
	10, // 9: siemens.industrialassethub.discover.v1.DeviceDiscoverApi.GetFilterTypes:output_type -> siemens.common.filters.v1.FilterTypesResponse
	11, // 10: siemens.industrialassethub.discover.v1.DeviceDiscoverApi.GetFilterOptions:output_type -> siemens.common.filters.v1.FilterOptionsResponse
	1,  // 11: siemens.industrialassethub.discover.v1.DeviceDiscoverApi.DiscoverDevices:output_type -> siemens.industrialassethub.discover.v1.DiscoverResponse
	9,  // [9:12] is the sub-list for method output_type
	6,  // [6:9] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_iah_discover_proto_init() }
func file_iah_discover_proto_init() {
	if File_iah_discover_proto != nil {
		return
	}
	file_common_address_proto_init()
	file_common_filters_proto_init()
	file_common_identifiers_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_iah_discover_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoverRequest); i {
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
		file_iah_discover_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoverResponse); i {
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
		file_iah_discover_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DiscoveredDevice); i {
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
			RawDescriptor: file_iah_discover_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_iah_discover_proto_goTypes,
		DependencyIndexes: file_iah_discover_proto_depIdxs,
		MessageInfos:      file_iah_discover_proto_msgTypes,
	}.Build()
	File_iah_discover_proto = out.File
	file_iah_discover_proto_rawDesc = nil
	file_iah_discover_proto_goTypes = nil
	file_iah_discover_proto_depIdxs = nil
}
