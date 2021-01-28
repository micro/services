// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/posts.proto

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
	Id                   string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title                string            `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Slug                 string            `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	Content              string            `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Created              int64             `protobuf:"varint,5,opt,name=created,proto3" json:"created,omitempty"`
	Updated              int64             `protobuf:"varint,6,opt,name=updated,proto3" json:"updated,omitempty"`
	Author               string            `protobuf:"bytes,7,opt,name=author,proto3" json:"author,omitempty"`
	Tags                 []string          `protobuf:"bytes,8,rep,name=tags,proto3" json:"tags,omitempty"`
	Metadata             map[string]string `protobuf:"bytes,9,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Image                string            `protobuf:"bytes,19,opt,name=image,proto3" json:"image,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Post) Reset()         { *m = Post{} }
func (m *Post) String() string { return proto.CompactTextString(m) }
func (*Post) ProtoMessage()    {}
func (*Post) Descriptor() ([]byte, []int) {
	return fileDescriptor_e93dc7d934d9dc10, []int{0}
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

func (m *Post) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Post) GetImage() string {
	if m != nil {
		return m.Image
	}
	return ""
}

// Query posts. Acts as a listing when no id or slug provided.
// Gets a single post by id or slug if any of them provided.
type QueryRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Slug                 string   `protobuf:"bytes,2,opt,name=slug,proto3" json:"slug,omitempty"`
	Tag                  string   `protobuf:"bytes,3,opt,name=tag,proto3" json:"tag,omitempty"`
	Offset               int64    `protobuf:"varint,4,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int64    `protobuf:"varint,5,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryRequest) Reset()         { *m = QueryRequest{} }
func (m *QueryRequest) String() string { return proto.CompactTextString(m) }
func (*QueryRequest) ProtoMessage()    {}
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e93dc7d934d9dc10, []int{1}
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

func (m *QueryRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *QueryRequest) GetSlug() string {
	if m != nil {
		return m.Slug
	}
	return ""
}

