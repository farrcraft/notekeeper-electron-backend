// Code generated by protoc-gen-go. DO NOT EDIT.
// source: title.proto

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

type Title struct {
	Text                 string   `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Bold                 bool     `protobuf:"varint,2,opt,name=bold,proto3" json:"bold,omitempty"`
	Italics              bool     `protobuf:"varint,3,opt,name=italics,proto3" json:"italics,omitempty"`
	Underscore           bool     `protobuf:"varint,4,opt,name=underscore,proto3" json:"underscore,omitempty"`
	Strike               bool     `protobuf:"varint,5,opt,name=strike,proto3" json:"strike,omitempty"`
	Color                string   `protobuf:"bytes,6,opt,name=color,proto3" json:"color,omitempty"`
	Background           string   `protobuf:"bytes,7,opt,name=background,proto3" json:"background,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Title) Reset()         { *m = Title{} }
func (m *Title) String() string { return proto.CompactTextString(m) }
func (*Title) ProtoMessage()    {}
func (*Title) Descriptor() ([]byte, []int) {
	return fileDescriptor_c4b226b0f73ca670, []int{0}
}

func (m *Title) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Title.Unmarshal(m, b)
}
func (m *Title) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Title.Marshal(b, m, deterministic)
}
func (m *Title) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Title.Merge(m, src)
}
func (m *Title) XXX_Size() int {
	return xxx_messageInfo_Title.Size(m)
}
func (m *Title) XXX_DiscardUnknown() {
	xxx_messageInfo_Title.DiscardUnknown(m)
}

var xxx_messageInfo_Title proto.InternalMessageInfo

func (m *Title) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Title) GetBold() bool {
	if m != nil {
		return m.Bold
	}
	return false
}

func (m *Title) GetItalics() bool {
	if m != nil {
		return m.Italics
	}
	return false
}

func (m *Title) GetUnderscore() bool {
	if m != nil {
		return m.Underscore
	}
	return false
}

func (m *Title) GetStrike() bool {
	if m != nil {
		return m.Strike
	}
	return false
}

func (m *Title) GetColor() string {
	if m != nil {
		return m.Color
	}
	return ""
}

func (m *Title) GetBackground() string {
	if m != nil {
		return m.Background
	}
	return ""
}

func init() {
	proto.RegisterType((*Title)(nil), "notekeeper.Title")
}

func init() { proto.RegisterFile("title.proto", fileDescriptor_c4b226b0f73ca670) }

var fileDescriptor_c4b226b0f73ca670 = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x4d, 0x6a, 0xc3, 0x40,
	0x0c, 0x85, 0x3b, 0xad, 0x7f, 0x5a, 0xb5, 0xab, 0xa1, 0x14, 0xd1, 0x45, 0x31, 0x5d, 0x79, 0x95,
	0x4d, 0x2e, 0x12, 0x4c, 0x2e, 0x60, 0x8f, 0x45, 0x18, 0x3c, 0x1e, 0x99, 0xb1, 0x0c, 0xb9, 0x59,
	0xae, 0x17, 0x2c, 0x27, 0xe0, 0xdd, 0x7b, 0xdf, 0x27, 0x24, 0x04, 0x9f, 0xe2, 0x25, 0xd0, 0x61,
	0x4a, 0x2c, 0x6c, 0x21, 0xb2, 0xd0, 0x40, 0x34, 0x51, 0xfa, 0xfd, 0x72, 0x3c, 0x8e, 0x1c, 0x37,
	0xf3, 0x7f, 0x33, 0x90, 0x9f, 0xd7, 0x49, 0x6b, 0x21, 0x13, 0xba, 0x0a, 0x9a, 0xca, 0xd4, 0x1f,
	0x8d, 0xe6, 0x95, 0x75, 0x1c, 0x7a, 0x7c, 0xad, 0x4c, 0xfd, 0xde, 0x68, 0xb6, 0x08, 0xa5, 0x97,
	0x36, 0x78, 0x37, 0xe3, 0x9b, 0xe2, 0x67, 0xb5, 0x7f, 0x00, 0x4b, 0xec, 0x29, 0xcd, 0x8e, 0x13,
	0x61, 0xa6, 0x72, 0x47, 0xec, 0x0f, 0x14, 0xb3, 0x24, 0x3f, 0x10, 0xe6, 0xea, 0x1e, 0xcd, 0x7e,
	0x43, 0xee, 0x38, 0x70, 0xc2, 0x42, 0x4f, 0x6f, 0x65, 0xdd, 0xd6, 0xb5, 0x6e, 0xb8, 0x24, 0x5e,
	0x62, 0x8f, 0xa5, 0xaa, 0x1d, 0x39, 0xbd, 0x74, 0x85, 0xbe, 0x70, 0xbc, 0x07, 0x00, 0x00, 0xff,
	0xff, 0xe3, 0xdb, 0xb7, 0x56, 0xeb, 0x00, 0x00, 0x00,
}
