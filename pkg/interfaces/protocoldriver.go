package interfaces

import "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

type ProtocolDriver interface {
	// Initialize 驱动全局一次性初始化
	Initialize(sdk NorthServiceSDK) error
	// Start 启动该通道实例
	Start(nodeName string, cfg models.ProtocolProperties) error
	// Stop 停止该通道实例
	Stop(nodeName string, force bool) error
	// UpdateNode 通道配置热更新
	UpdateNode(nodeName string, cfg models.ProtocolProperties) error
	// RemoveNode 销毁通道资源
	RemoveNode(nodeName string) error
	// ValidateNode 校验配置
	ValidateNode(cfg models.ProtocolProperties) error
	// HandleSend SDK传入Event，驱动负责协议发送
	HandleSend(nodeName string, event models.Event) error
}