func (m *QueryRequest) GetTag() string {
	if m != nil {
		return m.Tag
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
	return fileDescriptor_e93dc7d934d9dc10, []int{2}
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
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title     string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Slug      string `protobuf:"bytes,3,opt,name=slug,proto3" json:"slug,omitempty"`
	Content   string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Timestamp int64  `protobuf:"varint,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// When updating a post and wanting to delete all tags,
	// send a list of tags with only one member being an empty string [""]
	Tags                 []string          `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty"`
	Metadata             map[string]string `protobuf:"bytes,7,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Image                string            `protobuf:"bytes,8,opt,name=image,proto3" json:"image,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SaveRequest) Reset()         { *m = SaveRequest{} }
func (m *SaveRequest) String() string { return proto.CompactTextString(m) }
func (*SaveRequest) ProtoMessage()    {}
func (*SaveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e93dc7d934d9dc10, []int{3}
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

func (m *SaveRequest) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *SaveRequest) GetImage() string {
	if m != nil {
		return m.Image
	}
	return ""
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
	return fileDescriptor_e93dc7d934d9dc10, []int{4}
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
	return fileDescriptor_e93dc7d934d9dc10, []int{5}
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
	return fileDescriptor_e93dc7d934d9dc10, []int{6}
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
	proto.RegisterMapType((map[string]string)(nil), "posts.Post.MetadataEntry")
	proto.RegisterType((*QueryRequest)(nil), "posts.QueryRequest")
	proto.RegisterType((*QueryResponse)(nil), "posts.QueryResponse")
	proto.RegisterType((*SaveRequest)(nil), "posts.SaveRequest")
	proto.RegisterMapType((map[string]string)(nil), "posts.SaveRequest.MetadataEntry")
	proto.RegisterType((*SaveResponse)(nil), "posts.SaveResponse")
	proto.RegisterType((*DeleteRequest)(nil), "posts.DeleteRequest")
	proto.RegisterType((*DeleteResponse)(nil), "posts.DeleteResponse")
}

func init() { proto.RegisterFile("proto/posts.proto", fileDescriptor_e93dc7d934d9dc10) }

var fileDescriptor_e93dc7d934d9dc10 = []byte{
	// 452 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xcb, 0x8e, 0xd3, 0x40,
	0x10, 0xc4, 0x76, 0xec, 0x24, 0x9d, 0xcd, 0x2a, 0x74, 0x16, 0x34, 0x44, 0x08, 0x8c, 0x4f, 0x39,
	0x05, 0x11, 0x40, 0x20, 0xe0, 0x08, 0x47, 0x24, 0x30, 0x5f, 0x30, 0xe0, 0xde, 0x60, 0x61, 0xc7,
	0xc6, 0xd3, 0x5e, 0x29, 0xff, 0xc3, 0x85, 0xff, 0xe0, 0xc3, 0xd0, 0x3c, 0x1c, 0xec, 0x85, 0x15,
	0x97, 0xdc, 0xba, 0xaa, 0xdd, 0xd3, 0xd5, 0x55, 0x51, 0xe0, 0x76, 0xdd, 0x54, 0x5c, 0x3d, 0xae,
	0x2b, 0xc5, 0x6a, 0x63, 0x6a, 0x0c, 0x0d, 0x48, 0x7e, 0xf9, 0x30, 0xfa, 0x50, 0x29, 0xc6, 0x73,
	0xf0, 0xf3, 0x4c, 0x78, 0xb1, 0xb7, 0x9e, 0xa6, 0x7e, 0x9e, 0xe1, 0x05, 0x84, 0x9c, 0x73, 0x41,
	0xc2, 0x37, 0x94, 0x05, 0x88, 0x30, 0x52, 0x45, 0xbb, 0x13, 0x81, 0x21, 0x4d, 0x8d, 0x02, 0xc6,
	0x5f, 0xaa, 0x3d, 0xd3, 0x9e, 0xc5, 0xc8, 0xd0, 0x1d, 0x34, 0x9d, 0x86, 0x24, 0x53, 0x26, 0xc2,
	0xd8, 0x5b, 0x07, 0x69, 0x07, 0x75, 0xa7, 0xad, 0x33, 0xd3, 0x89, 0x6c, 0xc7, 0x41, 0xbc, 0x0b,
	0x91, 0x6c, 0xf9, 0x6b, 0xd5, 0x88, 0xb1, 0x79, 0xcc, 0x21, 0xbd, 0x99, 0xe5, 0x4e, 0x89, 0x49,
	0x1c, 0xe8, 0xcd, 0xba, 0xc6, 0xe7, 0x30, 0x29, 0x89, 0x65, 0x26, 0x59, 0x8a, 0x69, 0x1c, 0xac,
	0x67, 0xdb, 0x7b, 0x1b, 0x7b, 0xa3, 0x3e, 0x69, 0xf3, 0xde, 0xf5, 0xde, 0xed, 0xb9, 0x39, 0xa4,
	0xc7, 0x4f, 0xf5, 0x69, 0x79, 0x29, 0x77, 0x24, 0x96, 0xf6, 0x34, 0x03, 0x56, 0xaf, 0x61, 0x3e,
	0x18, 0xc0, 0x05, 0x04, 0xdf, 0xe8, 0xe0, 0x2c, 0xd1, 0xa5, 0x1e, 0xbc, 0x92, 0x45, 0x7b, 0xf4,
	0xc4, 0x80, 0x57, 0xfe, 0x4b, 0x2f, 0x69, 0xe0, 0xec, 0x63, 0x4b, 0xcd, 0x21, 0xa5, 0xef, 0x2d,
	0xfd, 0xc3, 0xcd, 0xce, 0x37, 0xbf, 0xe7, 0xdb, 0x02, 0x02, 0x96, 0x9d, 0x95, 0xba, 0xd4, 0xb7,
	0x57, 0x97, 0x97, 0x8a, 0xac, 0x91, 0x41, 0xea, 0x90, 0xde, 0x5b, 0xe4, 0x65, 0xce, 0xce, 0x45,
	0x0b, 0x92, 0x2d, 0xcc, 0xdd, 0x4e, 0x55, 0x57, 0x7b, 0x45, 0xf8, 0x08, 0x6c, 0xa8, 0xc2, 0x33,
	0x5e, 0xcc, 0x7a, 0x5e, 0xa4, 0x2e, 0xee, 0x1f, 0x3e, 0xcc, 0x3e, 0xc9, 0x2b, 0xba, 0x49, 0xe7,
	0x29, 0x52, 0xbf, 0x0f, 0x53, 0xce, 0x4b, 0x52, 0x2c, 0xcb, 0xda, 0x29, 0xfe, 0x43, 0x1c, 0x73,
	0x8c, 0x7a, 0x39, 0xbe, 0xe9, 0xe5, 0x38, 0x36, 0xda, 0x63, 0xa7, 0xbd, 0xa7, 0xf5, 0xff, 0x71,
	0x4e, 0x4e, 0x16, 0xe7, 0x03, 0x38, 0xb3, 0x9b, 0x9d, 0xb3, 0xd7, 0x6c, 0x4a, 0x1e, 0xc2, 0xfc,
	0x2d, 0x15, 0xc4, 0x37, 0xf9, 0x98, 0x2c, 0xe0, 0xbc, 0xfb, 0xc0, 0x3e, 0xb1, 0xfd, 0xe9, 0x41,
	0xa8, 0x93, 0x50, 0xf8, 0x0c, 0x42, 0x93, 0x1b, 0x2e, 0xdd, 0x91, 0xfd, 0x5f, 0xce, 0xea, 0x62,
	0x48, 0xda, 0xe9, 0xe4, 0x16, 0x3e, 0x81, 0x91, 0x96, 0x84, 0xf8, 0xb7, 0x33, 0xab, 0xe5, 0x80,
	0x3b, 0x8e, 0xbc, 0x80, 0xc8, 0x8a, 0xc0, 0xee, 0xd1, 0x81, 0xe8, 0xd5, 0x9d, 0x6b, 0x6c, 0x37,
	0xf8, 0x39, 0x32, 0x7f, 0x11, 0x4f, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0x17, 0x84, 0x7a, 0xa4,
	0x37, 0x04, 0x00, 0x00,
}
