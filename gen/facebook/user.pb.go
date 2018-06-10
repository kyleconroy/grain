// Code generated by protoc-gen-go. DO NOT EDIT.
// source: facebook/user.proto

package facebookpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type User struct {
	Id       string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Birthday string `protobuf:"bytes,2,opt,name=birthday" json:"birthday,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *User) GetBirthday() string {
	if m != nil {
		return m.Birthday
	}
	return ""
}

func init() {
	proto.RegisterType((*User)(nil), "grain.facebook.User")
}

func init() { proto.RegisterFile("facebook/user.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 110 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4e, 0x4b, 0x4c, 0x4e,
	0x4d, 0xca, 0xcf, 0xcf, 0xd6, 0x2f, 0x2d, 0x4e, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0xe2, 0x4b, 0x2f, 0x4a, 0xcc, 0xcc, 0xd3, 0x83, 0x49, 0x29, 0x19, 0x71, 0xb1, 0x84, 0x16, 0xa7,
	0x16, 0x09, 0xf1, 0x71, 0x31, 0x65, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x31, 0x65,
	0xa6, 0x08, 0x49, 0x71, 0x71, 0x24, 0x65, 0x16, 0x95, 0x64, 0xa4, 0x24, 0x56, 0x4a, 0x30, 0x81,
	0x45, 0xe1, 0x7c, 0x27, 0x9e, 0x28, 0x2e, 0x98, 0xfe, 0x82, 0xa4, 0x24, 0x36, 0xb0, 0xc1, 0xc6,
	0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x79, 0xa0, 0xaf, 0x69, 0x6f, 0x00, 0x00, 0x00,
}
