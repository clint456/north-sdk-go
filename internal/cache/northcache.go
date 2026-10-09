// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
// Package cache provides the in-memory metadata cache used by the north SDK.
package cache

import (
	"fmt"
	"sync"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

// NorthCache holds the northbound channel (North) and profile (NorthProfile)
// definitions managed by the SDK. All accessors are safe for concurrent use.
//
// A North is the northbound equivalent of a southbound Device and a NorthProfile
// is the equivalent of a DeviceProfile. North channels reference a profile via
// North.ProfileName, which is how resource/command lookups are resolved.
type NorthCache struct {
	mu          sync.RWMutex
	northMap    map[string]*models.North
	profileMap  map[string]*models.NorthProfile
	resourceMap map[string]map[string]models.NorthResource // profileName -> resourceName -> resource
	commandMap  map[string]map[string]models.NorthCommand  // profileName -> commandName -> command
}

// NewNorthCache creates an empty NorthCache.
func NewNorthCache() *NorthCache {
	return &NorthCache{
		northMap:    make(map[string]*models.North),
		profileMap:  make(map[string]*models.NorthProfile),
		resourceMap: make(map[string]map[string]models.NorthResource),
		commandMap:  make(map[string]map[string]models.NorthCommand),
	}
}

// ============================ North ============================

// ForName returns the north channel with the given name.
func (c *NorthCache) ForName(name string) (models.North, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	north, ok := c.northMap[name]
	if !ok {
		return models.North{}, false
	}
	return *north, true
}

// All returns all north channels currently in the cache.
func (c *NorthCache) All() []models.North {
	c.mu.RLock()
	defer c.mu.RUnlock()

	norths := make([]models.North, 0, len(c.northMap))
	for _, north := range c.northMap {
		norths = append(norths, *north)
	}
	return norths
}

// Add inserts a new north channel into the cache.
func (c *NorthCache) Add(north models.North) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.add(north)
}

func (c *NorthCache) add(north models.North) error {
	if _, ok := c.northMap[north.Name]; ok {
		return errors.NewCommonEdgeX(errors.KindDuplicateName, fmt.Sprintf("North %s already exists in cache", north.Name), nil)
	}
	c.northMap[north.Name] = &north
	return nil
}

// Update replaces the north channel in the cache, keyed by its name.
func (c *NorthCache) Update(north models.North) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.removeByName(north.Name); err != nil {
		return err
	}
	return c.add(north)
}

// RemoveByName removes the north channel with the given name.
func (c *NorthCache) RemoveByName(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.removeByName(name)
}

func (c *NorthCache) removeByName(name string) error {
	if _, ok := c.northMap[name]; !ok {
		return errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("North %s not found in cache", name), nil)
	}
	delete(c.northMap, name)
	return nil
}

// UpdateOperatingState updates the operating state of the north channel in cache.
func (c *NorthCache) UpdateOperatingState(name string, state models.OperatingState) error {
	if state != models.Up && state != models.Down {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "invalid OperatingState", nil)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	north, ok := c.northMap[name]
	if !ok {
		return errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("North %s not found in cache", name), nil)
	}
	north.OperatingState = state
	return nil
}

// ============================ NorthProfile ============================

// ProfileForName returns the north profile with the given name.
func (c *NorthCache) ProfileForName(name string) (models.NorthProfile, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	profile, ok := c.profileMap[name]
	if !ok {
		return models.NorthProfile{}, false
	}
	return *profile, true
}

// AllProfiles returns all north profiles currently in the cache.
func (c *NorthCache) AllProfiles() []models.NorthProfile {
	c.mu.RLock()
	defer c.mu.RUnlock()

	profiles := make([]models.NorthProfile, 0, len(c.profileMap))
	for _, profile := range c.profileMap {
		profiles = append(profiles, *profile)
	}
	return profiles
}

// AddProfile inserts a new north profile into the cache.
func (c *NorthCache) AddProfile(profile models.NorthProfile) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.addProfile(profile)
}

func (c *NorthCache) addProfile(profile models.NorthProfile) error {
	if _, ok := c.profileMap[profile.Name]; ok {
		return errors.NewCommonEdgeX(errors.KindDuplicateName, fmt.Sprintf("NorthProfile %s already exists in cache", profile.Name), nil)
	}
	c.profileMap[profile.Name] = &profile
	c.resourceMap[profile.Name] = northResourceSliceToMap(profile.NorthResources)
	c.commandMap[profile.Name] = northCommandSliceToMap(profile.NorthCommands)
	return nil
}

// UpdateProfile replaces the north profile in the cache, keyed by its name.
func (c *NorthCache) UpdateProfile(profile models.NorthProfile) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.removeProfileByName(profile.Name); err != nil {
		return err
	}
	return c.addProfile(profile)
}

// RemoveProfileByName removes the north profile with the given name.
func (c *NorthCache) RemoveProfileByName(name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.removeProfileByName(name)
}

func (c *NorthCache) removeProfileByName(name string) error {
	if _, ok := c.profileMap[name]; !ok {
		return errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("NorthProfile %s not found in cache", name), nil)
	}
	delete(c.profileMap, name)
	delete(c.resourceMap, name)
	delete(c.commandMap, name)
	return nil
}

// ============================ Lookups ============================

// NorthResource returns the NorthResource identified by the north channel name and the resource name.
// The resource is resolved through the profile referenced by the north channel.
func (c *NorthCache) NorthResource(northName string, resourceName string) (models.NorthResource, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	profileName, ok := c.profileNameForNorth(northName)
	if !ok {
		return models.NorthResource{}, false
	}
	resources, ok := c.resourceMap[profileName]
	if !ok {
		return models.NorthResource{}, false
	}
	resource, ok := resources[resourceName]
	return resource, ok
}

// NorthCommand returns the NorthCommand identified by the north channel name and the command name.
// The command is resolved through the profile referenced by the north channel.
func (c *NorthCache) NorthCommand(northName string, cmdName string) (models.NorthCommand, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	profileName, ok := c.profileNameForNorth(northName)
	if !ok {
		return models.NorthCommand{}, false
	}
	commands, ok := c.commandMap[profileName]
	if !ok {
		return models.NorthCommand{}, false
	}
	command, ok := commands[cmdName]
	return command, ok
}

// profileNameForNorth returns the profile name referenced by the given north channel.
// The caller must hold at least the read lock.
func (c *NorthCache) profileNameForNorth(northName string) (string, bool) {
	north, ok := c.northMap[northName]
	if !ok || north.ProfileName == "" {
		return "", false
	}
	return north.ProfileName, true
}

func northResourceSliceToMap(resources []models.NorthResource) map[string]models.NorthResource {
	result := make(map[string]models.NorthResource, len(resources))
	for _, resource := range resources {
		result[resource.Name] = resource
	}
	return result
}

func northCommandSliceToMap(commands []models.NorthCommand) map[string]models.NorthCommand {
	result := make(map[string]models.NorthCommand, len(commands))
	for _, command := range commands {
		result[command.Name] = command
	}
	return result
}
