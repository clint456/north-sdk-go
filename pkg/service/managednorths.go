// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package service

import (
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos/requests"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

// AddNorth adds a new north channel to the cache and, when configured, to the
// northbound metadata service. It returns the assigned id.
func (s *northService) AddNorth(north models.North) (string, error) {
	if north.Name == "" {
		return "", errors.NewCommonEdgeX(errors.KindContractInvalid, "north name cannot be empty", nil)
	}
	if _, ok := s.northCache.ForName(north.Name); ok {
		return "", errors.NewCommonEdgeX(errors.KindDuplicateName, fmt.Sprintf("North %s already exists", north.Name), nil)
	}

	north.ServiceName = s.serviceKey

	client := s.northClient()
	if client == nil {
		return "", errors.NewCommonEdgeX(errors.KindServerError, "north metadata client is not available", nil)
	}
	res, err := client.Add(s.correlationContext(), []requests.NorthRequest{requests.NewNorthRequest(dtos.FromNorthModelToDTO(north))})
	if err != nil {
		return "", err
	}
	north.Id = res[0].Id

	if err := s.northCache.Add(north); err != nil {
		return "", err
	}
	s.logDebugf("Added north channel %s", north.Name)
	return north.Id, nil
}

// Norths returns all managed north channels from the cache.
func (s *northService) Norths() []models.North {
	return s.northCache.All()
}

// GetNorthByName returns the north channel with the given name from the cache.
func (s *northService) GetNorthByName(name string) (models.North, error) {
	north, ok := s.northCache.ForName(name)
	if !ok {
		return models.North{}, errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("failed to find North %s in cache", name), nil)
	}
	return north, nil
}

// UpdateNorth updates the north channel in the cache and persists it to the
// metadata service when configured.
func (s *northService) UpdateNorth(north models.North) error {
	if north.Name == "" {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "north name cannot be empty", nil)
	}
	if _, err := s.GetNorthByName(north.Name); err != nil {
		return err
	}
	if north.ServiceName == "" {
		north.ServiceName = s.serviceKey
	}

	client := s.northClient()
	if client == nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "north metadata client is not available", nil)
	}
	_, err := client.Update(s.correlationContext(), []requests.NorthRequest{requests.NewNorthRequest(dtos.FromNorthModelToDTO(north))})
	if err != nil {
		return err
	}
	if err := s.northCache.Update(north); err != nil {
		return err
	}
	s.logDebugf("Updated north channel %s", north.Name)
	return nil
}

// RemoveNorthByName removes the north channel from the cache and the metadata service.
func (s *northService) RemoveNorthByName(name string) error {
	if _, err := s.GetNorthByName(name); err != nil {
		return err
	}

	client := s.northClient()
	if client == nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "north metadata client is not available", nil)
	}
	if _, err := client.DeleteByName(s.correlationContext(), name); err != nil {
		return err
	}
	if err := s.northCache.RemoveByName(name); err != nil {
		return err
	}
	s.logDebugf("Removed north channel %s", name)
	return nil
}

// UpdateNorthOperatingState updates the operating state of the north channel.
func (s *northService) UpdateNorthOperatingState(name string, state models.OperatingState) error {
	north, err := s.GetNorthByName(name)
	if err != nil {
		return err
	}

	north.OperatingState = state
	client := s.northClient()
	if client == nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "north metadata client is not available", nil)
	}
	if _, err := client.Update(s.correlationContext(), []requests.NorthRequest{requests.NewNorthRequest(dtos.FromNorthModelToDTO(north))}); err != nil {
		return err
	}
	return s.northCache.UpdateOperatingState(name, state)
}
