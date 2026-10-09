// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package container

import (
	"github.com/clint456/north-sdk-go/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/di"
)

// NorthServiceName contains the name of the north service SDK implementation in the DIC.
var NorthServiceName = di.TypeInstanceToName((*interfaces.NorthServiceSDK)(nil))

// NorthServiceFrom helper function queries the DIC and returns the north service SDK implementation.
func NorthServiceFrom(get di.Get) interfaces.NorthServiceSDK {
	service, ok := get(NorthServiceName).(interfaces.NorthServiceSDK)
	if !ok {
		return nil
	}
	return service
}
