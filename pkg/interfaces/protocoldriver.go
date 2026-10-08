package interfaces

import "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

type ProtocolDriver interface {
	// 初始化
	Initialize(sdk NorthServiceSDK) error
	// HandleReceived 处理接收到的消息
	HandleReceived(protocol map[string]models.ProtocolProperties, message []models.ReceivedMessage) error
	// HandleSend 处理发送的消息
	HandleSend(protocol map[string]models.ProtocolProperties, message []models.SendMessage) error
	//
	Stop(force bool) error

	Start() error
	// 更新
	UpdateNorth(protocol map[string]models.ProtocolProperties) error
	// 移除
	RemoveNorth(protocol map[string]models.ProtocolProperties) error
	// 校验
	ValidateNorth(protocol map[string]models.ProtocolProperties) error
}
