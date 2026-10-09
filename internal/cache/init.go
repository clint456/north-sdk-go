// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package cache

import (
	"context"
	"fmt"

	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/v4/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/di"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/dtos"
)

// InitCache builds the north metadata cache from the EdgeX core-metadata service.
//
// The North and NorthProfile clients must have been registered in the DIC by the
// bootstrap ClientsBootstrap handler, which requires the core-metadata entry in
// the service's Clients configuration. As with device-sdk-go, metadata is
// mandatory: the cache cannot be initialised without it.
func InitCache(serviceName string, dic *di.Container) (*NorthCache, error) {
	northClient := bootstrapContainer.NorthClientFrom(dic.Get)
	profileClient := bootstrapContainer.NorthProfileClientFrom(dic.Get)
	if northClient == nil || profileClient == nil {
		return nil, fmt.Errorf("north metadata clients are not available, check the Clients configuration for core-metadata")
	}

	northCache := NewNorthCache()
	ctx := context.Background()

	northRes, err := northClient.NorthsByServiceName(ctx, serviceName, 0, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to load north channels from core-metadata: %w", err)
	}
	for _, north := range northRes.Norths {
		if addErr := northCache.Add(dtos.ToNorthModel(north)); addErr != nil {
			return nil, addErr
		}
	}

	// Load the profile referenced by each north channel.
	for _, north := range northCache.All() {
		if north.ProfileName == "" {
			continue
		}
		if _, ok := northCache.ProfileForName(north.ProfileName); ok {
			continue
		}
		res, err := profileClient.NorthProfileByName(ctx, north.ProfileName)
		if err != nil {
			return nil, fmt.Errorf("failed to load north profile %s from core-metadata: %w", north.ProfileName, err)
		}
		if addErr := northCache.AddProfile(dtos.ToNorthProfileModel(res.Profile)); addErr != nil {
			return nil, addErr
		}
	}

	return northCache, nil
}
