// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package service

import (
	"context"

	sdkModels "github.com/clint456/north-sdk-go/pkg/models"
)

// defaultAsyncEventBufferSize is used when the configured buffer size is not
// positive.
const defaultAsyncEventBufferSize = 100

// processUpstreamEvents consumes the asynchronous upstream channel and dispatches
// every event to the protocol driver for northbound delivery.
func (s *northService) processUpstreamEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.asyncUpCh:
			if event == nil {
				continue
			}
			s.sendUpstreamEvent(event)
		}
	}
}

// sendUpstreamEvent resolves the target north channel and hands the event to the
// protocol driver.
func (s *northService) sendUpstreamEvent(event *sdkModels.UpstreamEvent) {
	northName := event.TargetNorth()
	if northName == "" {
		s.logErrorf("skip upstream event because the target north channel is empty")
		return
	}
	if _, ok := s.northCache.ForName(northName); !ok {
		s.logErrorf("skip upstream event because north channel %s does not exist", northName)
		return
	}
	if err := s.driver.HandleSend(northName, event.ToEvent()); err != nil {
		s.logErrorf("failed to send upstream event to north channel %s: %v", northName, err)
	}
}
