// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package service

import (
	"context"
	"errors"
	"sync"

	"github.com/bytedance/gopkg/util/logger"
	sdkCommon "github.com/clint456/north-sdk-go/internal/common"
	"github.com/clint456/north-sdk-go/internal/config"
	"github.com/clint456/north-sdk-go/pkg/interfaces"
	sdkModels "github.com/clint456/north-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
)

type northService struct {
	serviceKey  string
	baseService string
	lc          logger.LoggingClient
	// 若单进程单协议：
	driver interfaces.ProtocolDriver
	// 若单进程支持多协议插件，改用下面map
	// drivers map[string]interfaces.ProtocolDriver

	// 上行数据通道
	AsyncuUpCh chan *sdkModels.UpstreamEvent

	// 内存元数据缓存
	northCache *cache.NorthCache

	// 下行命令回调存储 + 并发锁
	handlerMu   sync.RWMutex
	cmdHandlers map[string]func(cmd models.NorthCommand) (models.CommandResult, error)

	// Metadata客户端，持久化元数据
	metaClient metadata.Client

	// 生命周期上下文，优雅关闭
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// SDK全局配置
	svcConfig *config.NorthServiceConfig

	// 通道增删改操作锁
	northMu sync.RWMutex
}

// NewNorthService 创建北向服务实例
func NewNorthService(serviceKey string, serviceVersion string, driver interfaces.ProtocolDriver) (interfaces.NorthServiceSDK, error) {
	var service northService
	if serviceKey == "" {
		return nil, errors.New("please specify north service name")
	}
	service.serviceKey = serviceKey

	if serviceVersion == "" {
		return nil, errors.New("please specify device service version")
	}
	sdkCommon.ServiceVersion = serviceVersion

	service.driver = driver

	service.svcConfig = &config.NorthServiceConfig
	return interfaces.NorthServiceSDK(&service)
}
