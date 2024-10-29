// ------------------------------------------------------------------
// Connectivity Suite Registry
// ------------------------------------------------------------------
//
// Naming convention according:
// https://cloud.google.com/apis/design/naming_convention
//
// ------------------------------------------------------------------

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v4.25.2
// source: conn_suite_registry.proto

package conn_suite_registry

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

// The key of ServiceInfo is the app_instance_id:
type ServiceInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Mandatory: unique ID of the application
	// Examples: "s7p1", "pnioc1"
	AppInstanceId string `protobuf:"bytes,7,opt,name=app_instance_id,json=appInstanceId,proto3" json:"app_instance_id,omitempty"`
	// Mandatory: Either 'dns_domainname' or 'ipv4_address' needs to be provided.
	// Intentionally there is no 'ipv6_address', because ipv6 without DNS will not really work.
	//
	// Types that are assignable to GrpcIp:
	//
	//	*ServiceInfo_DnsDomainname
	//	*ServiceInfo_Ipv4Address
	GrpcIp isServiceInfo_GrpcIp `protobuf_oneof:"grpc_ip"`
	// Mandatory: port number of gRPC interface
	GrpcIpPortNumber uint32 `protobuf:"varint,3,opt,name=grpc_ip_port_number,json=grpcIpPortNumber,proto3" json:"grpc_ip_port_number,omitempty"`
	// Mandatory: application type
	// +----------------------------------------------------
	// | Examples
	// +-----------------------+----------------------------
	// | "cs-driver"           | Connectivity Suite Driver
	// | "cs-gateway"          | Connectivity Suite Gateway which is configured
	// |                       | by the IIH Configurator
	// | "noncs-driver"        | Driver which supports all interfaces to be configurable
	// |                       | via IIH Configurator, but doesn't support DataApi
	// | "cs-import-converter" | Service which can convert a device-specific configuration
	// |                       | file to a 'Connectivity Suite compatible' configuration file.
	// |                       | This app is e.g. used by the IIH Configurator to let convert
	// |                       | specific config files provided by the user to compatible format
	// | "cs-tagbrowser"       | Service which can browse tags. See the 'BrowsingApi'
	// | "iah-discover"        | The Connector is acting as a Asset Link (AL) and returns IAH
	// |                       | compliant discover results.
	// |                       | The interface siemens.industrialassethub.discover.v1 must be supported.
	//
	// For more info see the Connectivity Suite documentation.
	AppTypes []string `protobuf:"bytes,4,rep,name=app_types,json=appTypes,proto3" json:"app_types,omitempty"`
	// List of all interfaces the app supports.
	// An interface is identified by the package name of the proto file / API service
	// Examples "siemens.connectivitysuite.browsing.v2" or "siemens.connectivitysuite.data.v1"
	// When this list is empty, the service is considered to support all mandatory interfaces
	// of the listed app_types and none of the optional ones. Those are then automatically
	// added by the Registry Service.
	Interfaces []string `protobuf:"bytes,8,rep,name=interfaces,proto3" json:"interfaces,omitempty"`
	// Mandatory: the schema identification(s) of the driver
	DriverSchemaUris []string `protobuf:"bytes,5,rep,name=driver_schema_uris,json=driverSchemaUris,proto3" json:"driver_schema_uris,omitempty"`
}

func (x *ServiceInfo) Reset() {
	*x = ServiceInfo{}
	mi := &file_conn_suite_registry_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServiceInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceInfo) ProtoMessage() {}

func (x *ServiceInfo) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_registry_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceInfo.ProtoReflect.Descriptor instead.
func (*ServiceInfo) Descriptor() ([]byte, []int) {
	return file_conn_suite_registry_proto_rawDescGZIP(), []int{0}
}

func (x *ServiceInfo) GetAppInstanceId() string {
	if x != nil {
		return x.AppInstanceId
	}
	return ""
}

func (m *ServiceInfo) GetGrpcIp() isServiceInfo_GrpcIp {
	if m != nil {
		return m.GrpcIp
	}
	return nil
}

