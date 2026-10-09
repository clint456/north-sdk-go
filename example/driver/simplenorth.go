// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
// Package driver provides a minimal example ProtocolDriver for the north SDK.
package driver

import (
	"fmt"
	"sync"

	"github.com/clint456/north-sdk-go/pkg/interfaces"
	sdkModels "github.com/clint456/north-sdk-go/pkg/models"
	edgexErr "github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
)

const (
	demoProfileName = "Simple-North-Profile"
	demoNorthName   = "north-mqtt"
)

// SimpleDriver is a trivial northbound ProtocolDriver used for demonstration.
// It does not talk to a real broker; it just logs what it would send.
type SimpleDriver struct {
	mu     sync.Mutex
	sdk    interfaces.NorthServiceSDK
	norths map[string]models.ProtocolProperties
}

// Initialize seeds a demo north profile and channel into core-metadata (via the
// SDK, which persists them), then remembers the SDK reference. Existing
// entities are tolerated so the example can be re-run.
func (d *SimpleDriver) Initialize(sdk interfaces.NorthServiceSDK) error {
	d.sdk = sdk
	d.norths = make(map[string]models.ProtocolProperties)

	if _, err := sdk.AddNorthProfile(demoProfile()); err != nil && edgexErr.Kind(err) != edgexErr.KindDuplicateName {
		return err
	}
	if _, err := sdk.AddNorth(demoNorth()); err != nil && edgexErr.Kind(err) != edgexErr.KindDuplicateName {
		return err
	}
	return nil
}

// Start is called by the SDK for every north channel present in the cache.
func (d *SimpleDriver) Start(northName string, cfg models.ProtocolProperties) error {
	d.mu.Lock()
	d.norths[northName] = cfg
	d.mu.Unlock()

	fmt.Printf("[simple-driver] north channel %s started with config %v\n", northName, cfg)

	// Register a downstream command handler for this channel.
	return d.sdk.RegisterCommandHandler(northName, d.handleCommand)
}

// Stop unregisters the command handler and forgets the channel.
func (d *SimpleDriver) Stop(northName string, force bool) error {
	_ = d.sdk.UnregisterCommandHandler(northName)

	d.mu.Lock()
	delete(d.norths, northName)
	d.mu.Unlock()

	fmt.Printf("[simple-driver] north channel %s stopped (force=%v)\n", northName, force)
	return nil
}

// UpdateNorth applies a hot configuration update.
func (d *SimpleDriver) UpdateNorth(northName string, cfg models.ProtocolProperties) error {
	d.mu.Lock()
	d.norths[northName] = cfg
	d.mu.Unlock()
	return nil
}

// RemoveNorth disposes any resources held for the north channel.
func (d *SimpleDriver) RemoveNorth(northName string) error {
	d.mu.Lock()
	delete(d.norths, northName)
	d.mu.Unlock()
	return nil
}

// ValidateNorth validates the protocol configuration.
func (d *SimpleDriver) ValidateNorth(cfg models.ProtocolProperties) error {
	if _, ok := cfg["broker"]; !ok {
		return fmt.Errorf("missing required 'broker' property")
	}
	return nil
}

// HandleSend is where the real protocol send would happen.
func (d *SimpleDriver) HandleSend(northName string, event models.Event) error {
	fmt.Printf("[simple-driver] send event %s (%d readings) to north channel %s\n",
		event.Id, len(event.Readings), northName)
	return nil
}

func (d *SimpleDriver) handleCommand(cmd models.NorthCommand) (sdkModels.CommandResult, error) {
	fmt.Printf("[simple-driver] received downstream command %s\n", cmd.Name)
	return sdkModels.NewCommandResult(map[string]interface{}{
		"command": cmd.Name,
		"status":  "ok",
	}), nil
}

func demoProfile() models.NorthProfile {
	return models.NorthProfile{
		Name:        demoProfileName,
		Description: "Demo north profile",
		NorthResources: []models.NorthResource{
			{
				Name:        "temperature",
				Description: "temperature reading",
			},
		},
		NorthCommands: []models.NorthCommand{
			{
				Name:      "setPoint",
				ReadWrite: "RW",
			},
		},
	}
}

func demoNorth() models.North {
	return models.North{
		Name:           demoNorthName,
		Description:    "Demo north channel",
		ProfileName:    demoProfileName,
		OperatingState: models.Up,
		Protocols: map[string]models.ProtocolProperties{
			"mqtt": {
				"broker": "tcp://localhost:1883",
				"topic":  "edge/data",
			},
		},
	}
}
