// ------------------------------------------------------------------
// Common Definition of Device Identifiers
// ------------------------------------------------------------------

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: common_identifiers.proto

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

type IdentifierRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Identifier:
	//
	//	*IdentifierRequest_Semantic
	//	*IdentifierRequest_Name
	Identifier isIdentifierRequest_Identifier `protobuf_oneof:"identifier"`
}

func (x *IdentifierRequest) Reset() {
	*x = IdentifierRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_identifiers_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdentifierRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentifierRequest) ProtoMessage() {}

func (x *IdentifierRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_identifiers_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdentifierRequest.ProtoReflect.Descriptor instead.
func (*IdentifierRequest) Descriptor() ([]byte, []int) {
	return file_common_identifiers_proto_rawDescGZIP(), []int{0}
}

func (m *IdentifierRequest) GetIdentifier() isIdentifierRequest_Identifier {
	if m != nil {
		return m.Identifier
	}
	return nil
}

func (x *IdentifierRequest) GetSemantic() *SemanticClassifier {
	if x, ok := x.GetIdentifier().(*IdentifierRequest_Semantic); ok {
		return x.Semantic
	}
	return nil
}

func (x *IdentifierRequest) GetName() string {
	if x, ok := x.GetIdentifier().(*IdentifierRequest_Name); ok {
		return x.Name
	}
	return ""
}

type isIdentifierRequest_Identifier interface {
	isIdentifierRequest_Identifier()
}

type IdentifierRequest_Semantic struct {
	// The specific identifier which is requested
	Semantic *SemanticClassifier `protobuf:"bytes,1,opt,name=semantic,proto3,oneof"`
}

type IdentifierRequest_Name struct {
	// if a Identifier does not have a semantic mapping we can use the name of it as fallback
	Name string `protobuf:"bytes,2,opt,name=name,proto3,oneof"`
}

func (*IdentifierRequest_Semantic) isIdentifierRequest_Identifier() {}

func (*IdentifierRequest_Name) isIdentifierRequest_Identifier() {}

type GetIdentifiersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Target      *Destination         `protobuf:"bytes,1,opt,name=target,proto3" json:"target,omitempty"`           //the device to which the nodes/data points belong
	Identifiers []*IdentifierRequest `protobuf:"bytes,2,rep,name=identifiers,proto3" json:"identifiers,omitempty"` // if no specific identifiers are provided, all supported identifiers (see GetSupportedIdentifiers) are returned
}

func (x *GetIdentifiersRequest) Reset() {
	*x = GetIdentifiersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_identifiers_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetIdentifiersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetIdentifiersRequest) ProtoMessage() {}

func (x *GetIdentifiersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_identifiers_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetIdentifiersRequest.ProtoReflect.Descriptor instead.
func (*GetIdentifiersRequest) Descriptor() ([]byte, []int) {
	return file_common_identifiers_proto_rawDescGZIP(), []int{1}
}

func (x *GetIdentifiersRequest) GetTarget() *Destination {
	if x != nil {
		return x.Target
	}
	return nil
}

func (x *GetIdentifiersRequest) GetIdentifiers() []*IdentifierRequest {
	if x != nil {
		return x.Identifiers
	}
	return nil
}

type GetIdentifiersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Identifiers []*DeviceIdentifier `protobuf:"bytes,1,rep,name=identifiers,proto3" json:"identifiers,omitempty"`
}

func (x *GetIdentifiersResponse) Reset() {
	*x = GetIdentifiersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_identifiers_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetIdentifiersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetIdentifiersResponse) ProtoMessage() {}

func (x *GetIdentifiersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_identifiers_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetIdentifiersResponse.ProtoReflect.Descriptor instead.
func (*GetIdentifiersResponse) Descriptor() ([]byte, []int) {
	return file_common_identifiers_proto_rawDescGZIP(), []int{2}
}

func (x *GetIdentifiersResponse) GetIdentifiers() []*DeviceIdentifier {
	if x != nil {
		return x.Identifiers
	}
	return nil
}

type GetSupportedSemanticsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetSupportedSemanticsRequest) Reset() {
	*x = GetSupportedSemanticsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_identifiers_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSupportedSemanticsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSupportedSemanticsRequest) ProtoMessage() {}