func (x *ServiceInfo) GetDnsDomainname() string {
	if x, ok := x.GetGrpcIp().(*ServiceInfo_DnsDomainname); ok {
		return x.DnsDomainname
	}
	return ""
}

func (x *ServiceInfo) GetIpv4Address() string {
	if x, ok := x.GetGrpcIp().(*ServiceInfo_Ipv4Address); ok {
		return x.Ipv4Address
	}
	return ""
}

func (x *ServiceInfo) GetGrpcIpPortNumber() uint32 {
	if x != nil {
		return x.GrpcIpPortNumber
	}
	return 0
}

func (x *ServiceInfo) GetAppTypes() []string {
	if x != nil {
		return x.AppTypes
	}
	return nil
}

func (x *ServiceInfo) GetInterfaces() []string {
	if x != nil {
		return x.Interfaces
	}
	return nil
}

func (x *ServiceInfo) GetDriverSchemaUris() []string {
	if x != nil {
		return x.DriverSchemaUris
	}
	return nil
}

type isServiceInfo_GrpcIp interface {
	isServiceInfo_GrpcIp()
}

type ServiceInfo_DnsDomainname struct {
	// DNS host name of the driver
	DnsDomainname string `protobuf:"bytes,1,opt,name=dns_domainname,json=dnsDomainname,proto3,oneof"`
}

type ServiceInfo_Ipv4Address struct {
	// ipv4 address of the driver
	Ipv4Address string `protobuf:"bytes,2,opt,name=ipv4_address,json=ipv4Address,proto3,oneof"`
}

func (*ServiceInfo_DnsDomainname) isServiceInfo_GrpcIp() {}

func (*ServiceInfo_Ipv4Address) isServiceInfo_GrpcIp() {}

// RegisterService
type RegisterServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *ServiceInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *RegisterServiceRequest) Reset() {
	*x = RegisterServiceRequest{}
	mi := &file_conn_suite_registry_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterServiceRequest) ProtoMessage() {}

func (x *RegisterServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_registry_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterServiceRequest.ProtoReflect.Descriptor instead.
func (*RegisterServiceRequest) Descriptor() ([]byte, []int) {
	return file_conn_suite_registry_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterServiceRequest) GetInfo() *ServiceInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type RegisterServiceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// maximum time between ServiceInfo refresh in [sec]
	ExpireTime uint32 `protobuf:"varint,2,opt,name=expire_time,json=expireTime,proto3" json:"expire_time,omitempty"`
}

func (x *RegisterServiceResponse) Reset() {
	*x = RegisterServiceResponse{}
	mi := &file_conn_suite_registry_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterServiceResponse) ProtoMessage() {}

func (x *RegisterServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_registry_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterServiceResponse.ProtoReflect.Descriptor instead.
func (*RegisterServiceResponse) Descriptor() ([]byte, []int) {
	return file_conn_suite_registry_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterServiceResponse) GetExpireTime() uint32 {
	if x != nil {
		return x.ExpireTime
	}
	return 0
}

// UnregisterService
// Only the following fields of 'info' are considered, the rest is ignored:
//
//	app_instance_id
type UnregisterServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *ServiceInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *UnregisterServiceRequest) Reset() {
	*x = UnregisterServiceRequest{}
	mi := &file_conn_suite_registry_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnregisterServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterServiceRequest) ProtoMessage() {}

func (x *UnregisterServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_registry_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterServiceRequest.ProtoReflect.Descriptor instead.
func (*UnregisterServiceRequest) Descriptor() ([]byte, []int) {
	return file_conn_suite_registry_proto_rawDescGZIP(), []int{3}
}

func (x *UnregisterServiceRequest) GetInfo() *ServiceInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type UnregisterServiceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UnregisterServiceResponse) Reset() {
	*x = UnregisterServiceResponse{}
	mi := &file_conn_suite_registry_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnregisterServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnregisterServiceResponse) ProtoMessage() {}

