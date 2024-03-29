// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.12.4
// source: proto/movie.proto

package movie

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MovieInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PosterPath       string  `protobuf:"bytes,1,opt,name=poster_path,json=posterPath,proto3" json:"poster_path,omitempty"`
	Adult            bool    `protobuf:"varint,2,opt,name=adult,proto3" json:"adult,omitempty"`
	Overview         string  `protobuf:"bytes,3,opt,name=overview,proto3" json:"overview,omitempty"`
	ReleaseDate      string  `protobuf:"bytes,4,opt,name=release_date,json=releaseDate,proto3" json:"release_date,omitempty"`
	GenreIds         []int32 `protobuf:"varint,5,rep,packed,name=genre_ids,json=genreIds,proto3" json:"genre_ids,omitempty"`
	Id               int32   `protobuf:"varint,6,opt,name=id,proto3" json:"id,omitempty"`
	OriginalTitle    string  `protobuf:"bytes,7,opt,name=original_title,json=originalTitle,proto3" json:"original_title,omitempty"`
	OriginalLanguage string  `protobuf:"bytes,8,opt,name=original_language,json=originalLanguage,proto3" json:"original_language,omitempty"`
	Title            string  `protobuf:"bytes,9,opt,name=title,proto3" json:"title,omitempty"`
	BackdropPath     string  `protobuf:"bytes,10,opt,name=backdrop_path,json=backdropPath,proto3" json:"backdrop_path,omitempty"`
	Popularity       float64 `protobuf:"fixed64,11,opt,name=popularity,proto3" json:"popularity,omitempty"`
	VoteCount        int32   `protobuf:"varint,12,opt,name=vote_count,json=voteCount,proto3" json:"vote_count,omitempty"`
	Video            bool    `protobuf:"varint,13,opt,name=video,proto3" json:"video,omitempty"`
	VoteAverage      float64 `protobuf:"fixed64,14,opt,name=vote_average,json=voteAverage,proto3" json:"vote_average,omitempty"`
}

func (x *MovieInfo) Reset() {
	*x = MovieInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_movie_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MovieInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MovieInfo) ProtoMessage() {}

func (x *MovieInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_movie_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MovieInfo.ProtoReflect.Descriptor instead.
func (*MovieInfo) Descriptor() ([]byte, []int) {
	return file_proto_movie_proto_rawDescGZIP(), []int{0}
}

func (x *MovieInfo) GetPosterPath() string {
	if x != nil {
		return x.PosterPath
	}
	return ""
}

func (x *MovieInfo) GetAdult() bool {
	if x != nil {
		return x.Adult
	}
	return false
}

func (x *MovieInfo) GetOverview() string {
	if x != nil {
		return x.Overview
	}
	return ""
}

func (x *MovieInfo) GetReleaseDate() string {
	if x != nil {
		return x.ReleaseDate
	}
	return ""
}

func (x *MovieInfo) GetGenreIds() []int32 {
	if x != nil {
		return x.GenreIds
	}
	return nil
}

func (x *MovieInfo) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MovieInfo) GetOriginalTitle() string {
	if x != nil {
		return x.OriginalTitle
	}
	return ""
}

func (x *MovieInfo) GetOriginalLanguage() string {
	if x != nil {
		return x.OriginalLanguage
	}
	return ""
}

func (x *MovieInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *MovieInfo) GetBackdropPath() string {
	if x != nil {
		return x.BackdropPath
	}
	return ""
}

func (x *MovieInfo) GetPopularity() float64 {
	if x != nil {
		return x.Popularity
	}
	return 0
}

func (x *MovieInfo) GetVoteCount() int32 {
	if x != nil {
		return x.VoteCount
	}
	return 0
}

func (x *MovieInfo) GetVideo() bool {
	if x != nil {
		return x.Video
	}
	return false
}

func (x *MovieInfo) GetVoteAverage() float64 {
	if x != nil {
		return x.VoteAverage
	}
	return 0
}

