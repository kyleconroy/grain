// Code generated by protoc-gen-go. DO NOT EDIT.
// source: facebook/photo.proto

package facebookpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Owner struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Id   string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
}

func (m *Owner) Reset()                    { *m = Owner{} }
func (m *Owner) String() string            { return proto.CompactTextString(m) }
func (*Owner) ProtoMessage()               {}
func (*Owner) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *Owner) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Owner) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Location struct {
	City      string  `protobuf:"bytes,1,opt,name=city" json:"city,omitempty"`
	Country   string  `protobuf:"bytes,2,opt,name=country" json:"country,omitempty"`
	Latitude  float32 `protobuf:"fixed32,3,opt,name=latitude" json:"latitude,omitempty"`
	Longitude float32 `protobuf:"fixed32,4,opt,name=longitude" json:"longitude,omitempty"`
}

func (m *Location) Reset()                    { *m = Location{} }
func (m *Location) String() string            { return proto.CompactTextString(m) }
func (*Location) ProtoMessage()               {}
func (*Location) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *Location) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *Location) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *Location) GetLatitude() float32 {
	if m != nil {
		return m.Latitude
	}
	return 0
}

func (m *Location) GetLongitude() float32 {
	if m != nil {
		return m.Longitude
	}
	return 0
}

type Place struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// The name of the Page
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	// The location of this place. Applicable to all Places
	Location *Location `protobuf:"bytes,3,opt,name=location" json:"location,omitempty"`
}

func (m *Place) Reset()                    { *m = Place{} }
func (m *Place) String() string            { return proto.CompactTextString(m) }
func (*Place) ProtoMessage()               {}
func (*Place) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *Place) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Place) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Place) GetLocation() *Location {
	if m != nil {
		return m.Location
	}
	return nil
}

type Album struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// A user-specified time for when this object was created
	BackdatedTime string `protobuf:"bytes,2,opt,name=backdated_time,json=backdatedTime" json:"backdated_time,omitempty"`
	// How accurate the backdated time is
	BackdatedTimeGranularity string `protobuf:"bytes,3,opt,name=backdated_time_granularity,json=backdatedTimeGranularity" json:"backdated_time_granularity,omitempty"`
	// Whether the viewer can backdate this album
	CanBackdate bool `protobuf:"varint,4,opt,name=can_backdate,json=canBackdate" json:"can_backdate,omitempty"`
	// Whether the viewer can upload photos to this album
	CanUpload bool `protobuf:"varint,5,opt,name=can_upload,json=canUpload" json:"can_upload,omitempty"`
	// The approximate number of photos in the album. This is not necessarily an exact count
	Count      uint32 `protobuf:"varint,6,opt,name=count" json:"count,omitempty"`
	PhotoCount uint32 `protobuf:"varint,18,opt,name=photo_count,json=photoCount" json:"photo_count,omitempty"`
	// Album cover photo id
	CoverPhoto *Photo `protobuf:"bytes,7,opt,name=cover_photo,json=coverPhoto" json:"cover_photo,omitempty"`
	// The time the album was initially created
	CreatedTime string `protobuf:"bytes,8,opt,name=created_time,json=createdTime" json:"created_time,omitempty"`
	// The description of the album
	Description string `protobuf:"bytes,9,opt,name=description" json:"description,omitempty"`
	// The URL for editing this album
	EditLink string `protobuf:"bytes,10,opt,name=edit_link,json=editLink" json:"edit_link,omitempty"`
	// If this object has a place, the event associated with the place
	// string event = 11;
	// The profile that created the album
	From *Owner `protobuf:"bytes,12,opt,name=from" json:"from,omitempty"`
	// Determines whether or not the album should be shown to users
	IsUserFacing bool `protobuf:"varint,13,opt,name=is_user_facing,json=isUserFacing" json:"is_user_facing,omitempty"`
	// A link to this album on Facebook
	Link string `protobuf:"bytes,14,opt,name=link" json:"link,omitempty"`
	// The textual location of the album"
	Location string `protobuf:"bytes,15,opt,name=location" json:"location,omitempty"`
	// Time of the last major update (e.g. addition of photos) expressed as UNIX time
	ModifiedMajor string `protobuf:"bytes,16,opt,name=modified_major,json=modifiedMajor" json:"modified_major,omitempty"`
	// The title of the album
	Name string `protobuf:"bytes,17,opt,name=name" json:"name,omitempty"`
	// The place associated with this album
	Place *Place `protobuf:"bytes,19,opt,name=place" json:"place,omitempty"`
	// The privacy settings for the album
	Privacy string `protobuf:"bytes,20,opt,name=privacy" json:"privacy,omitempty"`
	// The type of the album: profile, mobile, wall, normal or album
	Type string `protobuf:"bytes,21,opt,name=type" json:"type,omitempty"`
	// The last time the album was updated
	UpdatedTime string `protobuf:"bytes,22,opt,name=updated_time,json=updatedTime" json:"updated_time,omitempty"`
	// The approximate number of videos in the album. This is not necessarily an exact count
	VideoCount uint32 `protobuf:"varint,23,opt,name=video_count,json=videoCount" json:"video_count,omitempty"`
}

