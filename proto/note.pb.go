// Code generated by protoc-gen-go. DO NOT EDIT.
// source: note.proto

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

// This is just note metadata
// Note content is treated separately
type Note struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	NotebookId           string   `protobuf:"bytes,2,opt,name=notebookId,proto3" json:"notebookId,omitempty"`
	OwnerId              string   `protobuf:"bytes,3,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	StoreId              string   `protobuf:"bytes,4,opt,name=storeId,proto3" json:"storeId,omitempty"`
	Scope                string   `protobuf:"bytes,5,opt,name=scope,proto3" json:"scope,omitempty"`
	Store                string   `protobuf:"bytes,6,opt,name=store,proto3" json:"store,omitempty"`
	Name                 *Title   `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
	Type                 string   `protobuf:"bytes,8,opt,name=type,proto3" json:"type,omitempty"`
	Revisions            int32    `protobuf:"varint,9,opt,name=revisions,proto3" json:"revisions,omitempty"`
	Locked               bool     `protobuf:"varint,10,opt,name=locked,proto3" json:"locked,omitempty"`
	Created              string   `protobuf:"bytes,11,opt,name=created,proto3" json:"created,omitempty"`
	Updated              string   `protobuf:"bytes,12,opt,name=updated,proto3" json:"updated,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Note) Reset()         { *m = Note{} }
func (m *Note) String() string { return proto.CompactTextString(m) }
func (*Note) ProtoMessage()    {}
func (*Note) Descriptor() ([]byte, []int) {
	return fileDescriptor_640dafe07df50d4e, []int{0}
}

func (m *Note) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Note.Unmarshal(m, b)
}
func (m *Note) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Note.Marshal(b, m, deterministic)
}
func (m *Note) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Note.Merge(m, src)
}
func (m *Note) XXX_Size() int {
	return xxx_messageInfo_Note.Size(m)
}
func (m *Note) XXX_DiscardUnknown() {
	xxx_messageInfo_Note.DiscardUnknown(m)
}

var xxx_messageInfo_Note proto.InternalMessageInfo

func (m *Note) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Note) GetNotebookId() string {
	if m != nil {
		return m.NotebookId
	}
	return ""
}

func (m *Note) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *Note) GetStoreId() string {
	if m != nil {
		return m.StoreId
	}
	return ""
}

func (m *Note) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *Note) GetStore() string {
	if m != nil {
		return m.Store
	}
	return ""
}

func (m *Note) GetName() *Title {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *Note) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Note) GetRevisions() int32 {
	if m != nil {
		return m.Revisions
	}
	return 0
}

func (m *Note) GetLocked() bool {
	if m != nil {
		return m.Locked
	}
	return false
}

func (m *Note) GetCreated() string {
	if m != nil {
		return m.Created
	}
	return ""
}

func (m *Note) GetUpdated() string {
	if m != nil {
		return m.Updated
	}
	return ""
}

type CreateNoteRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	NotebookId           string         `protobuf:"bytes,2,opt,name=notebookId,proto3" json:"notebookId,omitempty"`
	StoreId              string         `protobuf:"bytes,3,opt,name=storeId,proto3" json:"storeId,omitempty"`
	OwnerId              string         `protobuf:"bytes,4,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Scope                string         `protobuf:"bytes,5,opt,name=scope,proto3" json:"scope,omitempty"`
	Store                string         `protobuf:"bytes,6,opt,name=store,proto3" json:"store,omitempty"`
	Name                 *Title         `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *CreateNoteRequest) Reset()         { *m = CreateNoteRequest{} }
func (m *CreateNoteRequest) String() string { return proto.CompactTextString(m) }
func (*CreateNoteRequest) ProtoMessage()    {}
func (*CreateNoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_640dafe07df50d4e, []int{1}
}

func (m *CreateNoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateNoteRequest.Unmarshal(m, b)
}
func (m *CreateNoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateNoteRequest.Marshal(b, m, deterministic)
}
func (m *CreateNoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateNoteRequest.Merge(m, src)
}
func (m *CreateNoteRequest) XXX_Size() int {
	return xxx_messageInfo_CreateNoteRequest.Size(m)
}
func (m *CreateNoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateNoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateNoteRequest proto.InternalMessageInfo

func (m *CreateNoteRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *CreateNoteRequest) GetNotebookId() string {
	if m != nil {
		return m.NotebookId
	}
	return ""
}

func (m *CreateNoteRequest) GetStoreId() string {
	if m != nil {
		return m.StoreId
	}
	return ""
}

func (m *CreateNoteRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *CreateNoteRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *CreateNoteRequest) GetStore() string {
	if m != nil {
		return m.Store
	}
	return ""
}

func (m *CreateNoteRequest) GetName() *Title {
	if m != nil {
		return m.Name
	}
	return nil
}

type SaveNoteRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	NotebookId           string         `protobuf:"bytes,3,opt,name=notebookId,proto3" json:"notebookId,omitempty"`
	StoreId              string         `protobuf:"bytes,4,opt,name=storeId,proto3" json:"storeId,omitempty"`
	OwnerId              string         `protobuf:"bytes,5,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Scope                string         `protobuf:"bytes,6,opt,name=scope,proto3" json:"scope,omitempty"`
	Store                string         `protobuf:"bytes,7,opt,name=store,proto3" json:"store,omitempty"`
	Name                 *Title         `protobuf:"bytes,8,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *SaveNoteRequest) Reset()         { *m = SaveNoteRequest{} }
func (m *SaveNoteRequest) String() string { return proto.CompactTextString(m) }
func (*SaveNoteRequest) ProtoMessage()    {}
func (*SaveNoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_640dafe07df50d4e, []int{2}
}

func (m *SaveNoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SaveNoteRequest.Unmarshal(m, b)
}
func (m *SaveNoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SaveNoteRequest.Marshal(b, m, deterministic)
}
func (m *SaveNoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SaveNoteRequest.Merge(m, src)
}
func (m *SaveNoteRequest) XXX_Size() int {
	return xxx_messageInfo_SaveNoteRequest.Size(m)
}
func (m *SaveNoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SaveNoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SaveNoteRequest proto.InternalMessageInfo

func (m *SaveNoteRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *SaveNoteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SaveNoteRequest) GetNotebookId() string {
	if m != nil {
		return m.NotebookId
	}
	return ""
}

func (m *SaveNoteRequest) GetStoreId() string {
	if m != nil {
		return m.StoreId
	}
	return ""
}

func (m *SaveNoteRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *SaveNoteRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *SaveNoteRequest) GetStore() string {
	if m != nil {
		return m.Store
	}
	return ""
}

func (m *SaveNoteRequest) GetName() *Title {
	if m != nil {
		return m.Name
	}
	return nil
}

type DeleteNoteRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	NotebookId           string         `protobuf:"bytes,3,opt,name=notebookId,proto3" json:"notebookId,omitempty"`
	StoreId              string         `protobuf:"bytes,4,opt,name=storeId,proto3" json:"storeId,omitempty"`
	OwnerId              string         `protobuf:"bytes,5,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Scope                string         `protobuf:"bytes,6,opt,name=scope,proto3" json:"scope,omitempty"`
	Store                string         `protobuf:"bytes,7,opt,name=store,proto3" json:"store,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *DeleteNoteRequest) Reset()         { *m = DeleteNoteRequest{} }
func (m *DeleteNoteRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteNoteRequest) ProtoMessage()    {}
func (*DeleteNoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_640dafe07df50d4e, []int{3}
}

func (m *DeleteNoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteNoteRequest.Unmarshal(m, b)
}
func (m *DeleteNoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteNoteRequest.Marshal(b, m, deterministic)
}
func (m *DeleteNoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteNoteRequest.Merge(m, src)
}
func (m *DeleteNoteRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteNoteRequest.Size(m)
}
func (m *DeleteNoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteNoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteNoteRequest proto.InternalMessageInfo

func (m *DeleteNoteRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *DeleteNoteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *DeleteNoteRequest) GetNotebookId() string {
	if m != nil {
		return m.NotebookId
	}
	return ""
}

func (m *DeleteNoteRequest) GetStoreId() string {
	if m != nil {
		return m.StoreId
	}
	return ""
}

func (m *DeleteNoteRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *DeleteNoteRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *DeleteNoteRequest) GetStore() string {
	if m != nil {
		return m.Store
	}
	return ""
}

type LoadNoteRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Id                   string         `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	StoreId              string         `protobuf:"bytes,3,opt,name=storeId,proto3" json:"storeId,omitempty"`
	OwnerId              string         `protobuf:"bytes,4,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Scope                string         `protobuf:"bytes,5,opt,name=scope,proto3" json:"scope,omitempty"`
	Store                string         `protobuf:"bytes,6,opt,name=store,proto3" json:"store,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *LoadNoteRequest) Reset()         { *m = LoadNoteRequest{} }
