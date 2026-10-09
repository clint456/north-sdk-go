// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package service

import (
	"context"
	"sort"
	"sync"

	sdkCache "github.com/clint456/north-sdk-go/internal/cache"
	sdkModels "github.com/clint456/north-sdk-go/pkg/models"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

// Bootstrap contains references to dependencies required by the BootstrapHandler.
type Bootstrap struct {
	northService *northService
}

// NewBootstrap is a factory method that returns an initialized Bootstrap receiver struct.
func NewBootstrap(s *northService) *Bootstrap {
	return &Bootstrap{northService: s}
}

// BootstrapHandler fulfills the BootstrapHandler contract. It wires up logging,
// the metadata client, the metadata cache, the async upstream channel and the
// protocol driver, then starts every managed north channel.
func (b *Bootstrap) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup, _ startup.Timer, dic *di.Container) bool {
	s := b.northService
	s.ctx = ctx
	s.wg = wg
	s.dic = dic
	s.lc = bootstrapContainer.LoggingClientFrom(dic.Get)

	// Initialise the metadata cache from core-metadata. As with device-sdk-go,
	// metadata is mandatory and the Clients configuration must contain core-metadata.
	northCache, err := sdkCache.InitCache(s.serviceKey, dic)
	if err != nil {
		s.lc.Errorf("failed to init north cache: %s", err.Error())
		return false
	}
	s.northCache = northCache

	// Create the async upstream channel and start its consumer.
	bufferSize := s.svcConfig.North.AsyncEventChanBufferSize
	if bufferSize <= 0 {
		bufferSize = defaultAsyncEventBufferSize
	}
	s.asyncUpCh = make(chan *sdkModels.UpstreamEvent, bufferSize)
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.processUpstreamEvents(ctx)
	}()

	// Initialise the protocol driver.
	if err := s.driver.Initialize(s); err != nil {
		s.lc.Errorf("protocol driver initialize failed: %s", err.Error())
		return false
	}

	// Start the north channels already present in the cache.
	for _, north := range s.northCache.All() {
		if err := s.driver.Start(north.Name, firstProtocolProperties(north)); err != nil {
			s.lc.Errorf("failed to start north channel %s: %s", north.Name, err.Error())
			return false
		}
		s.lc.Infof("north channel %s started", north.Name)
	}

	return true
}

// firstProtocolProperties returns the protocol properties of the (deterministically)
// first protocol configured on a north channel. A north channel normally declares
// a single protocol.
func firstProtocolProperties(north models.North) models.ProtocolProperties {
	if len(north.Protocols) == 0 {
		return models.ProtocolProperties{}
	}
	keys := make([]string, 0, len(north.Protocols))
	for key := range north.Protocols {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return north.Protocols[keys[0]]
}
