package interfaces

import (
	sdkModels "github.com/clint456/north-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

// NorthServiceSDK defines the minimal interface for Northbound Service SDK
// Core: manage north channel metadata, publish upstream event, handle platform down command
type NorthServiceSDK interface {
	// ========= North通道实例（对应南向Device） =========
	AddNorth(north models.North) (string, error)
	Norths() []models.North
	GetNorthByName(name string) (models.North, error)
	UpdateNorth(north models.North) error
	RemoveNorthByName(name string) error

	// ========= NorthProfile北向点位元数据（对应南向DeviceProfile） =========
	AddNorthProfile(profile models.NorthProfile) (string, error)
	NorthProfiles() []models.NorthProfile
	GetNorthProfileByName(name string) (models.NorthProfile, error)
	UpdateNorthProfile(profile models.NorthProfile) error
	RemoveNorthProfileByName(name string) error

	// ========= 点位、命令查询缓存 =========
	NorthResource(northName string, resourceName string) (models.NorthResource, bool)
	NorthCommand(northName string, cmdName string) (models.NorthCommand, bool)

	// ========= 核心：北向发布数据流（数据出） =========
	// PublishEvent 基础事件上送，所有协议统一入口，复用EdgeX Event模型
	PublishEvent(event models.Event) error
	// AsyncEventsChannel 异步事件通道，Mapdata写入，协议驱动消费
	AsyncEventsChannel() chan *sdkModels.UpstreamEvent

	// ========= 核心：平台下行命令（指令入） =========
	RegisterCommandHandler(northName string, handler func(cmd models.NorthCommand) (sdkModels.CommandResult, error)) error
	UnregisterCommandHandler(northName string) error
	// DispatchCommand routes an inbound platform command to the handler registered
	// for the given north channel. It is the entry point used by an inbound
	// transport (message bus / REST) to deliver commands to the SDK.
	DispatchCommand(northName string, cmd models.NorthCommand) (sdkModels.CommandResult, error)

	// ========= 通道状态运维 =========
	UpdateNorthOperatingState(name string, state models.OperatingState) error

	// ========= 生命周期与基础信息 =========
	// Run starts the north service. It should not be called directly; call
	// startup.Bootstrap instead.
	Run() error
	// Name returns the name of this North Service
	Name() string
	// Version returns the version number of this North Service
	Version() string
	// LoggingClient returns the logger.LoggingClient
	LoggingClient() logger.LoggingClient
}
