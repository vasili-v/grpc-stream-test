// Code generated by protoc-gen-go. DO NOT EDIT.
// source: stream.proto

/*
Package stream is a generated protocol buffer package.

It is generated from these files:
	stream.proto

It has these top-level messages:
	Attribute
	Request
	Response
*/
package stream

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

type Attribute struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Type  string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Value string `protobuf:"bytes,3,opt,name=value" json:"value,omitempty"`
}

func (m *Attribute) Reset()                    { *m = Attribute{} }
func (m *Attribute) String() string            { return proto.CompactTextString(m) }
func (*Attribute) ProtoMessage()               {}
func (*Attribute) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Attribute) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Attribute) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Attribute) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Request struct {
	Id         uint32       `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Attributes []*Attribute `protobuf:"bytes,2,rep,name=attributes" json:"attributes,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Request) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Request) GetAttributes() []*Attribute {
	if m != nil {
		return m.Attributes
	}
	return nil
}

type Response struct {
	Id          uint32       `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Status      string       `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	Obligations []*Attribute `protobuf:"bytes,3,rep,name=obligations" json:"obligations,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Response) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Response) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Response) GetObligations() []*Attribute {
	if m != nil {
		return m.Obligations
	}
	return nil
}

func init() {
	proto.RegisterType((*Attribute)(nil), "stream.Attribute")
	proto.RegisterType((*Request)(nil), "stream.Request")
	proto.RegisterType((*Response)(nil), "stream.Response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Stream service

type StreamClient interface {
	Test(ctx context.Context, opts ...grpc.CallOption) (Stream_TestClient, error)
}

type streamClient struct {
	cc *grpc.ClientConn
}

func NewStreamClient(cc *grpc.ClientConn) StreamClient {
	return &streamClient{cc}
}

func (c *streamClient) Test(ctx context.Context, opts ...grpc.CallOption) (Stream_TestClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Stream_serviceDesc.Streams[0], c.cc, "/stream.Stream/Test", opts...)
	if err != nil {
		return nil, err
	}
	x := &streamTestClient{stream}
	return x, nil
}

type Stream_TestClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type streamTestClient struct {
	grpc.ClientStream
}

func (x *streamTestClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *streamTestClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Stream service

type StreamServer interface {
	Test(Stream_TestServer) error
}

func RegisterStreamServer(s *grpc.Server, srv StreamServer) {
	s.RegisterService(&_Stream_serviceDesc, srv)
}

func _Stream_Test_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StreamServer).Test(&streamTestServer{stream})
}

type Stream_TestServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type streamTestServer struct {
	grpc.ServerStream
}

func (x *streamTestServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *streamTestServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Stream_serviceDesc = grpc.ServiceDesc{
	ServiceName: "stream.Stream",
	HandlerType: (*StreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Test",
			Handler:       _Stream_Test_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "stream.proto",
}

func init() { proto.RegisterFile("stream.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 223 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x3f, 0x4f, 0xc3, 0x30,
	0x14, 0xc4, 0xb1, 0x53, 0x0c, 0x7d, 0xe5, 0xaf, 0x85, 0x90, 0xc5, 0x54, 0x79, 0xca, 0x54, 0xa0,
	0x9d, 0x18, 0x19, 0xd8, 0x98, 0x0c, 0x5f, 0xc0, 0x51, 0x9f, 0x2a, 0x4b, 0x25, 0x0e, 0x79, 0xcf,
	0x48, 0x7c, 0x7b, 0x84, 0xe3, 0x40, 0x24, 0xc4, 0xe6, 0x3b, 0x9d, 0xef, 0x7e, 0x36, 0x9c, 0x10,
	0xf7, 0xe8, 0xdf, 0x56, 0x5d, 0x1f, 0x39, 0x6a, 0x35, 0x28, 0xfb, 0x04, 0xf3, 0x47, 0xe6, 0x3e,
	0x34, 0x89, 0x51, 0x9f, 0x81, 0x0c, 0x5b, 0x23, 0x96, 0xa2, 0x9e, 0x3b, 0x19, 0xb6, 0x5a, 0xc3,
	0x8c, 0x3f, 0x3b, 0x34, 0x32, 0x3b, 0xf9, 0xac, 0xaf, 0xe0, 0xf0, 0xc3, 0xef, 0x13, 0x9a, 0x2a,
	0x9b, 0x83, 0xb0, 0xcf, 0x70, 0xe4, 0xf0, 0x3d, 0x21, 0xf1, 0xa4, 0xe4, 0x34, 0x97, 0xdc, 0x03,
	0xf8, 0x71, 0x81, 0x8c, 0x5c, 0x56, 0xf5, 0x62, 0x7d, 0xb9, 0x2a, 0x30, 0x3f, 0xdb, 0x6e, 0x12,
	0xb2, 0x3b, 0x38, 0x76, 0x48, 0x5d, 0x6c, 0x09, 0xff, 0xd4, 0x5d, 0x83, 0x22, 0xf6, 0x9c, 0xa8,
	0x50, 0x15, 0xa5, 0x37, 0xb0, 0x88, 0xcd, 0x3e, 0xec, 0x3c, 0x87, 0xd8, 0x92, 0xa9, 0xfe, 0xdb,
	0x99, 0xa6, 0xd6, 0x0f, 0xa0, 0x5e, 0x72, 0x40, 0xdf, 0xc2, 0xec, 0xf5, 0x9b, 0xfe, 0x7c, 0xbc,
	0x51, 0x9e, 0x73, 0x73, 0xf1, 0x6b, 0x0c, 0x44, 0xf6, 0xa0, 0x16, 0x77, 0xa2, 0x51, 0xf9, 0x1f,
	0x37, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xf9, 0x81, 0x9a, 0x57, 0x57, 0x01, 0x00, 0x00,
}