func (m *Album) Reset()                    { *m = Album{} }
func (m *Album) String() string            { return proto.CompactTextString(m) }
func (*Album) ProtoMessage()               {}
func (*Album) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *Album) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Album) GetBackdatedTime() string {
	if m != nil {
		return m.BackdatedTime
	}
	return ""
}

func (m *Album) GetBackdatedTimeGranularity() string {
	if m != nil {
		return m.BackdatedTimeGranularity
	}
	return ""
}

func (m *Album) GetCanBackdate() bool {
	if m != nil {
		return m.CanBackdate
	}
	return false
}

func (m *Album) GetCanUpload() bool {
	if m != nil {
		return m.CanUpload
	}
	return false
}

func (m *Album) GetCount() uint32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *Album) GetPhotoCount() uint32 {
	if m != nil {
		return m.PhotoCount
	}
	return 0
}

func (m *Album) GetCoverPhoto() *Photo {
	if m != nil {
		return m.CoverPhoto
	}
	return nil
}

func (m *Album) GetCreatedTime() string {
	if m != nil {
		return m.CreatedTime
	}
	return ""
}

func (m *Album) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Album) GetEditLink() string {
	if m != nil {
		return m.EditLink
	}
	return ""
}

func (m *Album) GetFrom() *Owner {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Album) GetIsUserFacing() bool {
	if m != nil {
		return m.IsUserFacing
	}
	return false
}

func (m *Album) GetLink() string {
	if m != nil {
		return m.Link
	}
	return ""
}

func (m *Album) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *Album) GetModifiedMajor() string {
	if m != nil {
		return m.ModifiedMajor
	}
	return ""
}

func (m *Album) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Album) GetPlace() *Place {
	if m != nil {
		return m.Place
	}
	return nil
}

func (m *Album) GetPrivacy() string {
	if m != nil {
		return m.Privacy
	}
	return ""
}

func (m *Album) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Album) GetUpdatedTime() string {
	if m != nil {
		return m.UpdatedTime
	}
	return ""
}

func (m *Album) GetVideoCount() uint32 {
	if m != nil {
		return m.VideoCount
	}
	return 0
}

type Image struct {
	Height int32  `protobuf:"varint,1,opt,name=height" json:"height,omitempty"`
	Width  int32  `protobuf:"varint,2,opt,name=width" json:"width,omitempty"`
	Source string `protobuf:"bytes,3,opt,name=source" json:"source,omitempty"`
}

func (m *Image) Reset()                    { *m = Image{} }
func (m *Image) String() string            { return proto.CompactTextString(m) }
func (*Image) ProtoMessage()               {}
func (*Image) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func (m *Image) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Image) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *Image) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

