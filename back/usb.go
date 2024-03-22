package back

import (
	"github.com/gotmc/libusb/v2"
	"github.com/gotomicro/ego/core/elog"
)

func init() {
	ctx, err := libusb.NewContext()
	if err != nil {
		elog.Panic("Couldn't create USB context. Ending now.")
	}
	defer ctx.Close()

	ctx.HotplugRegisterCallbackEvent(0, 0, libusb.HotplugArrived|libusb.HotplugLeft, USBEvents)
	defer ctx.HotplugDeregisterAllCallbacks()
}

func USBEvents(vID, pID uint16, eventType libusb.HotPlugEventType) {
	elog.Infof("VendorID: %04x, ProductID: %04x, EventType: %d\r\n", vID, pID, eventType)
	// h5 sse
}
