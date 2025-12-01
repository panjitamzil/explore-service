package explorepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	ExploreService_ListLikedYou_FullMethodName    = "/explore.ExploreService/ListLikedYou"
	ExploreService_ListNewLikedYou_FullMethodName = "/explore.ExploreService/ListNewLikedYou"
	ExploreService_CountLikedYou_FullMethodName   = "/explore.ExploreService/CountLikedYou"
	ExploreService_PutDecision_FullMethodName     = "/explore.ExploreService/PutDecision"
)

type ExploreServiceClient interface {
	ListLikedYou(ctx context.Context, in *ListLikedYouRequest, opts ...grpc.CallOption) (*ListLikedYouResponse, error)
	ListNewLikedYou(ctx context.Context, in *ListLikedYouRequest, opts ...grpc.CallOption) (*ListLikedYouResponse, error)
	CountLikedYou(ctx context.Context, in *CountLikedYouRequest, opts ...grpc.CallOption) (*CountLikedYouResponse, error)
	PutDecision(ctx context.Context, in *PutDecisionRequest, opts ...grpc.CallOption) (*PutDecisionResponse, error)
}

type exploreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExploreServiceClient(cc grpc.ClientConnInterface) ExploreServiceClient {
	return &exploreServiceClient{cc}
}

func (c *exploreServiceClient) ListLikedYou(ctx context.Context, in *ListLikedYouRequest, opts ...grpc.CallOption) (*ListLikedYouResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListLikedYouResponse)
	err := c.cc.Invoke(ctx, ExploreService_ListLikedYou_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exploreServiceClient) ListNewLikedYou(ctx context.Context, in *ListLikedYouRequest, opts ...grpc.CallOption) (*ListLikedYouResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListLikedYouResponse)
	err := c.cc.Invoke(ctx, ExploreService_ListNewLikedYou_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exploreServiceClient) CountLikedYou(ctx context.Context, in *CountLikedYouRequest, opts ...grpc.CallOption) (*CountLikedYouResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CountLikedYouResponse)
	err := c.cc.Invoke(ctx, ExploreService_CountLikedYou_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exploreServiceClient) PutDecision(ctx context.Context, in *PutDecisionRequest, opts ...grpc.CallOption) (*PutDecisionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PutDecisionResponse)
	err := c.cc.Invoke(ctx, ExploreService_PutDecision_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ExploreServiceServer interface {
	ListLikedYou(context.Context, *ListLikedYouRequest) (*ListLikedYouResponse, error)
	ListNewLikedYou(context.Context, *ListLikedYouRequest) (*ListLikedYouResponse, error)
	CountLikedYou(context.Context, *CountLikedYouRequest) (*CountLikedYouResponse, error)
	PutDecision(context.Context, *PutDecisionRequest) (*PutDecisionResponse, error)
	mustEmbedUnimplementedExploreServiceServer()
}

type UnimplementedExploreServiceServer struct{}

func (UnimplementedExploreServiceServer) ListLikedYou(context.Context, *ListLikedYouRequest) (*ListLikedYouResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLikedYou not implemented")
}
func (UnimplementedExploreServiceServer) ListNewLikedYou(context.Context, *ListLikedYouRequest) (*ListLikedYouResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNewLikedYou not implemented")
}
func (UnimplementedExploreServiceServer) CountLikedYou(context.Context, *CountLikedYouRequest) (*CountLikedYouResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CountLikedYou not implemented")
}
func (UnimplementedExploreServiceServer) PutDecision(context.Context, *PutDecisionRequest) (*PutDecisionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutDecision not implemented")
}
func (UnimplementedExploreServiceServer) mustEmbedUnimplementedExploreServiceServer() {}
func (UnimplementedExploreServiceServer) testEmbeddedByValue()                        {}

type UnsafeExploreServiceServer interface {
	mustEmbedUnimplementedExploreServiceServer()
}

func RegisterExploreServiceServer(s grpc.ServiceRegistrar, srv ExploreServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ExploreService_ServiceDesc, srv)
}

func _ExploreService_ListLikedYou_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLikedYouRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExploreServiceServer).ListLikedYou(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExploreService_ListLikedYou_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExploreServiceServer).ListLikedYou(ctx, req.(*ListLikedYouRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExploreService_ListNewLikedYou_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLikedYouRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExploreServiceServer).ListNewLikedYou(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExploreService_ListNewLikedYou_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExploreServiceServer).ListNewLikedYou(ctx, req.(*ListLikedYouRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExploreService_CountLikedYou_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CountLikedYouRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExploreServiceServer).CountLikedYou(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExploreService_CountLikedYou_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExploreServiceServer).CountLikedYou(ctx, req.(*CountLikedYouRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExploreService_PutDecision_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutDecisionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExploreServiceServer).PutDecision(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExploreService_PutDecision_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExploreServiceServer).PutDecision(ctx, req.(*PutDecisionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var ExploreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "explore.ExploreService",
	HandlerType: (*ExploreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListLikedYou",
			Handler:    _ExploreService_ListLikedYou_Handler,
		},
		{
			MethodName: "ListNewLikedYou",
			Handler:    _ExploreService_ListNewLikedYou_Handler,
		},
		{
			MethodName: "CountLikedYou",
			Handler:    _ExploreService_CountLikedYou_Handler,
		},
		{
			MethodName: "PutDecision",
			Handler:    _ExploreService_PutDecision_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/explore-service.proto",
}
