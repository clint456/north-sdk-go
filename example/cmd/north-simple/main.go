// @version: 0.0.1
// @contact.name: Clint
// @contact.email: clinton_luo@163.com
// This package provides a simple example of a northbound service.
package main

import (
	"github.com/clint456/north-sdk-go/example/driver"
	"github.com/clint456/north-sdk-go/pkg/startup"
)

const (
	serviceName    = "north-simple"
	serviceVersion = "0.0.1"
)

func main() {
	sd := &driver.SimpleDriver{}
	startup.Bootstrap(serviceName, serviceVersion, sd)
}
