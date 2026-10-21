package interfaces

import (
	sdkModels "github.com/clint456/north-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

// NorthServiceSDK defines the minimal interface for Northbound Service SDK
// Core: manage north channel metadata, publish upstream event, handle platform down command
type NorthServiceSDK interface {
	// ========= Node通道实例（对应南向Node） =========
	AddNode(node models.Node) (string, error)
	Nodes() []models.Node
	GetNodeByName(name string) (models.Node, error)
	UpdateNode(node models.Node) error
	RemoveNodeByName(name string) error

	// ========= Profile点位元数据（对应南向Profile） =========
	AddProfile(profile models.Profile) (string, error)
	Profiles() []models.Profile
	GetProfileByName(name string) (models.Profile, error)
	UpdateProfile(profile models.Profile) error
	RemoveProfileByName(name string) error

	// ========= 点位、命令查询缓存 =========
	Resource(nodeName string, resourceName string) (models.Resource, bool)
	Command(nodeName string, cmdName string) (models.Command, bool)

	// ========= 核心：北向发布数据流（数据出） =========
	// PublishEvent 基础事件上送，所有协议统一入口，复用EdgeX Event模型
	PublishEvent(event models.Event) error
	// AsyncEventsChannel 异步事件通道，Mapdata写入，协议驱动消费
	AsyncEventsChannel() chan *sdkModels.UpstreamEvent

	// ========= 核心：平台下行命令（指令入） =========
	RegisterCommandHandler(nodeName string, handler func(cmd models.Command) (sdkModels.CommandResult, error)) error
	UnregisterCommandHandler(nodeName string) error

	// ========= 通道状态运维 =========
	UpdateNodeOperatingState(name string, state models.OperatingState) error
}
