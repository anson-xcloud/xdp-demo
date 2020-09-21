// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Cmd int32

const (
	Cmd_CmdNone             Cmd = 0
	Cmd_CmdHandshake        Cmd = 1
	Cmd_CmdSessionOnConnect Cmd = 10
	Cmd_CmdSessionOnRecv    Cmd = 11
	Cmd_CmdSessionOnClose   Cmd = 12
	Cmd_CmdSessionSend      Cmd = 20
	Cmd_CmdSessionMultiSend Cmd = 21
	Cmd_CmdSessionClose     Cmd = 22
)

var Cmd_name = map[int32]string{
	0:  "CmdNone",
	1:  "CmdHandshake",
	10: "CmdSessionOnConnect",
	11: "CmdSessionOnRecv",
	12: "CmdSessionOnClose",
	20: "CmdSessionSend",
	21: "CmdSessionMultiSend",
	22: "CmdSessionClose",
}

var Cmd_value = map[string]int32{
	"CmdNone":             0,
	"CmdHandshake":        1,
	"CmdSessionOnConnect": 10,
	"CmdSessionOnRecv":    11,
	"CmdSessionOnClose":   12,
	"CmdSessionSend":      20,
	"CmdSessionMultiSend": 21,
	"CmdSessionClose":     22,
}

func (x Cmd) String() string {
	return proto.EnumName(Cmd_name, int32(x))
}

func (Cmd) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

type HandshakeRequest struct {
	AppID                string   `protobuf:"bytes,1,opt,name=appID,proto3" json:"appID,omitempty"`
	AccessKey            string   `protobuf:"bytes,2,opt,name=accessKey,proto3" json:"accessKey,omitempty"`
	Config               string   `protobuf:"bytes,3,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HandshakeRequest) Reset()         { *m = HandshakeRequest{} }
func (m *HandshakeRequest) String() string { return proto.CompactTextString(m) }
func (*HandshakeRequest) ProtoMessage()    {}
func (*HandshakeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *HandshakeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HandshakeRequest.Unmarshal(m, b)
}
func (m *HandshakeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HandshakeRequest.Marshal(b, m, deterministic)
}
func (m *HandshakeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HandshakeRequest.Merge(m, src)
}
func (m *HandshakeRequest) XXX_Size() int {
	return xxx_messageInfo_HandshakeRequest.Size(m)
}
func (m *HandshakeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HandshakeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HandshakeRequest proto.InternalMessageInfo

func (m *HandshakeRequest) GetAppID() string {
	if m != nil {
		return m.AppID
	}
	return ""
}

func (m *HandshakeRequest) GetAccessKey() string {
	if m != nil {
		return m.AccessKey
	}
	return ""
}

func (m *HandshakeRequest) GetConfig() string {
	if m != nil {
		return m.Config
	}
	return ""
}

type SessionOnConnectNotify struct {
	Sid                  string   `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty"`
	OpenID               string   `protobuf:"bytes,2,opt,name=openID,proto3" json:"openID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionOnConnectNotify) Reset()         { *m = SessionOnConnectNotify{} }
func (m *SessionOnConnectNotify) String() string { return proto.CompactTextString(m) }
func (*SessionOnConnectNotify) ProtoMessage()    {}
func (*SessionOnConnectNotify) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *SessionOnConnectNotify) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionOnConnectNotify.Unmarshal(m, b)
}
func (m *SessionOnConnectNotify) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionOnConnectNotify.Marshal(b, m, deterministic)
}
func (m *SessionOnConnectNotify) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionOnConnectNotify.Merge(m, src)
}
func (m *SessionOnConnectNotify) XXX_Size() int {
	return xxx_messageInfo_SessionOnConnectNotify.Size(m)
}
func (m *SessionOnConnectNotify) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionOnConnectNotify.DiscardUnknown(m)
}

var xxx_messageInfo_SessionOnConnectNotify proto.InternalMessageInfo

func (m *SessionOnConnectNotify) GetSid() string {
	if m != nil {
		return m.Sid
	}
	return ""
}

func (m *SessionOnConnectNotify) GetOpenID() string {
	if m != nil {
		return m.OpenID
	}
	return ""
}

type SessionOnRecvNotify struct {
	Sid                  string            `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty"`
	Api                  string            `protobuf:"bytes,2,opt,name=api,proto3" json:"api,omitempty"`
	Headers              map[string]string `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body                 []byte            `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SessionOnRecvNotify) Reset()         { *m = SessionOnRecvNotify{} }
