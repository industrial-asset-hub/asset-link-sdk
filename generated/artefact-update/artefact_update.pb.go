// Artefact update interface
// This is the interface for pushing and pulling
// artefacts to and from drivers.
// The driver is responsible for the actual
// transfer of the artefact to the target device.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.3
// source: artefact_update.proto

package artefact_update

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

// Prefix with AS == Artefact status
// to be globally unique
type TransferStatus int32

const (
	TransferStatus_AS_OK                  TransferStatus = 0
	TransferStatus_AS_FAIL                TransferStatus = 1
	TransferStatus_AS_IDENTITY_CHECK_FAIL TransferStatus = 2
	TransferStatus_AS_INCOMPATIBLE        TransferStatus = 3
)

// Enum value maps for TransferStatus.
var (
	TransferStatus_name = map[int32]string{
		0: "AS_OK",
		1: "AS_FAIL",
		2: "AS_IDENTITY_CHECK_FAIL",
		3: "AS_INCOMPATIBLE",
	}
	TransferStatus_value = map[string]int32{
		"AS_OK":                  0,
		"AS_FAIL":                1,
		"AS_IDENTITY_CHECK_FAIL": 2,
		"AS_INCOMPATIBLE":        3,
	}
)

func (x TransferStatus) Enum() *TransferStatus {
	p := new(TransferStatus)
	*p = x
	return p
}

func (x TransferStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TransferStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_artefact_update_proto_enumTypes[0].Descriptor()
}

func (TransferStatus) Type() protoreflect.EnumType {
	return &file_artefact_update_proto_enumTypes[0]
}