func (x *UnregisterServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_registry_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnregisterServiceResponse.ProtoReflect.Descriptor instead.
func (*UnregisterServiceResponse) Descriptor() ([]byte, []int) {
	return file_conn_suite_registry_proto_rawDescGZIP(), []int{4}
}

// QueryRegisteredServices
// Any property which is not set, is ignored
// When a property has multiple values, any entry which matches at least one of the fields matches
// When multiple properties are set, all have to match
//
// Example1: return all apps which have either app_type "cs-driver" or "cs-gateway":
//
//	.app_types = ["cs-driver", "cs-gateway"]
//
// Example2: return all apps which have app_type "cs-import-converter"
//
//	and uri="https://siemens.com/connectivity_suite/schemas/s7plus/1.0.0/config.json":
//	 .app_types = ["cs-import-converter"]
//	 .driver_schema_uris = ["https://siemens.com/connectivity_suite/schemas/s7plus/1.0.0/config.json"]
type ServiceInfoFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// unique ID(s) of the application
	// Examples: "s7p1", "pnioc1"
	AppInstanceIds []string `protobuf:"bytes,2,rep,name=app_instance_ids,json=appInstanceIds,proto3" json:"app_instance_ids,omitempty"`
	// application types
	// Examples "cs-driver" or "cs-gateway"
	// (see 'ServiceInfo.app_types')
	AppTypes []string `protobuf:"bytes,3,rep,name=app_types,json=appTypes,proto3" json:"app_types,omitempty"`
	// the schema identification(s) of the driver
	DriverSchemaUris []string `protobuf:"bytes,4,rep,name=driver_schema_uris,json=driverSchemaUris,proto3" json:"driver_schema_uris,omitempty"`
	// interfaces (actually the package names of the interfaces)
	// Examples "siemens.connectivitysuite.browsing.v2" or "siemens.connectivitysuite.data.v1"
	// (see 'ServiceInfo.interfaces')
	Interfaces []string `protobuf:"bytes,5,rep,name=interfaces,proto3" json:"interfaces,omitempty"`
}

func (x *ServiceInfoFilter) Reset() {
	*x = ServiceInfoFilter{}
	mi := &file_conn_suite_registry_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServiceInfoFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceInfoFilter) ProtoMessage() {}

func (x *ServiceInfoFilter) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_registry_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceInfoFilter.ProtoReflect.Descriptor instead.
func (*ServiceInfoFilter) Descriptor() ([]byte, []int) {
	return file_conn_suite_registry_proto_rawDescGZIP(), []int{5}
}

func (x *ServiceInfoFilter) GetAppInstanceIds() []string {
	if x != nil {
		return x.AppInstanceIds
	}
	return nil
}

func (x *ServiceInfoFilter) GetAppTypes() []string {
	if x != nil {
		return x.AppTypes
	}
	return nil
}

func (x *ServiceInfoFilter) GetDriverSchemaUris() []string {
	if x != nil {
		return x.DriverSchemaUris
	}
	return nil
}

func (x *ServiceInfoFilter) GetInterfaces() []string {
	if x != nil {
		return x.Interfaces
	}
	return nil
}

type QueryRegisteredServicesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filter *ServiceInfoFilter `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *QueryRegisteredServicesRequest) Reset() {
	*x = QueryRegisteredServicesRequest{}
	mi := &file_conn_suite_registry_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryRegisteredServicesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRegisteredServicesRequest) ProtoMessage() {}

func (x *QueryRegisteredServicesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_registry_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRegisteredServicesRequest.ProtoReflect.Descriptor instead.
func (*QueryRegisteredServicesRequest) Descriptor() ([]byte, []int) {
	return file_conn_suite_registry_proto_rawDescGZIP(), []int{6}
}

func (x *QueryRegisteredServicesRequest) GetFilter() *ServiceInfoFilter {
	if x != nil {
		return x.Filter
	}
	return nil
}

type QueryRegisteredServicesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Infos []*ServiceInfo `protobuf:"bytes,1,rep,name=infos,proto3" json:"infos,omitempty"`
}

func (x *QueryRegisteredServicesResponse) Reset() {
	*x = QueryRegisteredServicesResponse{}
	mi := &file_conn_suite_registry_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryRegisteredServicesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRegisteredServicesResponse) ProtoMessage() {}

func (x *QueryRegisteredServicesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_conn_suite_registry_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRegisteredServicesResponse.ProtoReflect.Descriptor instead.
func (*QueryRegisteredServicesResponse) Descriptor() ([]byte, []int) {
	return file_conn_suite_registry_proto_rawDescGZIP(), []int{7}
}

func (x *QueryRegisteredServicesResponse) GetInfos() []*ServiceInfo {
	if x != nil {
		return x.Infos
	}
	return nil
}

var File_conn_suite_registry_proto protoreflect.FileDescriptor

var file_conn_suite_registry_proto_rawDesc = []byte{
	0x0a, 0x19, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x73, 0x75, 0x69, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x25, 0x73, 0x69, 0x65,
	0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74,
	0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e,
	0x76, 0x31, 0x22, 0xa8, 0x02, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x70, 0x70,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0e, 0x64, 0x6e,
	0x73, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x0d, 0x64, 0x6e, 0x73, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x69, 0x70, 0x76, 0x34, 0x5f, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x69, 0x70, 0x76,
	0x34, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2d, 0x0a, 0x13, 0x67, 0x72, 0x70, 0x63,
	0x5f, 0x69, 0x70, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x67, 0x72, 0x70, 0x63, 0x49, 0x70, 0x50, 0x6f, 0x72,
	0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x70, 0x70, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x61, 0x70, 0x70, 0x54,
	0x79, 0x70, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63,
	0x65, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x73, 0x12, 0x2c, 0x0a, 0x12, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x5f, 0x73,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x5f, 0x75, 0x72, 0x69, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x10, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x55, 0x72,
	0x69, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x67, 0x72, 0x70, 0x63, 0x5f, 0x69, 0x70, 0x22, 0x60, 0x0a,
	0x16, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x46, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74,
	0x65, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22,
	0x3a, 0x0a, 0x17, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x62, 0x0a, 0x18, 0x55,
	0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x46, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74,
	0x65, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22,
	0x1b, 0x0a, 0x19, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xa8, 0x01, 0x0a,
	0x11, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x12, 0x28, 0x0a, 0x10, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x70,
	0x70, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x73, 0x12, 0x1b, 0x0a, 0x09,
	0x61, 0x70, 0x70, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x08, 0x61, 0x70, 0x70, 0x54, 0x79, 0x70, 0x65, 0x73, 0x12, 0x2c, 0x0a, 0x12, 0x64, 0x72, 0x69,
	0x76, 0x65, 0x72, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x5f, 0x75, 0x72, 0x69, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x55, 0x72, 0x69, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x66, 0x61, 0x63, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x22, 0x72, 0x0a, 0x1e, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x50, 0x0a, 0x06, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x73, 0x69, 0x65, 0x6d,
	0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22, 0x6b, 0x0a, 0x1f, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48,
	0x0a, 0x05, 0x69, 0x6e, 0x66, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e,
	0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x05, 0x69, 0x6e, 0x66, 0x6f, 0x73, 0x32, 0xea, 0x03, 0x0a, 0x0b, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x72, 0x79, 0x41, 0x70, 0x69, 0x12, 0x92, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d, 0x2e, 0x73,
	0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76,
	0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3e, 0x2e, 0x73, 0x69,
	0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x98, 0x01,
	0x0a, 0x11, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x3f, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x6e, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x40, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65,
	0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x6e, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0xaa, 0x01, 0x0a, 0x17, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x12, 0x45, 0x2e, 0x73, 0x69, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65,
	0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x65, 0x64, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x46, 0x2e, 0x73, 0x69,
	0x65, 0x6d, 0x65, 0x6e, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x79, 0x73, 0x75, 0x69, 0x74, 0x65, 0x2e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x65, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_conn_suite_registry_proto_rawDescOnce sync.Once
	file_conn_suite_registry_proto_rawDescData = file_conn_suite_registry_proto_rawDesc
)

func file_conn_suite_registry_proto_rawDescGZIP() []byte {
	file_conn_suite_registry_proto_rawDescOnce.Do(func() {
		file_conn_suite_registry_proto_rawDescData = protoimpl.X.CompressGZIP(file_conn_suite_registry_proto_rawDescData)
	})
	return file_conn_suite_registry_proto_rawDescData
}

var file_conn_suite_registry_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_conn_suite_registry_proto_goTypes = []any{
	(*ServiceInfo)(nil),                     // 0: siemens.connectivitysuite.registry.v1.ServiceInfo
	(*RegisterServiceRequest)(nil),          // 1: siemens.connectivitysuite.registry.v1.RegisterServiceRequest
	(*RegisterServiceResponse)(nil),         // 2: siemens.connectivitysuite.registry.v1.RegisterServiceResponse
	(*UnregisterServiceRequest)(nil),        // 3: siemens.connectivitysuite.registry.v1.UnregisterServiceRequest
	(*UnregisterServiceResponse)(nil),       // 4: siemens.connectivitysuite.registry.v1.UnregisterServiceResponse
	(*ServiceInfoFilter)(nil),               // 5: siemens.connectivitysuite.registry.v1.ServiceInfoFilter
	(*QueryRegisteredServicesRequest)(nil),  // 6: siemens.connectivitysuite.registry.v1.QueryRegisteredServicesRequest
	(*QueryRegisteredServicesResponse)(nil), // 7: siemens.connectivitysuite.registry.v1.QueryRegisteredServicesResponse
}
var file_conn_suite_registry_proto_depIdxs = []int32{
	0, // 0: siemens.connectivitysuite.registry.v1.RegisterServiceRequest.info:type_name -> siemens.connectivitysuite.registry.v1.ServiceInfo
	0, // 1: siemens.connectivitysuite.registry.v1.UnregisterServiceRequest.info:type_name -> siemens.connectivitysuite.registry.v1.ServiceInfo
	5, // 2: siemens.connectivitysuite.registry.v1.QueryRegisteredServicesRequest.filter:type_name -> siemens.connectivitysuite.registry.v1.ServiceInfoFilter
	0, // 3: siemens.connectivitysuite.registry.v1.QueryRegisteredServicesResponse.infos:type_name -> siemens.connectivitysuite.registry.v1.ServiceInfo
	1, // 4: siemens.connectivitysuite.registry.v1.RegistryApi.RegisterService:input_type -> siemens.connectivitysuite.registry.v1.RegisterServiceRequest
	3, // 5: siemens.connectivitysuite.registry.v1.RegistryApi.UnregisterService:input_type -> siemens.connectivitysuite.registry.v1.UnregisterServiceRequest
	6, // 6: siemens.connectivitysuite.registry.v1.RegistryApi.QueryRegisteredServices:input_type -> siemens.connectivitysuite.registry.v1.QueryRegisteredServicesRequest
	2, // 7: siemens.connectivitysuite.registry.v1.RegistryApi.RegisterService:output_type -> siemens.connectivitysuite.registry.v1.RegisterServiceResponse
	4, // 8: siemens.connectivitysuite.registry.v1.RegistryApi.UnregisterService:output_type -> siemens.connectivitysuite.registry.v1.UnregisterServiceResponse
	7, // 9: siemens.connectivitysuite.registry.v1.RegistryApi.QueryRegisteredServices:output_type -> siemens.connectivitysuite.registry.v1.QueryRegisteredServicesResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_conn_suite_registry_proto_init() }
func file_conn_suite_registry_proto_init() {
	if File_conn_suite_registry_proto != nil {
		return
	}
	file_conn_suite_registry_proto_msgTypes[0].OneofWrappers = []any{
		(*ServiceInfo_DnsDomainname)(nil),
		(*ServiceInfo_Ipv4Address)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_conn_suite_registry_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_conn_suite_registry_proto_goTypes,
		DependencyIndexes: file_conn_suite_registry_proto_depIdxs,
		MessageInfos:      file_conn_suite_registry_proto_msgTypes,
	}.Build()
	File_conn_suite_registry_proto = out.File
	file_conn_suite_registry_proto_rawDesc = nil
	file_conn_suite_registry_proto_goTypes = nil
	file_conn_suite_registry_proto_depIdxs = nil
}
