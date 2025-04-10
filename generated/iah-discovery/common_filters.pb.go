// Filter type and option definitions

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v4.25.2
// source: common_filters.proto

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

type SupportedFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string      `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Datatype VariantType `protobuf:"varint,2,opt,name=datatype,proto3,enum=siemens.common.types.v1.VariantType" json:"datatype,omitempty"`
}

func (x *SupportedFilter) Reset() {
	*x = SupportedFilter{}
	mi := &file_common_filters_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SupportedFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SupportedFilter) ProtoMessage() {}

func (x *SupportedFilter) ProtoReflect() protoreflect.Message {
	mi := &file_common_filters_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SupportedFilter.ProtoReflect.Descriptor instead.
func (*SupportedFilter) Descriptor() ([]byte, []int) {
	return file_common_filters_proto_rawDescGZIP(), []int{0}
}

func (x *SupportedFilter) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SupportedFilter) GetDatatype() VariantType {
	if x != nil {
		return x.Datatype
	}
	return VariantType_VT_UNSPECIFIED
}

type SupportedOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key      string      `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Datatype VariantType `protobuf:"varint,2,opt,name=datatype,proto3,enum=siemens.common.types.v1.VariantType" json:"datatype,omitempty"`
}

func (x *SupportedOption) Reset() {
	*x = SupportedOption{}
	mi := &file_common_filters_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SupportedOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SupportedOption) ProtoMessage() {}

func (x *SupportedOption) ProtoReflect() protoreflect.Message {
	mi := &file_common_filters_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SupportedOption.ProtoReflect.Descriptor instead.
func (*SupportedOption) Descriptor() ([]byte, []int) {
	return file_common_filters_proto_rawDescGZIP(), []int{1}
}

func (x *SupportedOption) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SupportedOption) GetDatatype() VariantType {
	if x != nil {
		return x.Datatype
	}
	return VariantType_VT_UNSPECIFIED
}

type ActiveFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"` // e.g. Timeout
	// array of raw-data
	Operator ComparisonOperator `protobuf:"varint,2,opt,name=operator,proto3,enum=siemens.common.operators.v1.ComparisonOperator" json:"operator,omitempty"` //optional
	Value    *Variant           `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ActiveFilter) Reset() {
	*x = ActiveFilter{}
	mi := &file_common_filters_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActiveFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActiveFilter) ProtoMessage() {}

func (x *ActiveFilter) ProtoReflect() protoreflect.Message {
	mi := &file_common_filters_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActiveFilter.ProtoReflect.Descriptor instead.
func (*ActiveFilter) Descriptor() ([]byte, []int) {
	return file_common_filters_proto_rawDescGZIP(), []int{2}
}

func (x *ActiveFilter) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ActiveFilter) GetOperator() ComparisonOperator {
	if x != nil {
		return x.Operator
	}
	return ComparisonOperator_EQUAL
}

func (x *ActiveFilter) GetValue() *Variant {
	if x != nil {
		return x.Value
	}
	return nil
}

type ActiveOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"` // e.g. Timeout
	// array of raw-data
	Operator ComparisonOperator `protobuf:"varint,2,opt,name=operator,proto3,enum=siemens.common.operators.v1.ComparisonOperator" json:"operator,omitempty"` //optional (when missing it means EQUAL)
	Value    *Variant           `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ActiveOption) Reset() {
	*x = ActiveOption{}
	mi := &file_common_filters_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ActiveOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActiveOption) ProtoMessage() {}

func (x *ActiveOption) ProtoReflect() protoreflect.Message {
	mi := &file_common_filters_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActiveOption.ProtoReflect.Descriptor instead.
func (*ActiveOption) Descriptor() ([]byte, []int) {
	return file_common_filters_proto_rawDescGZIP(), []int{3}
}

func (x *ActiveOption) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *ActiveOption) GetOperator() ComparisonOperator {
	if x != nil {
		return x.Operator
	}
	return ComparisonOperator_EQUAL
}

func (x *ActiveOption) GetValue() *Variant {
	if x != nil {
		return x.Value
	}
	return nil
}

type FilterTypesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilterTypesRequest) Reset() {
	*x = FilterTypesRequest{}
	mi := &file_common_filters_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FilterTypesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterTypesRequest) ProtoMessage() {}

func (x *FilterTypesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_filters_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterTypesRequest.ProtoReflect.Descriptor instead.
func (*FilterTypesRequest) Descriptor() ([]byte, []int) {
	return file_common_filters_proto_rawDescGZIP(), []int{4}
}

type FilterTypesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilterTypes []*SupportedFilter `protobuf:"bytes,1,rep,name=filter_types,json=filterTypes,proto3" json:"filter_types,omitempty"`
}

func (x *FilterTypesResponse) Reset() {
	*x = FilterTypesResponse{}
	mi := &file_common_filters_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FilterTypesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterTypesResponse) ProtoMessage() {}

func (x *FilterTypesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_filters_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterTypesResponse.ProtoReflect.Descriptor instead.
func (*FilterTypesResponse) Descriptor() ([]byte, []int) {
	return file_common_filters_proto_rawDescGZIP(), []int{5}
}

func (x *FilterTypesResponse) GetFilterTypes() []*SupportedFilter {
	if x != nil {
		return x.FilterTypes
	}
	return nil
}

type FilterOptionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FilterOptionsRequest) Reset() {
	*x = FilterOptionsRequest{}
	mi := &file_common_filters_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FilterOptionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterOptionsRequest) ProtoMessage() {}

func (x *FilterOptionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_filters_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterOptionsRequest.ProtoReflect.Descriptor instead.
func (*FilterOptionsRequest) Descriptor() ([]byte, []int) {
	return file_common_filters_proto_rawDescGZIP(), []int{6}
}

type FilterOptionsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilterOptions []*SupportedOption `protobuf:"bytes,1,rep,name=filter_options,json=filterOptions,proto3" json:"filter_options,omitempty"`
}

func (x *FilterOptionsResponse) Reset() {
	*x = FilterOptionsResponse{}
	mi := &file_common_filters_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FilterOptionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterOptionsResponse) ProtoMessage() {}

func (x *FilterOptionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_filters_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterOptionsResponse.ProtoReflect.Descriptor instead.
func (*FilterOptionsResponse) Descriptor() ([]byte, []int) {
	return file_common_filters_proto_rawDescGZIP(), []int{7}
}

func (x *FilterOptionsResponse) GetFilterOptions() []*SupportedOption {
	if x != nil {
		return x.FilterOptions
	}
	return nil
}

var File_common_filters_proto protoreflect.FileDescriptor

var file_common_filters_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x76,
	0x31, 0x1a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x65, 0x0a, 0x0f, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x40, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x64, 0x61,
	0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x22, 0x65, 0x0a, 0x0f, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x65, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x40, 0x0a, 0x08, 0x64,
	0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e,
	0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x22, 0xa5, 0x01,
	0x0a, 0x0c, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x4b, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x2f, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x69, 0x73, 0x6f, 0x6e, 0x4f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x36, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x73,
	0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xa5, 0x01, 0x0a, 0x0c, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x4b, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2f, 0x2e, 0x73, 0x69, 0x65,
	0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x72, 0x69,
	0x73, 0x6f, 0x6e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x08, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x36, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x14, 0x0a,
	0x12, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x64, 0x0a, 0x13, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x54, 0x79, 0x70,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0c, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2a, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x70,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x0b, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x73, 0x22, 0x16, 0x0a, 0x14, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x6a, 0x0a, 0x15, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x0e, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_filters_proto_rawDescOnce sync.Once
	file_common_filters_proto_rawDescData = file_common_filters_proto_rawDesc
)

func file_common_filters_proto_rawDescGZIP() []byte {
	file_common_filters_proto_rawDescOnce.Do(func() {
		file_common_filters_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_filters_proto_rawDescData)
	})
	return file_common_filters_proto_rawDescData
}

var file_common_filters_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_common_filters_proto_goTypes = []any{
	(*SupportedFilter)(nil),       // 0: siemens.common.filters.v1.SupportedFilter
	(*SupportedOption)(nil),       // 1: siemens.common.filters.v1.SupportedOption
	(*ActiveFilter)(nil),          // 2: siemens.common.filters.v1.ActiveFilter
	(*ActiveOption)(nil),          // 3: siemens.common.filters.v1.ActiveOption
	(*FilterTypesRequest)(nil),    // 4: siemens.common.filters.v1.FilterTypesRequest
	(*FilterTypesResponse)(nil),   // 5: siemens.common.filters.v1.FilterTypesResponse
	(*FilterOptionsRequest)(nil),  // 6: siemens.common.filters.v1.FilterOptionsRequest
	(*FilterOptionsResponse)(nil), // 7: siemens.common.filters.v1.FilterOptionsResponse
	(VariantType)(0),              // 8: siemens.common.types.v1.VariantType
	(ComparisonOperator)(0),       // 9: siemens.common.operators.v1.ComparisonOperator
	(*Variant)(nil),               // 10: siemens.common.types.v1.Variant
}
var file_common_filters_proto_depIdxs = []int32{
	8,  // 0: siemens.common.filters.v1.SupportedFilter.datatype:type_name -> siemens.common.types.v1.VariantType
	8,  // 1: siemens.common.filters.v1.SupportedOption.datatype:type_name -> siemens.common.types.v1.VariantType
	9,  // 2: siemens.common.filters.v1.ActiveFilter.operator:type_name -> siemens.common.operators.v1.ComparisonOperator
	10, // 3: siemens.common.filters.v1.ActiveFilter.value:type_name -> siemens.common.types.v1.Variant
	9,  // 4: siemens.common.filters.v1.ActiveOption.operator:type_name -> siemens.common.operators.v1.ComparisonOperator
	10, // 5: siemens.common.filters.v1.ActiveOption.value:type_name -> siemens.common.types.v1.Variant
	0,  // 6: siemens.common.filters.v1.FilterTypesResponse.filter_types:type_name -> siemens.common.filters.v1.SupportedFilter
	1,  // 7: siemens.common.filters.v1.FilterOptionsResponse.filter_options:type_name -> siemens.common.filters.v1.SupportedOption
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_common_filters_proto_init() }
func file_common_filters_proto_init() {
	if File_common_filters_proto != nil {
		return
	}
	file_common_variant_proto_init()
	file_common_operators_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_filters_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_filters_proto_goTypes,
		DependencyIndexes: file_common_filters_proto_depIdxs,
		MessageInfos:      file_common_filters_proto_msgTypes,
	}.Build()
	File_common_filters_proto = out.File
	file_common_filters_proto_rawDesc = nil
	file_common_filters_proto_goTypes = nil
	file_common_filters_proto_depIdxs = nil
}
