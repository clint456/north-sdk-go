// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package models

// CommandResult 接收平台下行命令的返回结果
type CommandResult struct {
	// Data 业务返回数据，非协议配置
	Data map[string]interface{}
}

// NewCommandResult creates a CommandResult carrying the given business data.
func NewCommandResult(data map[string]interface{}) CommandResult {
	if data == nil {
		data = make(map[string]interface{})
	}
	return CommandResult{Data: data}
}

// Get returns the value stored under key.
func (r CommandResult) Get(key string) (interface{}, bool) {
	value, ok := r.Data[key]
	return value, ok
}

// Set stores value under key, initialising the underlying map when needed.
func (r *CommandResult) Set(key string, value interface{}) {
	if r.Data == nil {
		r.Data = make(map[string]interface{})
	}
	r.Data[key] = value
}
