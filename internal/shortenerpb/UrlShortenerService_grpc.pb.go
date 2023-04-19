// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: shortenerpb/UrlShortenerService.proto

package shortener

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

// UrlShortenerClient is the client API for UrlShortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UrlShortenerClient interface {
	ShortenUrl(ctx context.Context, in *ShortenURLRequest, opts ...grpc.CallOption) (*ShortenURLResponse, error)
	ExtractUrl(ctx context.Context, in *ExtractURLRequest, opts ...grpc.CallOption) (*ExtractURLResponse, error)
}

type urlShortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewUrlShortenerClient(cc grpc.ClientConnInterface) UrlShortenerClient {
	return &urlShortenerClient{cc}
}

func (c *urlShortenerClient) ShortenUrl(ctx context.Context, in *ShortenURLRequest, opts ...grpc.CallOption) (*ShortenURLResponse, error) {
	out := new(ShortenURLResponse)
	err := c.cc.Invoke(ctx, "/shortener.UrlShortener/ShortenUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *urlShortenerClient) ExtractUrl(ctx context.Context, in *ExtractURLRequest, opts ...grpc.CallOption) (*ExtractURLResponse, error) {
	out := new(ExtractURLResponse)
	err := c.cc.Invoke(ctx, "/shortener.UrlShortener/ExtractUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UrlShortenerServer is the server API for UrlShortener service.
// All implementations must embed UnimplementedUrlShortenerServer
// for forward compatibility
type UrlShortenerServer interface {
	ShortenUrl(context.Context, *ShortenURLRequest) (*ShortenURLResponse, error)
	ExtractUrl(context.Context, *ExtractURLRequest) (*ExtractURLResponse, error)
	mustEmbedUnimplementedUrlShortenerServer()
}

// UnimplementedUrlShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedUrlShortenerServer struct {
}

func (UnimplementedUrlShortenerServer) ShortenUrl(context.Context, *ShortenURLRequest) (*ShortenURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenUrl not implemented")
}
func (UnimplementedUrlShortenerServer) ExtractUrl(context.Context, *ExtractURLRequest) (*ExtractURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExtractUrl not implemented")
}
func (UnimplementedUrlShortenerServer) mustEmbedUnimplementedUrlShortenerServer() {}

// UnsafeUrlShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UrlShortenerServer will
// result in compilation errors.
type UnsafeUrlShortenerServer interface {
	mustEmbedUnimplementedUrlShortenerServer()
}

func RegisterUrlShortenerServer(s grpc.ServiceRegistrar, srv UrlShortenerServer) {
	s.RegisterService(&UrlShortener_ServiceDesc, srv)
}

func _UrlShortener_ShortenUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlShortenerServer).ShortenUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.UrlShortener/ShortenUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlShortenerServer).ShortenUrl(ctx, req.(*ShortenURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UrlShortener_ExtractUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtractURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlShortenerServer).ExtractUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.UrlShortener/ExtractUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlShortenerServer).ExtractUrl(ctx, req.(*ExtractURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UrlShortener_ServiceDesc is the grpc.ServiceDesc for UrlShortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UrlShortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.UrlShortener",
	HandlerType: (*UrlShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShortenUrl",
			Handler:    _UrlShortener_ShortenUrl_Handler,
		},
		{
			MethodName: "ExtractUrl",
			Handler:    _UrlShortener_ExtractUrl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shortenerpb/UrlShortenerService.proto",
}
