// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package container

import (
	"github.com/clint456/north-sdk-go/internal/config"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/di"
)

// ConfigurationName contains the name of the north service's ConfigurationStruct implementation in the DIC.
var ConfigurationName = di.TypeInstanceToName(config.NorthServiceConfig{})

// ConfigurationFrom helper function queries the DIC and returns north service's ConfigurationStruct implementation.
func ConfigurationFrom(get di.Get) *config.NorthServiceConfig {
	cfg, ok := get(ConfigurationName).(*config.NorthServiceConfig)
	if !ok {
		return nil
	}
	return cfg
}
