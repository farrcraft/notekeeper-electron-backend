// Code generated by protoc-gen-go. DO NOT EDIT.
// source: shelf.proto

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

type Shelf struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 *Title   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Scope                string   `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
	Default              bool     `protobuf:"varint,4,opt,name=default,proto3" json:"default,omitempty"`
	Trash                bool     `protobuf:"varint,5,opt,name=trash,proto3" json:"trash,omitempty"`
	Locked               bool     `protobuf:"varint,6,opt,name=locked,proto3" json:"locked,omitempty"`
	Created              string   `protobuf:"bytes,7,opt,name=created,proto3" json:"created,omitempty"`
	Updated              string   `protobuf:"bytes,8,opt,name=updated,proto3" json:"updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Shelf) Reset()         { *m = Shelf{} }
func (m *Shelf) String() string { return proto.CompactTextString(m) }
func (*Shelf) ProtoMessage()    {}
func (*Shelf) Descriptor() ([]byte, []int) {
	return fileDescriptor_997c08397bcb74ab, []int{0}
}

func (m *Shelf) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Shelf.Unmarshal(m, b)
}
func (m *Shelf) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Shelf.Marshal(b, m, deterministic)
}
func (m *Shelf) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Shelf.Merge(m, src)
}
func (m *Shelf) XXX_Size() int {
	return xxx_messageInfo_Shelf.Size(m)
}
func (m *Shelf) XXX_DiscardUnknown() {
	xxx_messageInfo_Shelf.DiscardUnknown(m)
}

var xxx_messageInfo_Shelf proto.InternalMessageInfo

func (m *Shelf) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Shelf) GetName() *Title {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *Shelf) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *Shelf) GetDefault() bool {
	if m != nil {
		return m.Default
	}
	return false
}

func (m *Shelf) GetTrash() bool {
	if m != nil {
		return m.Trash
	}
	return false
}

func (m *Shelf) GetLocked() bool {
	if m != nil {
		return m.Locked
	}
	return false
}

func (m *Shelf) GetCreated() string {
	if m != nil {
		return m.Created
	}
	return ""
}

func (m *Shelf) GetUpdated() string {
	if m != nil {
		return m.Updated
	}
	return ""
}

type GetShelvesRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Scope                string         `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetShelvesRequest) Reset()         { *m = GetShelvesRequest{} }
func (m *GetShelvesRequest) String() string { return proto.CompactTextString(m) }
func (*GetShelvesRequest) ProtoMessage()    {}
func (*GetShelvesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_997c08397bcb74ab, []int{1}
}

func (m *GetShelvesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetShelvesRequest.Unmarshal(m, b)
}
func (m *GetShelvesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetShelvesRequest.Marshal(b, m, deterministic)
}
func (m *GetShelvesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetShelvesRequest.Merge(m, src)
}
func (m *GetShelvesRequest) XXX_Size() int {
	return xxx_messageInfo_GetShelvesRequest.Size(m)
}
func (m *GetShelvesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetShelvesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetShelvesRequest proto.InternalMessageInfo

func (m *GetShelvesRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetShelvesRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GetShelvesRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

type GetShelvesResponse struct {
	Header               *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Shelves              []*Shelf        `protobuf:"bytes,2,rep,name=shelves,proto3" json:"shelves,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *GetShelvesResponse) Reset()         { *m = GetShelvesResponse{} }
func (m *GetShelvesResponse) String() string { return proto.CompactTextString(m) }
func (*GetShelvesResponse) ProtoMessage()    {}
func (*GetShelvesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_997c08397bcb74ab, []int{2}
}

