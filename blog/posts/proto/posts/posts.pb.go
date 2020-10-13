// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/posts/posts.proto

package posts

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

type Post struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title                string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Slug                 string   `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	Content              string   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Created              int64    `protobuf:"varint,5,opt,name=created,proto3" json:"created,omitempty"`
	Updated              int64    `protobuf:"varint,6,opt,name=updated,proto3" json:"updated,omitempty"`
	Author               string   `protobuf:"bytes,7,opt,name=author,proto3" json:"author,omitempty"`
	Tags                 []string `protobuf:"bytes,8,rep,name=tags,proto3" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Post) Reset()         { *m = Post{} }
func (m *Post) String() string { return proto.CompactTextString(m) }
func (*Post) ProtoMessage()    {}
func (*Post) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1e4efc789192621, []int{0}
}

func (m *Post) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Post.Unmarshal(m, b)
}
func (m *Post) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Post.Marshal(b, m, deterministic)
}
func (m *Post) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Post.Merge(m, src)
}
func (m *Post) XXX_Size() int {
	return xxx_messageInfo_Post.Size(m)
}
func (m *Post) XXX_DiscardUnknown() {
	xxx_messageInfo_Post.DiscardUnknown(m)
}

var xxx_messageInfo_Post proto.InternalMessageInfo

func (m *Post) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Post) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Post) GetSlug() string {
	if m != nil {
		return m.Slug
	}
	return ""
}

func (m *Post) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Post) GetCreated() int64 {
	if m != nil {
		return m.Created
	}
	return 0
}

func (m *Post) GetUpdated() int64 {
	if m != nil {
		return m.Updated
	}
	return 0
}

func (m *Post) GetAuthor() string {
	if m != nil {
		return m.Author
	}
	return ""
}

func (m *Post) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