func (x TransferStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TransferStatus.Descriptor instead.
func (TransferStatus) EnumDescriptor() ([]byte, []int) {
	return file_artefact_update_proto_rawDescGZIP(), []int{0}
}

type ArtefactType int32

const (
	ArtefactType_AT_FIRMWARE      ArtefactType = 0
	ArtefactType_AT_BACKUP        ArtefactType = 1
	ArtefactType_AT_CONFIGURATION ArtefactType = 2
)

// Enum value maps for ArtefactType.
var (
	ArtefactType_name = map[int32]string{
		0: "AT_FIRMWARE",
		1: "AT_BACKUP",
		2: "AT_CONFIGURATION",
	}
	ArtefactType_value = map[string]int32{
		"AT_FIRMWARE":      0,
		"AT_BACKUP":        1,
		"AT_CONFIGURATION": 2,
	}
)

func (x ArtefactType) Enum() *ArtefactType {
	p := new(ArtefactType)
	*p = x
	return p
}

func (x ArtefactType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ArtefactType) Descriptor() protoreflect.EnumDescriptor {
	return file_artefact_update_proto_enumTypes[1].Descriptor()
}

func (ArtefactType) Type() protoreflect.EnumType {
	return &file_artefact_update_proto_enumTypes[1]
}

func (x ArtefactType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ArtefactType.Descriptor instead.
func (ArtefactType) EnumDescriptor() ([]byte, []int) {
	return file_artefact_update_proto_rawDescGZIP(), []int{1}
}

type ArtefactUpdateState int32

const (
	ArtefactUpdateState_AUS_IDLE         ArtefactUpdateState = 0
	ArtefactUpdateState_AUS_DOWNLOAD     ArtefactUpdateState = 1
	ArtefactUpdateState_AUS_INSTALLATION ArtefactUpdateState = 2
	ArtefactUpdateState_AUS_ACTIVATION   ArtefactUpdateState = 3
)

// Enum value maps for ArtefactUpdateState.
var (
	ArtefactUpdateState_name = map[int32]string{
		0: "AUS_IDLE",
		1: "AUS_DOWNLOAD",
		2: "AUS_INSTALLATION",
		3: "AUS_ACTIVATION",
	}
	ArtefactUpdateState_value = map[string]int32{
		"AUS_IDLE":         0,
		"AUS_DOWNLOAD":     1,
		"AUS_INSTALLATION": 2,
		"AUS_ACTIVATION":   3,
	}
)

func (x ArtefactUpdateState) Enum() *ArtefactUpdateState {
	p := new(ArtefactUpdateState)
	*p = x
	return p
}

func (x ArtefactUpdateState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ArtefactUpdateState) Descriptor() protoreflect.EnumDescriptor {
	return file_artefact_update_proto_enumTypes[2].Descriptor()
}

func (ArtefactUpdateState) Type() protoreflect.EnumType {
	return &file_artefact_update_proto_enumTypes[2]
}

func (x ArtefactUpdateState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ArtefactUpdateState.Descriptor instead.
func (ArtefactUpdateState) EnumDescriptor() ([]byte, []int) {
	return file_artefact_update_proto_rawDescGZIP(), []int{2}
}

type ArtefactChunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*ArtefactChunk_Metadata
	//	*ArtefactChunk_FileContent
	//	*ArtefactChunk_Status
	Data isArtefactChunk_Data `protobuf_oneof:"data"`
}

func (x *ArtefactChunk) Reset() {
	*x = ArtefactChunk{}
	mi := &file_artefact_update_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ArtefactChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArtefactChunk) ProtoMessage() {}

func (x *ArtefactChunk) ProtoReflect() protoreflect.Message {
	mi := &file_artefact_update_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArtefactChunk.ProtoReflect.Descriptor instead.
func (*ArtefactChunk) Descriptor() ([]byte, []int) {
	return file_artefact_update_proto_rawDescGZIP(), []int{0}
}

func (m *ArtefactChunk) GetData() isArtefactChunk_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *ArtefactChunk) GetMetadata() *ArtefactMetaData {
	if x, ok := x.GetData().(*ArtefactChunk_Metadata); ok {
		return x.Metadata
	}
	return nil
}

func (x *ArtefactChunk) GetFileContent() []byte {
	if x, ok := x.GetData().(*ArtefactChunk_FileContent); ok {
		return x.FileContent
	}
	return nil
}

func (x *ArtefactChunk) GetStatus() *Status {
	if x, ok := x.GetData().(*ArtefactChunk_Status); ok {
		return x.Status
	}
	return nil
}

type isArtefactChunk_Data interface {
	isArtefactChunk_Data()
}

type ArtefactChunk_Metadata struct {
	Metadata *ArtefactMetaData `protobuf:"bytes,1,opt,name=metadata,proto3,oneof"` // first packet
}

type ArtefactChunk_FileContent struct {
	FileContent []byte `protobuf:"bytes,2,opt,name=file_content,json=fileContent,proto3,oneof"` // following
}

type ArtefactChunk_Status struct {
	Status *Status `protobuf:"bytes,3,opt,name=status,proto3,oneof"` // only for pull
}

func (*ArtefactChunk_Metadata) isArtefactChunk_Data() {}

func (*ArtefactChunk_FileContent) isArtefactChunk_Data() {}

func (*ArtefactChunk_Status) isArtefactChunk_Data() {}

type ArtefactMetaData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// oneof meta {
	Credential         *ArtefactCredentials `protobuf:"bytes,1,opt,name=credential,proto3" json:"credential,omitempty"`
	DeviceIdentifier   []byte               `protobuf:"bytes,2,opt,name=device_identifier,json=deviceIdentifier,proto3" json:"device_identifier,omitempty"`
	ArtefactIdentifier *ArtefactIdentifier  `protobuf:"bytes,3,opt,name=artefact_identifier,json=artefactIdentifier,proto3" json:"artefact_identifier,omitempty"` // }
}

func (x *ArtefactMetaData) Reset() {
	*x = ArtefactMetaData{}
	mi := &file_artefact_update_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ArtefactMetaData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArtefactMetaData) ProtoMessage() {}

func (x *ArtefactMetaData) ProtoReflect() protoreflect.Message {
	mi := &file_artefact_update_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArtefactMetaData.ProtoReflect.Descriptor instead.
func (*ArtefactMetaData) Descriptor() ([]byte, []int) {
	return file_artefact_update_proto_rawDescGZIP(), []int{1}
}

func (x *ArtefactMetaData) GetCredential() *ArtefactCredentials {
	if x != nil {
		return x.Credential
	}
	return nil
}

func (x *ArtefactMetaData) GetDeviceIdentifier() []byte {
	if x != nil {
		return x.DeviceIdentifier
	}
	return nil
}

func (x *ArtefactMetaData) GetArtefactIdentifier() *ArtefactIdentifier {
	if x != nil {
		return x.ArtefactIdentifier
	}
	return nil
}

type ArtefactCredentials struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CredentialType string `protobuf:"bytes,1,opt,name=credential_type,json=credentialType,proto3" json:"credential_type,omitempty"`
	Credentials    []byte `protobuf:"bytes,2,opt,name=credentials,proto3" json:"credentials,omitempty"`
}

func (x *ArtefactCredentials) Reset() {
	*x = ArtefactCredentials{}
	mi := &file_artefact_update_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ArtefactCredentials) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArtefactCredentials) ProtoMessage() {}

