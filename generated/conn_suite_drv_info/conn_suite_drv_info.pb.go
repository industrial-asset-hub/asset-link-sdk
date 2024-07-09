// ------------------------------------------------------------------
// Connectivity Suite Driver Version Service
// ------------------------------------------------------------------
//
// Naming convention according:
// https://cloud.google.com/apis/design/naming_convention
//
// ------------------------------------------------------------------

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.24.4
// source: conn_suite_drv_info.proto

// Include when using Protobuf-Any
// import "google/protobuf/any.proto";

package conn_suite_drv_info

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

// ==================================================================
// Driver Version
type GetVersionInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVersionInfoRequest) Reset() {
	*x = GetVersionInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_suite_drv_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersionInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersionInfoRequest) ProtoMessage() {}

func (x *GetVersionInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_drv_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersionInfoRequest.ProtoReflect.Descriptor instead.
func (*GetVersionInfoRequest) Descriptor() ([]byte, []int) {
	return file_conn_suite_drv_info_proto_rawDescGZIP(), []int{0}
}

// Version Info
type VersionInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Version numbering according 'Semantic Versioning'
	// (see https://semver.org/)
	// Major - increment for incompatible API changes
	Major uint32 `protobuf:"varint,1,opt,name=major,proto3" json:"major,omitempty"`
	// Minor - increment for added functionality in a backwards compatible manner
	Minor uint32 `protobuf:"varint,2,opt,name=minor,proto3" json:"minor,omitempty"`
	// Patch - increment for backwards compatible bug fixes
	Patch uint32 `protobuf:"varint,3,opt,name=patch,proto3" json:"patch,omitempty"`
	// Suffix - containing Build number and/or pre-release version.
	// According to the version definition of Industrial Edge OR to https://semver.org/
	// Don't expect the string to strictly follow semver, especially for checking
	// which version is newer!!!!!
	// Can be an empty string.
	// Industrial Edge always uses "-" as the first character which is a violation of semver!
	// Here some examples for Industrial Edge version suffixes:
	//   - ""
	//   - "-0"
	//   - "-1"
	//   - "-rc.1"
	//   - "-rc.1.alpha.23773115"
	//   - "-beta.2.rc.23652691"
	//   - "-3.0"
	Suffix string `protobuf:"bytes,7,opt,name=suffix,proto3" json:"suffix,omitempty"`
	// vendor name, e.g. "Siemens AG"
	VendorName string `protobuf:"bytes,4,opt,name=vendor_name,json=vendorName,proto3" json:"vendor_name,omitempty"`
	// product name, e.g. "CS S7-1500 Driver"
	ProductName string `protobuf:"bytes,5,opt,name=product_name,json=productName,proto3" json:"product_name,omitempty"`
	// the documentation URL of the driver
	// e.g. the company webpage with a deep link directly the docu
	DocuUrl string `protobuf:"bytes,6,opt,name=docu_url,json=docuUrl,proto3" json:"docu_url,omitempty"`
	// feedback url for customers, it's different for different products
	FeedbackUrl string `protobuf:"bytes,8,opt,name=feedback_url,json=feedbackUrl,proto3" json:"feedback_url,omitempty"`
}

func (x *VersionInfo) Reset() {
	*x = VersionInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_suite_drv_info_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VersionInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VersionInfo) ProtoMessage() {}

func (x *VersionInfo) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_drv_info_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VersionInfo.ProtoReflect.Descriptor instead.
func (*VersionInfo) Descriptor() ([]byte, []int) {
	return file_conn_suite_drv_info_proto_rawDescGZIP(), []int{1}
}

func (x *VersionInfo) GetMajor() uint32 {
	if x != nil {
		return x.Major
	}
	return 0
}

func (x *VersionInfo) GetMinor() uint32 {
	if x != nil {
		return x.Minor
	}
	return 0
}

func (x *VersionInfo) GetPatch() uint32 {
	if x != nil {
		return x.Patch
	}
	return 0
}

func (x *VersionInfo) GetSuffix() string {
	if x != nil {
		return x.Suffix
	}
	return ""
}

func (x *VersionInfo) GetVendorName() string {
	if x != nil {
		return x.VendorName
	}
	return ""
}

func (x *VersionInfo) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}

