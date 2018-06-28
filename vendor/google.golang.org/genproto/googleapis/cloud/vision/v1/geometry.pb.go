// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/vision/v1/geometry.proto

package vision // import "google.golang.org/genproto/googleapis/cloud/vision/v1"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// A vertex represents a 2D point in the image.
// NOTE: the vertex coordinates are in the same scale as the original image.
type Vertex struct {
	// X coordinate.
	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	// Y coordinate.
	Y                    int32    `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vertex) Reset()         { *m = Vertex{} }
func (m *Vertex) String() string { return proto.CompactTextString(m) }
func (*Vertex) ProtoMessage()    {}
func (*Vertex) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_8551ce1bb04014e1, []int{0}
}
func (m *Vertex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Vertex.Unmarshal(m, b)
}
func (m *Vertex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Vertex.Marshal(b, m, deterministic)
}
func (dst *Vertex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vertex.Merge(dst, src)
}
func (m *Vertex) XXX_Size() int {
	return xxx_messageInfo_Vertex.Size(m)
}
func (m *Vertex) XXX_DiscardUnknown() {
	xxx_messageInfo_Vertex.DiscardUnknown(m)
}

var xxx_messageInfo_Vertex proto.InternalMessageInfo

func (m *Vertex) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Vertex) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

// A vertex represents a 2D point in the image.
// NOTE: the normalized vertex coordinates are relative to the original image
// and range from 0 to 1.
type NormalizedVertex struct {
	// X coordinate.
	X float32 `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	// Y coordinate.
	Y                    float32  `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NormalizedVertex) Reset()         { *m = NormalizedVertex{} }
func (m *NormalizedVertex) String() string { return proto.CompactTextString(m) }
func (*NormalizedVertex) ProtoMessage()    {}
func (*NormalizedVertex) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_8551ce1bb04014e1, []int{1}
}
func (m *NormalizedVertex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NormalizedVertex.Unmarshal(m, b)
}
func (m *NormalizedVertex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NormalizedVertex.Marshal(b, m, deterministic)
}
func (dst *NormalizedVertex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NormalizedVertex.Merge(dst, src)
}
func (m *NormalizedVertex) XXX_Size() int {
	return xxx_messageInfo_NormalizedVertex.Size(m)
}
func (m *NormalizedVertex) XXX_DiscardUnknown() {
	xxx_messageInfo_NormalizedVertex.DiscardUnknown(m)
}

var xxx_messageInfo_NormalizedVertex proto.InternalMessageInfo

func (m *NormalizedVertex) GetX() float32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *NormalizedVertex) GetY() float32 {
	if m != nil {
		return m.Y
	}
	return 0
}

// A bounding polygon for the detected image annotation.
type BoundingPoly struct {
	// The bounding polygon vertices.
	Vertices []*Vertex `protobuf:"bytes,1,rep,name=vertices,proto3" json:"vertices,omitempty"`
	// The bounding polygon normalized vertices.
	NormalizedVertices   []*NormalizedVertex `protobuf:"bytes,2,rep,name=normalized_vertices,json=normalizedVertices,proto3" json:"normalized_vertices,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *BoundingPoly) Reset()         { *m = BoundingPoly{} }
func (m *BoundingPoly) String() string { return proto.CompactTextString(m) }
func (*BoundingPoly) ProtoMessage()    {}
func (*BoundingPoly) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_8551ce1bb04014e1, []int{2}
}
func (m *BoundingPoly) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BoundingPoly.Unmarshal(m, b)
}
func (m *BoundingPoly) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BoundingPoly.Marshal(b, m, deterministic)
}
func (dst *BoundingPoly) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BoundingPoly.Merge(dst, src)
}
func (m *BoundingPoly) XXX_Size() int {
	return xxx_messageInfo_BoundingPoly.Size(m)
}
func (m *BoundingPoly) XXX_DiscardUnknown() {
	xxx_messageInfo_BoundingPoly.DiscardUnknown(m)
}

var xxx_messageInfo_BoundingPoly proto.InternalMessageInfo

func (m *BoundingPoly) GetVertices() []*Vertex {
	if m != nil {
		return m.Vertices
	}
	return nil
}

func (m *BoundingPoly) GetNormalizedVertices() []*NormalizedVertex {
	if m != nil {
		return m.NormalizedVertices
	}
	return nil
}

