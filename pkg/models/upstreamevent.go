// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package models

import "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

type UpstreamEvent struct {
	// NorthName 目标北向通道名称，为空时回退使用 Name
	NorthName string
	// 事件ID，唯一标识一个事件
	Id string
	// 事件名称
	Name string
	// 事件描述
	Description string
	// 事件标签，用于分类和检索
	Labels []string
	// 事件发生时间戳，单位为毫秒
	Timestamp int64
	// 事件关联的北向资源列表
	Resources []models.NorthResource
	// 事件关联的北向命令列表
	Commands []models.NorthCommand
	// Readings 事件携带的测点数据
	Readings []models.Reading
	// Tags 事件附加的业务标签
	Tags map[string]interface{}
}

// TargetNorth returns the north channel this event should be sent to. It falls
// back to Name when NorthName is not set, mirroring how an EdgeX Event reuses
// DeviceName as the northbound destination.
func (e *UpstreamEvent) TargetNorth() string {
	if e.NorthName != "" {
		return e.NorthName
	}
	return e.Name
}

// ToEvent converts the UpstreamEvent into an EdgeX Event so it can be handed to
// a ProtocolDriver for northbound delivery.
func (e *UpstreamEvent) ToEvent() models.Event {
	return models.Event{
		Id:          e.Id,
		DeviceName:  e.TargetNorth(),
		ProfileName: e.Description,
		SourceName:  e.Name,
		Origin:      e.Timestamp,
		Readings:    e.Readings,
		Tags:        e.Tags,
	}
}