func (x *VersionInfo) GetDocuUrl() string {
	if x != nil {
		return x.DocuUrl
	}
	return ""
}

func (x *VersionInfo) GetFeedbackUrl() string {
	if x != nil {
		return x.FeedbackUrl
	}
	return ""
}

type GetVersionInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// version information
	Version *VersionInfo `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetVersionInfoResponse) Reset() {
	*x = GetVersionInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_suite_drv_info_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersionInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersionInfoResponse) ProtoMessage() {}

func (x *GetVersionInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_drv_info_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersionInfoResponse.ProtoReflect.Descriptor instead.
func (*GetVersionInfoResponse) Descriptor() ([]byte, []int) {
	return file_conn_suite_drv_info_proto_rawDescGZIP(), []int{2}
}

func (x *GetVersionInfoResponse) GetVersion() *VersionInfo {
	if x != nil {
		return x.Version
	}
	return nil
}

// ==================================================================
// Config Schema
type GetConfigSchemaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetConfigSchemaRequest) Reset() {
	*x = GetConfigSchemaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_suite_drv_info_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetConfigSchemaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConfigSchemaRequest) ProtoMessage() {}

func (x *GetConfigSchemaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_drv_info_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConfigSchemaRequest.ProtoReflect.Descriptor instead.
func (*GetConfigSchemaRequest) Descriptor() ([]byte, []int) {
	return file_conn_suite_drv_info_proto_rawDescGZIP(), []int{3}
}

type ConfigSchema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// URI of the schema
	Uri string `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
	// JSON string with configuration schema
	Schema string `protobuf:"bytes,2,opt,name=schema,proto3" json:"schema,omitempty"`
}

func (x *ConfigSchema) Reset() {
	*x = ConfigSchema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_suite_drv_info_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigSchema) ProtoMessage() {}

func (x *ConfigSchema) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_drv_info_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigSchema.ProtoReflect.Descriptor instead.
func (*ConfigSchema) Descriptor() ([]byte, []int) {
	return file_conn_suite_drv_info_proto_rawDescGZIP(), []int{4}
}

func (x *ConfigSchema) GetUri() string {
	if x != nil {
		return x.Uri
	}
	return ""
}

func (x *ConfigSchema) GetSchema() string {
	if x != nil {
		return x.Schema
	}
	return ""
}

type GetConfigSchemaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// configuration schema(s)
	Schemas []*ConfigSchema `protobuf:"bytes,1,rep,name=schemas,proto3" json:"schemas,omitempty"`
}

func (x *GetConfigSchemaResponse) Reset() {
	*x = GetConfigSchemaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_suite_drv_info_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetConfigSchemaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConfigSchemaResponse) ProtoMessage() {}

func (x *GetConfigSchemaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_drv_info_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConfigSchemaResponse.ProtoReflect.Descriptor instead.
func (*GetConfigSchemaResponse) Descriptor() ([]byte, []int) {
	return file_conn_suite_drv_info_proto_rawDescGZIP(), []int{5}
}

func (x *GetConfigSchemaResponse) GetSchemas() []*ConfigSchema {
	if x != nil {
		return x.Schemas
	}
	return nil
}

// ==================================================================
// Get Icon
type GetAppIconRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetAppIconRequest) Reset() {
	*x = GetAppIconRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_suite_drv_info_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAppIconRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAppIconRequest) ProtoMessage() {}

func (x *GetAppIconRequest) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_drv_info_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAppIconRequest.ProtoReflect.Descriptor instead.
func (*GetAppIconRequest) Descriptor() ([]byte, []int) {
	return file_conn_suite_drv_info_proto_rawDescGZIP(), []int{6}
}

type GetAppIconResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Image format of the returned image.
	// The aspect rate of the image shall be 1:1 that means a square image.
	// Server should select the format with best fit for the client request
	// (best effort, don't scale the image to this size!).
	// The format names are according to image mime-type defined by
	//
	//	https://www.iana.org/assignments/media-types/media-types.xhtml#image
	//	but without the "image/" prefix because that is obligatory.
	//
	// e.g. "png" or "svg+xml"
	// At the moment the only supported format is "png". We will decide in future
	// if we support more formats and introduce 'supported_image_formats'
	// in the request then.
	ImageFormat string `protobuf:"bytes,1,opt,name=image_format,json=imageFormat,proto3" json:"image_format,omitempty"`
	// Byte array containing the image data
	ImageData []byte `protobuf:"bytes,2,opt,name=image_data,json=imageData,proto3" json:"image_data,omitempty"`
}