type Photo struct {
	Id    string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Album *Album `protobuf:"bytes,2,opt,name=album" json:"album,omitempty"`
	// A user-specified time for when this object was created
	BackdatedTime string `protobuf:"bytes,3,opt,name=backdated_time,json=backdatedTime" json:"backdated_time,omitempty"`
	// How accurate the backdated time is
	BackdatedTimeGranularity string   `protobuf:"bytes,4,opt,name=backdated_time_granularity,json=backdatedTimeGranularity" json:"backdated_time_granularity,omitempty"`
	CanBackdate              bool     `protobuf:"varint,5,opt,name=can_backdate,json=canBackdate" json:"can_backdate,omitempty"`
	CanDelete                bool     `protobuf:"varint,6,opt,name=can_delete,json=canDelete" json:"can_delete,omitempty"`
	CanTag                   bool     `protobuf:"varint,7,opt,name=can_tag,json=canTag" json:"can_tag,omitempty"`
	CreatedTime              string   `protobuf:"bytes,8,opt,name=created_time,json=createdTime" json:"created_time,omitempty"`
	From                     *Owner   `protobuf:"bytes,10,opt,name=from" json:"from,omitempty"`
	Height                   int32    `protobuf:"varint,11,opt,name=height" json:"height,omitempty"`
	Width                    int32    `protobuf:"varint,12,opt,name=width" json:"width,omitempty"`
	Icon                     string   `protobuf:"bytes,13,opt,name=icon" json:"icon,omitempty"`
	Images                   []*Image `protobuf:"bytes,14,rep,name=images" json:"images,omitempty"`
	Name                     string   `protobuf:"bytes,15,opt,name=name" json:"name,omitempty"`
}

func (m *Photo) Reset()                    { *m = Photo{} }
func (m *Photo) String() string            { return proto.CompactTextString(m) }
func (*Photo) ProtoMessage()               {}
func (*Photo) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{5} }

func (m *Photo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Photo) GetAlbum() *Album {
	if m != nil {
		return m.Album
	}
	return nil
}

func (m *Photo) GetBackdatedTime() string {
	if m != nil {
		return m.BackdatedTime
	}
	return ""
}

func (m *Photo) GetBackdatedTimeGranularity() string {
	if m != nil {
		return m.BackdatedTimeGranularity
	}
	return ""
}

func (m *Photo) GetCanBackdate() bool {
	if m != nil {
		return m.CanBackdate
	}
	return false
}

func (m *Photo) GetCanDelete() bool {
	if m != nil {
		return m.CanDelete
	}
	return false
}

func (m *Photo) GetCanTag() bool {
	if m != nil {
		return m.CanTag
	}
	return false
}

func (m *Photo) GetCreatedTime() string {
	if m != nil {
		return m.CreatedTime
	}
	return ""
}

func (m *Photo) GetFrom() *Owner {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Photo) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Photo) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *Photo) GetIcon() string {
	if m != nil {
		return m.Icon
	}
	return ""
}

func (m *Photo) GetImages() []*Image {
	if m != nil {
		return m.Images
	}
	return nil
}

func (m *Photo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Owner)(nil), "grain.facebook.Owner")
	proto.RegisterType((*Location)(nil), "grain.facebook.Location")
	proto.RegisterType((*Place)(nil), "grain.facebook.Place")
	proto.RegisterType((*Album)(nil), "grain.facebook.Album")
	proto.RegisterType((*Image)(nil), "grain.facebook.Image")
	proto.RegisterType((*Photo)(nil), "grain.facebook.Photo")
}

