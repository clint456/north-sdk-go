// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package models

import "github.com/edgexfoundry/go-mod-core-contracts/v4/models"

type UpstreamEvent struct {
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
}
