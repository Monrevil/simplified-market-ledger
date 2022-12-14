// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LedgerClient is the client API for Ledger service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LedgerClient interface {
	NewIssuer(ctx context.Context, in *NewIssuerReq, opts ...grpc.CallOption) (*NewIssuerResp, error)
	SellInvoice(ctx context.Context, in *SellInvoiceReq, opts ...grpc.CallOption) (*SellInvoiceResp, error)
	GetInvoice(ctx context.Context, in *GetInvoiceReq, opts ...grpc.CallOption) (*Invoice, error)
	ListInvoices(ctx context.Context, in *ListInvoicesReq, opts ...grpc.CallOption) (*ListInvoicesResp, error)
	NewInvestor(ctx context.Context, in *NewInvestorReq, opts ...grpc.CallOption) (*NewInvestorResp, error)
	GetInvestor(ctx context.Context, in *GetInvestorReq, opts ...grpc.CallOption) (*Investor, error)
	PlaceBid(ctx context.Context, in *PlaceBidReq, opts ...grpc.CallOption) (*PlaceBidResp, error)
	ApproveFinancing(ctx context.Context, in *ApproveReq, opts ...grpc.CallOption) (*ApproveResp, error)
	ReverseFinancing(ctx context.Context, in *ReverseReq, opts ...grpc.CallOption) (*ReverseResp, error)
	ListInvestors(ctx context.Context, in *ListInvestorsReq, opts ...grpc.CallOption) (*ListInvestorsResp, error)
}

type ledgerClient struct {
	cc grpc.ClientConnInterface
}

func NewLedgerClient(cc grpc.ClientConnInterface) LedgerClient {
	return &ledgerClient{cc}
}

