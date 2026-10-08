// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package service

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/clint456/north-sdk-go/pkg/interfaces"
	sdkModels "github.com/clint456/north-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
)

type northService struct {
	serviceKey  string
	baseService string
	lc          logger.LoggingClient
	driver      interfaces.ProtocolDriver
	upCh        chan *sdkModels.UpstreamEvent
}

// NewNorthService 创建北向服务实例
func NewNorthService(serviceKey string, serviceVersion string, driver interfaces.ProtocolDriver) (interfaces.NorthServiceSDK, error) {
	// 初始化日志客户端
	lc := logger.NewClient(serviceKey, serviceVersion, false, "")
	// 创建北向服务实例
	ns := &northService{
		serviceKey:  serviceKey,
		baseService: serviceVersion,
		lc:          lc,
		driver:      driver,
		upCh:        make(chan *sdkModels.UpstreamEvent, 100), // 初始化异步事件通道，缓冲区大小为100
	}
	return interfaces.NorthServiceSDK(ns), nil
}

func 