func init() { proto.RegisterFile("facebook/photo.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 710 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0x5d, 0x6f, 0xd3, 0x4a,
	0x10, 0x55, 0x3e, 0x9c, 0x26, 0xe3, 0x34, 0xbd, 0x77, 0x6f, 0x3f, 0x56, 0xbd, 0xa0, 0x86, 0x08,
	0xa4, 0xa2, 0x8a, 0x20, 0x01, 0xe2, 0x89, 0x17, 0x0a, 0x02, 0x21, 0xb5, 0xa2, 0xb2, 0xda, 0x17,
	0x5e, 0xac, 0xcd, 0x7a, 0xe3, 0x2c, 0x71, 0x76, 0xad, 0xcd, 0xba, 0x55, 0x7e, 0x29, 0xff, 0x81,
	0x5f, 0x81, 0x76, 0xd6, 0x4e, 0xd2, 0x92, 0x0a, 0xfa, 0xb6, 0x73, 0xe6, 0xd8, 0x3b, 0x3e, 0x73,
	0x8e, 0x0c, 0xbb, 0x63, 0xc6, 0xc5, 0x48, 0xeb, 0xe9, 0xcb, 0x7c, 0xa2, 0xad, 0x1e, 0xe6, 0x46,
	0x5b, 0x4d, 0x7a, 0xa9, 0x61, 0x52, 0x0d, 0xab, 0xde, 0xe0, 0x04, 0x82, 0xaf, 0x37, 0x4a, 0x18,
	0x42, 0xa0, 0xa9, 0xd8, 0x4c, 0xd0, 0x5a, 0xbf, 0x76, 0xdc, 0x89, 0xf0, 0x4c, 0x7a, 0x50, 0x97,
	0x09, 0xad, 0x23, 0x52, 0x97, 0xc9, 0xc0, 0x40, 0xfb, 0x4c, 0x73, 0x66, 0xa5, 0x56, 0x8e, 0xcf,
	0xa5, 0x5d, 0x54, 0x7c, 0x77, 0x26, 0x14, 0xb6, 0xb8, 0x2e, 0x94, 0x35, 0x8b, 0xf2, 0xa1, 0xaa,
	0x24, 0x87, 0xd0, 0xce, 0x98, 0x95, 0xb6, 0x48, 0x04, 0x6d, 0xf4, 0x6b, 0xc7, 0xf5, 0x68, 0x59,
	0x93, 0x47, 0xd0, 0xc9, 0xb4, 0x4a, 0x7d, 0xb3, 0x89, 0xcd, 0x15, 0x30, 0x60, 0x10, 0x5c, 0x64,
	0x8c, 0x57, 0xc3, 0xd4, 0xaa, 0x61, 0x96, 0x03, 0xd7, 0xd7, 0x06, 0x7e, 0x03, 0xed, 0xac, 0x1c,
	0x10, 0xaf, 0x09, 0x5f, 0xd1, 0xe1, 0xed, 0x0f, 0x1e, 0x56, 0x1f, 0x10, 0x2d, 0x99, 0x83, 0x9f,
	0x01, 0x04, 0xef, 0xb3, 0x51, 0x31, 0xfb, 0xed, 0x8e, 0x67, 0xd0, 0x1b, 0x31, 0x3e, 0x4d, 0x98,
	0x15, 0x49, 0x6c, 0xe5, 0xf2, 0xb6, 0xed, 0x25, 0x7a, 0x29, 0x67, 0x82, 0xbc, 0x83, 0xc3, 0xdb,
	0xb4, 0x38, 0x35, 0x4c, 0x15, 0x19, 0x33, 0x4e, 0xa1, 0x06, 0x3e, 0x42, 0x6f, 0x3d, 0xf2, 0x79,
	0xd5, 0x27, 0x4f, 0xa0, 0xcb, 0x99, 0x8a, 0xab, 0x3e, 0x4a, 0xd0, 0x8e, 0x42, 0xce, 0xd4, 0x69,
	0x09, 0x91, 0xc7, 0x00, 0x8e, 0x52, 0xe4, 0x99, 0x66, 0x09, 0x0d, 0x90, 0xd0, 0xe1, 0x4c, 0x5d,
	0x21, 0x40, 0x76, 0x21, 0x40, 0xa1, 0x69, 0xab, 0x5f, 0x3b, 0xde, 0x8e, 0x7c, 0x41, 0x8e, 0x20,
	0xc4, 0xcd, 0xc7, 0xbe, 0x47, 0xb0, 0x07, 0x08, 0x7d, 0x40, 0xc2, 0x5b, 0x08, 0xb9, 0xbe, 0x16,
	0x26, 0x46, 0x8c, 0x6e, 0xa1, 0x60, 0x7b, 0x77, 0x05, 0xbb, 0x70, 0xcd, 0x08, 0x90, 0x89, 0x67,
	0x1c, 0xd8, 0x88, 0x95, 0x26, 0x6d, 0xfc, 0xc0, 0xb0, 0xc4, 0x50, 0x91, 0x3e, 0x84, 0x89, 0x98,
	0x73, 0x23, 0x73, 0xdc, 0x45, 0xc7, 0x33, 0xd6, 0x20, 0xf2, 0x3f, 0x74, 0x44, 0x22, 0x6d, 0x9c,
	0x49, 0x35, 0xa5, 0x80, 0xfd, 0xb6, 0x03, 0xce, 0xa4, 0x9a, 0x92, 0xe7, 0xd0, 0x1c, 0x1b, 0x3d,
	0xa3, 0xdd, 0xcd, 0x23, 0xa1, 0x63, 0x23, 0xa4, 0x90, 0xa7, 0xd0, 0x93, 0xf3, 0xb8, 0x98, 0x0b,
	0x13, 0x8f, 0x19, 0x97, 0x2a, 0xa5, 0xdb, 0x28, 0x4f, 0x57, 0xce, 0xaf, 0xe6, 0xc2, 0x7c, 0x42,
	0xcc, 0x99, 0x05, 0x2f, 0xea, 0x79, 0xb3, 0xb8, 0x33, 0x7a, 0xb2, 0x32, 0xcb, 0x8e, 0x1f, 0xa0,
	0xaa, 0xdd, 0xe2, 0x67, 0x3a, 0x91, 0x63, 0x29, 0x92, 0x78, 0xc6, 0xbe, 0x6b, 0x43, 0xff, 0xf1,
	0x8b, 0xaf, 0xd0, 0x73, 0x07, 0x2e, 0x3d, 0xf8, 0xef, 0x9a, 0x07, 0x4f, 0x20, 0xc8, 0x9d, 0x61,
	0xe9, 0x7f, 0xf7, 0xe8, 0xe9, 0x9a, 0x91, 0xe7, 0xb8, 0xc4, 0xe4, 0x46, 0x5e, 0x33, 0xbe, 0xa0,
	0xbb, 0x3e, 0x31, 0x65, 0xe9, 0x5e, 0x6d, 0x17, 0xb9, 0xa0, 0x7b, 0xfe, 0xd5, 0xee, 0xec, 0x84,
	0x2f, 0xf2, 0x35, 0x33, 0xee, 0x7b, 0x59, 0x4b, 0x0c, 0x85, 0x3f, 0x82, 0xf0, 0x5a, 0x26, 0xa2,
	0x5a, 0xfa, 0x81, 0x5f, 0x3a, 0x42, 0xb8, 0xf4, 0xc1, 0x39, 0x04, 0x5f, 0x66, 0x2c, 0x15, 0x64,
	0x1f, 0x5a, 0x13, 0x21, 0xd3, 0x89, 0x45, 0xbf, 0x07, 0x51, 0x59, 0x39, 0x33, 0xdd, 0xc8, 0xc4,
	0x4e, 0xd0, 0xea, 0x41, 0xe4, 0x0b, 0xc7, 0x9e, 0xeb, 0xc2, 0x70, 0x51, 0xda, 0xb9, 0xac, 0x06,
	0x3f, 0x1a, 0x10, 0x78, 0x57, 0xdc, 0xcd, 0xce, 0x09, 0x04, 0xcc, 0x85, 0x0a, 0xdf, 0xb3, 0x41,
	0x07, 0x4c, 0x5c, 0xe4, 0x39, 0x1b, 0x82, 0xd6, 0x78, 0x78, 0xd0, 0x9a, 0x0f, 0x0c, 0x5a, 0x70,
	0x6f, 0xd0, 0x12, 0x91, 0x09, 0x2b, 0x30, 0x4e, 0x3e, 0x68, 0x1f, 0x11, 0x20, 0x07, 0xb0, 0xe5,
	0xda, 0x96, 0xa5, 0x98, 0x96, 0x76, 0xd4, 0xe2, 0x4c, 0x5d, 0xb2, 0xf4, 0x6f, 0x22, 0x51, 0x79,
	0x1a, 0xfe, 0xec, 0xe9, 0xd5, 0x6a, 0xc2, 0xcd, 0xab, 0xe9, 0xae, 0xaf, 0x86, 0x40, 0x53, 0x72,
	0xad, 0xd0, 0xf7, 0x9d, 0x08, 0xcf, 0xe4, 0x05, 0xb4, 0xa4, 0xdb, 0xf2, 0x9c, 0xf6, 0xfa, 0x8d,
	0x4d, 0xd7, 0xa1, 0x07, 0xa2, 0x92, 0xb4, 0xf4, 0xf1, 0xce, 0xca, 0xc7, 0xa7, 0xdd, 0x6f, 0x50,
	0xb1, 0xf3, 0xd1, 0xa8, 0x85, 0xbf, 0x8f, 0xd7, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4d, 0x81,
	0x7e, 0x12, 0x56, 0x06, 0x00, 0x00,
}
