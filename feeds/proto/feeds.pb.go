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
	// eg. a16z
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// rss feed url
	// eg. http://a16z.com/feed/
	Url string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	// category of the feed
	Category             string   `protobuf:"bytes,3,opt,name=category,proto3" json:"category,omitempty"`
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

func (m *Feed) GetCategory() string {
	if m != nil {
		return m.Category
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
	Category             string   `protobuf:"bytes,7,opt,name=category,proto3" json:"category,omitempty"`
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

func (m *Entry) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

type AddRequest struct {
	// rss feed name
	// eg. a16z
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// rss feed url
	// eg. http://a16z.com/feed/
	Url string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	// category to add
	Category             string   `protobuf:"bytes,3,opt,name=category,proto3" json:"category,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddRequest) Reset()         { *m = AddRequest{} }
func (m *AddRequest) String() string { return proto.CompactTextString(m) }
func (*AddRequest) ProtoMessage()    {}
func (*AddRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{2}
}

func (m *AddRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddRequest.Unmarshal(m, b)
}
func (m *AddRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddRequest.Marshal(b, m, deterministic)
}
func (m *AddRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddRequest.Merge(m, src)
}
func (m *AddRequest) XXX_Size() int {
	return xxx_messageInfo_AddRequest.Size(m)
}
func (m *AddRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddRequest proto.InternalMessageInfo

func (m *AddRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AddRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *AddRequest) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

type AddResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddResponse) Reset()         { *m = AddResponse{} }
func (m *AddResponse) String() string { return proto.CompactTextString(m) }
func (*AddResponse) ProtoMessage()    {}
func (*AddResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{3}
}

func (m *AddResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddResponse.Unmarshal(m, b)
}
func (m *AddResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddResponse.Marshal(b, m, deterministic)
}
func (m *AddResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddResponse.Merge(m, src)
}
func (m *AddResponse) XXX_Size() int {
	return xxx_messageInfo_AddResponse.Size(m)
}
func (m *AddResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddResponse proto.InternalMessageInfo

type EntriesRequest struct {
	// rss feed url
	// eg. http://a16z.com/feed/
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EntriesRequest) Reset()         { *m = EntriesRequest{} }
func (m *EntriesRequest) String() string { return proto.CompactTextString(m) }
func (*EntriesRequest) ProtoMessage()    {}
func (*EntriesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{4}
}

func (m *EntriesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EntriesRequest.Unmarshal(m, b)
}
func (m *EntriesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EntriesRequest.Marshal(b, m, deterministic)
}
func (m *EntriesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EntriesRequest.Merge(m, src)
}
func (m *EntriesRequest) XXX_Size() int {
	return xxx_messageInfo_EntriesRequest.Size(m)
}
func (m *EntriesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EntriesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EntriesRequest proto.InternalMessageInfo

func (m *EntriesRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type EntriesResponse struct {
	Entries              []*Entry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EntriesResponse) Reset()         { *m = EntriesResponse{} }
func (m *EntriesResponse) String() string { return proto.CompactTextString(m) }
func (*EntriesResponse) ProtoMessage()    {}
func (*EntriesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{5}
}

func (m *EntriesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EntriesResponse.Unmarshal(m, b)
}
func (m *EntriesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EntriesResponse.Marshal(b, m, deterministic)
}
func (m *EntriesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EntriesResponse.Merge(m, src)
}
func (m *EntriesResponse) XXX_Size() int {
	return xxx_messageInfo_EntriesResponse.Size(m)
}
func (m *EntriesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EntriesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EntriesResponse proto.InternalMessageInfo

func (m *EntriesResponse) GetEntries() []*Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

type ListRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRequest) Reset()         { *m = ListRequest{} }
func (m *ListRequest) String() string { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()    {}
func (*ListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{6}
}

func (m *ListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRequest.Unmarshal(m, b)
}
func (m *ListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRequest.Marshal(b, m, deterministic)
}
func (m *ListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRequest.Merge(m, src)
}
func (m *ListRequest) XXX_Size() int {
	return xxx_messageInfo_ListRequest.Size(m)
}
func (m *ListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRequest proto.InternalMessageInfo

type ListResponse struct {
	Feeds                []*Feed  `protobuf:"bytes,1,rep,name=feeds,proto3" json:"feeds,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListResponse) Reset()         { *m = ListResponse{} }
func (m *ListResponse) String() string { return proto.CompactTextString(m) }
func (*ListResponse) ProtoMessage()    {}
func (*ListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{7}
}

func (m *ListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListResponse.Unmarshal(m, b)
}
func (m *ListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListResponse.Marshal(b, m, deterministic)
}
func (m *ListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListResponse.Merge(m, src)
}
func (m *ListResponse) XXX_Size() int {
	return xxx_messageInfo_ListResponse.Size(m)
}
func (m *ListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListResponse proto.InternalMessageInfo

func (m *ListResponse) GetFeeds() []*Feed {
	if m != nil {
		return m.Feeds
	}
	return nil
}

type RemoveRequest struct {
	// rss feed name
	// eg. a16z
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveRequest) Reset()         { *m = RemoveRequest{} }
func (m *RemoveRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveRequest) ProtoMessage()    {}
func (*RemoveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{8}
}

func (m *RemoveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveRequest.Unmarshal(m, b)
}
func (m *RemoveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveRequest.Marshal(b, m, deterministic)
}
func (m *RemoveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveRequest.Merge(m, src)
}
func (m *RemoveRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveRequest.Size(m)
}
func (m *RemoveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveRequest proto.InternalMessageInfo

func (m *RemoveRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type RemoveResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveResponse) Reset()         { *m = RemoveResponse{} }
func (m *RemoveResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveResponse) ProtoMessage()    {}
func (*RemoveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd517c38176c13bf, []int{9}
}

func (m *RemoveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveResponse.Unmarshal(m, b)
}
func (m *RemoveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveResponse.Marshal(b, m, deterministic)
}
func (m *RemoveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveResponse.Merge(m, src)
}
func (m *RemoveResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveResponse.Size(m)
}
func (m *RemoveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Feed)(nil), "feeds.Feed")
	proto.RegisterType((*Entry)(nil), "feeds.Entry")
	proto.RegisterType((*AddRequest)(nil), "feeds.AddRequest")
	proto.RegisterType((*AddResponse)(nil), "feeds.AddResponse")
	proto.RegisterType((*EntriesRequest)(nil), "feeds.EntriesRequest")
	proto.RegisterType((*EntriesResponse)(nil), "feeds.EntriesResponse")
	proto.RegisterType((*ListRequest)(nil), "feeds.ListRequest")
	proto.RegisterType((*ListResponse)(nil), "feeds.ListResponse")
	proto.RegisterType((*RemoveRequest)(nil), "feeds.RemoveRequest")
	proto.RegisterType((*RemoveResponse)(nil), "feeds.RemoveResponse")
}

func init() { proto.RegisterFile("proto/feeds.proto", fileDescriptor_dd517c38176c13bf) }

var fileDescriptor_dd517c38176c13bf = []byte{
	// 384 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0xbb, 0x4e, 0xf3, 0x40,
	0x10, 0x85, 0x7f, 0xc7, 0xb7, 0x9f, 0x71, 0x12, 0x92, 0x21, 0x89, 0x56, 0xae, 0x82, 0x91, 0x50,
	0xaa, 0xa0, 0x84, 0x02, 0x01, 0x55, 0x90, 0x40, 0x14, 0x88, 0xc2, 0x25, 0x9d, 0xc9, 0x0e, 0xc8,
	0x52, 0x62, 0x07, 0x7b, 0x83, 0x94, 0xc7, 0xe1, 0xc9, 0x78, 0x15, 0xe4, 0xbd, 0x38, 0x97, 0x82,
	0x8a, 0x6e, 0xce, 0xd9, 0x9d, 0x6f, 0xc7, 0x67, 0x64, 0xe8, 0xae, 0x8a, 0x5c, 0xe4, 0x17, 0x6f,
	0x44, 0xbc, 0x1c, 0xcb, 0x1a, 0x5d, 0x29, 0xa2, 0x47, 0x70, 0x1e, 0x88, 0x38, 0x22, 0x38, 0x59,
	0xb2, 0x24, 0x66, 0x0d, 0xad, 0xd1, 0x51, 0x2c, 0x6b, 0xec, 0x80, 0xbd, 0x2e, 0x16, 0xac, 0x21,
	0xad, 0xaa, 0xc4, 0x10, 0xfe, 0xcf, 0x13, 0x41, 0xef, 0x79, 0xb1, 0x61, 0xb6, 0xb4, 0x6b, 0x1d,
	0x7d, 0x59, 0xe0, 0xde, 0x67, 0xa2, 0xd8, 0x60, 0x1b, 0x1a, 0x29, 0xd7, 0xa4, 0x46, 0xca, 0x71,
	0x00, 0x1e, 0xcf, 0x97, 0x49, 0x9a, 0x69, 0x94, 0x56, 0x86, 0x6f, 0x6f, 0xf9, 0x3d, 0x70, 0x45,
	0x2a, 0x16, 0xc4, 0x1c, 0xe9, 0x29, 0x81, 0x0c, 0xfc, 0x79, 0x9e, 0x09, 0xca, 0x04, 0x73, 0xa5,
	0x6f, 0x64, 0x35, 0x35, 0x4f, 0x04, 0x31, 0x6f, 0x68, 0x8d, 0xec, 0x58, 0xd6, 0x7b, 0x33, 0xfa,
	0x07, 0x33, 0x3e, 0x03, 0xcc, 0x38, 0x8f, 0xe9, 0x63, 0x4d, 0xa5, 0xf8, 0x83, 0x6f, 0x6e, 0x41,
	0x20, 0x79, 0xe5, 0x2a, 0xcf, 0x4a, 0x8a, 0x22, 0x68, 0x57, 0x09, 0xa4, 0x54, 0x9a, 0x27, 0x34,
	0xce, 0xaa, 0x71, 0xd1, 0x35, 0x1c, 0xd7, 0x77, 0x54, 0x1b, 0x9e, 0x83, 0x4f, 0xca, 0x62, 0xd6,
	0xd0, 0x1e, 0x05, 0xd3, 0xe6, 0x58, 0x6d, 0x4a, 0xc6, 0x19, 0x9b, 0xc3, 0xea, 0xb5, 0xa7, 0xb4,
	0x14, 0x9a, 0x1d, 0x4d, 0xa0, 0xa9, 0xa4, 0xc6, 0x9c, 0x82, 0xda, 0xa9, 0x86, 0x04, 0x1a, 0x52,
	0xad, 0x37, 0xd6, 0xdb, 0x3e, 0x83, 0x56, 0x4c, 0xcb, 0xfc, 0x93, 0x7e, 0x89, 0x20, 0xea, 0x40,
	0xdb, 0x5c, 0x52, 0xe4, 0xe9, 0xb7, 0x05, 0x6e, 0x85, 0x29, 0x71, 0x0c, 0xf6, 0x8c, 0x73, 0xec,
	0x6a, 0xf6, 0x36, 0xcc, 0x10, 0x77, 0x2d, 0x9d, 0xc7, 0x3f, 0xbc, 0x02, 0x4f, 0xb1, 0xb0, 0xa7,
	0xcf, 0xf7, 0xde, 0x0f, 0xfb, 0x07, 0x6e, 0xdd, 0x78, 0x03, 0xbe, 0x8e, 0x09, 0xfb, 0x3b, 0x69,
	0x6c, 0xa3, 0x0d, 0x07, 0x87, 0x76, 0xdd, 0x3b, 0x01, 0xa7, 0x0a, 0x06, 0xcd, 0x48, 0x3b, 0xa1,
	0x85, 0x27, 0x7b, 0x9e, 0x69, 0xb9, 0x6b, 0xbd, 0x04, 0xf2, 0xb7, 0xb8, 0x95, 0xa7, 0xaf, 0x9e,
	0x14, 0x97, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe4, 0x42, 0x5d, 0x93, 0x38, 0x03, 0x00, 0x00,
}
