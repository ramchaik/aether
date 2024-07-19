// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.1
// source: project_log.proto

package project_log

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

type PushLogsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId string `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Log       string `protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`
	Timestamp int64  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"` // Epoch timestamp in seconds
}

func (x *PushLogsRequest) Reset() {
	*x = PushLogsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_project_log_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushLogsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushLogsRequest) ProtoMessage() {}

func (x *PushLogsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_project_log_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushLogsRequest.ProtoReflect.Descriptor instead.
func (*PushLogsRequest) Descriptor() ([]byte, []int) {
	return file_project_log_proto_rawDescGZIP(), []int{0}
}

func (x *PushLogsRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *PushLogsRequest) GetLog() string {
	if x != nil {
		return x.Log
	}
	return ""
}

func (x *PushLogsRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type PushLogsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PushLogsResponse) Reset() {
	*x = PushLogsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_project_log_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushLogsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushLogsResponse) ProtoMessage() {}

func (x *PushLogsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_project_log_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushLogsResponse.ProtoReflect.Descriptor instead.
func (*PushLogsResponse) Descriptor() ([]byte, []int) {
	return file_project_log_proto_rawDescGZIP(), []int{1}
}

func (x *PushLogsResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *PushLogsResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_project_log_proto protoreflect.FileDescriptor

var file_project_log_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6c, 0x6f, 0x67,
	0x22, 0x60, 0x0a, 0x0f, 0x50, 0x75, 0x73, 0x68, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6c, 0x6f, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x22, 0x46, 0x0a, 0x10, 0x50, 0x75, 0x73, 0x68, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x5e, 0x0a, 0x11, 0x50, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4c, 0x6f, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x49, 0x0a, 0x08, 0x50, 0x75, 0x73, 0x68, 0x4c, 0x6f, 0x67, 0x73, 0x12, 0x1c, 0x2e, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4c, 0x6f,
	0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4c, 0x6f, 0x67, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_project_log_proto_rawDescOnce sync.Once
	file_project_log_proto_rawDescData = file_project_log_proto_rawDesc
)

func file_project_log_proto_rawDescGZIP() []byte {
	file_project_log_proto_rawDescOnce.Do(func() {
		file_project_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_project_log_proto_rawDescData)
	})
	return file_project_log_proto_rawDescData
}

var file_project_log_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_project_log_proto_goTypes = []interface{}{
	(*PushLogsRequest)(nil),  // 0: project_log.PushLogsRequest
	(*PushLogsResponse)(nil), // 1: project_log.PushLogsResponse
}
var file_project_log_proto_depIdxs = []int32{
	0, // 0: project_log.ProjectLogService.PushLogs:input_type -> project_log.PushLogsRequest
	1, // 1: project_log.ProjectLogService.PushLogs:output_type -> project_log.PushLogsResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_project_log_proto_init() }
func file_project_log_proto_init() {
	if File_project_log_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_project_log_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushLogsRequest); i {
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
		file_project_log_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushLogsResponse); i {
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
			RawDescriptor: file_project_log_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_project_log_proto_goTypes,
		DependencyIndexes: file_project_log_proto_depIdxs,
		MessageInfos:      file_project_log_proto_msgTypes,
	}.Build()
	File_project_log_proto = out.File
	file_project_log_proto_rawDesc = nil
	file_project_log_proto_goTypes = nil
	file_project_log_proto_depIdxs = nil
}
