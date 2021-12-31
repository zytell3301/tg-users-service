// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0-devel
// 	protoc        v3.12.4
// source: api/pb/UsersService/users-service.proto

package UsersService

import (
	error1 "github.com/zytell3301/tg-users-service/pkg/error"
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

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Lastname     string `protobuf:"bytes,3,opt,name=Lastname,proto3" json:"Lastname,omitempty"`
	Bio          string `protobuf:"bytes,4,opt,name=Bio,proto3" json:"Bio,omitempty"`
	Username     string `protobuf:"bytes,5,opt,name=Username,proto3" json:"Username,omitempty"`
	Phone        string `protobuf:"bytes,6,opt,name=Phone,proto3" json:"Phone,omitempty"`
	OnlineStatus bool   `protobuf:"varint,7,opt,name=Online_status,json=OnlineStatus,proto3" json:"Online_status,omitempty"`
	CreatedAt    int64  `protobuf:"varint,8,opt,name=Created_at,json=CreatedAt,proto3" json:"Created_at,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_pb_UsersService_users_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_api_pb_UsersService_users_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_api_pb_UsersService_users_service_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetLastname() string {
	if x != nil {
		return x.Lastname
	}
	return ""
}

func (x *User) GetBio() string {
	if x != nil {
		return x.Bio
	}
	return ""
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *User) GetOnlineStatus() bool {
	if x != nil {
		return x.OnlineStatus
	}
	return false
}

func (x *User) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

var File_api_pb_UsersService_users_service_proto protoreflect.FileDescriptor

var file_api_pb_UsersService_users_service_proto_rawDesc = []byte{
	0x0a, 0x27, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62, 0x2f, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2d, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x7a, 0x79, 0x74, 0x65, 0x6c,
	0x33, 0x33, 0x30, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x1a, 0x1f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62, 0x2f, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xce, 0x01, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x4c, 0x61, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x4c, 0x61, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x42,
	0x69, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x42, 0x69, 0x6f, 0x12, 0x1a, 0x0a,
	0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x68, 0x6f,
	0x6e, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x12,
	0x23, 0x0a, 0x0d, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x32, 0x4f, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x3f, 0x0a, 0x07, 0x4e, 0x65, 0x77, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1c,
	0x2e, 0x7a, 0x79, 0x74, 0x65, 0x6c, 0x33, 0x33, 0x30, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x73,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x7a,
	0x79, 0x74, 0x65, 0x6c, 0x33, 0x33, 0x30, 0x31, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x7a, 0x79, 0x74, 0x65, 0x6c, 0x6c, 0x33, 0x33, 0x30, 0x31, 0x2f, 0x74, 0x67,
	0x2d, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_pb_UsersService_users_service_proto_rawDescOnce sync.Once
	file_api_pb_UsersService_users_service_proto_rawDescData = file_api_pb_UsersService_users_service_proto_rawDesc
)

func file_api_pb_UsersService_users_service_proto_rawDescGZIP() []byte {
	file_api_pb_UsersService_users_service_proto_rawDescOnce.Do(func() {
		file_api_pb_UsersService_users_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_pb_UsersService_users_service_proto_rawDescData)
	})
	return file_api_pb_UsersService_users_service_proto_rawDescData
}

var file_api_pb_UsersService_users_service_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_pb_UsersService_users_service_proto_goTypes = []interface{}{
	(*User)(nil),         // 0: zytel3301.UsersService.User
	(*error1.Error)(nil), // 1: zytel3301.error.Error
}
var file_api_pb_UsersService_users_service_proto_depIdxs = []int32{
	0, // 0: zytel3301.UsersService.UsersService.NewUser:input_type -> zytel3301.UsersService.User
	1, // 1: zytel3301.UsersService.UsersService.NewUser:output_type -> zytel3301.error.Error
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_pb_UsersService_users_service_proto_init() }
func file_api_pb_UsersService_users_service_proto_init() {
	if File_api_pb_UsersService_users_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_pb_UsersService_users_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
			RawDescriptor: file_api_pb_UsersService_users_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_pb_UsersService_users_service_proto_goTypes,
		DependencyIndexes: file_api_pb_UsersService_users_service_proto_depIdxs,
		MessageInfos:      file_api_pb_UsersService_users_service_proto_msgTypes,
	}.Build()
	File_api_pb_UsersService_users_service_proto = out.File
	file_api_pb_UsersService_users_service_proto_rawDesc = nil
	file_api_pb_UsersService_users_service_proto_goTypes = nil
	file_api_pb_UsersService_users_service_proto_depIdxs = nil
}