// A 3D position in the image, used primarily for Face detection landmarks.
// A valid Position must have both x and y coordinates.
// The position coordinates are in the same scale as the original image.
type Position struct {
	// X coordinate.
	X float32 `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	// Y coordinate.
	Y float32 `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
	// Z coordinate (or depth).
	Z                    float32  `protobuf:"fixed32,3,opt,name=z,proto3" json:"z,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Position) Reset()         { *m = Position{} }
func (m *Position) String() string { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()    {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_8551ce1bb04014e1, []int{3}
}
func (m *Position) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Position.Unmarshal(m, b)
}
func (m *Position) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Position.Marshal(b, m, deterministic)
}
func (dst *Position) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Position.Merge(dst, src)
}
func (m *Position) XXX_Size() int {
	return xxx_messageInfo_Position.Size(m)
}
func (m *Position) XXX_DiscardUnknown() {
	xxx_messageInfo_Position.DiscardUnknown(m)
}

var xxx_messageInfo_Position proto.InternalMessageInfo

func (m *Position) GetX() float32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Position) GetY() float32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Position) GetZ() float32 {
	if m != nil {
		return m.Z
	}
	return 0
}

func init() {
	proto.RegisterType((*Vertex)(nil), "google.cloud.vision.v1.Vertex")
	proto.RegisterType((*NormalizedVertex)(nil), "google.cloud.vision.v1.NormalizedVertex")
	proto.RegisterType((*BoundingPoly)(nil), "google.cloud.vision.v1.BoundingPoly")
	proto.RegisterType((*Position)(nil), "google.cloud.vision.v1.Position")
}

func init() {
	proto.RegisterFile("google/cloud/vision/v1/geometry.proto", fileDescriptor_geometry_8551ce1bb04014e1)
}

var fileDescriptor_geometry_8551ce1bb04014e1 = []byte{
	// 276 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0xc7, 0x49, 0x87, 0x63, 0xc4, 0x09, 0x52, 0x41, 0x8a, 0x07, 0x19, 0x45, 0xa1, 0xa7, 0x84,
	0xa9, 0x27, 0xf5, 0xd4, 0x8b, 0x37, 0x29, 0x3d, 0x08, 0x7a, 0x91, 0xda, 0x86, 0x10, 0x48, 0xdf,
	0x1b, 0x69, 0x56, 0xd6, 0x7e, 0x1f, 0xbf, 0xa3, 0x47, 0x69, 0x32, 0x0a, 0x2b, 0xce, 0x5b, 0xfe,
	0xc9, 0x8f, 0xdf, 0x7b, 0x79, 0x8f, 0xde, 0x4a, 0x44, 0xa9, 0x05, 0x2f, 0x35, 0x6e, 0x2b, 0xde,
	0xaa, 0x46, 0x21, 0xf0, 0x76, 0xcd, 0xa5, 0xc0, 0x5a, 0x58, 0xd3, 0xb1, 0x8d, 0x41, 0x8b, 0xe1,
	0xa5, 0xc7, 0x98, 0xc3, 0x98, 0xc7, 0x58, 0xbb, 0x8e, 0x6f, 0xe8, 0xfc, 0x4d, 0x18, 0x2b, 0x76,
	0xe1, 0x92, 0x92, 0x5d, 0x44, 0x56, 0x24, 0x39, 0xc9, 0x89, 0x4b, 0x5d, 0x14, 0xf8, 0xd4, 0xc5,
	0x8c, 0x9e, 0xbf, 0xa2, 0xa9, 0x0b, 0xad, 0x7a, 0x51, 0x4d, 0xf9, 0xe0, 0x80, 0x0f, 0x06, 0xfe,
	0x9b, 0xd0, 0x65, 0x8a, 0x5b, 0xa8, 0x14, 0xc8, 0x0c, 0x75, 0x17, 0x3e, 0xd2, 0x45, 0x2b, 0x8c,
	0x55, 0xa5, 0x68, 0x22, 0xb2, 0x9a, 0x25, 0xa7, 0x77, 0xd7, 0xec, 0xef, 0x8e, 0x98, 0xd7, 0xe7,
	0x23, 0x1f, 0xbe, 0xd3, 0x0b, 0x18, 0x8b, 0x7f, 0x8e, 0x9a, 0xc0, 0x69, 0x92, 0x63, 0x9a, 0x69,
	0xbf, 0x79, 0x08, 0x07, 0x37, 0x83, 0x23, 0x7e, 0xa0, 0x8b, 0x0c, 0x1b, 0x65, 0x15, 0xc2, 0x7f,
	0xff, 0x19, 0x52, 0x1f, 0xcd, 0x7c, 0xea, 0x53, 0xa0, 0x57, 0x25, 0xd6, 0x47, 0x0a, 0xa7, 0x67,
	0x2f, 0xfb, 0xc9, 0x67, 0xc3, 0xe0, 0x33, 0xf2, 0xf1, 0xbc, 0x07, 0x25, 0xea, 0x02, 0x24, 0x43,
	0x23, 0xb9, 0x14, 0xe0, 0xd6, 0xc2, 0xfd, 0x53, 0xb1, 0x51, 0xcd, 0x74, 0x81, 0x4f, 0xfe, 0xf4,
	0x43, 0xc8, 0xd7, 0xdc, 0xb1, 0xf7, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x6f, 0xd6, 0x4f, 0x3a,
	0xeb, 0x01, 0x00, 0x00,
}
