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

// AddNorthProfile adds a new north profile to the cache and, when configured, to
// the northbound metadata service. It returns the assigned id.
func (s *northService) AddNorthProfile(profile models.NorthProfile) (string, error) {
	if profile.Name == "" {
		return "", errors.NewCommonEdgeX(errors.KindContractInvalid, "north profile name cannot be empty", nil)
	}
	if _, ok := s.northCache.ProfileForName(profile.Name); ok {
		return "", errors.NewCommonEdgeX(errors.KindDuplicateName, fmt.Sprintf("NorthProfile %s already exists", profile.Name), nil)
	}

	client := s.northProfileClient()
	if client == nil {
		return "", errors.NewCommonEdgeX(errors.KindServerError, "north profile metadata client is not available", nil)
	}
	res, err := client.Add(s.correlationContext(), []requests.NorthProfileRequest{requests.NewNorthProfileRequest(dtos.FromNorthProfileModelToDTO(profile))})
	if err != nil {
		return "", err
	}
	profile.Id = res[0].Id

	if err := s.northCache.AddProfile(profile); err != nil {
		return "", err
	}
	s.logDebugf("Added north profile %s", profile.Name)
	return profile.Id, nil
}

// NorthProfiles returns all managed north profiles from the cache.
func (s *northService) NorthProfiles() []models.NorthProfile {
	return s.northCache.AllProfiles()
}

// GetNorthProfileByName returns the north profile with the given name from the cache.
func (s *northService) GetNorthProfileByName(name string) (models.NorthProfile, error) {
	profile, ok := s.northCache.ProfileForName(name)
	if !ok {
		return models.NorthProfile{}, errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("failed to find NorthProfile %s in cache", name), nil)
	}
	return profile, nil
}

// UpdateNorthProfile updates the north profile in the cache and persists it to the
// metadata service when configured.
func (s *northService) UpdateNorthProfile(profile models.NorthProfile) error {
	if profile.Name == "" {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "north profile name cannot be empty", nil)
	}
	if _, err := s.GetNorthProfileByName(profile.Name); err != nil {
		return err
	}

	client := s.northProfileClient()
	if client == nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "north profile metadata client is not available", nil)
	}
	_, err := client.Update(s.correlationContext(), []requests.NorthProfileRequest{requests.NewNorthProfileRequest(dtos.FromNorthProfileModelToDTO(profile))})
	if err != nil {
		return err
	}
	if err := s.northCache.UpdateProfile(profile); err != nil {
		return err
	}
	s.logDebugf("Updated north profile %s", profile.Name)
	return nil
}

// RemoveNorthProfileByName removes the north profile from the cache and the metadata service.
func (s *northService) RemoveNorthProfileByName(name string) error {
	if _, err := s.GetNorthProfileByName(name); err != nil {
		return err
	}

	client := s.northProfileClient()
	if client == nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "north profile metadata client is not available", nil)
	}
	if _, err := client.DeleteByName(s.correlationContext(), name); err != nil {
		return err
	}
	if err := s.northCache.RemoveProfileByName(name); err != nil {
		return err
	}
	s.logDebugf("Removed north profile %s", name)
	return nil
}

// NorthResource returns the NorthResource identified by the north channel name and
// the resource name. The resource is resolved through the channel's profile.
func (s *northService) NorthResource(northName string, resourceName string) (models.NorthResource, bool) {
	return s.northCache.NorthResource(northName, resourceName)
}

// NorthCommand returns the NorthCommand identified by the north channel name and
// the command name. The command is resolved through the channel's profile.
func (s *northService) NorthCommand(northName string, cmdName string) (models.NorthCommand, bool) {
	return s.northCache.NorthCommand(northName, cmdName)
}
