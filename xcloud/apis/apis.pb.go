// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.15.1
// source: apis.proto

package apis

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ServiceStatus int32

const (
	ServiceStatus_ServiceStatusNone    ServiceStatus = 0
	ServiceStatus_ServiceStatusRunning ServiceStatus = 1
	ServiceStatus_ServiceStatusSuspend ServiceStatus = 2
)

// Enum value maps for ServiceStatus.
var (
	ServiceStatus_name = map[int32]string{
		0: "ServiceStatusNone",
		1: "ServiceStatusRunning",
		2: "ServiceStatusSuspend",
	}
	ServiceStatus_value = map[string]int32{
		"ServiceStatusNone":    0,
		"ServiceStatusRunning": 1,
		"ServiceStatusSuspend": 2,
	}
)

func (x ServiceStatus) Enum() *ServiceStatus {
	p := new(ServiceStatus)
	*p = x
	return p
}

func (x ServiceStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServiceStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_apis_proto_enumTypes[0].Descriptor()
}

func (ServiceStatus) Type() protoreflect.EnumType {
	return &file_apis_proto_enumTypes[0]
}

func (x ServiceStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ServiceStatus.Descriptor instead.
func (ServiceStatus) EnumDescriptor() ([]byte, []int) {
	return file_apis_proto_rawDescGZIP(), []int{0}
}

type Peer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sid   string `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty"`
	Appid string `protobuf:"bytes,2,opt,name=appid,proto3" json:"appid,omitempty"`
	// below invalid for user session
	Openid     string `protobuf:"bytes,3,opt,name=openid,proto3" json:"openid,omitempty"`
	Network    string `protobuf:"bytes,4,opt,name=network,proto3" json:"network,omitempty"`
	Addr       string `protobuf:"bytes,5,opt,name=addr,proto3" json:"addr,omitempty"`
	Authorized bool   `protobuf:"varint,6,opt,name=authorized,proto3" json:"authorized,omitempty"`
}

func (x *Peer) Reset() {
	*x = Peer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Peer) ProtoMessage() {}

func (x *Peer) ProtoReflect() protoreflect.Message {
	mi := &file_apis_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Peer.ProtoReflect.Descriptor instead.
func (*Peer) Descriptor() ([]byte, []int) {
	return file_apis_proto_rawDescGZIP(), []int{0}
}

func (x *Peer) GetSid() string {
	if x != nil {
		return x.Sid
	}
	return ""
}

func (x *Peer) GetAppid() string {
	if x != nil {
		return x.Appid
	}
	return ""
}

func (x *Peer) GetOpenid() string {
	if x != nil {
		return x.Openid
	}
	return ""
}

func (x *Peer) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Peer) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Peer) GetAuthorized() bool {
	if x != nil {
		return x.Authorized
	}
	return false
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Api     string            `protobuf:"bytes,1,opt,name=api,proto3" json:"api,omitempty"`
	Version string            `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Source  *Peer             `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"` // source and target once will used only one
	Target  *Peer             `protobuf:"bytes,4,opt,name=target,proto3" json:"target,omitempty"`
	Headers map[string]string `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body    []byte            `protobuf:"bytes,6,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_apis_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_apis_proto_rawDescGZIP(), []int{1}
}

func (x *Request) GetApi() string {
	if x != nil {
		return x.Api
	}
	return ""
}

func (x *Request) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Request) GetSource() *Peer {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *Request) GetTarget() *Peer {
	if x != nil {
		return x.Target
	}
	return nil
}

func (x *Request) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *Request) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body []byte `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_apis_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_apis_proto_rawDescGZIP(), []int{2}
}

func (x *Response) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type ServiceRegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Rid    int32  `protobuf:"varint,2,opt,name=rid,proto3" json:"rid,omitempty"`
	Token  string `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	Config string `protobuf:"bytes,4,opt,name=config,proto3" json:"config,omitempty"`
}

func (x *ServiceRegisterRequest) Reset() {
	*x = ServiceRegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceRegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceRegisterRequest) ProtoMessage() {}

func (x *ServiceRegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apis_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceRegisterRequest.ProtoReflect.Descriptor instead.
func (*ServiceRegisterRequest) Descriptor() ([]byte, []int) {
	return file_apis_proto_rawDescGZIP(), []int{3}
}

func (x *ServiceRegisterRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ServiceRegisterRequest) GetRid() int32 {
	if x != nil {
		return x.Rid
	}
	return 0
}

func (x *ServiceRegisterRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ServiceRegisterRequest) GetConfig() string {
	if x != nil {
		return x.Config
	}
	return ""
}

type ServiceRegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rid int32 `protobuf:"varint,1,opt,name=rid,proto3" json:"rid,omitempty"`
}

func (x *ServiceRegisterResponse) Reset() {
	*x = ServiceRegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceRegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceRegisterResponse) ProtoMessage() {}

func (x *ServiceRegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_apis_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceRegisterResponse.ProtoReflect.Descriptor instead.
func (*ServiceRegisterResponse) Descriptor() ([]byte, []int) {
	return file_apis_proto_rawDescGZIP(), []int{4}
}

func (x *ServiceRegisterResponse) GetRid() int32 {
	if x != nil {
		return x.Rid
	}
	return 0
}

type ServiceSuspendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status uint32 `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ServiceSuspendRequest) Reset() {
	*x = ServiceSuspendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceSuspendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceSuspendRequest) ProtoMessage() {}

func (x *ServiceSuspendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apis_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceSuspendRequest.ProtoReflect.Descriptor instead.
func (*ServiceSuspendRequest) Descriptor() ([]byte, []int) {
	return file_apis_proto_rawDescGZIP(), []int{5}
}

func (x *ServiceSuspendRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ServiceSuspendRequest) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type ServiceResumeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status uint32 `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ServiceResumeRequest) Reset() {
	*x = ServiceResumeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceResumeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceResumeRequest) ProtoMessage() {}

