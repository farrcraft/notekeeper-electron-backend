// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tag.proto

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

type Tag struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 *Title   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Scope                string   `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
	Created              string   `protobuf:"bytes,4,opt,name=created,proto3" json:"created,omitempty"`
	Updated              string   `protobuf:"bytes,5,opt,name=updated,proto3" json:"updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()    {}
func (*Tag) Descriptor() ([]byte, []int) {
	return fileDescriptor_27f545bcde37ecb5, []int{0}
}

func (m *Tag) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tag.Unmarshal(m, b)
}
func (m *Tag) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tag.Marshal(b, m, deterministic)
}
func (m *Tag) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tag.Merge(m, src)
}
func (m *Tag) XXX_Size() int {
	return xxx_messageInfo_Tag.Size(m)
}
func (m *Tag) XXX_DiscardUnknown() {
	xxx_messageInfo_Tag.DiscardUnknown(m)
}

var xxx_messageInfo_Tag proto.InternalMessageInfo

func (m *Tag) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Tag) GetName() *Title {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *Tag) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *Tag) GetCreated() string {
	if m != nil {
		return m.Created
	}
	return ""
}

func (m *Tag) GetUpdated() string {
	if m != nil {
		return m.Updated
	}
	return ""
}

type GetTagsRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Scope                string         `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetTagsRequest) Reset()         { *m = GetTagsRequest{} }
func (m *GetTagsRequest) String() string { return proto.CompactTextString(m) }
func (*GetTagsRequest) ProtoMessage()    {}
func (*GetTagsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_27f545bcde37ecb5, []int{1}
}

func (m *GetTagsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTagsRequest.Unmarshal(m, b)
}
func (m *GetTagsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTagsRequest.Marshal(b, m, deterministic)
}
func (m *GetTagsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTagsRequest.Merge(m, src)
}
func (m *GetTagsRequest) XXX_Size() int {
	return xxx_messageInfo_GetTagsRequest.Size(m)
}
func (m *GetTagsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTagsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTagsRequest proto.InternalMessageInfo

func (m *GetTagsRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetTagsRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *GetTagsRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

type GetTagsResponse struct {
	Header               *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Tags                 []*Tag          `protobuf:"bytes,2,rep,name=tags,proto3" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *GetTagsResponse) Reset()         { *m = GetTagsResponse{} }
func (m *GetTagsResponse) String() string { return proto.CompactTextString(m) }
func (*GetTagsResponse) ProtoMessage()    {}
func (*GetTagsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_27f545bcde37ecb5, []int{2}
}

func (m *GetTagsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTagsResponse.Unmarshal(m, b)
}
func (m *GetTagsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTagsResponse.Marshal(b, m, deterministic)
}
func (m *GetTagsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTagsResponse.Merge(m, src)
}
func (m *GetTagsResponse) XXX_Size() int {
	return xxx_messageInfo_GetTagsResponse.Size(m)
}
func (m *GetTagsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTagsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTagsResponse proto.InternalMessageInfo

func (m *GetTagsResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetTagsResponse) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

type CreateTagRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Name                 *Title         `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Id                   string         `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	Scope                string         `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *CreateTagRequest) Reset()         { *m = CreateTagRequest{} }
func (m *CreateTagRequest) String() string { return proto.CompactTextString(m) }
func (*CreateTagRequest) ProtoMessage()    {}
func (*CreateTagRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_27f545bcde37ecb5, []int{3}
}

func (m *CreateTagRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateTagRequest.Unmarshal(m, b)
}
func (m *CreateTagRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateTagRequest.Marshal(b, m, deterministic)
}
func (m *CreateTagRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateTagRequest.Merge(m, src)
}
func (m *CreateTagRequest) XXX_Size() int {
	return xxx_messageInfo_CreateTagRequest.Size(m)
}
func (m *CreateTagRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateTagRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateTagRequest proto.InternalMessageInfo

func (m *CreateTagRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *CreateTagRequest) GetName() *Title {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *CreateTagRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *CreateTagRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

type SaveTagRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	OwnerId              string         `protobuf:"bytes,3,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Scope                string         `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
	Name                 *Title         `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *SaveTagRequest) Reset()         { *m = SaveTagRequest{} }
func (m *SaveTagRequest) String() string { return proto.CompactTextString(m) }
func (*SaveTagRequest) ProtoMessage()    {}
func (*SaveTagRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_27f545bcde37ecb5, []int{4}
}

func (m *SaveTagRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SaveTagRequest.Unmarshal(m, b)
}
func (m *SaveTagRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SaveTagRequest.Marshal(b, m, deterministic)
}
func (m *SaveTagRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SaveTagRequest.Merge(m, src)
}
func (m *SaveTagRequest) XXX_Size() int {
	return xxx_messageInfo_SaveTagRequest.Size(m)
}
func (m *SaveTagRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SaveTagRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SaveTagRequest proto.InternalMessageInfo

func (m *SaveTagRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *SaveTagRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SaveTagRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *SaveTagRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *SaveTagRequest) GetName() *Title {
	if m != nil {
		return m.Name
	}
	return nil
}

type DeleteTagRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	OwnerId              string         `protobuf:"bytes,3,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Scope                string         `protobuf:"bytes,4,opt,name=scope,proto3" json:"scope,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *DeleteTagRequest) Reset()         { *m = DeleteTagRequest{} }
func (m *DeleteTagRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteTagRequest) ProtoMessage()    {}
func (*DeleteTagRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_27f545bcde37ecb5, []int{5}
}

