// ------------------------------------------------------------------
// Common definition of types used by several APIs
// ------------------------------------------------------------------
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: common_variant.proto

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

type VariantType int32

const (
	VariantType_VT_UNSPECIFIED VariantType = 0
	VariantType_VT_BOOL        VariantType = 1
	VariantType_VT_INT64       VariantType = 2
	VariantType_VT_UINT64      VariantType = 3
	VariantType_VT_DOUBLE      VariantType = 4
	VariantType_VT_STRING      VariantType = 5
	VariantType_VT_BYTES       VariantType = 6
)

// Enum value maps for VariantType.
var (
	VariantType_name = map[int32]string{
		0: "VT_UNSPECIFIED",
		1: "VT_BOOL",
		2: "VT_INT64",
		3: "VT_UINT64",
		4: "VT_DOUBLE",
		5: "VT_STRING",
		6: "VT_BYTES",
	}
	VariantType_value = map[string]int32{
		"VT_UNSPECIFIED": 0,
		"VT_BOOL":        1,
		"VT_INT64":       2,
		"VT_UINT64":      3,
		"VT_DOUBLE":      4,
		"VT_STRING":      5,
		"VT_BYTES":       6,
	}
)

func (x VariantType) Enum() *VariantType {
	p := new(VariantType)
	*p = x
	return p
}

func (x VariantType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VariantType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_variant_proto_enumTypes[0].Descriptor()
}

func (VariantType) Type() protoreflect.EnumType {
	return &file_common_variant_proto_enumTypes[0]
}

func (x VariantType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VariantType.Descriptor instead.
func (VariantType) EnumDescriptor() ([]byte, []int) {
	return file_common_variant_proto_rawDescGZIP(), []int{0}
}

type Variant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//
	//	*Variant_BoolValue
	//	*Variant_Int64Value
	//	*Variant_Uint64Value
	//	*Variant_Float64Value
	//	*Variant_Text
	//	*Variant_RawData
	Value isVariant_Value `protobuf_oneof:"value"`
}