func (m *SessionOnRecvNotify) String() string { return proto.CompactTextString(m) }
func (*SessionOnRecvNotify) ProtoMessage()    {}
func (*SessionOnRecvNotify) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *SessionOnRecvNotify) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionOnRecvNotify.Unmarshal(m, b)
}
func (m *SessionOnRecvNotify) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionOnRecvNotify.Marshal(b, m, deterministic)
}
func (m *SessionOnRecvNotify) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionOnRecvNotify.Merge(m, src)
}
func (m *SessionOnRecvNotify) XXX_Size() int {
	return xxx_messageInfo_SessionOnRecvNotify.Size(m)
}
func (m *SessionOnRecvNotify) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionOnRecvNotify.DiscardUnknown(m)
}

var xxx_messageInfo_SessionOnRecvNotify proto.InternalMessageInfo

func (m *SessionOnRecvNotify) GetSid() string {
	if m != nil {
		return m.Sid
	}
	return ""
}

func (m *SessionOnRecvNotify) GetApi() string {
	if m != nil {
		return m.Api
	}
	return ""
}

func (m *SessionOnRecvNotify) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *SessionOnRecvNotify) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type SessionOnRecvNotifyResponse struct {
	Status               uint32            `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Headers              map[string]string `protobuf:"bytes,2,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body                 []byte            `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SessionOnRecvNotifyResponse) Reset()         { *m = SessionOnRecvNotifyResponse{} }
func (m *SessionOnRecvNotifyResponse) String() string { return proto.CompactTextString(m) }
func (*SessionOnRecvNotifyResponse) ProtoMessage()    {}
func (*SessionOnRecvNotifyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *SessionOnRecvNotifyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionOnRecvNotifyResponse.Unmarshal(m, b)
}
func (m *SessionOnRecvNotifyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionOnRecvNotifyResponse.Marshal(b, m, deterministic)
}
func (m *SessionOnRecvNotifyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionOnRecvNotifyResponse.Merge(m, src)
}
func (m *SessionOnRecvNotifyResponse) XXX_Size() int {
	return xxx_messageInfo_SessionOnRecvNotifyResponse.Size(m)
}
func (m *SessionOnRecvNotifyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionOnRecvNotifyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SessionOnRecvNotifyResponse proto.InternalMessageInfo

func (m *SessionOnRecvNotifyResponse) GetStatus() uint32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *SessionOnRecvNotifyResponse) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *SessionOnRecvNotifyResponse) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type SessionOnCloseNotify struct {
	Sid                  string   `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionOnCloseNotify) Reset()         { *m = SessionOnCloseNotify{} }
func (m *SessionOnCloseNotify) String() string { return proto.CompactTextString(m) }
func (*SessionOnCloseNotify) ProtoMessage()    {}
func (*SessionOnCloseNotify) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}

func (m *SessionOnCloseNotify) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionOnCloseNotify.Unmarshal(m, b)
}
func (m *SessionOnCloseNotify) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionOnCloseNotify.Marshal(b, m, deterministic)
}
func (m *SessionOnCloseNotify) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionOnCloseNotify.Merge(m, src)
}
func (m *SessionOnCloseNotify) XXX_Size() int {
	return xxx_messageInfo_SessionOnCloseNotify.Size(m)
}
func (m *SessionOnCloseNotify) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionOnCloseNotify.DiscardUnknown(m)
}

var xxx_messageInfo_SessionOnCloseNotify proto.InternalMessageInfo

func (m *SessionOnCloseNotify) GetSid() string {
	if m != nil {
		return m.Sid
	}
	return ""
}

type SessionSendRequest struct {
	Sid                  string   `protobuf:"bytes,1,opt,name=sid,proto3" json:"sid,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionSendRequest) Reset()         { *m = SessionSendRequest{} }
func (m *SessionSendRequest) String() string { return proto.CompactTextString(m) }
func (*SessionSendRequest) ProtoMessage()    {}
func (*SessionSendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}

func (m *SessionSendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionSendRequest.Unmarshal(m, b)
}
func (m *SessionSendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionSendRequest.Marshal(b, m, deterministic)
}
func (m *SessionSendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionSendRequest.Merge(m, src)
}
func (m *SessionSendRequest) XXX_Size() int {
	return xxx_messageInfo_SessionSendRequest.Size(m)
}
func (m *SessionSendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionSendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SessionSendRequest proto.InternalMessageInfo

func (m *SessionSendRequest) GetSid() string {
	if m != nil {
		return m.Sid
	}
	return ""
}

func (m *SessionSendRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type SessionMultiSendRequest struct {
	Sids                 []string `protobuf:"bytes,1,rep,name=sids,proto3" json:"sids,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionMultiSendRequest) Reset()         { *m = SessionMultiSendRequest{} }
func (m *SessionMultiSendRequest) String() string { return proto.CompactTextString(m) }
func (*SessionMultiSendRequest) ProtoMessage()    {}
func (*SessionMultiSendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{6}
}

func (m *SessionMultiSendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionMultiSendRequest.Unmarshal(m, b)
}
func (m *SessionMultiSendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionMultiSendRequest.Marshal(b, m, deterministic)
}
func (m *SessionMultiSendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionMultiSendRequest.Merge(m, src)
}
func (m *SessionMultiSendRequest) XXX_Size() int {
	return xxx_messageInfo_SessionMultiSendRequest.Size(m)
}
func (m *SessionMultiSendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionMultiSendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SessionMultiSendRequest proto.InternalMessageInfo

func (m *SessionMultiSendRequest) GetSids() []string {
	if m != nil {
		return m.Sids
	}
	return nil
}

func (m *SessionMultiSendRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("api.Cmd", Cmd_name, Cmd_value)
	proto.RegisterType((*HandshakeRequest)(nil), "api.HandshakeRequest")
	proto.RegisterType((*SessionOnConnectNotify)(nil), "api.SessionOnConnectNotify")
	proto.RegisterType((*SessionOnRecvNotify)(nil), "api.SessionOnRecvNotify")
	proto.RegisterMapType((map[string]string)(nil), "api.SessionOnRecvNotify.HeadersEntry")
	proto.RegisterType((*SessionOnRecvNotifyResponse)(nil), "api.SessionOnRecvNotifyResponse")
	proto.RegisterMapType((map[string]string)(nil), "api.SessionOnRecvNotifyResponse.HeadersEntry")
	proto.RegisterType((*SessionOnCloseNotify)(nil), "api.SessionOnCloseNotify")
	proto.RegisterType((*SessionSendRequest)(nil), "api.SessionSendRequest")
	proto.RegisterType((*SessionMultiSendRequest)(nil), "api.SessionMultiSendRequest")
}

func init() {
	proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c)
}

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 449 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0xc1, 0x6e, 0x13, 0x31,
	0x10, 0x86, 0xd9, 0x6e, 0x68, 0x95, 0xc9, 0x02, 0x66, 0x92, 0xa6, 0x2b, 0xe0, 0x10, 0xad, 0x84,
	0xb4, 0x42, 0x22, 0x07, 0xb8, 0xa0, 0x5c, 0x10, 0x6c, 0x11, 0xad, 0x10, 0x45, 0x72, 0xef, 0x48,
	0xee, 0x7a, 0x4a, 0x57, 0x49, 0xec, 0x25, 0x76, 0x2a, 0xed, 0x3b, 0xf1, 0x24, 0x1c, 0x78, 0x26,
	0x64, 0xc7, 0xc9, 0xa6, 0x90, 0x88, 0x0b, 0xb7, 0x19, 0xcf, 0xf8, 0x1b, 0xff, 0xbf, 0x6d, 0xe8,
	0x8a, 0xba, 0x1a, 0xd7, 0x0b, 0x6d, 0x35, 0xc6, 0xa2, 0xae, 0xb2, 0xaf, 0xc0, 0xce, 0x84, 0x92,
	0xe6, 0x46, 0x4c, 0x89, 0xd3, 0xf7, 0x25, 0x19, 0x8b, 0x03, 0xb8, 0x2f, 0xea, 0xfa, 0xfc, 0x34,
	0x8d, 0x46, 0x51, 0xde, 0xe5, 0xab, 0x04, 0x9f, 0x41, 0x57, 0x94, 0x25, 0x19, 0xf3, 0x89, 0x9a,
	0xf4, 0xc0, 0x57, 0xda, 0x05, 0x1c, 0xc2, 0x61, 0xa9, 0xd5, 0x75, 0xf5, 0x2d, 0x8d, 0x7d, 0x29,
	0x64, 0xd9, 0x7b, 0x18, 0x5e, 0x92, 0x31, 0x95, 0x56, 0x5f, 0x54, 0xa1, 0x95, 0xa2, 0xd2, 0x5e,
	0x68, 0x5b, 0x5d, 0x37, 0xc8, 0x20, 0x36, 0x95, 0x0c, 0x33, 0x5c, 0xe8, 0x18, 0xba, 0x26, 0x75,
	0x7e, 0x1a, 0xf0, 0x21, 0xcb, 0x7e, 0x46, 0xd0, 0xdf, 0x40, 0x38, 0x95, 0xb7, 0x7b, 0x09, 0x0c,
	0x9c, 0xa8, 0xb0, 0xdd, 0x85, 0xf8, 0x16, 0x8e, 0x6e, 0x48, 0x48, 0x5a, 0x98, 0x34, 0x1e, 0xc5,
	0x79, 0xef, 0xd5, 0xf3, 0xb1, 0x73, 0x60, 0x07, 0x6e, 0x7c, 0xb6, 0xea, 0xfb, 0xa0, 0xec, 0xa2,
	0xe1, 0xeb, 0x5d, 0x88, 0xd0, 0xb9, 0xd2, 0xb2, 0x49, 0x3b, 0xa3, 0x28, 0x4f, 0xb8, 0x8f, 0x9f,
	0x4c, 0x20, 0xd9, 0x6e, 0x76, 0x63, 0xa7, 0xd4, 0xac, 0x0f, 0x32, 0xa5, 0xc6, 0x59, 0x78, 0x2b,
	0x66, 0x4b, 0x0a, 0x47, 0x59, 0x25, 0x93, 0x83, 0x37, 0x51, 0xf6, 0x2b, 0x82, 0xa7, 0x3b, 0xa6,
	0x73, 0x32, 0xb5, 0x56, 0x86, 0x9c, 0x09, 0xc6, 0x0a, 0xbb, 0x34, 0x1e, 0xf7, 0x80, 0x87, 0x0c,
	0x3f, 0xb6, 0x42, 0x0e, 0xbc, 0x90, 0x97, 0xfb, 0x84, 0xac, 0x51, 0xff, 0x10, 0x14, 0xff, 0x27,
	0x41, 0x39, 0x0c, 0xda, 0x1b, 0x9e, 0x69, 0x43, 0xfb, 0x6e, 0x27, 0x9b, 0x00, 0x86, 0xce, 0x4b,
	0x52, 0x72, 0xfd, 0xda, 0xfe, 0xbe, 0x45, 0x84, 0x8e, 0x14, 0x56, 0xf8, 0x51, 0x09, 0xf7, 0x71,
	0xf6, 0x0e, 0x4e, 0xc2, 0xde, 0xcf, 0xcb, 0x99, 0xad, 0xb6, 0x01, 0x08, 0x1d, 0x53, 0x49, 0xe7,
	0x57, 0x9c, 0x77, 0xb9, 0x8f, 0x77, 0x21, 0x5e, 0xfc, 0x88, 0x20, 0x2e, 0xe6, 0x12, 0x7b, 0x70,
	0x54, 0xcc, 0xe5, 0x85, 0x56, 0xc4, 0xee, 0x21, 0x83, 0xa4, 0x98, 0xcb, 0xcd, 0x17, 0x60, 0x11,
	0x9e, 0x40, 0xbf, 0x98, 0xcb, 0x3f, 0x1f, 0x2d, 0x03, 0x1c, 0x00, 0xdb, 0x2e, 0x38, 0xc3, 0x59,
	0x0f, 0x8f, 0xe1, 0xf1, 0x9d, 0x76, 0xe7, 0x00, 0x4b, 0x10, 0xe1, 0x61, 0xbb, 0xec, 0x4e, 0xcb,
	0x06, 0x77, 0xc9, 0x1b, 0x19, 0xec, 0x18, 0xfb, 0xf0, 0xa8, 0x2d, 0xac, 0x08, 0xc3, 0xab, 0x43,
	0xff, 0x4b, 0x5f, 0xff, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xeb, 0x4f, 0x23, 0xbf, 0xb2, 0x03, 0x00,
	0x00,
}
