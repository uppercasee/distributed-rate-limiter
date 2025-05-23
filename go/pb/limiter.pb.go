// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: limiter.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CheckRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ClientId      string                 `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	mi := &file_limiter_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_limiter_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_limiter_proto_rawDescGZIP(), []int{0}
}

func (x *CheckRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

type CheckResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Allowed       bool                   `protobuf:"varint,1,opt,name=allowed,proto3" json:"allowed,omitempty"`
	RetryAfter    int64                  `protobuf:"varint,2,opt,name=retry_after,json=retryAfter,proto3" json:"retry_after,omitempty"` // seconds
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	mi := &file_limiter_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_limiter_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_limiter_proto_rawDescGZIP(), []int{1}
}

func (x *CheckResponse) GetAllowed() bool {
	if x != nil {
		return x.Allowed
	}
	return false
}

func (x *CheckResponse) GetRetryAfter() int64 {
	if x != nil {
		return x.RetryAfter
	}
	return 0
}

var File_limiter_proto protoreflect.FileDescriptor

const file_limiter_proto_rawDesc = "" +
	"\n" +
	"\rlimiter.proto\x12\alimiter\"+\n" +
	"\fCheckRequest\x12\x1b\n" +
	"\tclient_id\x18\x01 \x01(\tR\bclientId\"J\n" +
	"\rCheckResponse\x12\x18\n" +
	"\aallowed\x18\x01 \x01(\bR\aallowed\x12\x1f\n" +
	"\vretry_after\x18\x02 \x01(\x03R\n" +
	"retryAfter2L\n" +
	"\x12RateLimiterService\x126\n" +
	"\x05Check\x12\x15.limiter.CheckRequest\x1a\x16.limiter.CheckResponseB\x06Z\x04./pbb\x06proto3"

var (
	file_limiter_proto_rawDescOnce sync.Once
	file_limiter_proto_rawDescData []byte
)

func file_limiter_proto_rawDescGZIP() []byte {
	file_limiter_proto_rawDescOnce.Do(func() {
		file_limiter_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_limiter_proto_rawDesc), len(file_limiter_proto_rawDesc)))
	})
	return file_limiter_proto_rawDescData
}

var file_limiter_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_limiter_proto_goTypes = []any{
	(*CheckRequest)(nil),  // 0: limiter.CheckRequest
	(*CheckResponse)(nil), // 1: limiter.CheckResponse
}
var file_limiter_proto_depIdxs = []int32{
	0, // 0: limiter.RateLimiterService.Check:input_type -> limiter.CheckRequest
	1, // 1: limiter.RateLimiterService.Check:output_type -> limiter.CheckResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_limiter_proto_init() }
func file_limiter_proto_init() {
	if File_limiter_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_limiter_proto_rawDesc), len(file_limiter_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_limiter_proto_goTypes,
		DependencyIndexes: file_limiter_proto_depIdxs,
		MessageInfos:      file_limiter_proto_msgTypes,
	}.Build()
	File_limiter_proto = out.File
	file_limiter_proto_goTypes = nil
	file_limiter_proto_depIdxs = nil
}