func (m *DeleteTagRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteTagRequest.Unmarshal(m, b)
}
func (m *DeleteTagRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteTagRequest.Marshal(b, m, deterministic)
}
func (m *DeleteTagRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteTagRequest.Merge(m, src)
}
func (m *DeleteTagRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteTagRequest.Size(m)
}
func (m *DeleteTagRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteTagRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteTagRequest proto.InternalMessageInfo

func (m *DeleteTagRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *DeleteTagRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DeleteTagRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *DeleteTagRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func init() {
	proto.RegisterType((*Tag)(nil), "notekeeper.Tag")
	proto.RegisterType((*GetTagsRequest)(nil), "notekeeper.GetTagsRequest")
	proto.RegisterType((*GetTagsResponse)(nil), "notekeeper.GetTagsResponse")
	proto.RegisterType((*CreateTagRequest)(nil), "notekeeper.CreateTagRequest")
	proto.RegisterType((*SaveTagRequest)(nil), "notekeeper.SaveTagRequest")
	proto.RegisterType((*DeleteTagRequest)(nil), "notekeeper.DeleteTagRequest")
}

func init() { proto.RegisterFile("tag.proto", fileDescriptor_27f545bcde37ecb5) }

var fileDescriptor_27f545bcde37ecb5 = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x53, 0x4d, 0x4e, 0x83, 0x40,
	0x18, 0x75, 0xf8, 0x69, 0xd3, 0x8f, 0x84, 0xd6, 0x89, 0x8b, 0xb1, 0x2b, 0x82, 0x31, 0x61, 0x45,
	0x22, 0x1e, 0x41, 0x13, 0x75, 0x67, 0x90, 0x0b, 0x8c, 0xe5, 0x0b, 0xa2, 0x95, 0x41, 0x66, 0xaa,
	0x27, 0xd0, 0x03, 0x78, 0x0c, 0x4f, 0x69, 0x18, 0x7e, 0x6c, 0xb1, 0x31, 0x46, 0x17, 0x5d, 0xbe,
	0x79, 0xef, 0xcb, 0xfb, 0x21, 0xc0, 0x44, 0xf1, 0x2c, 0x2c, 0x2b, 0xa1, 0x04, 0x85, 0x42, 0x28,
	0x7c, 0x40, 0x2c, 0xb1, 0x9a, 0x3b, 0x2a, 0x57, 0x4b, 0x6c, 0x08, 0xff, 0x8d, 0x80, 0x99, 0xf0,
	0x8c, 0xba, 0x60, 0xe4, 0x29, 0x23, 0x1e, 0x09, 0x26, 0xb1, 0x91, 0xa7, 0xf4, 0x18, 0xac, 0x82,
	0x3f, 0x22, 0x33, 0x3c, 0x12, 0x38, 0xd1, 0x7e, 0xf8, 0x75, 0x1f, 0x26, 0xf5, 0x79, 0xac, 0x69,
	0x7a, 0x00, 0xb6, 0x5c, 0x88, 0x12, 0x99, 0xa9, 0x2f, 0x1b, 0x40, 0x19, 0x8c, 0x17, 0x15, 0x72,
	0x85, 0x29, 0xb3, 0xf4, 0x7b, 0x07, 0x6b, 0x66, 0x55, 0xa6, 0x9a, 0xb1, 0x1b, 0xa6, 0x85, 0x7e,
	0x0e, 0xee, 0x05, 0xaa, 0x84, 0x67, 0x32, 0xc6, 0xa7, 0x15, 0x4a, 0x45, 0x4f, 0x60, 0x74, 0x87,
	0x3c, 0xc5, 0x4a, 0xc7, 0x72, 0xa2, 0xc3, 0xf5, 0x10, 0xad, 0xe8, 0x52, 0x0b, 0xe2, 0x56, 0xd8,
	0xb6, 0x30, 0xfa, 0x16, 0x5b, 0xe3, 0xf9, 0xf7, 0x30, 0xed, 0xad, 0x64, 0x29, 0x0a, 0x89, 0x34,
	0x1a, 0x78, 0xcd, 0x37, 0xbd, 0x1a, 0xd5, 0xc0, 0xec, 0x08, 0x2c, 0xc5, 0x33, 0xc9, 0x0c, 0xcf,
	0x0c, 0x9c, 0x68, 0xba, 0x31, 0x11, 0xcf, 0x62, 0x4d, 0xfa, 0xef, 0x04, 0x66, 0x67, 0xba, 0x7c,
	0xfd, 0xf6, 0xf7, 0x66, 0xbf, 0xfc, 0x1e, 0xcd, 0x00, 0xe6, 0xf7, 0x01, 0xac, 0xf5, 0x01, 0x3e,
	0x08, 0xb8, 0x37, 0xfc, 0xf9, 0x9f, 0x91, 0x86, 0x63, 0x33, 0x18, 0x8b, 0x97, 0x02, 0xab, 0xab,
	0x2e, 0x40, 0x07, 0xb7, 0xa7, 0xe8, 0x2b, 0xd9, 0x3f, 0x56, 0xf2, 0x5f, 0x09, 0xcc, 0xce, 0x71,
	0x89, 0x6a, 0xb7, 0x71, 0xaf, 0xf7, 0x6e, 0x47, 0xfa, 0x97, 0x39, 0xfd, 0x0c, 0x00, 0x00, 0xff,
	0xff, 0x4f, 0x18, 0x1a, 0x21, 0x58, 0x03, 0x00, 0x00,
}