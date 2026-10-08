package startup

import (
	"fmt"
	"os"

	"github.com/clint456/north-sdk-go/pkg/interfaces"
)

func Bootstrap(serviceKey string, serviceVersion string, driver interfaces.ProtocolDriver) {
	northService, err := service.NewNorthService(serviceKey, serviceVersion, driver)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(-1)
	}

	err = northService.Run()
	if err != nil {
		northService.LoggingClient().Errorf("North Service %s", err.Error())
		os.Exit(-1)
	}

	os.Exit(0)
}
