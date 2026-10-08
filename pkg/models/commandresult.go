// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package models

// 接收平台下行命令的返回结果
type CommandResult struct {
	data map[string]interface{} // 业务数据，非协议配置
}
