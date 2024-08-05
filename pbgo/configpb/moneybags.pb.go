// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: moneybags.proto

package configpb

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

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PostgresDb     *PostgresDB     `protobuf:"bytes,1,opt,name=postgres_db,json=postgresDb,proto3" json:"postgres_db,omitempty"`
	TelegramConfig *TelegramConfig `protobuf:"bytes,2,opt,name=telegram_config,json=telegramConfig,proto3" json:"telegram_config,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_moneybags_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_moneybags_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_moneybags_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetPostgresDb() *PostgresDB {
	if x != nil {
		return x.PostgresDb
	}
	return nil
}

func (x *Config) GetTelegramConfig() *TelegramConfig {
	if x != nil {
		return x.TelegramConfig
	}
	return nil
}

type PostgresDB struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host     string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	User     string `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	DbName   string `protobuf:"bytes,4,opt,name=db_name,json=dbName,proto3" json:"db_name,omitempty"`
	Port     uint32 `protobuf:"varint,5,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *PostgresDB) Reset() {
	*x = PostgresDB{}
	if protoimpl.UnsafeEnabled {
		mi := &file_moneybags_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostgresDB) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostgresDB) ProtoMessage() {}

func (x *PostgresDB) ProtoReflect() protoreflect.Message {
	mi := &file_moneybags_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostgresDB.ProtoReflect.Descriptor instead.
func (*PostgresDB) Descriptor() ([]byte, []int) {
	return file_moneybags_proto_rawDescGZIP(), []int{1}
}

func (x *PostgresDB) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *PostgresDB) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *PostgresDB) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *PostgresDB) GetDbName() string {
	if x != nil {
		return x.DbName
	}
	return ""
}

func (x *PostgresDB) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type TelegramConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApiKey      string `protobuf:"bytes,1,opt,name=api_key,json=apiKey,proto3" json:"api_key,omitempty"`
	Debug       bool   `protobuf:"varint,2,opt,name=debug,proto3" json:"debug,omitempty"`
	AllowedUser string `protobuf:"bytes,3,opt,name=allowed_user,json=allowedUser,proto3" json:"allowed_user,omitempty"`
}

func (x *TelegramConfig) Reset() {
	*x = TelegramConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_moneybags_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TelegramConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TelegramConfig) ProtoMessage() {}

func (x *TelegramConfig) ProtoReflect() protoreflect.Message {
	mi := &file_moneybags_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TelegramConfig.ProtoReflect.Descriptor instead.
func (*TelegramConfig) Descriptor() ([]byte, []int) {
	return file_moneybags_proto_rawDescGZIP(), []int{2}
}

func (x *TelegramConfig) GetApiKey() string {
	if x != nil {
		return x.ApiKey
	}
	return ""
}

func (x *TelegramConfig) GetDebug() bool {
	if x != nil {
		return x.Debug
	}
	return false
}

func (x *TelegramConfig) GetAllowedUser() string {
	if x != nil {
		return x.AllowedUser
	}
	return ""
}

var File_moneybags_proto protoreflect.FileDescriptor

var file_moneybags_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x62, 0x61, 0x67, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x7e, 0x0a, 0x06, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x12, 0x33, 0x0a, 0x0b, 0x70, 0x6f, 0x73, 0x74, 0x67, 0x72, 0x65, 0x73, 0x5f,
	0x64, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x67, 0x72, 0x65, 0x73, 0x44, 0x42, 0x52, 0x0a, 0x70, 0x6f,
	0x73, 0x74, 0x67, 0x72, 0x65, 0x73, 0x44, 0x62, 0x12, 0x3f, 0x0a, 0x0f, 0x74, 0x65, 0x6c, 0x65,
	0x67, 0x72, 0x61, 0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x54, 0x65, 0x6c, 0x65, 0x67,
	0x72, 0x61, 0x6d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x0e, 0x74, 0x65, 0x6c, 0x65, 0x67,
	0x72, 0x61, 0x6d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x7d, 0x0a, 0x0a, 0x50, 0x6f, 0x73,
	0x74, 0x67, 0x72, 0x65, 0x73, 0x44, 0x42, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x64,
	0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x62,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x62, 0x0a, 0x0e, 0x54, 0x65, 0x6c, 0x65,
	0x67, 0x72, 0x61, 0x6d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x70,
	0x69, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x70, 0x69,
	0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x6c, 0x6c,
	0x6f, 0x77, 0x65, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x42, 0x11, 0x5a, 0x0f,
	0x2e, 0x2f, 0x70, 0x62, 0x67, 0x6f, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_moneybags_proto_rawDescOnce sync.Once
	file_moneybags_proto_rawDescData = file_moneybags_proto_rawDesc
)

func file_moneybags_proto_rawDescGZIP() []byte {
	file_moneybags_proto_rawDescOnce.Do(func() {
		file_moneybags_proto_rawDescData = protoimpl.X.CompressGZIP(file_moneybags_proto_rawDescData)
	})
	return file_moneybags_proto_rawDescData
}

var file_moneybags_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_moneybags_proto_goTypes = []any{
	(*Config)(nil),         // 0: config.Config
	(*PostgresDB)(nil),     // 1: config.PostgresDB
	(*TelegramConfig)(nil), // 2: config.TelegramConfig
}
var file_moneybags_proto_depIdxs = []int32{
	1, // 0: config.Config.postgres_db:type_name -> config.PostgresDB
	2, // 1: config.Config.telegram_config:type_name -> config.TelegramConfig
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_moneybags_proto_init() }
func file_moneybags_proto_init() {
	if File_moneybags_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_moneybags_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Config); i {
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
		file_moneybags_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*PostgresDB); i {
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
		file_moneybags_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*TelegramConfig); i {
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
			RawDescriptor: file_moneybags_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_moneybags_proto_goTypes,
		DependencyIndexes: file_moneybags_proto_depIdxs,
		MessageInfos:      file_moneybags_proto_msgTypes,
	}.Build()
	File_moneybags_proto = out.File
	file_moneybags_proto_rawDesc = nil
	file_moneybags_proto_goTypes = nil
	file_moneybags_proto_depIdxs = nil
}
