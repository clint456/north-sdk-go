// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package service

import (
	"fmt"

	sdkModels "github.com/clint456/north-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

// RegisterCommandHandler registers the downstream command handler for a north
// channel. The handler is invoked by DispatchCommand when a platform command
// arrives for that channel.
func (s *northService) RegisterCommandHandler(northName string, handler func(cmd models.NorthCommand) (sdkModels.CommandResult, error)) error {
	if handler == nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "command handler cannot be nil", nil)
	}
	if _, ok := s.northCache.ForName(northName); !ok {
		return errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("failed to find North %s in cache", northName), nil)
	}

	s.handlerMu.Lock()
	defer s.handlerMu.Unlock()

	if _, ok := s.cmdHandlers[northName]; ok {
		return errors.NewCommonEdgeX(errors.KindDuplicateName, fmt.Sprintf("command handler for north %s already registered", northName), nil)
	}
	s.cmdHandlers[northName] = handler
	s.logDebugf("Registered command handler for north channel %s", northName)
	return nil
}

// UnregisterCommandHandler removes the downstream command handler for a north channel.
func (s *northService) UnregisterCommandHandler(northName string) error {
	s.handlerMu.Lock()
	defer s.handlerMu.Unlock()

	if _, ok := s.cmdHandlers[northName]; !ok {
		return errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("no command handler registered for north %s", northName), nil)
	}
	delete(s.cmdHandlers, northName)
	s.logDebugf("Unregistered command handler for north channel %s", northName)
	return nil
}

// DispatchCommand invokes the command handler registered for the given north
// channel. It is the entry point an inbound transport uses to deliver commands
// to the SDK.
func (s *northService) DispatchCommand(northName string, cmd models.NorthCommand) (sdkModels.CommandResult, error) {
	s.handlerMu.RLock()
	handler, ok := s.cmdHandlers[northName]
	s.handlerMu.RUnlock()

	if !ok {
		return sdkModels.CommandResult{}, errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("no command handler registered for north %s", northName), nil)
	}
	return handler(cmd)
}
