// Code generated by protoc-gen-go.
// source: backend.proto
// DO NOT EDIT!

/*
Package notekeeper is a generated protocol buffer package.

It is generated from these files:
	backend.proto

It has these top-level messages:
	UIStateRequest
	UIStateResponse
	AccountStateRequest
	AccountStateResponse
	OpenMasterDbRequest
	OpenMasterDbResponse
	CreateAccountRequest
	CreateAccountResponse
	UnlockAccountRequest
	UnlockAccountResponse
	SigninAccountRequest
	SigninAccountResponse
	SignoutAccountRequest
	SignoutAccountResponse
	LockAccountRequest
	LockAccountResponse
	CreateNotebookRequest
	CreateNotebookResponse
*/
package notekeeper

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type UIStateRequest struct {
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *UIStateRequest) Reset()                    { *m = UIStateRequest{} }
func (m *UIStateRequest) String() string            { return proto.CompactTextString(m) }
func (*UIStateRequest) ProtoMessage()               {}
func (*UIStateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *UIStateRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type UIStateResponse struct {
	WindowWidth  int32 `protobuf:"varint,1,opt,name=windowWidth" json:"windowWidth,omitempty"`
	WindowHeight int32 `protobuf:"varint,2,opt,name=windowHeight" json:"windowHeight,omitempty"`
}

func (m *UIStateResponse) Reset()                    { *m = UIStateResponse{} }
func (m *UIStateResponse) String() string            { return proto.CompactTextString(m) }
func (*UIStateResponse) ProtoMessage()               {}
func (*UIStateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UIStateResponse) GetWindowWidth() int32 {
	if m != nil {
		return m.WindowWidth
	}
	return 0
}

func (m *UIStateResponse) GetWindowHeight() int32 {
	if m != nil {
		return m.WindowHeight
	}
	return 0
}

type AccountStateRequest struct {
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *AccountStateRequest) Reset()                    { *m = AccountStateRequest{} }
func (m *AccountStateRequest) String() string            { return proto.CompactTextString(m) }
func (*AccountStateRequest) ProtoMessage()               {}
func (*AccountStateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AccountStateRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type AccountStateResponse struct {
	SignedIn bool `protobuf:"varint,1,opt,name=signedIn" json:"signedIn,omitempty"`
	Locked   bool `protobuf:"varint,2,opt,name=locked" json:"locked,omitempty"`
}

func (m *AccountStateResponse) Reset()                    { *m = AccountStateResponse{} }
func (m *AccountStateResponse) String() string            { return proto.CompactTextString(m) }
func (*AccountStateResponse) ProtoMessage()               {}
func (*AccountStateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AccountStateResponse) GetSignedIn() bool {
	if m != nil {
		return m.SignedIn
	}
	return false
}

func (m *AccountStateResponse) GetLocked() bool {
	if m != nil {
		return m.Locked
	}
	return false
}

type OpenMasterDbRequest struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *OpenMasterDbRequest) Reset()                    { *m = OpenMasterDbRequest{} }
func (m *OpenMasterDbRequest) String() string            { return proto.CompactTextString(m) }
func (*OpenMasterDbRequest) ProtoMessage()               {}
func (*OpenMasterDbRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *OpenMasterDbRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type OpenMasterDbResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *OpenMasterDbResponse) Reset()                    { *m = OpenMasterDbResponse{} }
func (m *OpenMasterDbResponse) String() string            { return proto.CompactTextString(m) }
func (*OpenMasterDbResponse) ProtoMessage()               {}
func (*OpenMasterDbResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *OpenMasterDbResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type CreateAccountRequest struct {
	Name       string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Email      string `protobuf:"bytes,2,opt,name=email" json:"email,omitempty"`
	Passphrase string `protobuf:"bytes,3,opt,name=passphrase" json:"passphrase,omitempty"`
}

func (m *CreateAccountRequest) Reset()                    { *m = CreateAccountRequest{} }
func (m *CreateAccountRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateAccountRequest) ProtoMessage()               {}
func (*CreateAccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CreateAccountRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateAccountRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateAccountRequest) GetPassphrase() string {
	if m != nil {
		return m.Passphrase
	}
	return ""
}

type CreateAccountResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Id     string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
}

func (m *CreateAccountResponse) Reset()                    { *m = CreateAccountResponse{} }
func (m *CreateAccountResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateAccountResponse) ProtoMessage()               {}
func (*CreateAccountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *CreateAccountResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *CreateAccountResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type UnlockAccountRequest struct {
	Id         string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Passphrase string `protobuf:"bytes,2,opt,name=passphrase" json:"passphrase,omitempty"`
}

func (m *UnlockAccountRequest) Reset()                    { *m = UnlockAccountRequest{} }
func (m *UnlockAccountRequest) String() string            { return proto.CompactTextString(m) }
func (*UnlockAccountRequest) ProtoMessage()               {}
func (*UnlockAccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *UnlockAccountRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UnlockAccountRequest) GetPassphrase() string {
	if m != nil {
		return m.Passphrase
	}
	return ""
}

type UnlockAccountResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *UnlockAccountResponse) Reset()                    { *m = UnlockAccountResponse{} }
func (m *UnlockAccountResponse) String() string            { return proto.CompactTextString(m) }
func (*UnlockAccountResponse) ProtoMessage()               {}
func (*UnlockAccountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *UnlockAccountResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type SigninAccountRequest struct {
	Name       string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Email      string `protobuf:"bytes,2,opt,name=email" json:"email,omitempty"`
	Passphrase string `protobuf:"bytes,3,opt,name=passphrase" json:"passphrase,omitempty"`
}

func (m *SigninAccountRequest) Reset()                    { *m = SigninAccountRequest{} }
func (m *SigninAccountRequest) String() string            { return proto.CompactTextString(m) }
func (*SigninAccountRequest) ProtoMessage()               {}
func (*SigninAccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *SigninAccountRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *SigninAccountRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *SigninAccountRequest) GetPassphrase() string {
	if m != nil {
		return m.Passphrase
	}
	return ""
}

type SigninAccountResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Id     string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
}

func (m *SigninAccountResponse) Reset()                    { *m = SigninAccountResponse{} }
func (m *SigninAccountResponse) String() string            { return proto.CompactTextString(m) }
func (*SigninAccountResponse) ProtoMessage()               {}
func (*SigninAccountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *SigninAccountResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *SigninAccountResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type SignoutAccountRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *SignoutAccountRequest) Reset()                    { *m = SignoutAccountRequest{} }
func (m *SignoutAccountRequest) String() string            { return proto.CompactTextString(m) }
func (*SignoutAccountRequest) ProtoMessage()               {}
func (*SignoutAccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *SignoutAccountRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type SignoutAccountResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *SignoutAccountResponse) Reset()                    { *m = SignoutAccountResponse{} }
func (m *SignoutAccountResponse) String() string            { return proto.CompactTextString(m) }
func (*SignoutAccountResponse) ProtoMessage()               {}
func (*SignoutAccountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *SignoutAccountResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type LockAccountRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *LockAccountRequest) Reset()                    { *m = LockAccountRequest{} }
func (m *LockAccountRequest) String() string            { return proto.CompactTextString(m) }
func (*LockAccountRequest) ProtoMessage()               {}
func (*LockAccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *LockAccountRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type LockAccountResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
}

func (m *LockAccountResponse) Reset()                    { *m = LockAccountResponse{} }
func (m *LockAccountResponse) String() string            { return proto.CompactTextString(m) }
func (*LockAccountResponse) ProtoMessage()               {}
func (*LockAccountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *LockAccountResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type CreateNotebookRequest struct {
	Name    string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	ShelfId string `protobuf:"bytes,3,opt,name=shelf_id,json=shelfId" json:"shelf_id,omitempty"`
}

func (m *CreateNotebookRequest) Reset()                    { *m = CreateNotebookRequest{} }
func (m *CreateNotebookRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateNotebookRequest) ProtoMessage()               {}
func (*CreateNotebookRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *CreateNotebookRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateNotebookRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *CreateNotebookRequest) GetShelfId() string {
	if m != nil {
		return m.ShelfId
	}
	return ""
}

type CreateNotebookResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Id     string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
}

func (m *CreateNotebookResponse) Reset()                    { *m = CreateNotebookResponse{} }
func (m *CreateNotebookResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateNotebookResponse) ProtoMessage()               {}
func (*CreateNotebookResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *CreateNotebookResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *CreateNotebookResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*UIStateRequest)(nil), "notekeeper.UIStateRequest")
	proto.RegisterType((*UIStateResponse)(nil), "notekeeper.UIStateResponse")
	proto.RegisterType((*AccountStateRequest)(nil), "notekeeper.AccountStateRequest")
	proto.RegisterType((*AccountStateResponse)(nil), "notekeeper.AccountStateResponse")
	proto.RegisterType((*OpenMasterDbRequest)(nil), "notekeeper.OpenMasterDbRequest")
	proto.RegisterType((*OpenMasterDbResponse)(nil), "notekeeper.OpenMasterDbResponse")
	proto.RegisterType((*CreateAccountRequest)(nil), "notekeeper.CreateAccountRequest")
	proto.RegisterType((*CreateAccountResponse)(nil), "notekeeper.CreateAccountResponse")
	proto.RegisterType((*UnlockAccountRequest)(nil), "notekeeper.UnlockAccountRequest")
	proto.RegisterType((*UnlockAccountResponse)(nil), "notekeeper.UnlockAccountResponse")
	proto.RegisterType((*SigninAccountRequest)(nil), "notekeeper.SigninAccountRequest")
	proto.RegisterType((*SigninAccountResponse)(nil), "notekeeper.SigninAccountResponse")
	proto.RegisterType((*SignoutAccountRequest)(nil), "notekeeper.SignoutAccountRequest")
	proto.RegisterType((*SignoutAccountResponse)(nil), "notekeeper.SignoutAccountResponse")
	proto.RegisterType((*LockAccountRequest)(nil), "notekeeper.LockAccountRequest")
	proto.RegisterType((*LockAccountResponse)(nil), "notekeeper.LockAccountResponse")
	proto.RegisterType((*CreateNotebookRequest)(nil), "notekeeper.CreateNotebookRequest")
	proto.RegisterType((*CreateNotebookResponse)(nil), "notekeeper.CreateNotebookResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Backend service

type BackendClient interface {
	OpenMasterDb(ctx context.Context, in *OpenMasterDbRequest, opts ...grpc.CallOption) (*OpenMasterDbResponse, error)
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
	UnlockAccount(ctx context.Context, in *UnlockAccountRequest, opts ...grpc.CallOption) (*UnlockAccountResponse, error)
	SigninAccount(ctx context.Context, in *SigninAccountRequest, opts ...grpc.CallOption) (*SigninAccountResponse, error)
	SignoutAccount(ctx context.Context, in *SignoutAccountRequest, opts ...grpc.CallOption) (*SignoutAccountResponse, error)
	LockAccount(ctx context.Context, in *LockAccountRequest, opts ...grpc.CallOption) (*LockAccountResponse, error)
	UIState(ctx context.Context, in *UIStateRequest, opts ...grpc.CallOption) (*UIStateResponse, error)
	AccountState(ctx context.Context, in *AccountStateRequest, opts ...grpc.CallOption) (*AccountStateResponse, error)
	CreateNotebook(ctx context.Context, in *CreateNotebookRequest, opts ...grpc.CallOption) (*CreateNotebookResponse, error)
}

type backendClient struct {
	cc *grpc.ClientConn
}

func NewBackendClient(cc *grpc.ClientConn) BackendClient {
	return &backendClient{cc}
}

func (c *backendClient) OpenMasterDb(ctx context.Context, in *OpenMasterDbRequest, opts ...grpc.CallOption) (*OpenMasterDbResponse, error) {
	out := new(OpenMasterDbResponse)
	err := grpc.Invoke(ctx, "/notekeeper.Backend/OpenMasterDb", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	out := new(CreateAccountResponse)
	err := grpc.Invoke(ctx, "/notekeeper.Backend/CreateAccount", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) UnlockAccount(ctx context.Context, in *UnlockAccountRequest, opts ...grpc.CallOption) (*UnlockAccountResponse, error) {
	out := new(UnlockAccountResponse)
	err := grpc.Invoke(ctx, "/notekeeper.Backend/UnlockAccount", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) SigninAccount(ctx context.Context, in *SigninAccountRequest, opts ...grpc.CallOption) (*SigninAccountResponse, error) {
	out := new(SigninAccountResponse)
	err := grpc.Invoke(ctx, "/notekeeper.Backend/SigninAccount", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) SignoutAccount(ctx context.Context, in *SignoutAccountRequest, opts ...grpc.CallOption) (*SignoutAccountResponse, error) {
	out := new(SignoutAccountResponse)
	err := grpc.Invoke(ctx, "/notekeeper.Backend/SignoutAccount", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) LockAccount(ctx context.Context, in *LockAccountRequest, opts ...grpc.CallOption) (*LockAccountResponse, error) {
	out := new(LockAccountResponse)
	err := grpc.Invoke(ctx, "/notekeeper.Backend/LockAccount", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) UIState(ctx context.Context, in *UIStateRequest, opts ...grpc.CallOption) (*UIStateResponse, error) {
	out := new(UIStateResponse)
	err := grpc.Invoke(ctx, "/notekeeper.Backend/UIState", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) AccountState(ctx context.Context, in *AccountStateRequest, opts ...grpc.CallOption) (*AccountStateResponse, error) {
	out := new(AccountStateResponse)
	err := grpc.Invoke(ctx, "/notekeeper.Backend/AccountState", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) CreateNotebook(ctx context.Context, in *CreateNotebookRequest, opts ...grpc.CallOption) (*CreateNotebookResponse, error) {
	out := new(CreateNotebookResponse)
	err := grpc.Invoke(ctx, "/notekeeper.Backend/CreateNotebook", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Backend service

type BackendServer interface {
	OpenMasterDb(context.Context, *OpenMasterDbRequest) (*OpenMasterDbResponse, error)
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
	UnlockAccount(context.Context, *UnlockAccountRequest) (*UnlockAccountResponse, error)
	SigninAccount(context.Context, *SigninAccountRequest) (*SigninAccountResponse, error)
	SignoutAccount(context.Context, *SignoutAccountRequest) (*SignoutAccountResponse, error)
	LockAccount(context.Context, *LockAccountRequest) (*LockAccountResponse, error)
	UIState(context.Context, *UIStateRequest) (*UIStateResponse, error)
	AccountState(context.Context, *AccountStateRequest) (*AccountStateResponse, error)
	CreateNotebook(context.Context, *CreateNotebookRequest) (*CreateNotebookResponse, error)
}

func RegisterBackendServer(s *grpc.Server, srv BackendServer) {
	s.RegisterService(&_Backend_serviceDesc, srv)
}

func _Backend_OpenMasterDb_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenMasterDbRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).OpenMasterDb(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notekeeper.Backend/OpenMasterDb",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).OpenMasterDb(ctx, req.(*OpenMasterDbRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notekeeper.Backend/CreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_UnlockAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnlockAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).UnlockAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notekeeper.Backend/UnlockAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).UnlockAccount(ctx, req.(*UnlockAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_SigninAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SigninAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).SigninAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notekeeper.Backend/SigninAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).SigninAccount(ctx, req.(*SigninAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_SignoutAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignoutAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).SignoutAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notekeeper.Backend/SignoutAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).SignoutAccount(ctx, req.(*SignoutAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_LockAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LockAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).LockAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notekeeper.Backend/LockAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).LockAccount(ctx, req.(*LockAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_UIState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UIStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).UIState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notekeeper.Backend/UIState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).UIState(ctx, req.(*UIStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_AccountState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).AccountState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notekeeper.Backend/AccountState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).AccountState(ctx, req.(*AccountStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_CreateNotebook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNotebookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).CreateNotebook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/notekeeper.Backend/CreateNotebook",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).CreateNotebook(ctx, req.(*CreateNotebookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Backend_serviceDesc = grpc.ServiceDesc{
	ServiceName: "notekeeper.Backend",
	HandlerType: (*BackendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OpenMasterDb",
			Handler:    _Backend_OpenMasterDb_Handler,
		},
		{
			MethodName: "CreateAccount",
			Handler:    _Backend_CreateAccount_Handler,
		},
		{
			MethodName: "UnlockAccount",
			Handler:    _Backend_UnlockAccount_Handler,
		},
		{
			MethodName: "SigninAccount",
			Handler:    _Backend_SigninAccount_Handler,
		},
		{
			MethodName: "SignoutAccount",
			Handler:    _Backend_SignoutAccount_Handler,
		},
		{
			MethodName: "LockAccount",
			Handler:    _Backend_LockAccount_Handler,
		},
		{
			MethodName: "UIState",
			Handler:    _Backend_UIState_Handler,
		},
		{
			MethodName: "AccountState",
			Handler:    _Backend_AccountState_Handler,
		},
		{
			MethodName: "CreateNotebook",
			Handler:    _Backend_CreateNotebook_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backend.proto",
}

func init() { proto.RegisterFile("backend.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 576 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x55, 0xdb, 0x6e, 0xd3, 0x40,
	0x10, 0x4d, 0x02, 0xcd, 0x65, 0x9a, 0x06, 0x69, 0x93, 0x86, 0x60, 0xa4, 0x36, 0xac, 0x10, 0x17,
	0x21, 0x02, 0x82, 0x0f, 0xe0, 0x56, 0x21, 0x82, 0xca, 0x45, 0x8e, 0x4a, 0xc5, 0x53, 0xd9, 0xc4,
	0x43, 0x62, 0x25, 0xdd, 0x35, 0xde, 0xb5, 0xfa, 0x0d, 0xfc, 0x35, 0xb2, 0xd7, 0x31, 0x5e, 0x67,
	0xeb, 0x90, 0x07, 0xde, 0x3c, 0xbb, 0x67, 0xcf, 0xcc, 0x1c, 0x9d, 0xf1, 0xc0, 0xc1, 0x94, 0xcd,
	0x96, 0xc8, 0xbd, 0x51, 0x10, 0x0a, 0x25, 0x08, 0x70, 0xa1, 0x70, 0x89, 0x18, 0x60, 0x48, 0x1f,
	0x40, 0xe7, 0x6c, 0x3c, 0x51, 0x4c, 0xa1, 0x8b, 0xbf, 0x22, 0x94, 0x8a, 0xf4, 0x60, 0x4f, 0x89,
	0x25, 0xf2, 0x41, 0x75, 0x58, 0x7d, 0xd4, 0x72, 0x75, 0x40, 0xcf, 0xe1, 0x56, 0x86, 0x93, 0x81,
	0xe0, 0x12, 0xc9, 0x10, 0xf6, 0xaf, 0x7c, 0xee, 0x89, 0xab, 0x73, 0xdf, 0x53, 0x8b, 0x04, 0xbe,
	0xe7, 0xe6, 0x8f, 0x08, 0x85, 0xb6, 0x0e, 0x3f, 0xa0, 0x3f, 0x5f, 0xa8, 0x41, 0x2d, 0x81, 0x18,
	0x67, 0xf4, 0x09, 0x74, 0xdf, 0xcc, 0x66, 0x22, 0xe2, 0xea, 0x1f, 0xaa, 0xf8, 0x08, 0x3d, 0x13,
	0x9c, 0x96, 0xe2, 0x40, 0x53, 0xfa, 0x73, 0x8e, 0xde, 0x58, 0x3f, 0x68, 0xba, 0x59, 0x4c, 0xfa,
	0x50, 0x5f, 0x89, 0xd9, 0x12, 0xbd, 0x24, 0x7d, 0xd3, 0x4d, 0x23, 0xfa, 0x18, 0xba, 0x5f, 0x02,
	0xe4, 0x9f, 0x98, 0x54, 0x18, 0x9e, 0x4c, 0xd7, 0x89, 0x09, 0xdc, 0x0c, 0x58, 0xda, 0x4e, 0xcb,
	0x4d, 0xbe, 0xe9, 0x08, 0x7a, 0x26, 0x34, 0x4d, 0xdb, 0x87, 0xba, 0x54, 0x4c, 0x45, 0x32, 0x45,
	0xa7, 0x11, 0xfd, 0x01, 0xbd, 0x77, 0x21, 0x32, 0x85, 0x69, 0xb1, 0x39, 0x6e, 0xce, 0x2e, 0x71,
	0xcd, 0x1d, 0x7f, 0xc7, 0x8d, 0xe2, 0x25, 0xf3, 0x57, 0x49, 0x75, 0x2d, 0x57, 0x07, 0xe4, 0x08,
	0x20, 0x60, 0x52, 0x06, 0x8b, 0x90, 0x49, 0x1c, 0xdc, 0x48, 0xae, 0x72, 0x27, 0xf4, 0x15, 0x1c,
	0x16, 0x32, 0x94, 0x97, 0x44, 0x3a, 0x50, 0xf3, 0xbd, 0x34, 0x47, 0xcd, 0xf7, 0xe8, 0x7b, 0xe8,
	0x9d, 0xf1, 0x58, 0x89, 0x42, 0x89, 0x1a, 0x57, 0x5d, 0xe3, 0x0a, 0x85, 0xd4, 0x36, 0x0a, 0x79,
	0x06, 0x87, 0x05, 0x9e, 0xed, 0xda, 0x4c, 0xfc, 0x39, 0xf7, 0xf9, 0xff, 0xd4, 0xa6, 0x90, 0x61,
	0x47, 0x6d, 0x1e, 0x6a, 0x02, 0x11, 0xa9, 0x72, 0x71, 0xe8, 0x73, 0xe8, 0x17, 0x81, 0x5b, 0xba,
	0xbf, 0x0f, 0xe4, 0x74, 0xab, 0xe8, 0xf4, 0x29, 0x74, 0x4f, 0x77, 0x90, 0xf4, 0x62, 0x6d, 0x86,
	0xcf, 0x42, 0xe1, 0x54, 0x88, 0x65, 0x99, 0xa6, 0xb7, 0xa1, 0x11, 0x49, 0x0c, 0x2f, 0xb2, 0x8e,
	0xeb, 0x71, 0x38, 0xf6, 0xc8, 0x1d, 0x68, 0xca, 0x05, 0xae, 0x7e, 0xc6, 0x37, 0x5a, 0xd4, 0x46,
	0x12, 0x8f, 0x3d, 0xfa, 0x1a, 0xfa, 0xc5, 0x04, 0xbb, 0x49, 0xfa, 0xe2, 0x77, 0x1d, 0x1a, 0x6f,
	0xf5, 0x4f, 0x88, 0x4c, 0xa0, 0x9d, 0x9f, 0x26, 0x72, 0x3c, 0xfa, 0xfb, 0x3f, 0x1a, 0x59, 0x46,
	0xd2, 0x19, 0x5e, 0x0f, 0xd0, 0x65, 0xd0, 0x0a, 0xf9, 0x06, 0x07, 0xc6, 0x40, 0x10, 0xe3, 0x91,
	0x6d, 0x1a, 0x9d, 0x7b, 0x25, 0x88, 0x3c, 0xaf, 0xe1, 0x6f, 0x93, 0xd7, 0x36, 0x42, 0x26, 0xaf,
	0x75, 0x38, 0x34, 0xaf, 0x61, 0x52, 0x93, 0xd7, 0x36, 0x21, 0x26, 0xaf, 0xd5, 0xe1, 0xb4, 0x42,
	0xbe, 0x43, 0xc7, 0xb4, 0x24, 0xd9, 0x78, 0xb6, 0xe1, 0x6b, 0x87, 0x96, 0x41, 0x32, 0xea, 0xaf,
	0xb0, 0x9f, 0x73, 0x25, 0x39, 0xca, 0x3f, 0xda, 0x34, 0xb5, 0x73, 0x7c, 0xed, 0x7d, 0xc6, 0x78,
	0x02, 0x8d, 0x74, 0xa9, 0x10, 0xc7, 0x10, 0xcd, 0xd8, 0x48, 0xce, 0x5d, 0xeb, 0x5d, 0xc6, 0x32,
	0x81, 0x76, 0x7e, 0x29, 0x98, 0x7e, 0xb2, 0xec, 0x16, 0xd3, 0x4f, 0xb6, 0x7d, 0xa2, 0x75, 0x34,
	0x2d, 0x4f, 0x2c, 0x76, 0x29, 0xcc, 0x9b, 0xa9, 0xa3, 0x7d, 0x62, 0x68, 0x65, 0x5a, 0x4f, 0xb6,
	0xf0, 0xcb, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xba, 0x48, 0x44, 0x07, 0x96, 0x07, 0x00, 0x00,
}