func (x *GetAppIconResponse) Reset() {
	*x = GetAppIconResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_suite_drv_info_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAppIconResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAppIconResponse) ProtoMessage() {}

func (x *GetAppIconResponse) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_drv_info_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAppIconResponse.ProtoReflect.Descriptor instead.
func (*GetAppIconResponse) Descriptor() ([]byte, []int) {
	return file_conn_suite_drv_info_proto_rawDescGZIP(), []int{7}
}

func (x *GetAppIconResponse) GetImageFormat() string {
	if x != nil {
		return x.ImageFormat
	}
	return ""
}

func (x *GetAppIconResponse) GetImageData() []byte {
	if x != nil {
		return x.ImageData
	}
	return nil
}

var File_conn_suite_drv_info_proto protoreflect.FileDescriptor

var file_conn_suite_drv_info_proto_rawDesc = []byte{
	0x0a, 0x19, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x73, 0x75, 0x69, 0x74, 0x65, 0x5f, 0x64, 0x72, 0x76,
	0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x24, 0x73, 0x69, 0x65,
	0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74,
	0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x64, 0x72, 0x76, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x76,
	0x31, 0x22, 0x17, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xe9, 0x01, 0x0a, 0x0b, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x61,
	0x6a, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x6d, 0x61, 0x6a, 0x6f, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x05, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x61, 0x74, 0x63, 0x68, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x70, 0x61, 0x74, 0x63, 0x68, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x75, 0x66, 0x66, 0x69, 0x78, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x75,
	0x66, 0x66, 0x69, 0x78, 0x12, 0x1f, 0x0a, 0x0b, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x76, 0x65, 0x6e, 0x64, 0x6f,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x75,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x6f, 0x63, 0x75,
	0x55, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x65, 0x65, 0x64, 0x62, 0x61, 0x63, 0x6b, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x66, 0x65, 0x65, 0x64, 0x62,
	0x61, 0x63, 0x6b, 0x55, 0x72, 0x6c, 0x22, 0x65, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x4b, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x31, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x64, 0x72,
	0x76, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x18, 0x0a,
	0x16, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x38, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x69, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x69, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x22, 0x67, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x07,
	0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e,
	0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x64, 0x72, 0x76, 0x69, 0x6e, 0x66,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x52, 0x07, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x73, 0x22, 0x13, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x41, 0x70, 0x70, 0x49, 0x63, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x56, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x49, 0x63, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x32, 0xb6, 0x03, 0x0a, 0x0d, 0x44, 0x72, 0x69, 0x76,
	0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x41, 0x70, 0x69, 0x12, 0x8d, 0x01, 0x0a, 0x0e, 0x47, 0x65,
	0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x3b, 0x2e, 0x73,
	0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76,
	0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x64, 0x72, 0x76, 0x69, 0x6e, 0x66, 0x6f,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3c, 0x2e, 0x73, 0x69, 0x65, 0x6d,
	0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x64, 0x72, 0x76, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x90, 0x01, 0x0a, 0x0f, 0x47, 0x65,
	0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x3c, 0x2e,
	0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x64, 0x72, 0x76, 0x69, 0x6e, 0x66,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3d, 0x2e, 0x73, 0x69,
	0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x64, 0x72, 0x76, 0x69, 0x6e, 0x66, 0x6f, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x81, 0x01, 0x0a,
	0x0a, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x49, 0x63, 0x6f, 0x6e, 0x12, 0x37, 0x2e, 0x73, 0x69,
	0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x64, 0x72, 0x76, 0x69, 0x6e, 0x66, 0x6f, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x49, 0x63, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x38, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65,
	0x2e, 0x64, 0x72, 0x76, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x70, 0x70, 0x49, 0x63, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_conn_suite_drv_info_proto_rawDescOnce sync.Once
	file_conn_suite_drv_info_proto_rawDescData = file_conn_suite_drv_info_proto_rawDesc
)

func file_conn_suite_drv_info_proto_rawDescGZIP() []byte {
	file_conn_suite_drv_info_proto_rawDescOnce.Do(func() {
		file_conn_suite_drv_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_conn_suite_drv_info_proto_rawDescData)
	})
	return file_conn_suite_drv_info_proto_rawDescData
}