func (m *GetShelvesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetShelvesResponse.Unmarshal(m, b)
}
func (m *GetShelvesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetShelvesResponse.Marshal(b, m, deterministic)
}
func (m *GetShelvesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetShelvesResponse.Merge(m, src)
}
func (m *GetShelvesResponse) XXX_Size() int {
	return xxx_messageInfo_GetShelvesResponse.Size(m)
}
func (m *GetShelvesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetShelvesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetShelvesResponse proto.InternalMessageInfo

func (m *GetShelvesResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetShelvesResponse) GetShelves() []*Shelf {
	if m != nil {
		return m.Shelves
	}
	return nil
}

type CreateShelfRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Name                 *Title         `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Id                   string         `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Scope                string         `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *CreateShelfRequest) Reset()         { *m = CreateShelfRequest{} }
func (m *CreateShelfRequest) String() string { return proto.CompactTextString(m) }
func (*CreateShelfRequest) ProtoMessage()    {}
func (*CreateShelfRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_997c08397bcb74ab, []int{3}
}

func (m *CreateShelfRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateShelfRequest.Unmarshal(m, b)
}
func (m *CreateShelfRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateShelfRequest.Marshal(b, m, deterministic)
}
func (m *CreateShelfRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateShelfRequest.Merge(m, src)
}
func (m *CreateShelfRequest) XXX_Size() int {
	return xxx_messageInfo_CreateShelfRequest.Size(m)
}
func (m *CreateShelfRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateShelfRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateShelfRequest proto.InternalMessageInfo

func (m *CreateShelfRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *CreateShelfRequest) GetName() *Title {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *CreateShelfRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *CreateShelfRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

type SaveShelfRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	OwnerId              string         `protobuf:"bytes,3,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Scope                string         `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
	Name                 *Title         `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Locked               bool           `protobuf:"varint,6,opt,name=locked,proto3" json:"locked,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *SaveShelfRequest) Reset()         { *m = SaveShelfRequest{} }
func (m *SaveShelfRequest) String() string { return proto.CompactTextString(m) }
func (*SaveShelfRequest) ProtoMessage()    {}
func (*SaveShelfRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_997c08397bcb74ab, []int{4}
}

func (m *SaveShelfRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SaveShelfRequest.Unmarshal(m, b)
}
func (m *SaveShelfRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SaveShelfRequest.Marshal(b, m, deterministic)
}
func (m *SaveShelfRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SaveShelfRequest.Merge(m, src)
}
func (m *SaveShelfRequest) XXX_Size() int {
	return xxx_messageInfo_SaveShelfRequest.Size(m)
}
func (m *SaveShelfRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SaveShelfRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SaveShelfRequest proto.InternalMessageInfo

func (m *SaveShelfRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *SaveShelfRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SaveShelfRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *SaveShelfRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *SaveShelfRequest) GetName() *Title {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *SaveShelfRequest) GetLocked() bool {
	if m != nil {
		return m.Locked
	}
	return false
}

type DeleteShelfRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	OwnerId              string         `protobuf:"bytes,3,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Scope                string         `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *DeleteShelfRequest) Reset()         { *m = DeleteShelfRequest{} }
func (m *DeleteShelfRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteShelfRequest) ProtoMessage()    {}
func (*DeleteShelfRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_997c08397bcb74ab, []int{5}
}

func (m *DeleteShelfRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteShelfRequest.Unmarshal(m, b)
}
func (m *DeleteShelfRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteShelfRequest.Marshal(b, m, deterministic)
}
func (m *DeleteShelfRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteShelfRequest.Merge(m, src)
}
func (m *DeleteShelfRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteShelfRequest.Size(m)
}
func (m *DeleteShelfRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteShelfRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteShelfRequest proto.InternalMessageInfo

func (m *DeleteShelfRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *DeleteShelfRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DeleteShelfRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *DeleteShelfRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func init() {
	proto.RegisterType((*Shelf)(nil), "notekeeper.Shelf")
	proto.RegisterType((*GetShelvesRequest)(nil), "notekeeper.GetShelvesRequest")
	proto.RegisterType((*GetShelvesResponse)(nil), "notekeeper.GetShelvesResponse")
	proto.RegisterType((*CreateShelfRequest)(nil), "notekeeper.CreateShelfRequest")
	proto.RegisterType((*SaveShelfRequest)(nil), "notekeeper.SaveShelfRequest")
	proto.RegisterType((*DeleteShelfRequest)(nil), "notekeeper.DeleteShelfRequest")
}

func init() { proto.RegisterFile("shelf.proto", fileDescriptor_997c08397bcb74ab) }

var fileDescriptor_997c08397bcb74ab = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x53, 0xd1, 0x4a, 0xf3, 0x30,
	0x18, 0xfd, 0xd3, 0xad, 0xed, 0xfe, 0xaf, 0x22, 0x2e, 0x88, 0xc4, 0x5d, 0x95, 0x82, 0x30, 0x10,
	0x06, 0xce, 0x47, 0x50, 0x50, 0xef, 0x46, 0xe7, 0x0b, 0xd4, 0xe6, 0x1b, 0x1b, 0xeb, 0x9a, 0xda,
	0xa4, 0xf3, 0x11, 0x7c, 0x01, 0x1f, 0xc9, 0x07, 0xf0, 0x91, 0x24, 0x49, 0x3b, 0x67, 0x9d, 0x22,
	0x7a, 0xe1, 0xe5, 0xc9, 0x39, 0x39, 0xdf, 0x77, 0x4e, 0x5a, 0x08, 0xe4, 0x1c, 0xb3, 0xd9, 0xa8,
	0x28, 0x85, 0x12, 0x14, 0x72, 0xa1, 0x70, 0x89, 0x58, 0x60, 0x39, 0xd8, 0x4b, 0xc5, 0x6a, 0x25,
	0x72, 0xcb, 0x0c, 0x02, 0xb5, 0x50, 0x19, 0x5a, 0x10, 0xbd, 0x10, 0x70, 0xa7, 0xfa, 0x1a, 0xdd,
	0x07, 0x67, 0xc1, 0x19, 0x09, 0xc9, 0xf0, 0x7f, 0xec, 0x2c, 0x38, 0x3d, 0x81, 0x6e, 0x9e, 0xac,
	0x90, 0x39, 0x21, 0x19, 0x06, 0xe3, 0xfe, 0xe8, 0xcd, 0x6f, 0x74, 0xab, 0x0d, 0x62, 0x43, 0xd3,
	0x43, 0x70, 0x65, 0x2a, 0x0a, 0x64, 0x1d, 0x73, 0xd3, 0x02, 0xca, 0xc0, 0xe7, 0x38, 0x4b, 0xaa,
	0x4c, 0xb1, 0x6e, 0x48, 0x86, 0xbd, 0xb8, 0x81, 0x5a, 0xaf, 0xca, 0x44, 0xce, 0x99, 0x6b, 0xce,
	0x2d, 0xa0, 0x47, 0xe0, 0x65, 0x22, 0x5d, 0x22, 0x67, 0x9e, 0x39, 0xae, 0x91, 0xf6, 0x49, 0x4b,
	0x4c, 0x14, 0x72, 0xe6, 0x1b, 0xff, 0x06, 0x6a, 0xa6, 0x2a, 0xb8, 0x61, 0x7a, 0x96, 0xa9, 0x61,
	0x94, 0x41, 0xff, 0x0a, 0x95, 0x0e, 0xb5, 0x46, 0x19, 0xe3, 0x7d, 0x85, 0x52, 0xd1, 0x33, 0xf0,
	0xe6, 0x98, 0x70, 0x2c, 0x4d, 0xc2, 0x60, 0x7c, 0xbc, 0x9d, 0xa7, 0x16, 0x5d, 0x1b, 0x41, 0x5c,
	0x0b, 0xeb, 0x42, 0x9c, 0x4d, 0x21, 0x3b, 0x93, 0x46, 0x15, 0xd0, 0xed, 0x69, 0xb2, 0x10, 0xb9,
	0x44, 0x3a, 0x6e, 0x8d, 0x1b, 0xbc, 0x1f, 0x67, 0x55, 0xad, 0x79, 0xa7, 0xe0, 0x4b, 0x6b, 0xc3,
	0x9c, 0xb0, 0xd3, 0xee, 0xdc, 0x3c, 0x52, 0xdc, 0x28, 0xa2, 0x27, 0x02, 0xf4, 0xc2, 0x54, 0x61,
	0x89, 0x9f, 0xc7, 0xfc, 0xe6, 0x3b, 0xdb, 0x36, 0x3a, 0x1f, 0xdb, 0xe8, 0x6e, 0xb7, 0xf1, 0x4c,
	0xe0, 0x60, 0x9a, 0xac, 0x7f, 0xbd, 0x54, 0xbb, 0x7b, 0x06, 0xbe, 0x78, 0xc8, 0xb1, 0xbc, 0x69,
	0x56, 0x68, 0xe0, 0xee, 0x3d, 0x36, 0xa1, 0xdc, 0xaf, 0x43, 0x7d, 0xf2, 0xd9, 0x45, 0x8f, 0x04,
	0xe8, 0x25, 0x66, 0xa8, 0xfe, 0x3a, 0xc8, 0xe4, 0xdf, 0x84, 0xdc, 0x79, 0xe6, 0x57, 0x3d, 0x7f,
	0x0d, 0x00, 0x00, 0xff, 0xff, 0x3a, 0x11, 0x69, 0xa0, 0xe0, 0x03, 0x00, 0x00,
}