func (c *ledgerClient) NewIssuer(ctx context.Context, in *NewIssuerReq, opts ...grpc.CallOption) (*NewIssuerResp, error) {
	out := new(NewIssuerResp)
	err := c.cc.Invoke(ctx, "/api.Ledger/NewIssuer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) SellInvoice(ctx context.Context, in *SellInvoiceReq, opts ...grpc.CallOption) (*SellInvoiceResp, error) {
	out := new(SellInvoiceResp)
	err := c.cc.Invoke(ctx, "/api.Ledger/SellInvoice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) GetInvoice(ctx context.Context, in *GetInvoiceReq, opts ...grpc.CallOption) (*Invoice, error) {
	out := new(Invoice)
	err := c.cc.Invoke(ctx, "/api.Ledger/GetInvoice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) ListInvoices(ctx context.Context, in *ListInvoicesReq, opts ...grpc.CallOption) (*ListInvoicesResp, error) {
	out := new(ListInvoicesResp)
	err := c.cc.Invoke(ctx, "/api.Ledger/ListInvoices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) NewInvestor(ctx context.Context, in *NewInvestorReq, opts ...grpc.CallOption) (*NewInvestorResp, error) {
	out := new(NewInvestorResp)
	err := c.cc.Invoke(ctx, "/api.Ledger/NewInvestor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) GetInvestor(ctx context.Context, in *GetInvestorReq, opts ...grpc.CallOption) (*Investor, error) {
	out := new(Investor)
	err := c.cc.Invoke(ctx, "/api.Ledger/GetInvestor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) PlaceBid(ctx context.Context, in *PlaceBidReq, opts ...grpc.CallOption) (*PlaceBidResp, error) {
	out := new(PlaceBidResp)
	err := c.cc.Invoke(ctx, "/api.Ledger/PlaceBid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) ApproveFinancing(ctx context.Context, in *ApproveReq, opts ...grpc.CallOption) (*ApproveResp, error) {
	out := new(ApproveResp)
	err := c.cc.Invoke(ctx, "/api.Ledger/ApproveFinancing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) ReverseFinancing(ctx context.Context, in *ReverseReq, opts ...grpc.CallOption) (*ReverseResp, error) {
	out := new(ReverseResp)
	err := c.cc.Invoke(ctx, "/api.Ledger/ReverseFinancing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) ListInvestors(ctx context.Context, in *ListInvestorsReq, opts ...grpc.CallOption) (*ListInvestorsResp, error) {
	out := new(ListInvestorsResp)
	err := c.cc.Invoke(ctx, "/api.Ledger/ListInvestors", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LedgerServer is the server API for Ledger service.
// All implementations must embed UnimplementedLedgerServer
// for forward compatibility
type LedgerServer interface {
	NewIssuer(context.Context, *NewIssuerReq) (*NewIssuerResp, error)
	SellInvoice(context.Context, *SellInvoiceReq) (*SellInvoiceResp, error)
	GetInvoice(context.Context, *GetInvoiceReq) (*Invoice, error)
	ListInvoices(context.Context, *ListInvoicesReq) (*ListInvoicesResp, error)
	NewInvestor(context.Context, *NewInvestorReq) (*NewInvestorResp, error)
	GetInvestor(context.Context, *GetInvestorReq) (*Investor, error)
	PlaceBid(context.Context, *PlaceBidReq) (*PlaceBidResp, error)
	ApproveFinancing(context.Context, *ApproveReq) (*ApproveResp, error)
	ReverseFinancing(context.Context, *ReverseReq) (*ReverseResp, error)
	ListInvestors(context.Context, *ListInvestorsReq) (*ListInvestorsResp, error)
	mustEmbedUnimplementedLedgerServer()
}

// UnimplementedLedgerServer must be embedded to have forward compatible implementations.
type UnimplementedLedgerServer struct {
}

func (UnimplementedLedgerServer) NewIssuer(context.Context, *NewIssuerReq) (*NewIssuerResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewIssuer not implemented")
}
func (UnimplementedLedgerServer) SellInvoice(context.Context, *SellInvoiceReq) (*SellInvoiceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SellInvoice not implemented")
}
func (UnimplementedLedgerServer) GetInvoice(context.Context, *GetInvoiceReq) (*Invoice, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInvoice not implemented")
}
func (UnimplementedLedgerServer) ListInvoices(context.Context, *ListInvoicesReq) (*ListInvoicesResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInvoices not implemented")
}
func (UnimplementedLedgerServer) NewInvestor(context.Context, *NewInvestorReq) (*NewInvestorResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewInvestor not implemented")
}
func (UnimplementedLedgerServer) GetInvestor(context.Context, *GetInvestorReq) (*Investor, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInvestor not implemented")
}
func (UnimplementedLedgerServer) PlaceBid(context.Context, *PlaceBidReq) (*PlaceBidResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceBid not implemented")
}
func (UnimplementedLedgerServer) ApproveFinancing(context.Context, *ApproveReq) (*ApproveResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApproveFinancing not implemented")
}
func (UnimplementedLedgerServer) ReverseFinancing(context.Context, *ReverseReq) (*ReverseResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReverseFinancing not implemented")
}
func (UnimplementedLedgerServer) ListInvestors(context.Context, *ListInvestorsReq) (*ListInvestorsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInvestors not implemented")
}
func (UnimplementedLedgerServer) mustEmbedUnimplementedLedgerServer() {}

// UnsafeLedgerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LedgerServer will
// result in compilation errors.
type UnsafeLedgerServer interface {
	mustEmbedUnimplementedLedgerServer()
}

func RegisterLedgerServer(s grpc.ServiceRegistrar, srv LedgerServer) {
	s.RegisterService(&Ledger_ServiceDesc, srv)
}

func _Ledger_NewIssuer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewIssuerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).NewIssuer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/NewIssuer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).NewIssuer(ctx, req.(*NewIssuerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_SellInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellInvoiceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).SellInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/SellInvoice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).SellInvoice(ctx, req.(*SellInvoiceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_GetInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInvoiceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).GetInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/GetInvoice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).GetInvoice(ctx, req.(*GetInvoiceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_ListInvoices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInvoicesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).ListInvoices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/ListInvoices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).ListInvoices(ctx, req.(*ListInvoicesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_NewInvestor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewInvestorReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).NewInvestor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/NewInvestor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).NewInvestor(ctx, req.(*NewInvestorReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_GetInvestor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInvestorReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).GetInvestor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/GetInvestor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).GetInvestor(ctx, req.(*GetInvestorReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_PlaceBid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaceBidReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).PlaceBid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/PlaceBid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).PlaceBid(ctx, req.(*PlaceBidReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_ApproveFinancing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApproveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).ApproveFinancing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/ApproveFinancing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).ApproveFinancing(ctx, req.(*ApproveReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_ReverseFinancing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReverseReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).ReverseFinancing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/ReverseFinancing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).ReverseFinancing(ctx, req.(*ReverseReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_ListInvestors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInvestorsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).ListInvestors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/ListInvestors",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).ListInvestors(ctx, req.(*ListInvestorsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Ledger_ServiceDesc is the grpc.ServiceDesc for Ledger service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ledger_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Ledger",
	HandlerType: (*LedgerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewIssuer",
			Handler:    _Ledger_NewIssuer_Handler,
		},
		{
			MethodName: "SellInvoice",
			Handler:    _Ledger_SellInvoice_Handler,
		},
		{
			MethodName: "GetInvoice",
			Handler:    _Ledger_GetInvoice_Handler,
		},
		{
			MethodName: "ListInvoices",
			Handler:    _Ledger_ListInvoices_Handler,
		},
		{
			MethodName: "NewInvestor",
			Handler:    _Ledger_NewInvestor_Handler,
		},
		{
			MethodName: "GetInvestor",
			Handler:    _Ledger_GetInvestor_Handler,
		},
		{
			MethodName: "PlaceBid",
			Handler:    _Ledger_PlaceBid_Handler,
		},
		{
			MethodName: "ApproveFinancing",
			Handler:    _Ledger_ApproveFinancing_Handler,
		},
		{
			MethodName: "ReverseFinancing",
			Handler:    _Ledger_ReverseFinancing_Handler,
		},
		{
			MethodName: "ListInvestors",
			Handler:    _Ledger_ListInvestors_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/api.proto",
}
