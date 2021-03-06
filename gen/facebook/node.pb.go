// Code generated by protoc-gen-go. DO NOT EDIT.
// source: facebook/node.proto

package facebookpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Node struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Node) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*Node)(nil), "grain.facebook.Node")
}

func init() { proto.RegisterFile("facebook/node.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 92 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4e, 0x4b, 0x4c, 0x4e,
	0x4d, 0xca, 0xcf, 0xcf, 0xd6, 0xcf, 0xcb, 0x4f, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0xe2, 0x4b, 0x2f, 0x4a, 0xcc, 0xcc, 0xd3, 0x83, 0x49, 0x29, 0x89, 0x71, 0xb1, 0xf8, 0xe5, 0xa7,
	0xa4, 0x0a, 0xf1, 0x71, 0x31, 0x65, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x31, 0x65,
	0xa6, 0x38, 0xf1, 0x44, 0x71, 0xc1, 0xd4, 0x14, 0x24, 0x25, 0xb1, 0x81, 0x35, 0x1b, 0x03, 0x02,
	0x00, 0x00, 0xff, 0xff, 0x12, 0x0e, 0xa2, 0xa1, 0x53, 0x00, 0x00, 0x00,
}