var file_conn_suite_drv_info_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_conn_suite_drv_info_proto_goTypes = []interface{}{
	(*GetVersionInfoRequest)(nil),   // 0: siemens.connectivitysuite.drvinfo.v1.GetVersionInfoRequest
	(*VersionInfo)(nil),             // 1: siemens.connectivitysuite.drvinfo.v1.VersionInfo
	(*GetVersionInfoResponse)(nil),  // 2: siemens.connectivitysuite.drvinfo.v1.GetVersionInfoResponse
	(*GetConfigSchemaRequest)(nil),  // 3: siemens.connectivitysuite.drvinfo.v1.GetConfigSchemaRequest
	(*ConfigSchema)(nil),            // 4: siemens.connectivitysuite.drvinfo.v1.ConfigSchema
	(*GetConfigSchemaResponse)(nil), // 5: siemens.connectivitysuite.drvinfo.v1.GetConfigSchemaResponse
	(*GetAppIconRequest)(nil),       // 6: siemens.connectivitysuite.drvinfo.v1.GetAppIconRequest
	(*GetAppIconResponse)(nil),      // 7: siemens.connectivitysuite.drvinfo.v1.GetAppIconResponse
}
var file_conn_suite_drv_info_proto_depIdxs = []int32{
	1, // 0: siemens.connectivitysuite.drvinfo.v1.GetVersionInfoResponse.version:type_name -> siemens.connectivitysuite.drvinfo.v1.VersionInfo
	4, // 1: siemens.connectivitysuite.drvinfo.v1.GetConfigSchemaResponse.schemas:type_name -> siemens.connectivitysuite.drvinfo.v1.ConfigSchema
	0, // 2: siemens.connectivitysuite.drvinfo.v1.DriverInfoApi.GetVersionInfo:input_type -> siemens.connectivitysuite.drvinfo.v1.GetVersionInfoRequest
	3, // 3: siemens.connectivitysuite.drvinfo.v1.DriverInfoApi.GetConfigSchema:input_type -> siemens.connectivitysuite.drvinfo.v1.GetConfigSchemaRequest
	6, // 4: siemens.connectivitysuite.drvinfo.v1.DriverInfoApi.GetAppIcon:input_type -> siemens.connectivitysuite.drvinfo.v1.GetAppIconRequest
	2, // 5: siemens.connectivitysuite.drvinfo.v1.DriverInfoApi.GetVersionInfo:output_type -> siemens.connectivitysuite.drvinfo.v1.GetVersionInfoResponse
	5, // 6: siemens.connectivitysuite.drvinfo.v1.DriverInfoApi.GetConfigSchema:output_type -> siemens.connectivitysuite.drvinfo.v1.GetConfigSchemaResponse
	7, // 7: siemens.connectivitysuite.drvinfo.v1.DriverInfoApi.GetAppIcon:output_type -> siemens.connectivitysuite.drvinfo.v1.GetAppIconResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_conn_suite_drv_info_proto_init() }
func file_conn_suite_drv_info_proto_init() {
	if File_conn_suite_drv_info_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_conn_suite_drv_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersionInfoRequest); i {
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
		file_conn_suite_drv_info_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VersionInfo); i {
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
		file_conn_suite_drv_info_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersionInfoResponse); i {
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
		file_conn_suite_drv_info_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetConfigSchemaRequest); i {
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
		file_conn_suite_drv_info_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigSchema); i {
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
		file_conn_suite_drv_info_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetConfigSchemaResponse); i {
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
		file_conn_suite_drv_info_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAppIconRequest); i {
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
		file_conn_suite_drv_info_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAppIconResponse); i {
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
			RawDescriptor: file_conn_suite_drv_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_conn_suite_drv_info_proto_goTypes,
		DependencyIndexes: file_conn_suite_drv_info_proto_depIdxs,
		MessageInfos:      file_conn_suite_drv_info_proto_msgTypes,
	}.Build()
	File_conn_suite_drv_info_proto = out.File
	file_conn_suite_drv_info_proto_rawDesc = nil
	file_conn_suite_drv_info_proto_goTypes = nil
	file_conn_suite_drv_info_proto_depIdxs = nil
}
