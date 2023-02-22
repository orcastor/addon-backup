package main

import (
	"os/exec"
	"strings"

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

/*
adb devices -l
* daemon not running; starting now at tcp:5037
* daemon started successfully
List of devices attached

List of devices attached
4UX02211060009xx       unauthorized usb:338952192X transport_id:1

List of devices attached
4UX02211060009xx       device usb:338952192X product:NOP-AN00P model:NOP_AN00 device:HWNOP transport_id:2
*/
func ListAndroidDevices() (devices []map[string]string) {
	output, _ := exec.Command("adb", "devices", "-l").CombinedOutput()
	res := strings.Split(string(output), "List of devices attached\n")
	if len(res) <= 0 {
		return
	}
	for _, dev := range strings.Split(res[1], "\n") {
		row := strings.Split(dev, "       ")
		if len(row) <= 1 {
			continue
		}
		m := make(map[string]string, 0)
		m["id"] = row[0]
		for _, info := range strings.Split(row[1], " ") {
			kv := strings.Split(info, ":")
			switch len(kv) {
			case 1:
				m["type"] = kv[0]
			case 2:
				m[kv[0]] = kv[1]
			}
		}
		m["name"] = GetAndroidDeviceInfo(row[0], "ro.config.marketing_name")
		devices = append(devices, m)
	}
	return devices
}

/*
# > adb -s <id> shell getprop ro.config.marketing_name
PORSCHE DESIGN HUAWEI Mate 40
*/
func GetAndroidDeviceInfo(id, prop string) string {
	output, _ := exec.Command("adb", []string{"-s", id, "shell", "getprop", prop}...).CombinedOutput()
	return strings.TrimSpace(string(output))
}

/*
# > idevice_id
9dd16339e3fb8357f5954fdeb83383e0e97aabxx (USB)
00008030-00140C56143040xx (USB)
*/
func ListIOSDevices() (devices []map[string]string) {
	output, _ := exec.Command("idevice_id").CombinedOutput()
	devs := strings.Split(string(output), "\n")
	for _, dev := range devs {
		row := strings.Split(dev, " ")
		if len(row) <= 1 {
			continue
		}
		m := make(map[string]string, 0)
		m["id"] = row[0]
		m["source"] = row[1]
		// 基础信息
		GetIOSDeviceInfo(m, row[0])
		// 磁盘空间信息
		GetIOSDeviceInfo(m, row[0], "-q", "com.apple.disk_usage")
	}
	return devices
}

/*
# > ideviceinfo -u 9dd16339e3fb8357f5954fdeb83383e0e97aabxx
DeviceClass: iPhone
DeviceColor: #e4e7e8
DeviceName: iPhone
SerialNumber: F2LQT4JFGRX*
              F2LDC0HMN70*
PhoneNumber: +86 188 xxxx 8888
ProductType: iPhone8,2
ProductVersion: 15.7.1
ModelNumber: ML6J2
RegionInfo: CH/A
UniqueDeviceID: 9dd16339e3fb8357f5954fdeb83383e0e97aabxx

# > ideviceinfo -u 9dd16339e3fb8357f5954fdeb83383e0e97aabxx -q com.apple.disk_usage
TotalDiskCapacity: 31708938240
TotalDataAvailable: 4504629248
TotalSystemAvailable: 335544320
*/
func GetIOSDeviceInfo(m map[string]string, args ...string) {
	output, _ := exec.Command("ideviceinfo", append([]string{"-u"}, args...)...).CombinedOutput()
	for _, attr := range strings.Split(string(output), "\n") {
		kv := strings.Split(attr, ": ")
		if len(kv) <= 1 {
			continue
		}
		switch kv[0] {
		case "DeviceClass", "DeviceColor", "DeviceName", "SerialNumber", "PhoneNumber":
		case "ProductType", "ProductVersion", "ModelNumber", "RegionInfo":
		case "TotalDiskCapacity", "TotalDataAvailable", "TotalSystemAvailable":
			m[kv[0]] = kv[1]
		}
	}
}
