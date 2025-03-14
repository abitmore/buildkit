// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.11.4
// source: github.com/moby/buildkit/api/types/worker.proto

package moby_buildkit_v1_types

import (
	pb "github.com/moby/buildkit/solver/pb"
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

type WorkerRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID              string            `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Labels          map[string]string `protobuf:"bytes,2,rep,name=Labels,proto3" json:"Labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Platforms       []*pb.Platform    `protobuf:"bytes,3,rep,name=platforms,proto3" json:"platforms,omitempty"`
	GCPolicy        []*GCPolicy       `protobuf:"bytes,4,rep,name=GCPolicy,proto3" json:"GCPolicy,omitempty"`
	BuildkitVersion *BuildkitVersion  `protobuf:"bytes,5,opt,name=BuildkitVersion,proto3" json:"BuildkitVersion,omitempty"`
}

func (x *WorkerRecord) Reset() {
	*x = WorkerRecord{}
	mi := &file_github_com_moby_buildkit_api_types_worker_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WorkerRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerRecord) ProtoMessage() {}

func (x *WorkerRecord) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_moby_buildkit_api_types_worker_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerRecord.ProtoReflect.Descriptor instead.
func (*WorkerRecord) Descriptor() ([]byte, []int) {
	return file_github_com_moby_buildkit_api_types_worker_proto_rawDescGZIP(), []int{0}
}

func (x *WorkerRecord) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *WorkerRecord) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *WorkerRecord) GetPlatforms() []*pb.Platform {
	if x != nil {
		return x.Platforms
	}
	return nil
}

func (x *WorkerRecord) GetGCPolicy() []*GCPolicy {
	if x != nil {
		return x.GCPolicy
	}
	return nil
}

func (x *WorkerRecord) GetBuildkitVersion() *BuildkitVersion {
	if x != nil {
		return x.BuildkitVersion
	}
	return nil
}

type GCPolicy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	All          bool     `protobuf:"varint,1,opt,name=all,proto3" json:"all,omitempty"`
	KeepDuration int64    `protobuf:"varint,2,opt,name=keepDuration,proto3" json:"keepDuration,omitempty"`
	Filters      []string `protobuf:"bytes,4,rep,name=filters,proto3" json:"filters,omitempty"`
	// reservedSpace was renamed from freeBytes
	ReservedSpace int64 `protobuf:"varint,3,opt,name=reservedSpace,proto3" json:"reservedSpace,omitempty"`
	MaxUsedSpace  int64 `protobuf:"varint,5,opt,name=maxUsedSpace,proto3" json:"maxUsedSpace,omitempty"`
	MinFreeSpace  int64 `protobuf:"varint,6,opt,name=minFreeSpace,proto3" json:"minFreeSpace,omitempty"`
}

func (x *GCPolicy) Reset() {
	*x = GCPolicy{}
	mi := &file_github_com_moby_buildkit_api_types_worker_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GCPolicy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GCPolicy) ProtoMessage() {}

func (x *GCPolicy) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_moby_buildkit_api_types_worker_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GCPolicy.ProtoReflect.Descriptor instead.
func (*GCPolicy) Descriptor() ([]byte, []int) {
	return file_github_com_moby_buildkit_api_types_worker_proto_rawDescGZIP(), []int{1}
}

func (x *GCPolicy) GetAll() bool {
	if x != nil {
		return x.All
	}
	return false
}

func (x *GCPolicy) GetKeepDuration() int64 {
	if x != nil {
		return x.KeepDuration
	}
	return 0
}

func (x *GCPolicy) GetFilters() []string {
	if x != nil {
		return x.Filters
	}
	return nil
}

func (x *GCPolicy) GetReservedSpace() int64 {
	if x != nil {
		return x.ReservedSpace
	}
	return 0
}

func (x *GCPolicy) GetMaxUsedSpace() int64 {
	if x != nil {
		return x.MaxUsedSpace
	}
	return 0
}

func (x *GCPolicy) GetMinFreeSpace() int64 {
	if x != nil {
		return x.MinFreeSpace
	}
	return 0
}

type BuildkitVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Package  string `protobuf:"bytes,1,opt,name=package,proto3" json:"package,omitempty"`
	Version  string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Revision string `protobuf:"bytes,3,opt,name=revision,proto3" json:"revision,omitempty"`
}

func (x *BuildkitVersion) Reset() {
	*x = BuildkitVersion{}
	mi := &file_github_com_moby_buildkit_api_types_worker_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BuildkitVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildkitVersion) ProtoMessage() {}

func (x *BuildkitVersion) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_moby_buildkit_api_types_worker_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildkitVersion.ProtoReflect.Descriptor instead.
func (*BuildkitVersion) Descriptor() ([]byte, []int) {
	return file_github_com_moby_buildkit_api_types_worker_proto_rawDescGZIP(), []int{2}
}

func (x *BuildkitVersion) GetPackage() string {
	if x != nil {
		return x.Package
	}
	return ""
}

func (x *BuildkitVersion) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *BuildkitVersion) GetRevision() string {
	if x != nil {
		return x.Revision
	}
	return ""
}

var File_github_com_moby_buildkit_api_types_worker_proto protoreflect.FileDescriptor

