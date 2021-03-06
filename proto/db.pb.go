// Code generated by protoc-gen-go. DO NOT EDIT.
// source: db.proto

package notekeeper

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

type OpenMasterDbRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Path                 string         `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *OpenMasterDbRequest) Reset()         { *m = OpenMasterDbRequest{} }
func (m *OpenMasterDbRequest) String() string { return proto.CompactTextString(m) }
func (*OpenMasterDbRequest) ProtoMessage()    {}
func (*OpenMasterDbRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8817812184a13374, []int{0}
}

func (m *OpenMasterDbRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OpenMasterDbRequest.Unmarshal(m, b)
}
func (m *OpenMasterDbRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OpenMasterDbRequest.Marshal(b, m, deterministic)
}
func (m *OpenMasterDbRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OpenMasterDbRequest.Merge(m, src)
}
func (m *OpenMasterDbRequest) XXX_Size() int {
	return xxx_messageInfo_OpenMasterDbRequest.Size(m)
}
func (m *OpenMasterDbRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OpenMasterDbRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OpenMasterDbRequest proto.InternalMessageInfo

func (m *OpenMasterDbRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *OpenMasterDbRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func init() {
	proto.RegisterType((*OpenMasterDbRequest)(nil), "notekeeper.OpenMasterDbRequest")
}

func init() { proto.RegisterFile("db.proto", fileDescriptor_8817812184a13374) }

var fileDescriptor_8817812184a13374 = []byte{
	// 131 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x48, 0x49, 0xd2, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xca, 0xcb, 0x2f, 0x49, 0xcd, 0x4e, 0x4d, 0x2d, 0x48, 0x2d,
	0x92, 0xe2, 0x49, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0x83, 0xc8, 0x28, 0xc5, 0x70, 0x09, 0xfb, 0x17,
	0xa4, 0xe6, 0xf9, 0x26, 0x16, 0x97, 0xa4, 0x16, 0xb9, 0x24, 0x05, 0xa5, 0x16, 0x96, 0xa6, 0x16,
	0x97, 0x08, 0x19, 0x72, 0xb1, 0x65, 0xa4, 0x26, 0xa6, 0xa4, 0x16, 0x49, 0x30, 0x2a, 0x30, 0x6a,
	0x70, 0x1b, 0x49, 0xea, 0x21, 0x4c, 0xd0, 0x83, 0x2a, 0xf2, 0x00, 0x2b, 0x08, 0x82, 0x2a, 0x14,
	0x12, 0xe2, 0x62, 0x29, 0x48, 0x2c, 0xc9, 0x90, 0x60, 0x52, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3,
	0x93, 0xd8, 0xc0, 0x96, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x27, 0x80, 0xd6, 0x37, 0x8a,
	0x00, 0x00, 0x00,
}