func (x *ArtefactCredentials) ProtoReflect() protoreflect.Message {
	mi := &file_artefact_update_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArtefactCredentials.ProtoReflect.Descriptor instead.
func (*ArtefactCredentials) Descriptor() ([]byte, []int) {
	return file_artefact_update_proto_rawDescGZIP(), []int{2}
}

func (x *ArtefactCredentials) GetCredentialType() string {
	if x != nil {
		return x.CredentialType
	}
	return ""
}

func (x *ArtefactCredentials) GetCredentials() []byte {
	if x != nil {
		return x.Credentials
	}
	return nil
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  TransferStatus `protobuf:"varint,1,opt,name=status,proto3,enum=factory_x.artefact_update.v1.TransferStatus" json:"status,omitempty"`
	Message string         `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	mi := &file_artefact_update_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_artefact_update_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_artefact_update_proto_rawDescGZIP(), []int{3}
}

func (x *Status) GetStatus() TransferStatus {
	if x != nil {
		return x.Status
	}
	return TransferStatus_AS_OK
}

func (x *Status) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ArtefactIdentifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type ArtefactType `protobuf:"varint,1,opt,name=type,proto3,enum=factory_x.artefact_update.v1.ArtefactType" json:"type,omitempty"`
}

func (x *ArtefactIdentifier) Reset() {
	*x = ArtefactIdentifier{}
	mi := &file_artefact_update_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ArtefactIdentifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArtefactIdentifier) ProtoMessage() {}

func (x *ArtefactIdentifier) ProtoReflect() protoreflect.Message {
	mi := &file_artefact_update_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArtefactIdentifier.ProtoReflect.Descriptor instead.
func (*ArtefactIdentifier) Descriptor() ([]byte, []int) {
	return file_artefact_update_proto_rawDescGZIP(), []int{4}
}

func (x *ArtefactIdentifier) GetType() ArtefactType {
	if x != nil {
		return x.Type
	}
	return ArtefactType_AT_FIRMWARE
}

type ArtefactUpdateStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status   *Status             `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	State    ArtefactUpdateState `protobuf:"varint,2,opt,name=state,proto3,enum=factory_x.artefact_update.v1.ArtefactUpdateState" json:"state,omitempty"`
	Progress int32               `protobuf:"varint,3,opt,name=progress,proto3" json:"progress,omitempty"`
}

func (x *ArtefactUpdateStatus) Reset() {
	*x = ArtefactUpdateStatus{}
	mi := &file_artefact_update_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ArtefactUpdateStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArtefactUpdateStatus) ProtoMessage() {}

func (x *ArtefactUpdateStatus) ProtoReflect() protoreflect.Message {
	mi := &file_artefact_update_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArtefactUpdateStatus.ProtoReflect.Descriptor instead.
func (*ArtefactUpdateStatus) Descriptor() ([]byte, []int) {
	return file_artefact_update_proto_rawDescGZIP(), []int{5}
}

func (x *ArtefactUpdateStatus) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *ArtefactUpdateStatus) GetState() ArtefactUpdateState {
	if x != nil {
		return x.State
	}
	return ArtefactUpdateState_AUS_IDLE
}

func (x *ArtefactUpdateStatus) GetProgress() int32 {
	if x != nil {
		return x.Progress
	}
	return 0
}

var File_artefact_update_proto protoreflect.FileDescriptor

var file_artefact_update_proto_rawDesc = []byte{
	0x0a, 0x15, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x76, 0x31, 0x22, 0xca, 0x01, 0x0a, 0x0d, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61,
	0x63, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x4c, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x66, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x79, 0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63,
	0x74, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x08, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x23, 0x0a, 0x0c, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x0b, 0x66,
	0x69, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x3e, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x66, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x79, 0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0xf5, 0x01, 0x0a, 0x10, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x4d,
	0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x51, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x66, 0x61,
	0x63, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74,
	0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x74, 0x65, 0x66,
	0x61, 0x63, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x52, 0x0a,
	0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x12, 0x2b, 0x0a, 0x11, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x61, 0x0a, 0x13, 0x61, 0x72, 0x74, 0x65, 0x66,
	0x61, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x78,
	0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x49, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x12, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74,
	0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x22, 0x60, 0x0a, 0x13, 0x41, 0x72,
	0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x73, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x72,
	0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0b, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x22, 0x68, 0x0a, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x44, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2c, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79,
	0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x54, 0x0a, 0x12, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61,
	0x63, 0x74, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x3e, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2a, 0x2e, 0x66, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x79, 0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61,
	0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0xb9, 0x01, 0x0a,
	0x14, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x3c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x5f,
	0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x47, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x31, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x78, 0x2e, 0x61,
	0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x2a, 0x59, 0x0a, 0x0e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x53,
	0x5f, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x41, 0x53, 0x5f, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x49, 0x54,
	0x59, 0x5f, 0x43, 0x48, 0x45, 0x43, 0x4b, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x10, 0x02, 0x12, 0x13,
	0x0a, 0x0f, 0x41, 0x53, 0x5f, 0x49, 0x4e, 0x43, 0x4f, 0x4d, 0x50, 0x41, 0x54, 0x49, 0x42, 0x4c,
	0x45, 0x10, 0x03, 0x2a, 0x44, 0x0a, 0x0c, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x41, 0x54, 0x5f, 0x46, 0x49, 0x52, 0x4d, 0x57, 0x41,
	0x52, 0x45, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x41, 0x54, 0x5f, 0x42, 0x41, 0x43, 0x4b, 0x55,
	0x50, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x41, 0x54, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47,
	0x55, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02, 0x2a, 0x5f, 0x0a, 0x13, 0x41, 0x72, 0x74,
	0x65, 0x66, 0x61, 0x63, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x0c, 0x0a, 0x08, 0x41, 0x55, 0x53, 0x5f, 0x49, 0x44, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x10,
	0x0a, 0x0c, 0x41, 0x55, 0x53, 0x5f, 0x44, 0x4f, 0x57, 0x4e, 0x4c, 0x4f, 0x41, 0x44, 0x10, 0x01,
	0x12, 0x14, 0x0a, 0x10, 0x41, 0x55, 0x53, 0x5f, 0x49, 0x4e, 0x53, 0x54, 0x41, 0x4c, 0x4c, 0x41,
	0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x41, 0x55, 0x53, 0x5f, 0x41, 0x43,
	0x54, 0x49, 0x56, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x03, 0x32, 0xfb, 0x01, 0x0a, 0x11, 0x41,
	0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x70, 0x69,
	0x12, 0x75, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74,
	0x12, 0x2b, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74,
	0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x1a, 0x32, 0x2e,
	0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61,
	0x63, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x74,
	0x65, 0x66, 0x61, 0x63, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x6f, 0x0a, 0x0c, 0x50, 0x75, 0x6c, 0x6c, 0x41,
	0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x12, 0x2e, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x79, 0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x4d,
	0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x2b, 0x2e, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72,
	0x79, 0x5f, 0x78, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x74, 0x65, 0x66, 0x61, 0x63, 0x74, 0x43,
	0x68, 0x75, 0x6e, 0x6b, 0x22, 0x00, 0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_artefact_update_proto_rawDescOnce sync.Once
	file_artefact_update_proto_rawDescData = file_artefact_update_proto_rawDesc
)

func file_artefact_update_proto_rawDescGZIP() []byte {
	file_artefact_update_proto_rawDescOnce.Do(func() {
		file_artefact_update_proto_rawDescData = protoimpl.X.CompressGZIP(file_artefact_update_proto_rawDescData)
	})
	return file_artefact_update_proto_rawDescData
}

var file_artefact_update_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_artefact_update_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_artefact_update_proto_goTypes = []any{
	(TransferStatus)(0),          // 0: factory_x.artefact_update.v1.TransferStatus
	(ArtefactType)(0),            // 1: factory_x.artefact_update.v1.ArtefactType
	(ArtefactUpdateState)(0),     // 2: factory_x.artefact_update.v1.ArtefactUpdateState
	(*ArtefactChunk)(nil),        // 3: factory_x.artefact_update.v1.ArtefactChunk
	(*ArtefactMetaData)(nil),     // 4: factory_x.artefact_update.v1.ArtefactMetaData
	(*ArtefactCredentials)(nil),  // 5: factory_x.artefact_update.v1.ArtefactCredentials
	(*Status)(nil),               // 6: factory_x.artefact_update.v1.Status
	(*ArtefactIdentifier)(nil),   // 7: factory_x.artefact_update.v1.ArtefactIdentifier
	(*ArtefactUpdateStatus)(nil), // 8: factory_x.artefact_update.v1.ArtefactUpdateStatus
}
var file_artefact_update_proto_depIdxs = []int32{
	4,  // 0: factory_x.artefact_update.v1.ArtefactChunk.metadata:type_name -> factory_x.artefact_update.v1.ArtefactMetaData
	6,  // 1: factory_x.artefact_update.v1.ArtefactChunk.status:type_name -> factory_x.artefact_update.v1.Status
	5,  // 2: factory_x.artefact_update.v1.ArtefactMetaData.credential:type_name -> factory_x.artefact_update.v1.ArtefactCredentials
	7,  // 3: factory_x.artefact_update.v1.ArtefactMetaData.artefact_identifier:type_name -> factory_x.artefact_update.v1.ArtefactIdentifier
	0,  // 4: factory_x.artefact_update.v1.Status.status:type_name -> factory_x.artefact_update.v1.TransferStatus
	1,  // 5: factory_x.artefact_update.v1.ArtefactIdentifier.type:type_name -> factory_x.artefact_update.v1.ArtefactType
	6,  // 6: factory_x.artefact_update.v1.ArtefactUpdateStatus.status:type_name -> factory_x.artefact_update.v1.Status
	2,  // 7: factory_x.artefact_update.v1.ArtefactUpdateStatus.state:type_name -> factory_x.artefact_update.v1.ArtefactUpdateState
	3,  // 8: factory_x.artefact_update.v1.ArtefactUpdateApi.PushArtefact:input_type -> factory_x.artefact_update.v1.ArtefactChunk
	4,  // 9: factory_x.artefact_update.v1.ArtefactUpdateApi.PullArtefact:input_type -> factory_x.artefact_update.v1.ArtefactMetaData
	8,  // 10: factory_x.artefact_update.v1.ArtefactUpdateApi.PushArtefact:output_type -> factory_x.artefact_update.v1.ArtefactUpdateStatus
	3,  // 11: factory_x.artefact_update.v1.ArtefactUpdateApi.PullArtefact:output_type -> factory_x.artefact_update.v1.ArtefactChunk
	10, // [10:12] is the sub-list for method output_type
	8,  // [8:10] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_artefact_update_proto_init() }
func file_artefact_update_proto_init() {
	if File_artefact_update_proto != nil {
		return
	}
	file_artefact_update_proto_msgTypes[0].OneofWrappers = []any{
		(*ArtefactChunk_Metadata)(nil),
		(*ArtefactChunk_FileContent)(nil),
		(*ArtefactChunk_Status)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_artefact_update_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_artefact_update_proto_goTypes,
		DependencyIndexes: file_artefact_update_proto_depIdxs,
		EnumInfos:         file_artefact_update_proto_enumTypes,
		MessageInfos:      file_artefact_update_proto_msgTypes,
	}.Build()
	File_artefact_update_proto = out.File
	file_artefact_update_proto_rawDesc = nil
	file_artefact_update_proto_goTypes = nil
	file_artefact_update_proto_depIdxs = nil
}