func (x *GetSupportedSemanticsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_identifiers_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSupportedSemanticsRequest.ProtoReflect.Descriptor instead.
func (*GetSupportedSemanticsRequest) Descriptor() ([]byte, []int) {
	return file_common_identifiers_proto_rawDescGZIP(), []int{3}
}

type GetSupportedSemanticsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SupportedSemantic []*SupportedSemantic `protobuf:"bytes,1,rep,name=supportedSemantic,proto3" json:"supportedSemantic,omitempty"`
}

func (x *GetSupportedSemanticsResponse) Reset() {
	*x = GetSupportedSemanticsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_identifiers_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSupportedSemanticsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSupportedSemanticsResponse) ProtoMessage() {}

func (x *GetSupportedSemanticsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_common_identifiers_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSupportedSemanticsResponse.ProtoReflect.Descriptor instead.
func (*GetSupportedSemanticsResponse) Descriptor() ([]byte, []int) {
	return file_common_identifiers_proto_rawDescGZIP(), []int{4}
}

func (x *GetSupportedSemanticsResponse) GetSupportedSemantic() []*SupportedSemantic {
	if x != nil {
		return x.SupportedSemantic
	}
	return nil
}

// This message type is defined to work-around missing
// support for "repeated oneof" in gRPC. See
// https://github.com/protocolbuffers/protobuf/issues/2592
// for further reference.
type DeviceIdentifierValueList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []*DeviceIdentifier `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *DeviceIdentifierValueList) Reset() {
	*x = DeviceIdentifierValueList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_identifiers_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceIdentifierValueList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceIdentifierValueList) ProtoMessage() {}

func (x *DeviceIdentifierValueList) ProtoReflect() protoreflect.Message {
	mi := &file_common_identifiers_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceIdentifierValueList.ProtoReflect.Descriptor instead.
func (*DeviceIdentifierValueList) Descriptor() ([]byte, []int) {
	return file_common_identifiers_proto_rawDescGZIP(), []int{5}
}

func (x *DeviceIdentifierValueList) GetValue() []*DeviceIdentifier {
	if x != nil {
		return x.Value
	}
	return nil
}

type DeviceIdentifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The RAW value of the identifier, read from the device (connection)
	//
	// Types that are assignable to Value:
	//
	//	*DeviceIdentifier_Int64Value
	//	*DeviceIdentifier_Uint64Value
	//	*DeviceIdentifier_Float64Value
	//	*DeviceIdentifier_Text
	//	*DeviceIdentifier_RawData
	//	*DeviceIdentifier_Children
	Value isDeviceIdentifier_Value `protobuf_oneof:"value"`
	// List of semantic mappings for this identifier
	Classifiers []*SemanticClassifier `protobuf:"bytes,8,rep,name=classifiers,proto3" json:"classifiers,omitempty"`
}

func (x *DeviceIdentifier) Reset() {
	*x = DeviceIdentifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_identifiers_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceIdentifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceIdentifier) ProtoMessage() {}

func (x *DeviceIdentifier) ProtoReflect() protoreflect.Message {
	mi := &file_common_identifiers_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceIdentifier.ProtoReflect.Descriptor instead.
func (*DeviceIdentifier) Descriptor() ([]byte, []int) {
	return file_common_identifiers_proto_rawDescGZIP(), []int{6}
}

func (m *DeviceIdentifier) GetValue() isDeviceIdentifier_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *DeviceIdentifier) GetInt64Value() int64 {
	if x, ok := x.GetValue().(*DeviceIdentifier_Int64Value); ok {
		return x.Int64Value
	}
	return 0
}

func (x *DeviceIdentifier) GetUint64Value() uint64 {
	if x, ok := x.GetValue().(*DeviceIdentifier_Uint64Value); ok {
		return x.Uint64Value
	}
	return 0
}

func (x *DeviceIdentifier) GetFloat64Value() float64 {
	if x, ok := x.GetValue().(*DeviceIdentifier_Float64Value); ok {
		return x.Float64Value
	}
	return 0
}

func (x *DeviceIdentifier) GetText() string {
	if x, ok := x.GetValue().(*DeviceIdentifier_Text); ok {
		return x.Text
	}
	return ""
}

func (x *DeviceIdentifier) GetRawData() []byte {
	if x, ok := x.GetValue().(*DeviceIdentifier_RawData); ok {
		return x.RawData
	}
	return nil
}

func (x *DeviceIdentifier) GetChildren() *DeviceIdentifierValueList {
	if x, ok := x.GetValue().(*DeviceIdentifier_Children); ok {
		return x.Children
	}
	return nil
}

func (x *DeviceIdentifier) GetClassifiers() []*SemanticClassifier {
	if x != nil {
		return x.Classifiers
	}
	return nil
}

type isDeviceIdentifier_Value interface {
	isDeviceIdentifier_Value()
}

type DeviceIdentifier_Int64Value struct {
	// Transfer any integer value up to 64 bit
	Int64Value int64 `protobuf:"varint,2,opt,name=int64_value,json=int64Value,proto3,oneof"`
}

type DeviceIdentifier_Uint64Value struct {
	// Transfer any unsigned integer value up to 64 bit
	Uint64Value uint64 `protobuf:"varint,3,opt,name=uint64_value,json=uint64Value,proto3,oneof"`
}

type DeviceIdentifier_Float64Value struct {
	// Transfer any floating point value
	Float64Value float64 `protobuf:"fixed64,4,opt,name=float64_value,json=float64Value,proto3,oneof"`
}

type DeviceIdentifier_Text struct {
	// Transfer a UTF8-text
	Text string `protobuf:"bytes,5,opt,name=text,proto3,oneof"`
}

type DeviceIdentifier_RawData struct {
	// Raw data
	RawData []byte `protobuf:"bytes,6,opt,name=raw_data,json=rawData,proto3,oneof"`
}

type DeviceIdentifier_Children struct {
	// A list of child identifiers. This value type
	// can be used to transfer structured, hierarchical
	// information.
	// Example: Ethernet interface (parent identifier) with
	// its assigned IP addresses (child identifiers) and MAC
	// address (another child identifier).
	Children *DeviceIdentifierValueList `protobuf:"bytes,7,opt,name=children,proto3,oneof"`
}

func (*DeviceIdentifier_Int64Value) isDeviceIdentifier_Value() {}

func (*DeviceIdentifier_Uint64Value) isDeviceIdentifier_Value() {}

func (*DeviceIdentifier_Float64Value) isDeviceIdentifier_Value() {}

func (*DeviceIdentifier_Text) isDeviceIdentifier_Value() {}

func (*DeviceIdentifier_RawData) isDeviceIdentifier_Value() {}

func (*DeviceIdentifier_Children) isDeviceIdentifier_Value() {}

type SupportedSemantic struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Human readable name of the identifier, e.g. "Manufacturer"
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// List of semantic identifiers
	Classifiers []*SemanticClassifier `protobuf:"bytes,2,rep,name=classifiers,proto3" json:"classifiers,omitempty"`
}

func (x *SupportedSemantic) Reset() {
	*x = SupportedSemantic{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_identifiers_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SupportedSemantic) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SupportedSemantic) ProtoMessage() {}

func (x *SupportedSemantic) ProtoReflect() protoreflect.Message {
	mi := &file_common_identifiers_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SupportedSemantic.ProtoReflect.Descriptor instead.
func (*SupportedSemantic) Descriptor() ([]byte, []int) {
	return file_common_identifiers_proto_rawDescGZIP(), []int{7}
}

func (x *SupportedSemantic) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SupportedSemantic) GetClassifiers() []*SemanticClassifier {
	if x != nil {
		return x.Classifiers
	}
	return nil
}

type SemanticClassifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type of the semantic identifier
	// supported default classifier:
	// "IRDI": e.g. for "CDD" or "eClass" classifier)
	// "URI": e.g. for JSON-Link
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// The value of the semantic identifier, e.g. "0112/2///61987#ABA565#007" (IEC CDD/IRDI "Manufacturer")
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *SemanticClassifier) Reset() {
	*x = SemanticClassifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_identifiers_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SemanticClassifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SemanticClassifier) ProtoMessage() {}

func (x *SemanticClassifier) ProtoReflect() protoreflect.Message {
	mi := &file_common_identifiers_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SemanticClassifier.ProtoReflect.Descriptor instead.
func (*SemanticClassifier) Descriptor() ([]byte, []int) {
	return file_common_identifiers_proto_rawDescGZIP(), []int{8}
}

func (x *SemanticClassifier) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *SemanticClassifier) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_common_identifiers_proto protoreflect.FileDescriptor

var file_common_identifiers_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x73, 0x69, 0x65, 0x6d,
	0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x88, 0x01, 0x0a, 0x11, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x4f, 0x0a, 0x08, 0x73, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69,
	0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63,
	0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x48, 0x00, 0x52, 0x08, 0x73, 0x65,
	0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x0c, 0x0a, 0x0a,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x22, 0xab, 0x01, 0x0a, 0x15, 0x47,
	0x65, 0x74, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x3e, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x12, 0x52, 0x0a, 0x0b, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69,
	0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x73, 0x69, 0x65, 0x6d,
	0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x0b, 0x69, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x22, 0x6b, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x49,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x51, 0x0a, 0x0b, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x0b, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x73, 0x22, 0x1e, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x53, 0x75, 0x70, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x64, 0x53, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x7f, 0x0a, 0x1d, 0x47, 0x65, 0x74, 0x53, 0x75, 0x70, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x64, 0x53, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5e, 0x0a, 0x11, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x65, 0x64, 0x53, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x30, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x53, 0x65, 0x6d, 0x61, 0x6e,
	0x74, 0x69, 0x63, 0x52, 0x11, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x53, 0x65,
	0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x22, 0x62, 0x0a, 0x19, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x45, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xea, 0x02, 0x0a, 0x10, 0x44,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12,
	0x21, 0x0a, 0x0b, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x0b, 0x75, 0x69, 0x6e, 0x74,
	0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x25, 0x0a, 0x0d, 0x66, 0x6c, 0x6f, 0x61, 0x74,
	0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00,
	0x52, 0x0c, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14,
	0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x12, 0x1b, 0x0a, 0x08, 0x72, 0x61, 0x77, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x07, 0x72, 0x61, 0x77, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x56, 0x0a, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x00, 0x52,
	0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x12, 0x53, 0x0a, 0x0b, 0x63, 0x6c, 0x61,
	0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31,
	0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x52, 0x0b, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x42, 0x07,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x7c, 0x0a, 0x11, 0x53, 0x75, 0x70, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x64, 0x53, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x53, 0x0a, 0x0b, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x43, 0x6c,
	0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x0b, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x73, 0x22, 0x3e, 0x0a, 0x12, 0x53, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69,
	0x63, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0xa8, 0x02, 0x0a, 0x0e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x73, 0x41, 0x70, 0x69, 0x12, 0x7f, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x49,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x12, 0x34, 0x2e, 0x73, 0x69, 0x65,
	0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x35, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x94, 0x01, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x53, 0x65, 0x6d, 0x61, 0x6e, 0x74,
	0x69, 0x63, 0x73, 0x12, 0x3b, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64,
	0x53, 0x65, 0x6d, 0x61, 0x6e, 0x74, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x3c, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x53, 0x65, 0x6d,
	0x61, 0x6e, 0x74, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_identifiers_proto_rawDescOnce sync.Once
	file_common_identifiers_proto_rawDescData = file_common_identifiers_proto_rawDesc
)

func file_common_identifiers_proto_rawDescGZIP() []byte {
	file_common_identifiers_proto_rawDescOnce.Do(func() {
		file_common_identifiers_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_identifiers_proto_rawDescData)
	})
	return file_common_identifiers_proto_rawDescData
}

var file_common_identifiers_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_common_identifiers_proto_goTypes = []interface{}{
	(*IdentifierRequest)(nil),             // 0: siemens.common.identifiers.v1.IdentifierRequest
	(*GetIdentifiersRequest)(nil),         // 1: siemens.common.identifiers.v1.GetIdentifiersRequest
	(*GetIdentifiersResponse)(nil),        // 2: siemens.common.identifiers.v1.GetIdentifiersResponse
	(*GetSupportedSemanticsRequest)(nil),  // 3: siemens.common.identifiers.v1.GetSupportedSemanticsRequest
	(*GetSupportedSemanticsResponse)(nil), // 4: siemens.common.identifiers.v1.GetSupportedSemanticsResponse
	(*DeviceIdentifierValueList)(nil),     // 5: siemens.common.identifiers.v1.DeviceIdentifierValueList
	(*DeviceIdentifier)(nil),              // 6: siemens.common.identifiers.v1.DeviceIdentifier
	(*SupportedSemantic)(nil),             // 7: siemens.common.identifiers.v1.SupportedSemantic
	(*SemanticClassifier)(nil),            // 8: siemens.common.identifiers.v1.SemanticClassifier
	(*Destination)(nil),                   // 9: siemens.common.address.v1.Destination
}
var file_common_identifiers_proto_depIdxs = []int32{
	8,  // 0: siemens.common.identifiers.v1.IdentifierRequest.semantic:type_name -> siemens.common.identifiers.v1.SemanticClassifier
	9,  // 1: siemens.common.identifiers.v1.GetIdentifiersRequest.target:type_name -> siemens.common.address.v1.Destination
	0,  // 2: siemens.common.identifiers.v1.GetIdentifiersRequest.identifiers:type_name -> siemens.common.identifiers.v1.IdentifierRequest
	6,  // 3: siemens.common.identifiers.v1.GetIdentifiersResponse.identifiers:type_name -> siemens.common.identifiers.v1.DeviceIdentifier
	7,  // 4: siemens.common.identifiers.v1.GetSupportedSemanticsResponse.supportedSemantic:type_name -> siemens.common.identifiers.v1.SupportedSemantic
	6,  // 5: siemens.common.identifiers.v1.DeviceIdentifierValueList.value:type_name -> siemens.common.identifiers.v1.DeviceIdentifier
	5,  // 6: siemens.common.identifiers.v1.DeviceIdentifier.children:type_name -> siemens.common.identifiers.v1.DeviceIdentifierValueList
	8,  // 7: siemens.common.identifiers.v1.DeviceIdentifier.classifiers:type_name -> siemens.common.identifiers.v1.SemanticClassifier
	8,  // 8: siemens.common.identifiers.v1.SupportedSemantic.classifiers:type_name -> siemens.common.identifiers.v1.SemanticClassifier
	1,  // 9: siemens.common.identifiers.v1.IdentifiersApi.GetIdentifiers:input_type -> siemens.common.identifiers.v1.GetIdentifiersRequest
	3,  // 10: siemens.common.identifiers.v1.IdentifiersApi.GetSupportedSemantics:input_type -> siemens.common.identifiers.v1.GetSupportedSemanticsRequest
	2,  // 11: siemens.common.identifiers.v1.IdentifiersApi.GetIdentifiers:output_type -> siemens.common.identifiers.v1.GetIdentifiersResponse
	4,  // 12: siemens.common.identifiers.v1.IdentifiersApi.GetSupportedSemantics:output_type -> siemens.common.identifiers.v1.GetSupportedSemanticsResponse
	11, // [11:13] is the sub-list for method output_type
	9,  // [9:11] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_common_identifiers_proto_init() }
func file_common_identifiers_proto_init() {
	if File_common_identifiers_proto != nil {
		return
	}
	file_common_address_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_common_identifiers_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdentifierRequest); i {
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
		file_common_identifiers_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetIdentifiersRequest); i {
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
		file_common_identifiers_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetIdentifiersResponse); i {
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
		file_common_identifiers_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSupportedSemanticsRequest); i {
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
		file_common_identifiers_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSupportedSemanticsResponse); i {
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
		file_common_identifiers_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceIdentifierValueList); i {
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
		file_common_identifiers_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceIdentifier); i {
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
		file_common_identifiers_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SupportedSemantic); i {
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
		file_common_identifiers_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SemanticClassifier); i {
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
	file_common_identifiers_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*IdentifierRequest_Semantic)(nil),
		(*IdentifierRequest_Name)(nil),
	}
	file_common_identifiers_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*DeviceIdentifier_Int64Value)(nil),
		(*DeviceIdentifier_Uint64Value)(nil),
		(*DeviceIdentifier_Float64Value)(nil),
		(*DeviceIdentifier_Text)(nil),
		(*DeviceIdentifier_RawData)(nil),
		(*DeviceIdentifier_Children)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_identifiers_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_common_identifiers_proto_goTypes,
		DependencyIndexes: file_common_identifiers_proto_depIdxs,
		MessageInfos:      file_common_identifiers_proto_msgTypes,
	}.Build()
	File_common_identifiers_proto = out.File
	file_common_identifiers_proto_rawDesc = nil
	file_common_identifiers_proto_goTypes = nil
	file_common_identifiers_proto_depIdxs = nil
}