func (x *ServiceResumeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apis_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceResumeRequest.ProtoReflect.Descriptor instead.
func (*ServiceResumeRequest) Descriptor() ([]byte, []int) {
	return file_apis_proto_rawDescGZIP(), []int{6}
}

func (x *ServiceResumeRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ServiceResumeRequest) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type ServiceStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status uint32 `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ServiceStateRequest) Reset() {
	*x = ServiceStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_apis_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceStateRequest) ProtoMessage() {}

func (x *ServiceStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_apis_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceStateRequest.ProtoReflect.Descriptor instead.
func (*ServiceStateRequest) Descriptor() ([]byte, []int) {
	return file_apis_proto_rawDescGZIP(), []int{7}
}

func (x *ServiceStateRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ServiceStateRequest) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

var File_apis_proto protoreflect.FileDescriptor

var file_apis_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x78, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x61, 0x70, 0x69, 0x73, 0x22, 0x94, 0x01, 0x0a, 0x04, 0x50, 0x65, 0x65,
	0x72, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x73, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x70, 0x65,
	0x6e, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x70, 0x65, 0x6e, 0x69,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x61,
	0x64, 0x64, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12,
	0x1e, 0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x22,
	0x95, 0x02, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x61,
	0x70, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x70, 0x69, 0x12, 0x18, 0x0a,
	0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x78, 0x63, 0x6c, 0x6f, 0x75, 0x64,
	0x61, 0x70, 0x69, 0x73, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x28, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x78, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x50,
	0x65, 0x65, 0x72, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x3a, 0x0a, 0x07, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x78,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x1a, 0x3a, 0x0a, 0x0c, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x1e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x68, 0x0a, 0x16, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x72, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x22, 0x2b, 0x0a, 0x17, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x72, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x72, 0x69, 0x64, 0x22, 0x3f,
	0x0a, 0x15, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x75, 0x73, 0x70, 0x65, 0x6e, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x3e, 0x0a, 0x14, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6d, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22,
	0x3d, 0x0a, 0x13, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2a, 0x5a,
	0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x15, 0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x10, 0x01,
	0x12, 0x18, 0x0a, 0x14, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x53, 0x75, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x10, 0x02, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x73, 0x6f, 0x6e, 0x2d, 0x78,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x78, 0x64, 0x70, 0x2d, 0x64, 0x65, 0x6d, 0x6f, 0x2f, 0x78,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_apis_proto_rawDescOnce sync.Once
	file_apis_proto_rawDescData = file_apis_proto_rawDesc
)

func file_apis_proto_rawDescGZIP() []byte {
	file_apis_proto_rawDescOnce.Do(func() {
		file_apis_proto_rawDescData = protoimpl.X.CompressGZIP(file_apis_proto_rawDescData)
	})
	return file_apis_proto_rawDescData
}

var file_apis_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_apis_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_apis_proto_goTypes = []interface{}{
	(ServiceStatus)(0),              // 0: xcloudapis.ServiceStatus
	(*Peer)(nil),                    // 1: xcloudapis.Peer
	(*Request)(nil),                 // 2: xcloudapis.Request
	(*Response)(nil),                // 3: xcloudapis.Response
	(*ServiceRegisterRequest)(nil),  // 4: xcloudapis.ServiceRegisterRequest
	(*ServiceRegisterResponse)(nil), // 5: xcloudapis.ServiceRegisterResponse
	(*ServiceSuspendRequest)(nil),   // 6: xcloudapis.ServiceSuspendRequest
	(*ServiceResumeRequest)(nil),    // 7: xcloudapis.ServiceResumeRequest
	(*ServiceStateRequest)(nil),     // 8: xcloudapis.ServiceStateRequest
	nil,                             // 9: xcloudapis.Request.HeadersEntry
}
var file_apis_proto_depIdxs = []int32{
	1, // 0: xcloudapis.Request.source:type_name -> xcloudapis.Peer
	1, // 1: xcloudapis.Request.target:type_name -> xcloudapis.Peer
	9, // 2: xcloudapis.Request.headers:type_name -> xcloudapis.Request.HeadersEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_apis_proto_init() }
func file_apis_proto_init() {
	if File_apis_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_apis_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Peer); i {
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
		file_apis_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_apis_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_apis_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceRegisterRequest); i {
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
		file_apis_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceRegisterResponse); i {
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
		file_apis_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceSuspendRequest); i {
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
		file_apis_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceResumeRequest); i {
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
		file_apis_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceStateRequest); i {
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
			RawDescriptor: file_apis_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_apis_proto_goTypes,
		DependencyIndexes: file_apis_proto_depIdxs,
		EnumInfos:         file_apis_proto_enumTypes,
		MessageInfos:      file_apis_proto_msgTypes,
	}.Build()
	File_apis_proto = out.File
	file_apis_proto_rawDesc = nil
	file_apis_proto_goTypes = nil
	file_apis_proto_depIdxs = nil
}