var file_github_com_moby_buildkit_api_types_worker_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x62,
	0x79, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x16, 0x6d, 0x6f, 0x62, 0x79, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x1a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x62, 0x79, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64,
	0x6b, 0x69, 0x74, 0x2f, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x62, 0x2f, 0x6f, 0x70,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe0, 0x02, 0x0a, 0x0c, 0x57, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x48, 0x0a, 0x06, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x6d, 0x6f, 0x62, 0x79, 0x2e,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x4c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x12, 0x2a, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x52, 0x09, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x73, 0x12, 0x3c,
	0x0a, 0x08, 0x47, 0x43, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x6d, 0x6f, 0x62, 0x79, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x47, 0x43, 0x50, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x52, 0x08, 0x47, 0x43, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x51, 0x0a, 0x0f,
	0x42, 0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x6d, 0x6f, 0x62, 0x79, 0x2e, 0x62, 0x75, 0x69,
	0x6c, 0x64, 0x6b, 0x69, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x42,
	0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0f,
	0x42, 0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a,
	0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xc8, 0x01, 0x0a, 0x08, 0x47,
	0x43, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x6c, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x61, 0x6c, 0x6c, 0x12, 0x22, 0x0a, 0x0c, 0x6b, 0x65, 0x65,
	0x70, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0c, 0x6b, 0x65, 0x65, 0x70, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a,
	0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x64, 0x53, 0x70, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d,
	0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x53, 0x70, 0x61, 0x63, 0x65, 0x12, 0x22, 0x0a,
	0x0c, 0x6d, 0x61, 0x78, 0x55, 0x73, 0x65, 0x64, 0x53, 0x70, 0x61, 0x63, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0c, 0x6d, 0x61, 0x78, 0x55, 0x73, 0x65, 0x64, 0x53, 0x70, 0x61, 0x63,
	0x65, 0x12, 0x22, 0x0a, 0x0c, 0x6d, 0x69, 0x6e, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6d, 0x69, 0x6e, 0x46, 0x72, 0x65, 0x65,
	0x53, 0x70, 0x61, 0x63, 0x65, 0x22, 0x61, 0x0a, 0x0f, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69,
	0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x63, 0x6b,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08,
	0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x3b, 0x5a, 0x39, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x62, 0x79, 0x2f, 0x62, 0x75, 0x69, 0x6c,
	0x64, 0x6b, 0x69, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x3b, 0x6d,
	0x6f, 0x62, 0x79, 0x5f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x6b, 0x69, 0x74, 0x5f, 0x76, 0x31, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_moby_buildkit_api_types_worker_proto_rawDescOnce sync.Once
	file_github_com_moby_buildkit_api_types_worker_proto_rawDescData = file_github_com_moby_buildkit_api_types_worker_proto_rawDesc
)

func file_github_com_moby_buildkit_api_types_worker_proto_rawDescGZIP() []byte {
	file_github_com_moby_buildkit_api_types_worker_proto_rawDescOnce.Do(func() {
		file_github_com_moby_buildkit_api_types_worker_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_moby_buildkit_api_types_worker_proto_rawDescData)
	})
	return file_github_com_moby_buildkit_api_types_worker_proto_rawDescData
}

var file_github_com_moby_buildkit_api_types_worker_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_github_com_moby_buildkit_api_types_worker_proto_goTypes = []any{
	(*WorkerRecord)(nil),    // 0: moby.buildkit.v1.types.WorkerRecord
	(*GCPolicy)(nil),        // 1: moby.buildkit.v1.types.GCPolicy
	(*BuildkitVersion)(nil), // 2: moby.buildkit.v1.types.BuildkitVersion
	nil,                     // 3: moby.buildkit.v1.types.WorkerRecord.LabelsEntry
	(*pb.Platform)(nil),     // 4: pb.Platform
}
var file_github_com_moby_buildkit_api_types_worker_proto_depIdxs = []int32{
	3, // 0: moby.buildkit.v1.types.WorkerRecord.Labels:type_name -> moby.buildkit.v1.types.WorkerRecord.LabelsEntry
	4, // 1: moby.buildkit.v1.types.WorkerRecord.platforms:type_name -> pb.Platform
	1, // 2: moby.buildkit.v1.types.WorkerRecord.GCPolicy:type_name -> moby.buildkit.v1.types.GCPolicy
	2, // 3: moby.buildkit.v1.types.WorkerRecord.BuildkitVersion:type_name -> moby.buildkit.v1.types.BuildkitVersion
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_github_com_moby_buildkit_api_types_worker_proto_init() }
func file_github_com_moby_buildkit_api_types_worker_proto_init() {
	if File_github_com_moby_buildkit_api_types_worker_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_moby_buildkit_api_types_worker_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_moby_buildkit_api_types_worker_proto_goTypes,
		DependencyIndexes: file_github_com_moby_buildkit_api_types_worker_proto_depIdxs,
		MessageInfos:      file_github_com_moby_buildkit_api_types_worker_proto_msgTypes,
	}.Build()
	File_github_com_moby_buildkit_api_types_worker_proto = out.File
	file_github_com_moby_buildkit_api_types_worker_proto_rawDesc = nil
	file_github_com_moby_buildkit_api_types_worker_proto_goTypes = nil
	file_github_com_moby_buildkit_api_types_worker_proto_depIdxs = nil
}
