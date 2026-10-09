package config

import (
	bootstrapConfig "github.com/edgexfoundry/go-mod-bootstrap/v4/config"
)

// NorthServiceConfig 北向服务SDK的整体配置结构
type NorthServiceConfig struct {
	// Writable 可动态更新配置，支持注册中心热更新，无需重启服务
	Writable WritableInfo

	// Clients 依赖的EdgeX各服务客户端集合
	Clients bootstrapConfig.ClientsCollection

	// Registry 注册中心连接配置
	Registry bootstrapConfig.RegistryInfo

	// Service 北向服务基础信息
	Service bootstrapConfig.ServiceInfo

	// North 北向SDK专属配置项
	North NorthInfo

	// Driver 协议驱动全局自定义配置，键值字符串映射
	Driver map[string]string

	// MessageBus 消息总线连接信息
	MessageBus bootstrapConfig.MessageBusInfo

	// MaxEventSize 事件最大字节大小，用于消息总线或CoreData上报
	MaxEventSize int64

	// MaxConcurrentCommands 限制下行命令并发数量，防止北向通道响应缓慢/离线时
	// 耗尽命令消费协程资源；小于等于0时使用内置默认值
	MaxConcurrentCommands int
}

// NorthInfo 北向SDK专属配置信息
type NorthInfo struct {
	// AsyncEventChanBufferSize 异步上行事件通道的缓冲区长度
	AsyncEventChanBufferSize int
}

// UpdateFromRaw converts configuration received from the Configuration Provider to a service-specific
// configuration struct which is then used to overwrite the service's existing configuration struct.
func (c *NorthServiceConfig) UpdateFromRaw(rawConfig interface{}) bool {
	configuration, ok := rawConfig.(*NorthServiceConfig)
	if ok {
		*c = *configuration
	}
	return ok
}

// EmptyWritablePtr returns a pointer to a service-specific empty WritableInfo struct.
func (c *NorthServiceConfig) EmptyWritablePtr() interface{} {
	return &WritableInfo{}
}

// GetWritablePtr returns pointer to the writable section.
func (c *NorthServiceConfig) GetWritablePtr() any {
	return &c.Writable
}

// UpdateWritableFromRaw converts updated configuration received from the Configuration Provider to a
// service-specific WritableInfo struct.
func (c *NorthServiceConfig) UpdateWritableFromRaw(rawWritable interface{}) bool {
	writable, ok := rawWritable.(*WritableInfo)
	if ok {
		c.Writable = *writable
	}
	return ok
}

// GetBootstrap returns the configuration elements required by the bootstrap.
func (c *NorthServiceConfig) GetBootstrap() bootstrapConfig.BootstrapConfiguration {
	return bootstrapConfig.BootstrapConfiguration{
		Clients:    &c.Clients,
		Service:    &c.Service,
		Registry:   &c.Registry,
		MessageBus: &c.MessageBus,
	}
}

// GetLogLevel returns the current configuration's log level.
func (c *NorthServiceConfig) GetLogLevel() string {
	return c.Writable.LogLevel
}

// GetRegistryInfo gets the RegistryInfo field from the configuration.
func (c *NorthServiceConfig) GetRegistryInfo() bootstrapConfig.RegistryInfo {
	return c.Registry
}

// GetInsecureSecrets returns the service's InsecureSecrets.
func (c *NorthServiceConfig) GetInsecureSecrets() bootstrapConfig.InsecureSecrets {
	return c.Writable.InsecureSecrets
}

// GetTelemetryInfo returns the service's Telemetry settings.
func (c *NorthServiceConfig) GetTelemetryInfo() *bootstrapConfig.TelemetryInfo {
	return &c.Writable.Telemetry
}
