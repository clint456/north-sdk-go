package interfaces

import "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

type ProtocolDriver interface {
	// Initialize 驱动全局一次性初始化
	Initialize(sdk NorthServiceSDK) error
	// Start 启动该通道实例
	Start(northName string, cfg models.ProtocolProperties) error
	// Stop 停止该通道实例
	Stop(northName string, force bool) error
	// UpdateNorth 通道配置热更新
	UpdateNorth(northName string, cfg models.ProtocolProperties) error
	// RemoveNorth 销毁通道资源
	RemoveNorth(northName string) error
	// ValidateNorth 校验配置
	ValidateNorth(cfg models.ProtocolProperties) error
	// HandleSend SDK传入Event，驱动负责协议发送
	HandleSend(northName string, event models.Event) error
}
