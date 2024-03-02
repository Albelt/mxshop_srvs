// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.0--rc1
// source: config.proto

package config

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

	Server   *Server   `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
	Mysql    *Mysql    `protobuf:"bytes,2,opt,name=mysql,proto3" json:"mysql,omitempty"`
	Redis    *Redis    `protobuf:"bytes,3,opt,name=redis,proto3" json:"redis,omitempty"`
	Consul   *Consul   `protobuf:"bytes,4,opt,name=consul,proto3" json:"consul,omitempty"`
	RocketMq *RocketMq `protobuf:"bytes,5,opt,name=rocket_mq,json=rocketMq,proto3" json:"rocket_mq,omitempty" mapstructure:"rocket_mq"` //@gotags: mapstructure:"rocket_mq"
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[0]
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
	return file_config_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetServer() *Server {
	if x != nil {
		return x.Server
	}
	return nil
}

func (x *Config) GetMysql() *Mysql {
	if x != nil {
		return x.Mysql
	}
	return nil
}

func (x *Config) GetRedis() *Redis {
	if x != nil {
		return x.Redis
	}
	return nil
}

func (x *Config) GetConsul() *Consul {
	if x != nil {
		return x.Consul
	}
	return nil
}

func (x *Config) GetRocketMq() *RocketMq {
	if x != nil {
		return x.RocketMq
	}
	return nil
}

type Server struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip   string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Port int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	// 服务名称、标签、出口IP
	Name       string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Tags       []string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	ExternalIp string   `protobuf:"bytes,5,opt,name=external_ip,json=externalIp,proto3" json:"external_ip,omitempty" mapstructure:"external_ip"` //@gotags: mapstructure:"external_ip"
	Debug      bool     `protobuf:"varint,6,opt,name=debug,proto3" json:"debug,omitempty"`
}

func (x *Server) Reset() {
	*x = Server{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Server) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Server) ProtoMessage() {}

func (x *Server) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Server.ProtoReflect.Descriptor instead.
func (*Server) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{1}
}

func (x *Server) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *Server) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *Server) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Server) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Server) GetExternalIp() string {
	if x != nil {
		return x.ExternalIp
	}
	return ""
}

func (x *Server) GetDebug() bool {
	if x != nil {
		return x.Debug
	}
	return false
}

type Mysql struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr          string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	User          string `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	Password      string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	Db            string `protobuf:"bytes,4,opt,name=db,proto3" json:"db,omitempty"`
	SlowThreshold int32  `protobuf:"varint,5,opt,name=slow_threshold,json=slowThreshold,proto3" json:"slow_threshold,omitempty"`
}

func (x *Mysql) Reset() {
	*x = Mysql{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Mysql) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Mysql) ProtoMessage() {}

func (x *Mysql) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Mysql.ProtoReflect.Descriptor instead.
func (*Mysql) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{2}
}

func (x *Mysql) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Mysql) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *Mysql) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Mysql) GetDb() string {
	if x != nil {
		return x.Db
	}
	return ""
}

func (x *Mysql) GetSlowThreshold() int32 {
	if x != nil {
		return x.SlowThreshold
	}
	return 0
}

type Redis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr     string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Db       int64  `protobuf:"varint,3,opt,name=db,proto3" json:"db,omitempty"`
}

func (x *Redis) Reset() {
	*x = Redis{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Redis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Redis) ProtoMessage() {}

func (x *Redis) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Redis.ProtoReflect.Descriptor instead.
func (*Redis) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{3}
}

func (x *Redis) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Redis) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Redis) GetDb() int64 {
	if x != nil {
		return x.Db
	}
	return 0
}

type Consul struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr        string              `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	HealthCheck *Consul_HealthCheck `protobuf:"bytes,2,opt,name=health_check,json=healthCheck,proto3" json:"health_check,omitempty" mapstructure:"health_check"` //@gotags: mapstructure:"health_check"
}

func (x *Consul) Reset() {
	*x = Consul{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Consul) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Consul) ProtoMessage() {}

func (x *Consul) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Consul.ProtoReflect.Descriptor instead.
func (*Consul) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{4}
}

func (x *Consul) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Consul) GetHealthCheck() *Consul_HealthCheck {
	if x != nil {
		return x.HealthCheck
	}
	return nil
}

type RocketMq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NameServer        string `protobuf:"bytes,1,opt,name=name_server,json=nameServer,proto3" json:"name_server,omitempty" mapstructure:"name_server"` //@gotags: mapstructure:"name_server"
	Topic             string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	ConsumerGroupName string `protobuf:"bytes,3,opt,name=consumer_group_name,json=consumerGroupName,proto3" json:"consumer_group_name,omitempty" mapstructure:"consumer_group_name"` //@gotags: mapstructure:"consumer_group_name"
}

func (x *RocketMq) Reset() {
	*x = RocketMq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RocketMq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RocketMq) ProtoMessage() {}

func (x *RocketMq) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RocketMq.ProtoReflect.Descriptor instead.
func (*RocketMq) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{5}
}

func (x *RocketMq) GetNameServer() string {
	if x != nil {
		return x.NameServer
	}
	return ""
}

func (x *RocketMq) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *RocketMq) GetConsumerGroupName() string {
	if x != nil {
		return x.ConsumerGroupName
	}
	return ""
}

type Consul_HealthCheck struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Interval   string `protobuf:"bytes,1,opt,name=interval,proto3" json:"interval,omitempty"`
	Timeout    string `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Deregister string `protobuf:"bytes,3,opt,name=deregister,proto3" json:"deregister,omitempty"`
}

