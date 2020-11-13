// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/feeds.proto

package feeds

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Feed struct {
	// rss feed name
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// rss feed url
	Url                  string   `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Feed) Reset()         { *m = Feed{} }
func (m *Feed) String() string { return proto.CompactTextString(m) }
func (*Feed) ProtoMessage()    {}
func (*Feed) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{0}
}

func (m *Feed) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Feed.Unmarshal(m, b)
}
func (m *Feed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Feed.Marshal(b, m, deterministic)
}
func (m *Feed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Feed.Merge(m, src)
}
func (m *Feed) XXX_Size() int {
	return xxx_messageInfo_Feed.Size(m)
}
func (m *Feed) XXX_DiscardUnknown() {
	xxx_messageInfo_Feed.DiscardUnknown(m)
}

var xxx_messageInfo_Feed proto.InternalMessageInfo

func (m *Feed) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Feed) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type Entry struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Domain               string   `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	Url                  string   `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Title                string   `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Content              string   `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
	Date                 int64    `protobuf:"varint,6,opt,name=date,proto3" json:"date,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}
func (*Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{1}
}

func (m *Entry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Entry.Unmarshal(m, b)
}
func (m *Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Entry.Marshal(b, m, deterministic)
}
func (m *Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entry.Merge(m, src)
}
func (m *Entry) XXX_Size() int {
	return xxx_messageInfo_Entry.Size(m)
}
func (m *Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Entry proto.InternalMessageInfo

func (m *Entry) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Entry) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *Entry) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Entry) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Entry) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Entry) GetDate() int64 {
	if m != nil {
		return m.Date
	}
	return 0
}

type NewRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Url                  string   `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewRequest) Reset()         { *m = NewRequest{} }
func (m *NewRequest) String() string { return proto.CompactTextString(m) }
func (*NewRequest) ProtoMessage()    {}
func (*NewRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{2}
}

func (m *NewRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewRequest.Unmarshal(m, b)
}
func (m *NewRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewRequest.Marshal(b, m, deterministic)
}
func (m *NewRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewRequest.Merge(m, src)
}
func (m *NewRequest) XXX_Size() int {
	return xxx_messageInfo_NewRequest.Size(m)
}
func (m *NewRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NewRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NewRequest proto.InternalMessageInfo

func (m *NewRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NewRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type NewResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewResponse) Reset()         { *m = NewResponse{} }
func (m *NewResponse) String() string { return proto.CompactTextString(m) }
func (*NewResponse) ProtoMessage()    {}
func (*NewResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{3}
}

func (m *NewResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewResponse.Unmarshal(m, b)
}
func (m *NewResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewResponse.Marshal(b, m, deterministic)
}
func (m *NewResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewResponse.Merge(m, src)
}
func (m *NewResponse) XXX_Size() int {
	return xxx_messageInfo_NewResponse.Size(m)
}
func (m *NewResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NewResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NewResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Feed)(nil), "feeds.Feed")
	proto.RegisterType((*Entry)(nil), "feeds.Entry")
	proto.RegisterType((*NewRequest)(nil), "feeds.NewRequest")
	proto.RegisterType((*NewResponse)(nil), "feeds.NewResponse")
}

func init() { proto.RegisterFile("proto/feeds.proto", fileDescriptor_dd517c38176c13bf) }

var fileDescriptor_dd517c38176c13bf = []byte{
	// 233 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x71, 0x1c, 0x07, 0x71, 0x55, 0x11, 0x3d, 0x21, 0x64, 0x31, 0x55, 0x9e, 0x3a, 0xa0,
	0x20, 0x95, 0x81, 0x81, 0x0d, 0x09, 0xc6, 0x0e, 0x19, 0xd9, 0x02, 0x3e, 0x24, 0x4b, 0xad, 0x5d,
	0xe2, 0xab, 0x2a, 0x7e, 0x00, 0xff, 0x1b, 0xe5, 0x92, 0x00, 0x23, 0xdb, 0xfb, 0x9e, 0xef, 0x9e,
	0xce, 0x0f, 0x16, 0xfb, 0x2e, 0x71, 0xba, 0x7d, 0x27, 0xf2, 0xb9, 0x16, 0x8d, 0x46, 0xc0, 0xdd,
	0x40, 0xf9, 0x4c, 0xe4, 0x11, 0xa1, 0x8c, 0xed, 0x8e, 0xac, 0x5a, 0xaa, 0xd5, 0x59, 0x23, 0x1a,
	0x2f, 0x40, 0x1f, 0xba, 0xad, 0x2d, 0xc4, 0xea, 0xa5, 0xfb, 0x52, 0x60, 0x9e, 0x22, 0x77, 0x9f,
	0x78, 0x0e, 0x45, 0xf0, 0xe3, 0x74, 0x11, 0x3c, 0x5e, 0x41, 0xe5, 0xd3, 0xae, 0x0d, 0x71, 0x1c,
	0x1f, 0x69, 0xca, 0xd0, 0x3f, 0x19, 0x78, 0x09, 0x86, 0x03, 0x6f, 0xc9, 0x96, 0xe2, 0x0d, 0x80,
	0x16, 0x4e, 0xdf, 0x52, 0x64, 0x8a, 0x6c, 0x8d, 0xf8, 0x13, 0xf6, 0x97, 0xf9, 0x96, 0xc9, 0x56,
	0x4b, 0xb5, 0xd2, 0x8d, 0x68, 0xb7, 0x06, 0xd8, 0xd0, 0xb1, 0xa1, 0x8f, 0x03, 0x65, 0xfe, 0xe7,
	0xed, 0x73, 0x98, 0xc9, 0x4e, 0xde, 0xa7, 0x98, 0x69, 0x7d, 0x0f, 0xa6, 0xff, 0x78, 0xc6, 0x1a,
	0xf4, 0x86, 0x8e, 0xb8, 0xa8, 0x87, 0x76, 0x7e, 0x73, 0xaf, 0xf1, 0xaf, 0x35, 0xac, 0xb9, 0x93,
	0xc7, 0xf9, 0xcb, 0x4c, 0x1a, 0x7c, 0x90, 0xc7, 0xd7, 0x4a, 0xe0, 0xee, 0x3b, 0x00, 0x00, 0xff,
	0xff, 0xb0, 0x77, 0x9c, 0xb9, 0x63, 0x01, 0x00, 0x00,
}
