// Code generated by protoc-gen-go. DO NOT EDIT.
// source: facebook/download/about_you.proto

/*
Package grain_facebook_download is a generated protocol buffer package.

It is generated from these files:
	facebook/download/about_you.proto
	facebook/download/ads.proto
	facebook/download/apps.proto
	facebook/download/comments.proto
	facebook/download/events.proto
	facebook/download/following_and_followers.proto
	facebook/download/friends.proto
	facebook/download/groups.proto

It has these top-level messages:
	AddressBook
	Ads
	InstalledApps
	PostsFromApps
	YourApps
	Comment
	PhotoMetadata
	MediaMetadata
	Media
	AttachmentData
	Attachment
	Post
	YourPosts
	Coordinate
	Event
	EventInvitations
	EventResponses
	YourEvents
	FollowedPages
	Friend
	FriendsAdded
	RejectedRequests
	RemovedFriends
	SentRequests
	GroupsYouManage
	GroupMembershipActivity
	GroupPostsComments
*/
package grain_facebook_download

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

type AddressBook struct {
	AddressBook *AddressBook_AddressBookEntry `protobuf:"bytes,1,opt,name=address_book,json=addressBook" json:"address_book,omitempty"`
}

func (m *AddressBook) Reset()                    { *m = AddressBook{} }
func (m *AddressBook) String() string            { return proto.CompactTextString(m) }
func (*AddressBook) ProtoMessage()               {}
func (*AddressBook) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *AddressBook) GetAddressBook() *AddressBook_AddressBookEntry {
	if m != nil {
		return m.AddressBook
	}
	return nil
}

type AddressBook_Details struct {
	ContactPoint string `protobuf:"bytes,1,opt,name=contact_point,json=contactPoint" json:"contact_point,omitempty"`
}

func (m *AddressBook_Details) Reset()                    { *m = AddressBook_Details{} }
func (m *AddressBook_Details) String() string            { return proto.CompactTextString(m) }
func (*AddressBook_Details) ProtoMessage()               {}
func (*AddressBook_Details) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *AddressBook_Details) GetContactPoint() string {
	if m != nil {
		return m.ContactPoint
	}
	return ""
}

type AddressBook_Contact struct {
	Name    string                 `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Details []*AddressBook_Details `protobuf:"bytes,2,rep,name=details" json:"details,omitempty"`
}

func (m *AddressBook_Contact) Reset()                    { *m = AddressBook_Contact{} }
func (m *AddressBook_Contact) String() string            { return proto.CompactTextString(m) }
func (*AddressBook_Contact) ProtoMessage()               {}
func (*AddressBook_Contact) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *AddressBook_Contact) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AddressBook_Contact) GetDetails() []*AddressBook_Details {
	if m != nil {
		return m.Details
	}
	return nil
}

type AddressBook_AddressBookEntry struct {
	AddressBook []*AddressBook_Contact `protobuf:"bytes,1,rep,name=address_book,json=addressBook" json:"address_book,omitempty"`
}

func (m *AddressBook_AddressBookEntry) Reset()                    { *m = AddressBook_AddressBookEntry{} }
func (m *AddressBook_AddressBookEntry) String() string            { return proto.CompactTextString(m) }
func (*AddressBook_AddressBookEntry) ProtoMessage()               {}
func (*AddressBook_AddressBookEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 2} }

func (m *AddressBook_AddressBookEntry) GetAddressBook() []*AddressBook_Contact {
	if m != nil {
		return m.AddressBook
	}
	return nil
}

func init() {
	proto.RegisterType((*AddressBook)(nil), "grain.facebook.download.AddressBook")
	proto.RegisterType((*AddressBook_Details)(nil), "grain.facebook.download.AddressBook.Details")
	proto.RegisterType((*AddressBook_Contact)(nil), "grain.facebook.download.AddressBook.Contact")
	proto.RegisterType((*AddressBook_AddressBookEntry)(nil), "grain.facebook.download.AddressBook.AddressBookEntry")
}

func init() { proto.RegisterFile("facebook/download/about_you.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 227 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4c, 0x4b, 0x4c, 0x4e,
	0x4d, 0xca, 0xcf, 0xcf, 0xd6, 0x4f, 0xc9, 0x2f, 0xcf, 0xcb, 0xc9, 0x4f, 0x4c, 0xd1, 0x4f, 0x4c,
	0xca, 0x2f, 0x2d, 0x89, 0xaf, 0xcc, 0x2f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x4f,
	0x2f, 0x4a, 0xcc, 0xcc, 0xd3, 0x83, 0x29, 0xd4, 0x83, 0x29, 0x54, 0x7a, 0xcc, 0xc4, 0xc5, 0xed,
	0x98, 0x92, 0x52, 0x94, 0x5a, 0x5c, 0xec, 0x94, 0x9f, 0x9f, 0x2d, 0x14, 0xc1, 0xc5, 0x93, 0x08,
	0xe1, 0xc6, 0x83, 0x14, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b, 0x99, 0xea, 0xe1, 0xd0, 0xaf,
	0x87, 0xa4, 0x17, 0x99, 0xed, 0x9a, 0x57, 0x52, 0x54, 0x19, 0xc4, 0x9d, 0x88, 0x10, 0x91, 0xd2,
	0xe3, 0x62, 0x77, 0x49, 0x2d, 0x49, 0xcc, 0xcc, 0x29, 0x16, 0x52, 0xe6, 0xe2, 0x4d, 0xce, 0xcf,
	0x2b, 0x49, 0x4c, 0x2e, 0x89, 0x2f, 0xc8, 0xcf, 0xcc, 0x2b, 0x01, 0xdb, 0xc2, 0x19, 0xc4, 0x03,
	0x15, 0x0c, 0x00, 0x89, 0x49, 0xa5, 0x72, 0xb1, 0x3b, 0x43, 0xf8, 0x42, 0x42, 0x5c, 0x2c, 0x79,
	0x89, 0xb9, 0xa9, 0x50, 0x65, 0x60, 0xb6, 0x90, 0x1b, 0x17, 0x7b, 0x0a, 0xc4, 0x38, 0x09, 0x26,
	0x05, 0x66, 0x0d, 0x6e, 0x23, 0x1d, 0xa2, 0xdc, 0x08, 0x75, 0x42, 0x10, 0x4c, 0xb3, 0x54, 0x32,
	0x97, 0x00, 0xba, 0xbb, 0x85, 0xfc, 0x31, 0x02, 0x81, 0x78, 0x0b, 0xa0, 0x6e, 0x46, 0xf1, 0x7b,
	0x12, 0x1b, 0x38, 0x16, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8e, 0x54, 0xa6, 0xff, 0xaa,
	0x01, 0x00, 0x00,
}