func (x *Variant) Reset() {
	*x = Variant{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_variant_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Variant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Variant) ProtoMessage() {}

func (x *Variant) ProtoReflect() protoreflect.Message {
	mi := &file_common_variant_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Variant.ProtoReflect.Descriptor instead.
func (*Variant) Descriptor() ([]byte, []int) {
	return file_common_variant_proto_rawDescGZIP(), []int{0}
}

func (m *Variant) GetValue() isVariant_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *Variant) GetBoolValue() bool {
	if x, ok := x.GetValue().(*Variant_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

func (x *Variant) GetInt64Value() int64 {
	if x, ok := x.GetValue().(*Variant_Int64Value); ok {
		return x.Int64Value
	}
	return 0
}

func (x *Variant) GetUint64Value() uint64 {
	if x, ok := x.GetValue().(*Variant_Uint64Value); ok {
		return x.Uint64Value
	}
	return 0
}

func (x *Variant) GetFloat64Value() float64 {
	if x, ok := x.GetValue().(*Variant_Float64Value); ok {
		return x.Float64Value
	}
	return 0
}

func (x *Variant) GetText() string {
	if x, ok := x.GetValue().(*Variant_Text); ok {
		return x.Text
	}
	return ""
}

func (x *Variant) GetRawData() []byte {
	if x, ok := x.GetValue().(*Variant_RawData); ok {
		return x.RawData
	}
	return nil
}

type isVariant_Value interface {
	isVariant_Value()
}

type Variant_BoolValue struct {
	// Simple bool value
	BoolValue bool `protobuf:"varint,1,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

type Variant_Int64Value struct {
	// Transfer any integer value up to 64 bit
	Int64Value int64 `protobuf:"varint,2,opt,name=int64_value,json=int64Value,proto3,oneof"`
}

type Variant_Uint64Value struct {
	// Transfer any unsigned integer value up to 64 bit
	Uint64Value uint64 `protobuf:"varint,3,opt,name=uint64_value,json=uint64Value,proto3,oneof"`
}

type Variant_Float64Value struct {
	// Transfer any floating point value
	Float64Value float64 `protobuf:"fixed64,4,opt,name=float64_value,json=float64Value,proto3,oneof"`
}

type Variant_Text struct {
	// Transfer a UTF8-text
	Text string `protobuf:"bytes,5,opt,name=text,proto3,oneof"`
}

type Variant_RawData struct {
	// Transfer array of bytes
	// Example for raw-data: S7-1500 system diagnosis data
	RawData []byte `protobuf:"bytes,6,opt,name=raw_data,json=rawData,proto3,oneof"`
}

func (*Variant_BoolValue) isVariant_Value() {}

func (*Variant_Int64Value) isVariant_Value() {}

func (*Variant_Uint64Value) isVariant_Value() {}

func (*Variant_Float64Value) isVariant_Value() {}

func (*Variant_Text) isVariant_Value() {}

func (*Variant_RawData) isVariant_Value() {}

var File_common_variant_proto protoreflect.FileDescriptor

var file_common_variant_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x17, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x22,
	0xd5, 0x01, 0x0a, 0x07, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0a, 0x62,
	0x6f, 0x6f, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48,
	0x00, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x0b,
	0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x23, 0x0a, 0x0c, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00, 0x52, 0x0b, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x25, 0x0a, 0x0d, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x36, 0x34, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x0c, 0x66,
	0x6c, 0x6f, 0x61, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x12, 0x1b, 0x0a, 0x08, 0x72, 0x61, 0x77, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x07, 0x72, 0x61, 0x77, 0x44, 0x61, 0x74, 0x61, 0x42, 0x07,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0x77, 0x0a, 0x0b, 0x56, 0x61, 0x72, 0x69, 0x61,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x56, 0x54, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x56, 0x54,
	0x5f, 0x42, 0x4f, 0x4f, 0x4c, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x56, 0x54, 0x5f, 0x49, 0x4e,
	0x54, 0x36, 0x34, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x56, 0x54, 0x5f, 0x55, 0x49, 0x4e, 0x54,
	0x36, 0x34, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x56, 0x54, 0x5f, 0x44, 0x4f, 0x55, 0x42, 0x4c,
	0x45, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x56, 0x54, 0x5f, 0x53, 0x54, 0x52, 0x49, 0x4e, 0x47,
	0x10, 0x05, 0x12, 0x0c, 0x0a, 0x08, 0x56, 0x54, 0x5f, 0x42, 0x59, 0x54, 0x45, 0x53, 0x10, 0x06,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_variant_proto_rawDescOnce sync.Once
	file_common_variant_proto_rawDescData = file_common_variant_proto_rawDesc
)

func file_common_variant_proto_rawDescGZIP() []byte {
	file_common_variant_proto_rawDescOnce.Do(func() {
		file_common_variant_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_variant_proto_rawDescData)
	})
	return file_common_variant_proto_rawDescData
}

var file_common_variant_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_common_variant_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_variant_proto_goTypes = []interface{}{
	(VariantType)(0), // 0: siemens.common.types.v1.VariantType
	(*Variant)(nil),  // 1: siemens.common.types.v1.Variant
}
var file_common_variant_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_variant_proto_init() }
func file_common_variant_proto_init() {
	if File_common_variant_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_variant_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Variant); i {
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
	file_common_variant_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Variant_BoolValue)(nil),
		(*Variant_Int64Value)(nil),
		(*Variant_Uint64Value)(nil),
		(*Variant_Float64Value)(nil),
		(*Variant_Text)(nil),
		(*Variant_RawData)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_variant_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_variant_proto_goTypes,
		DependencyIndexes: file_common_variant_proto_depIdxs,
		EnumInfos:         file_common_variant_proto_enumTypes,
		MessageInfos:      file_common_variant_proto_msgTypes,
	}.Build()
	File_common_variant_proto = out.File
	file_common_variant_proto_rawDesc = nil
	file_common_variant_proto_goTypes = nil
	file_common_variant_proto_depIdxs = nil
}