func (m *LoadNoteRequest) String() string { return proto.CompactTextString(m) }
func (*LoadNoteRequest) ProtoMessage()    {}
func (*LoadNoteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_640dafe07df50d4e, []int{4}
}

func (m *LoadNoteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadNoteRequest.Unmarshal(m, b)
}
func (m *LoadNoteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadNoteRequest.Marshal(b, m, deterministic)
}
func (m *LoadNoteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadNoteRequest.Merge(m, src)
}
func (m *LoadNoteRequest) XXX_Size() int {
	return xxx_messageInfo_LoadNoteRequest.Size(m)
}
func (m *LoadNoteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadNoteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoadNoteRequest proto.InternalMessageInfo

func (m *LoadNoteRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *LoadNoteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *LoadNoteRequest) GetStoreId() string {
	if m != nil {
		return m.StoreId
	}
	return ""
}

func (m *LoadNoteRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *LoadNoteRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *LoadNoteRequest) GetStore() string {
	if m != nil {
		return m.Store
	}
	return ""
}

type LoadNoteResponse struct {
	Header               *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Note                 *Note           `protobuf:"bytes,2,opt,name=note,proto3" json:"note,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *LoadNoteResponse) Reset()         { *m = LoadNoteResponse{} }
func (m *LoadNoteResponse) String() string { return proto.CompactTextString(m) }
func (*LoadNoteResponse) ProtoMessage()    {}
func (*LoadNoteResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_640dafe07df50d4e, []int{5}
}

func (m *LoadNoteResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoadNoteResponse.Unmarshal(m, b)
}
func (m *LoadNoteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoadNoteResponse.Marshal(b, m, deterministic)
}
func (m *LoadNoteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadNoteResponse.Merge(m, src)
}
func (m *LoadNoteResponse) XXX_Size() int {
	return xxx_messageInfo_LoadNoteResponse.Size(m)
}
func (m *LoadNoteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadNoteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LoadNoteResponse proto.InternalMessageInfo

func (m *LoadNoteResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *LoadNoteResponse) GetNote() *Note {
	if m != nil {
		return m.Note
	}
	return nil
}

type GetNotesRequest struct {
	Header               *RequestHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	NotebookId           string         `protobuf:"bytes,2,opt,name=notebookId,proto3" json:"notebookId,omitempty"`
	StoreId              string         `protobuf:"bytes,3,opt,name=storeId,proto3" json:"storeId,omitempty"`
	OwnerId              string         `protobuf:"bytes,4,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Scope                string         `protobuf:"bytes,5,opt,name=scope,proto3" json:"scope,omitempty"`
	Store                string         `protobuf:"bytes,6,opt,name=store,proto3" json:"store,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetNotesRequest) Reset()         { *m = GetNotesRequest{} }
func (m *GetNotesRequest) String() string { return proto.CompactTextString(m) }
func (*GetNotesRequest) ProtoMessage()    {}
func (*GetNotesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_640dafe07df50d4e, []int{6}
}

func (m *GetNotesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetNotesRequest.Unmarshal(m, b)
}
func (m *GetNotesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetNotesRequest.Marshal(b, m, deterministic)
}
func (m *GetNotesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetNotesRequest.Merge(m, src)
}
func (m *GetNotesRequest) XXX_Size() int {
	return xxx_messageInfo_GetNotesRequest.Size(m)
}
func (m *GetNotesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetNotesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetNotesRequest proto.InternalMessageInfo

func (m *GetNotesRequest) GetHeader() *RequestHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetNotesRequest) GetNotebookId() string {
	if m != nil {
		return m.NotebookId
	}
	return ""
}

func (m *GetNotesRequest) GetStoreId() string {
	if m != nil {
		return m.StoreId
	}
	return ""
}

func (m *GetNotesRequest) GetOwnerId() string {
	if m != nil {
		return m.OwnerId
	}
	return ""
}

func (m *GetNotesRequest) GetScope() string {
	if m != nil {
		return m.Scope
	}
	return ""
}

func (m *GetNotesRequest) GetStore() string {
	if m != nil {
		return m.Store
	}
	return ""
}

type GetNotesResponse struct {
	Header               *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Notes                []*Note         `protobuf:"bytes,2,rep,name=notes,proto3" json:"notes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *GetNotesResponse) Reset()         { *m = GetNotesResponse{} }
func (m *GetNotesResponse) String() string { return proto.CompactTextString(m) }
func (*GetNotesResponse) ProtoMessage()    {}
func (*GetNotesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_640dafe07df50d4e, []int{7}
}

func (m *GetNotesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetNotesResponse.Unmarshal(m, b)
}
func (m *GetNotesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetNotesResponse.Marshal(b, m, deterministic)
}
func (m *GetNotesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetNotesResponse.Merge(m, src)
}
func (m *GetNotesResponse) XXX_Size() int {
	return xxx_messageInfo_GetNotesResponse.Size(m)
}
func (m *GetNotesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetNotesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetNotesResponse proto.InternalMessageInfo

func (m *GetNotesResponse) GetHeader() *ResponseHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *GetNotesResponse) GetNotes() []*Note {
	if m != nil {
		return m.Notes
	}
	return nil
}

func init() {
	proto.RegisterType((*Note)(nil), "notekeeper.Note")
	proto.RegisterType((*CreateNoteRequest)(nil), "notekeeper.CreateNoteRequest")
	proto.RegisterType((*SaveNoteRequest)(nil), "notekeeper.SaveNoteRequest")
	proto.RegisterType((*DeleteNoteRequest)(nil), "notekeeper.DeleteNoteRequest")
	proto.RegisterType((*LoadNoteRequest)(nil), "notekeeper.LoadNoteRequest")
	proto.RegisterType((*LoadNoteResponse)(nil), "notekeeper.LoadNoteResponse")
	proto.RegisterType((*GetNotesRequest)(nil), "notekeeper.GetNotesRequest")
	proto.RegisterType((*GetNotesResponse)(nil), "notekeeper.GetNotesResponse")
}

func init() { proto.RegisterFile("note.proto", fileDescriptor_640dafe07df50d4e) }

var fileDescriptor_640dafe07df50d4e = []byte{
	// 451 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x55, 0xd1, 0x8a, 0xd3, 0x40,
	0x14, 0x75, 0xd2, 0x24, 0x6d, 0x6f, 0x16, 0xdb, 0x0e, 0x22, 0x63, 0x11, 0x09, 0x41, 0xa5, 0x4f,
	0x05, 0xeb, 0x27, 0x28, 0x68, 0x41, 0x64, 0x89, 0xfe, 0x40, 0xb6, 0x73, 0xc1, 0xd0, 0x34, 0x13,
	0x33, 0xb3, 0x2b, 0xfe, 0x97, 0x6f, 0xbe, 0xfa, 0x0d, 0xfe, 0x85, 0x5f, 0xe0, 0x8b, 0xcc, 0x9d,
	0x59, 0x92, 0x2e, 0xed, 0x22, 0xec, 0xfa, 0xe0, 0xbe, 0xcd, 0x3d, 0xe7, 0xcc, 0x70, 0xce, 0x99,
	0x09, 0x01, 0xa8, 0x95, 0xc1, 0x65, 0xd3, 0x2a, 0xa3, 0x38, 0xad, 0xb7, 0x88, 0x0d, 0xb6, 0xf3,
	0x93, 0x8d, 0xda, 0xed, 0x54, 0xed, 0x98, 0x79, 0x62, 0x4a, 0x53, 0x79, 0x59, 0xf6, 0x3d, 0x80,
	0xf0, 0xbd, 0x32, 0xc8, 0xef, 0x43, 0x50, 0x4a, 0xc1, 0x52, 0xb6, 0x18, 0xe7, 0x41, 0x29, 0xf9,
	0x13, 0x77, 0xda, 0x99, 0x52, 0xdb, 0xb5, 0x14, 0x01, 0xe1, 0x3d, 0x84, 0x0b, 0x18, 0xaa, 0x2f,
	0x35, 0xb6, 0x6b, 0x29, 0x06, 0x44, 0x5e, 0x8e, 0x96, 0xd1, 0x46, 0xb5, 0xb8, 0x96, 0x22, 0x74,
	0x8c, 0x1f, 0xf9, 0x03, 0x88, 0xf4, 0x46, 0x35, 0x28, 0x22, 0xc2, 0xdd, 0x40, 0xa8, 0x15, 0x88,
	0xd8, 0xa3, 0x76, 0xe0, 0xcf, 0x20, 0xac, 0x8b, 0x1d, 0x8a, 0x61, 0xca, 0x16, 0xc9, 0x6a, 0xb6,
	0xec, 0xe2, 0x2c, 0x3f, 0x5a, 0xff, 0x39, 0xd1, 0x9c, 0x43, 0x68, 0xbe, 0x36, 0x28, 0x46, 0xb4,
	0x97, 0xd6, 0xfc, 0x31, 0x8c, 0x5b, 0xbc, 0x28, 0x75, 0xa9, 0x6a, 0x2d, 0xc6, 0x29, 0x5b, 0x44,
	0x79, 0x07, 0xf0, 0x87, 0x10, 0x57, 0x6a, 0xb3, 0x45, 0x29, 0x20, 0x65, 0x8b, 0x51, 0xee, 0x27,
	0x6b, 0x7b, 0xd3, 0x62, 0x61, 0x50, 0x8a, 0xc4, 0xd9, 0xf6, 0xa3, 0x65, 0xce, 0x1b, 0x49, 0xcc,
	0x89, 0x63, 0xfc, 0x98, 0xfd, 0x62, 0x30, 0x7b, 0x45, 0x2a, 0xdb, 0x61, 0x8e, 0x9f, 0xcf, 0x51,
	0x1b, 0xfe, 0x02, 0xe2, 0x4f, 0x58, 0x48, 0x6c, 0xa9, 0xce, 0x64, 0xf5, 0xa8, 0x6f, 0xde, 0x8b,
	0xde, 0x92, 0x20, 0xf7, 0xc2, 0xbf, 0x69, 0xfb, 0xb2, 0xd3, 0xc1, 0x7e, 0xa7, 0xbd, 0x7b, 0x08,
	0xf7, 0xef, 0xe1, 0xf6, 0xdb, 0xce, 0x7e, 0x33, 0x98, 0x7c, 0x28, 0x2e, 0x6e, 0x9a, 0xd6, 0xbd,
	0xb5, 0xe0, 0xc8, 0x5b, 0x1b, 0x5c, 0x97, 0x3e, 0x3c, 0x9a, 0x3e, 0x3a, 0x92, 0x3e, 0x3e, 0x98,
	0x7e, 0x78, 0x28, 0xfd, 0xe8, 0xfa, 0xf4, 0x3f, 0x19, 0xcc, 0x5e, 0x63, 0x85, 0xe6, 0x8e, 0xe5,
	0xcf, 0xbe, 0x31, 0x98, 0xbc, 0x53, 0x85, 0xbc, 0xe5, 0x58, 0xff, 0xf8, 0xd1, 0x66, 0x15, 0x4c,
	0x3b, 0xd7, 0xba, 0x51, 0xb5, 0x46, 0xbe, 0xba, 0x62, 0x7b, 0xbe, 0x6f, 0xdb, 0xa9, 0xae, 0xf8,
	0x7e, 0x0a, 0xa1, 0x15, 0x91, 0xf3, 0x64, 0x35, 0xed, 0xef, 0xa0, 0xb3, 0x89, 0xcd, 0x7e, 0x30,
	0x98, 0xbc, 0x41, 0x63, 0x11, 0xfd, 0xff, 0x7e, 0xe9, 0x59, 0x0d, 0xd3, 0x2e, 0xc5, 0x0d, 0x4a,
	0x7b, 0x0e, 0x91, 0x15, 0x69, 0x11, 0xa4, 0x83, 0x83, 0xad, 0x39, 0xfa, 0xf4, 0xde, 0x29, 0x3b,
	0x8b, 0xe9, 0x5f, 0xf3, 0xf2, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5f, 0x9b, 0x11, 0x9b, 0xa0,
	0x06, 0x00, 0x00,
}