func (x *Consul_HealthCheck) Reset() {
	*x = Consul_HealthCheck{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Consul_HealthCheck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Consul_HealthCheck) ProtoMessage() {}

func (x *Consul_HealthCheck) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Consul_HealthCheck.ProtoReflect.Descriptor instead.
func (*Consul_HealthCheck) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{4, 0}
}

func (x *Consul_HealthCheck) GetInterval() string {
	if x != nil {
		return x.Interval
	}
	return ""
}

func (x *Consul_HealthCheck) GetTimeout() string {
	if x != nil {
		return x.Timeout
	}
	return ""
}

func (x *Consul_HealthCheck) GetDeregister() string {
	if x != nil {
		return x.Deregister
	}
	return ""
}

var File_config_proto protoreflect.FileDescriptor

var file_config_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xae,
	0x01, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1f, 0x0a, 0x06, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x05, 0x6d, 0x79,
	0x73, 0x71, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x4d, 0x79, 0x73, 0x71,
	0x6c, 0x52, 0x05, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x12, 0x1c, 0x0a, 0x05, 0x72, 0x65, 0x64, 0x69,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x73, 0x52,
	0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x12, 0x1f, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6c, 0x52,
	0x06, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6c, 0x12, 0x26, 0x0a, 0x09, 0x72, 0x6f, 0x63, 0x6b, 0x65,
	0x74, 0x5f, 0x6d, 0x71, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x52, 0x6f, 0x63,
	0x6b, 0x65, 0x74, 0x4d, 0x71, 0x52, 0x08, 0x72, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x4d, 0x71, 0x22,
	0x8b, 0x01, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x5f, 0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x78, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x22, 0x82, 0x01,
	0x0a, 0x05, 0x4d, 0x79, 0x73, 0x71, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x64,
	0x62, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x64, 0x62, 0x12, 0x25, 0x0a, 0x0e, 0x73,
	0x6c, 0x6f, 0x77, 0x5f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0d, 0x73, 0x6c, 0x6f, 0x77, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f,
	0x6c, 0x64, 0x22, 0x47, 0x0a, 0x05, 0x52, 0x65, 0x64, 0x69, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x61,
	0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x64,
	0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x64, 0x62, 0x22, 0xb9, 0x01, 0x0a, 0x06,
	0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x36, 0x0a, 0x0c, 0x68, 0x65,
	0x61, 0x6c, 0x74, 0x68, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6c, 0x2e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x0b, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x1a, 0x63, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x18, 0x0a,
	0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x65, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x22, 0x71, 0x0a, 0x08, 0x52, 0x6f, 0x63, 0x6b, 0x65,
	0x74, 0x4d, 0x71, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x2e, 0x0a, 0x13, 0x63, 0x6f,
	0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65,
	0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_proto_rawDescOnce sync.Once
	file_config_proto_rawDescData = file_config_proto_rawDesc
)

func file_config_proto_rawDescGZIP() []byte {
	file_config_proto_rawDescOnce.Do(func() {
		file_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_proto_rawDescData)
	})
	return file_config_proto_rawDescData
}

var file_config_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_config_proto_goTypes = []interface{}{
	(*Config)(nil),             // 0: Config
	(*Server)(nil),             // 1: Server
	(*Mysql)(nil),              // 2: Mysql
	(*Redis)(nil),              // 3: Redis
	(*Consul)(nil),             // 4: Consul
	(*RocketMq)(nil),           // 5: RocketMq
	(*Consul_HealthCheck)(nil), // 6: Consul.HealthCheck
}
var file_config_proto_depIdxs = []int32{
	1, // 0: Config.server:type_name -> Server
	2, // 1: Config.mysql:type_name -> Mysql
	3, // 2: Config.redis:type_name -> Redis
	4, // 3: Config.consul:type_name -> Consul
	5, // 4: Config.rocket_mq:type_name -> RocketMq
	6, // 5: Consul.health_check:type_name -> Consul.HealthCheck
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_config_proto_init() }
func file_config_proto_init() {
	if File_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Server); i {
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
		file_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Mysql); i {
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
		file_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Redis); i {
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
		file_config_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Consul); i {
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
		file_config_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RocketMq); i {
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
		file_config_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Consul_HealthCheck); i {
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
			RawDescriptor: file_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_proto_goTypes,
		DependencyIndexes: file_config_proto_depIdxs,
		MessageInfos:      file_config_proto_msgTypes,
	}.Build()
	File_config_proto = out.File
	file_config_proto_rawDesc = nil
	file_config_proto_goTypes = nil
	file_config_proto_depIdxs = nil
}
