// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
package container

import (
	"github.com/clint456/north-sdk-go/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/v4/di"
)

// ProtocolDriverName contains the name of the protocol driver implementation in the DIC.
var ProtocolDriverName = di.TypeInstanceToName((*interfaces.ProtocolDriver)(nil))

// ProtocolDriverFrom helper function queries the DIC and returns the protocol driver implementation.
func ProtocolDriverFrom(get di.Get) interfaces.ProtocolDriver {
	driver, ok := get(ProtocolDriverName).(interfaces.ProtocolDriver)
	if !ok {
		return nil
	}
	return driver
}
