// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package service

import (
	"context"
	"fmt"
	"os"
	"sync"

	sdkCache "github.com/clint456/north-sdk-go/internal/cache"
	sdkCommon "github.com/clint456/north-sdk-go/internal/common"
	"github.com/clint456/north-sdk-go/internal/config"
	"github.com/clint456/north-sdk-go/internal/container"
	"github.com/clint456/north-sdk-go/pkg/interfaces"
	sdkModels "github.com/clint456/north-sdk-go/pkg/models"
	bootstrap "github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/flags"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/handlers"
	bootstrapInterfaces "github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/startup"
	bootstrapTypes "github.com/edgexfoundry/go-mod-bootstrap/v4/config"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/di"
	clientsInterfaces "github.com/edgexfoundry/go-mod-core-contracts/v4/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
	"github.com/google/uuid"
)

type northService struct {
	serviceKey  string
	baseService string
	lc          logger.LoggingClient

	// driver 协议驱动，负责实际的北向协议收发
	driver interfaces.ProtocolDriver

	// asyncUpCh 通道实例上行数据通道：外部（Mapdata）写入，SDK 消费后交由协议驱动发送
	asyncUpCh chan *sdkModels.UpstreamEvent

	// northCache 内存元数据缓存（北向通道 + 点位）
	northCache *sdkCache.NorthCache

	// handlerMu 保护 cmdHandlers 的并发访问
	handlerMu   sync.RWMutex
	cmdHandlers map[string]func(cmd models.NorthCommand) (sdkModels.CommandResult, error)

	// 生命周期上下文，用于优雅关闭
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	// 启动参数与依赖注入容器
	flags *flags.Default
	dic   *di.Container

	// SDK全局配置
	svcConfig *config.NorthServiceConfig
}

// NewNorthService 创建北向服务实例
func NewNorthService(serviceKey string, serviceVersion string, driver interfaces.ProtocolDriver) (interfaces.NorthServiceSDK, error) {
	var service northService
	if serviceKey == "" {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "please specify north service name", nil)
	}
	service.serviceKey = serviceKey
	service.baseService = serviceKey

	if serviceVersion == "" {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "please specify north service version", nil)
	}
	if driver == nil {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "please specify a protocol driver", nil)
	}
	sdkCommon.ServiceVersion = serviceVersion

	service.driver = driver
	service.cmdHandlers = make(map[string]func(cmd models.NorthCommand) (sdkModels.CommandResult, error))
	service.svcConfig = &config.NorthServiceConfig{}
	service.northCache = sdkCache.NewNorthCache()

	return interfaces.NorthServiceSDK(&service), nil
}

// Run starts this North Service. This should not be called directly by a north
// service. Instead, call startup.Bootstrap(...).
func (s *northService) Run() error {
	startupTimer := startup.NewStartUpTimer(s.serviceKey)

	s.flags = flags.New()
	s.flags.Parse(os.Args[1:])

	s.dic = di.NewContainer(di.ServiceConstructorMap{
		container.ConfigurationName: func(get di.Get) any {
			return s.svcConfig
		},
		container.ProtocolDriverName: func(get di.Get) any {
			return s.driver
		},
		container.NorthServiceName: func(get di.Get) any {
			return interfaces.NorthServiceSDK(s)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	wg, deferred, successful := bootstrap.RunAndReturnWaitGroup(
		ctx,
		cancel,
		s.flags,
		s.serviceKey,
		common.ConfigStemDevice,
		s.svcConfig,
		nil,
		startupTimer,
		s.dic,
		true,
		bootstrapTypes.ServiceTypeOther,
		[]bootstrapInterfaces.BootstrapHandler{
			handlers.NewClientsBootstrap().BootstrapHandler,
			handlers.NewStartMessage(s.serviceKey, sdkCommon.ServiceVersion).BootstrapHandler,
			NewBootstrap(s).BootstrapHandler,
		},
	)
	if !successful {
		deferred()
		cancel()
		return errors.NewCommonEdgeX(errors.KindServerError, "bootstrapping failed", nil)
	}

	s.ctx = ctx
	s.cancel = cancel
	s.wg = wg

	defer func() {
		deferred()
		s.Stop(false)
	}()

	wg.Wait()
	return nil
}

// Stop shuts down all managed north channels through the protocol driver.
func (s *northService) Stop(force bool) {
	if s.driver == nil || s.northCache == nil {
		return
	}
	for _, north := range s.northCache.All() {
		if err := s.driver.Stop(north.Name, force); err != nil {
			s.logErrorf("failed to stop north channel %s: %v", north.Name, err)
		}
	}
}

// Name returns the name of this North Service.
func (s *northService) Name() string {
	return s.serviceKey
}

// Version returns the version number of this North Service.
func (s *northService) Version() string {
	return sdkCommon.ServiceVersion
}

// LoggingClient returns the logger.LoggingClient.
func (s *northService) LoggingClient() logger.LoggingClient {
	if s.lc == nil && s.dic != nil {
		s.lc = bootstrapContainer.LoggingClientFrom(s.dic.Get)
	}
	return s.lc
}

// PublishEvent is the single synchronous entry point for northbound data.
// The target north channel is taken from event.DeviceName, mirroring how EdgeX
// reuses the DeviceName field to address a destination.
func (s *northService) PublishEvent(event models.Event) error {
	if s.driver == nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "protocol driver is not initialized", nil)
	}
	northName := event.DeviceName
	if northName == "" {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "event.DeviceName must carry the target north channel name", nil)
	}
	if s.northCache != nil {
		if _, ok := s.northCache.ForName(northName); !ok {
			return errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("failed to find north channel %s", northName), nil)
		}
	}
	if err := s.driver.HandleSend(northName, event); err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, fmt.Sprintf("failed to publish event to north channel %s", northName), err)
	}
	return nil
}

// AsyncEventsChannel returns the asynchronous upstream event channel. Producers
// (e.g. the data mapping layer) write UpstreamEvent values to it and the SDK
// dispatches each one to the protocol driver for northbound delivery.
func (s *northService) AsyncEventsChannel() chan *sdkModels.UpstreamEvent {
	return s.asyncUpCh
}

// northClient returns the core-metadata North client, or nil when not configured.
func (s *northService) northClient() clientsInterfaces.NorthClient {
	if s.dic == nil {
		return nil
	}
	return bootstrapContainer.NorthClientFrom(s.dic.Get)
}

// northProfileClient returns the core-metadata NorthProfile client, or nil when not configured.
func (s *northService) northProfileClient() clientsInterfaces.NorthProfileClient {
	if s.dic == nil {
		return nil
	}
	return bootstrapContainer.NorthProfileClientFrom(s.dic.Get)
}

// correlationContext returns a background context carrying a fresh correlation id
// for use with outbound metadata service calls.
func (s *northService) correlationContext() context.Context {
	return context.WithValue(context.Background(), sdkCommon.CorrelationHeaderKey, uuid.NewString())
}

// logErrorf logs through the logging client when available, otherwise to stderr.
func (s *northService) logErrorf(format string, args ...interface{}) {
	if lc := s.LoggingClient(); lc != nil {
		lc.Errorf(format, args...)
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func (s *northService) logWarnf(format string, args ...interface{}) {
	if lc := s.LoggingClient(); lc != nil {
		lc.Warnf(format, args...)
	}
}

func (s *northService) logInfof(format string, args ...interface{}) {
	if lc := s.LoggingClient(); lc != nil {
		lc.Infof(format, args...)
	}
}

func (s *northService) logDebugf(format string, args ...interface{}) {
	if lc := s.LoggingClient(); lc != nil {
		lc.Debugf(format, args...)
	}
}