// Search for movies by simple text search
type SearchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// a ISO 639-1 value to display translated data
	Language string `protobuf:"bytes,1,opt,name=language,proto3" json:"language,omitempty"`
	// a text query to search
	Query string `protobuf:"bytes,2,opt,name=query,proto3" json:"query,omitempty"`
	// page to query
	Page int32 `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
	// a ISO 3166-1 code to filter release dates.
	Region string `protobuf:"bytes,4,opt,name=region,proto3" json:"region,omitempty"`
	// year of making
	Year int32 `protobuf:"varint,5,opt,name=year,proto3" json:"year,omitempty"`
	// year of release
	PrimaryReleaseYear int32 `protobuf:"varint,6,opt,name=primary_release_year,json=primaryReleaseYear,proto3" json:"primary_release_year,omitempty"`
}

func (x *SearchRequest) Reset() {
	*x = SearchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_movie_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchRequest) ProtoMessage() {}

func (x *SearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_movie_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchRequest.ProtoReflect.Descriptor instead.
func (*SearchRequest) Descriptor() ([]byte, []int) {
	return file_proto_movie_proto_rawDescGZIP(), []int{1}
}

func (x *SearchRequest) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *SearchRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *SearchRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchRequest) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *SearchRequest) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *SearchRequest) GetPrimaryReleaseYear() int32 {
	if x != nil {
		return x.PrimaryReleaseYear
	}
	return 0
}

type SearchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalResults int32        `protobuf:"varint,1,opt,name=total_results,json=totalResults,proto3" json:"total_results,omitempty"`
	TotalPages   int32        `protobuf:"varint,2,opt,name=total_pages,json=totalPages,proto3" json:"total_pages,omitempty"`
	Page         int32        `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
	Results      []*MovieInfo `protobuf:"bytes,4,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *SearchResponse) Reset() {
	*x = SearchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_movie_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchResponse) ProtoMessage() {}

func (x *SearchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_movie_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchResponse.ProtoReflect.Descriptor instead.
func (*SearchResponse) Descriptor() ([]byte, []int) {
	return file_proto_movie_proto_rawDescGZIP(), []int{2}
}

func (x *SearchResponse) GetTotalResults() int32 {
	if x != nil {
		return x.TotalResults
	}
	return 0
}

func (x *SearchResponse) GetTotalPages() int32 {
	if x != nil {
		return x.TotalPages
	}
	return 0
}

func (x *SearchResponse) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchResponse) GetResults() []*MovieInfo {
	if x != nil {
		return x.Results
	}
	return nil
}

var File_proto_movie_proto protoreflect.FileDescriptor

var file_proto_movie_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x22, 0xb5, 0x03, 0x0a, 0x09, 0x4d,
	0x6f, 0x76, 0x69, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x6f, 0x73, 0x74,
	0x65, 0x72, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70,
	0x6f, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x74, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x64, 0x75,
	0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x61, 0x64, 0x75, 0x6c, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x6f, 0x76, 0x65, 0x72, 0x76, 0x69, 0x65, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6f, 0x76, 0x65, 0x72, 0x76, 0x69, 0x65, 0x77, 0x12, 0x21, 0x0a, 0x0c, 0x72,
	0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x44, 0x61, 0x74, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x05, 0x52, 0x08, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x49, 0x64, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x54, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x6c,
	0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x62, 0x61, 0x63, 0x6b, 0x64, 0x72, 0x6f,
	0x70, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x62, 0x61,
	0x63, 0x6b, 0x64, 0x72, 0x6f, 0x70, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x6f,
	0x70, 0x75, 0x6c, 0x61, 0x72, 0x69, 0x74, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a,
	0x70, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x72, 0x69, 0x74, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x6f,
	0x74, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x76, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x69, 0x64,
	0x65, 0x6f, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x12,
	0x21, 0x0a, 0x0c, 0x76, 0x6f, 0x74, 0x65, 0x5f, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x18,
	0x0e, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x41, 0x76, 0x65, 0x72, 0x61,
	0x67, 0x65, 0x22, 0xb3, 0x01, 0x0a, 0x0d, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x67, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x79, 0x65, 0x61, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x79, 0x65, 0x61, 0x72, 0x12, 0x30, 0x0a, 0x14, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72,
	0x79, 0x5f, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x79, 0x65, 0x61, 0x72, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x12, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x65, 0x6c,
	0x65, 0x61, 0x73, 0x65, 0x59, 0x65, 0x61, 0x72, 0x22, 0x96, 0x01, 0x0a, 0x0e, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73,
	0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x2a, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x4d,
	0x6f, 0x76, 0x69, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x73, 0x32, 0x40, 0x0a, 0x05, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x53, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x12, 0x14, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x6d,
	0x6f, 0x76, 0x69, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_movie_proto_rawDescOnce sync.Once
	file_proto_movie_proto_rawDescData = file_proto_movie_proto_rawDesc
)

func file_proto_movie_proto_rawDescGZIP() []byte {
	file_proto_movie_proto_rawDescOnce.Do(func() {
		file_proto_movie_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_movie_proto_rawDescData)
	})
	return file_proto_movie_proto_rawDescData
}

var file_proto_movie_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_movie_proto_goTypes = []interface{}{
	(*MovieInfo)(nil),      // 0: movie.MovieInfo
	(*SearchRequest)(nil),  // 1: movie.SearchRequest
	(*SearchResponse)(nil), // 2: movie.SearchResponse
}
var file_proto_movie_proto_depIdxs = []int32{
	0, // 0: movie.SearchResponse.results:type_name -> movie.MovieInfo
	1, // 1: movie.Movie.Search:input_type -> movie.SearchRequest
	2, // 2: movie.Movie.Search:output_type -> movie.SearchResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_movie_proto_init() }
func file_proto_movie_proto_init() {
	if File_proto_movie_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_movie_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MovieInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_movie_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_movie_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_movie_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_movie_proto_goTypes,
		DependencyIndexes: file_proto_movie_proto_depIdxs,
		MessageInfos:      file_proto_movie_proto_msgTypes,
	}.Build()
	File_proto_movie_proto = out.File
	file_proto_movie_proto_rawDesc = nil
	file_proto_movie_proto_goTypes = nil
	file_proto_movie_proto_depIdxs = nil
}