type QueryRequest struct {
	Slug                 string   `protobuf:"bytes,1,opt,name=slug,proto3" json:"slug,omitempty"`
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int64    `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryRequest) Reset()         { *m = QueryRequest{} }
func (m *QueryRequest) String() string { return proto.CompactTextString(m) }
func (*QueryRequest) ProtoMessage()    {}
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1e4efc789192621, []int{1}
}

func (m *QueryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryRequest.Unmarshal(m, b)
}
func (m *QueryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryRequest.Marshal(b, m, deterministic)
}
func (m *QueryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRequest.Merge(m, src)
}
func (m *QueryRequest) XXX_Size() int {
	return xxx_messageInfo_QueryRequest.Size(m)
}
func (m *QueryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRequest proto.InternalMessageInfo

func (m *QueryRequest) GetSlug() string {
	if m != nil {
		return m.Slug
	}
	return ""
}

func (m *QueryRequest) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *QueryRequest) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type QueryResponse struct {
	Posts                []*Post  `protobuf:"bytes,1,rep,name=posts,proto3" json:"posts,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryResponse) Reset()         { *m = QueryResponse{} }
func (m *QueryResponse) String() string { return proto.CompactTextString(m) }
func (*QueryResponse) ProtoMessage()    {}
func (*QueryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1e4efc789192621, []int{2}
}

func (m *QueryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResponse.Unmarshal(m, b)
}
func (m *QueryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResponse.Marshal(b, m, deterministic)
}
func (m *QueryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResponse.Merge(m, src)
}
func (m *QueryResponse) XXX_Size() int {
	return xxx_messageInfo_QueryResponse.Size(m)
}
func (m *QueryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResponse proto.InternalMessageInfo

func (m *QueryResponse) GetPosts() []*Post {
	if m != nil {
		return m.Posts
	}
	return nil
}

type SaveRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title                string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Slug                 string   `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	Content              string   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Timestamp            int64    `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Tags                 []string `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SaveRequest) Reset()         { *m = SaveRequest{} }
func (m *SaveRequest) String() string { return proto.CompactTextString(m) }
func (*SaveRequest) ProtoMessage()    {}
func (*SaveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1e4efc789192621, []int{3}
}

func (m *SaveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SaveRequest.Unmarshal(m, b)
}
func (m *SaveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SaveRequest.Marshal(b, m, deterministic)
}
func (m *SaveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SaveRequest.Merge(m, src)
}
func (m *SaveRequest) XXX_Size() int {
	return xxx_messageInfo_SaveRequest.Size(m)
}
func (m *SaveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SaveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SaveRequest proto.InternalMessageInfo

func (m *SaveRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SaveRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SaveRequest) GetSlug() string {
	if m != nil {
		return m.Slug
	}
	return ""
}

func (m *SaveRequest) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *SaveRequest) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *SaveRequest) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

type SaveResponse struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SaveResponse) Reset()         { *m = SaveResponse{} }
func (m *SaveResponse) String() string { return proto.CompactTextString(m) }
func (*SaveResponse) ProtoMessage()    {}
func (*SaveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1e4efc789192621, []int{4}
}

func (m *SaveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SaveResponse.Unmarshal(m, b)
}
func (m *SaveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SaveResponse.Marshal(b, m, deterministic)
}
func (m *SaveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SaveResponse.Merge(m, src)
}
func (m *SaveResponse) XXX_Size() int {
	return xxx_messageInfo_SaveResponse.Size(m)
}
func (m *SaveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SaveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SaveResponse proto.InternalMessageInfo

func (m *SaveResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type DeleteRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteRequest) Reset()         { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()    {}
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1e4efc789192621, []int{5}
}

func (m *DeleteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteRequest.Unmarshal(m, b)
}
func (m *DeleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteRequest.Marshal(b, m, deterministic)
}
func (m *DeleteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteRequest.Merge(m, src)
}
func (m *DeleteRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteRequest.Size(m)
}
func (m *DeleteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteRequest proto.InternalMessageInfo

func (m *DeleteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type DeleteResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteResponse) Reset()         { *m = DeleteResponse{} }
func (m *DeleteResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteResponse) ProtoMessage()    {}
func (*DeleteResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a1e4efc789192621, []int{6}
}

func (m *DeleteResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteResponse.Unmarshal(m, b)
}
func (m *DeleteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteResponse.Marshal(b, m, deterministic)
}
func (m *DeleteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteResponse.Merge(m, src)
}
func (m *DeleteResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteResponse.Size(m)
}
func (m *DeleteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Post)(nil), "posts.Post")
	proto.RegisterType((*QueryRequest)(nil), "posts.QueryRequest")
	proto.RegisterType((*QueryResponse)(nil), "posts.QueryResponse")
	proto.RegisterType((*SaveRequest)(nil), "posts.SaveRequest")
	proto.RegisterType((*SaveResponse)(nil), "posts.SaveResponse")
	proto.RegisterType((*DeleteRequest)(nil), "posts.DeleteRequest")
	proto.RegisterType((*DeleteResponse)(nil), "posts.DeleteResponse")
}

func init() {
	proto.RegisterFile("proto/posts/posts.proto", fileDescriptor_a1e4efc789192621)
}

var fileDescriptor_a1e4efc789192621 = []byte{
	// 363 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x52, 0x5d, 0x4e, 0xf3, 0x30,
	0x10, 0xfc, 0xd2, 0xfc, 0xf4, 0xeb, 0xf6, 0x47, 0x68, 0x5b, 0xc0, 0xaa, 0x10, 0x94, 0x3c, 0xf5,
	0xa9, 0x88, 0x82, 0xc4, 0x05, 0x38, 0x40, 0x09, 0x27, 0x08, 0xd4, 0x2d, 0x91, 0xd2, 0x26, 0xc4,
	0x0e, 0x12, 0xe7, 0xe0, 0x14, 0x5c, 0x81, 0xd3, 0x61, 0xaf, 0xed, 0x92, 0x56, 0xe2, 0x91, 0x97,
	0x68, 0x67, 0x36, 0x1e, 0xcf, 0x4c, 0x02, 0xa7, 0x65, 0x55, 0xc8, 0xe2, 0xaa, 0x2c, 0x84, 0x14,
	0xe6, 0x39, 0x23, 0x06, 0x43, 0x02, 0xf1, 0x97, 0x07, 0xc1, 0x42, 0x4d, 0x38, 0x80, 0x56, 0xb6,
	0x64, 0xde, 0xc4, 0x9b, 0x76, 0x12, 0x35, 0xe1, 0x08, 0x42, 0x99, 0xc9, 0x9c, 0xb3, 0x16, 0x51,
	0x06, 0x20, 0x42, 0x20, 0xf2, 0x7a, 0xcd, 0x7c, 0x22, 0x69, 0x46, 0x06, 0xed, 0xe7, 0x62, 0x2b,
	0xf9, 0x56, 0xb2, 0x80, 0x68, 0x07, 0x69, 0x53, 0xf1, 0x54, 0xf2, 0x25, 0x0b, 0xd5, 0xc6, 0x4f,
	0x1c, 0xd4, 0x9b, 0xba, 0x5c, 0xd2, 0x26, 0x32, 0x1b, 0x0b, 0xf1, 0x04, 0xa2, 0xb4, 0x96, 0x2f,
	0x45, 0xc5, 0xda, 0x24, 0x66, 0x91, 0xbe, 0x59, 0xa6, 0x6b, 0xc1, 0xfe, 0x4f, 0x7c, 0x7d, 0xb3,
	0x9e, 0xe3, 0x05, 0xf4, 0x1e, 0x6a, 0x5e, 0xbd, 0x27, 0xfc, 0xb5, 0xe6, 0x2a, 0x83, 0x73, 0xe7,
	0x35, 0xdc, 0x29, 0xbd, 0x62, 0xb5, 0x12, 0x5c, 0x52, 0x10, 0x3f, 0xb1, 0x48, 0xe7, 0xcb, 0xb3,
	0x4d, 0x26, 0x29, 0x8a, 0x9f, 0x18, 0x10, 0xcf, 0xa1, 0x6f, 0x15, 0x45, 0x59, 0x6c, 0x05, 0xc7,
	0x4b, 0x30, 0x45, 0x29, 0x4d, 0x7f, 0xda, 0x9d, 0x77, 0x67, 0xa6, 0x43, 0x5d, 0x59, 0x62, 0x2b,
	0xfc, 0xf0, 0xa0, 0xfb, 0x98, 0xbe, 0x71, 0xe7, 0xe2, 0x2f, 0x9a, 0x3c, 0x83, 0x8e, 0xcc, 0x36,
	0x4a, 0x3d, 0xdd, 0x94, 0xb6, 0xcb, 0x1f, 0x62, 0xd7, 0x4d, 0xd4, 0xe8, 0xe6, 0x1c, 0x7a, 0xc6,
	0x94, 0x0d, 0x72, 0xe0, 0x2a, 0xbe, 0x80, 0xfe, 0x3d, 0xcf, 0xb9, 0xfc, 0xcd, 0x76, 0x7c, 0x04,
	0x03, 0xf7, 0x82, 0x91, 0x98, 0x7f, 0x7a, 0x10, 0xea, 0xe0, 0x02, 0x6f, 0x21, 0xa4, 0x9a, 0x70,
	0x68, 0xfb, 0x68, 0x7e, 0x86, 0xf1, 0x68, 0x9f, 0x34, 0xa7, 0xe3, 0x7f, 0x78, 0x0d, 0x81, 0xb6,
	0x84, 0x68, 0xf7, 0x8d, 0xd2, 0xc6, 0xc3, 0x3d, 0x6e, 0x77, 0xe4, 0x0e, 0x22, 0x63, 0x02, 0x9d,
	0xe8, 0x9e, 0xe9, 0xf1, 0xf1, 0x01, 0xeb, 0x0e, 0x3e, 0x45, 0xf4, 0x97, 0xdf, 0x7c, 0x07, 0x00,
	0x00, 0xff, 0xff, 0xed, 0x03, 0x8f, 0x37, 0x00, 0x03, 0x00, 0x00,
}
