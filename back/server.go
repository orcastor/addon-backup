package main

import (
	"github.com/gotmc/libusb/v2"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"

	"github.com/orcastor/orcas/core"
	"github.com/orcastor/orcas/rpc/middleware"
)

// EGO_DEBUG=true EGO_LOG_EXTRA_KEYS=uid ORCAS_BASE=/tmp/test ORCAS_DATA=/tmp/test ORCAS_SECRET=xxxxxxxx egoctl run --runargs --config=config.toml
// go run server.go --config=config.toml
func main() {
	core.InitDB()

	ctx, err := libusb.NewContext()
	if err != nil {
		elog.Panic("Couldn't create USB context. Ending now.")
	}
	defer ctx.Close()

	ctx.HotplugRegisterCallbackEvent(0, 0, libusb.HotplugArrived|libusb.HotplugLeft, USBEvents)
	defer ctx.HotplugDeregisterAllCallbacks()

	if err := ego.New().Serve(func() *egin.Component {
		server := egin.Load("server.http").Build()

		server.Use(middleware.Metrics())
		server.Use(middleware.CORS())
		server.Use(middleware.JWT())

		_ = server.Group("/bak/api")
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.Any("err", err))
	}
}

func USBEvents(vID, pID uint16, eventType libusb.HotPlugEventType) {
	// elog.Infof("VendorID: %04x, ProductID: %04x, EventType: %d\r\n", vID, pID, eventType)
}
