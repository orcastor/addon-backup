package main

import (
	"os/exec"
	"regexp"
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
	res := strings.Split(string(output), "List of devices attached")
	if len(res) <= 1 {
		return
	}
	for _, dev := range strings.Split(res[1], "\n") {
		row := strings.Split(strings.Trim(dev, "\r\n"), "       ")
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
		GetAndroidDeviceDiskSpace(m)
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

func GetAndroidDeviceDiskSpace(m map[string]string) {
	output, _ := exec.Command("adb", "shell", "df", "/storage/emulated").CombinedOutput()
	dfs := strings.Split(string(output), "\n")
	if len(dfs) <= 1 {
		return
	}
	dfs[1] = regexp.MustCompile(`\s+`).ReplaceAllString(dfs[1], " ")
	row := strings.Split(strings.Trim(dfs[1], "\r\n "), " ")
	if len(row) <= 1 {
		return
	}

	m["total"] = row[1]
	m["used"] = row[2]
	m["available"] = row[3]
	m["use%"] = row[4]
	return
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
		m["name"] = GetIOSProductName(m["ProductType"])
		devices = append(devices, m)
	}
	return devices
}

/*
# > ideviceinfo -u 9dd16339e3fb8357f5954fdeb83383e0e97aabxx
DeviceClass: iPhone
DeviceColor: #e4e7e8
DeviceName: iPhone
SerialNumber: F2LQT4JFGRX* F2LDC0HMN70*
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

// https://github.com/lmirosevic/GBDeviceInfo/blob/master/GBDeviceInfo/GBDeviceInfo_iOS.m
func GetIOSProductName(ProductType string) string {
	if a, ok := strings.CutPrefix(ProductType, "iPhone"); ok {
		if n, okk := iPhoneNames[a]; okk {
			return "iPhone " + n
		}
	}
	if a, ok := strings.CutPrefix(ProductType, "iPad"); ok {
		if n, okk := iPadNames[a]; okk {
			return "iPad " + n
		}
	}
	if a, ok := strings.CutPrefix(ProductType, "iPod"); ok {
		if n, okk := iPodNames[a]; okk {
			return "iPod " + n
		}
	}
	return ProductType
}

var iPhoneNames = map[string]string{
	"1,1":  "1",
	"1,2":  "3G",
	"2,1":  "3GS",
	"3,1":  "4",
	"3,2":  "4",
	"3,3":  "4",
	"4,1":  "4S",
	"5,1":  "5",
	"5,2":  "5",
	"5,3":  "5c",
	"5,4":  "5c",
	"6,1":  "5s",
	"6,2":  "5s",
	"7,1":  "6 Plus",
	"7,2":  "6",
	"8,1":  "6s",
	"8,2":  "6s Plus",
	"8,4":  "SE",
	"9,1":  "7",
	"9,3":  "7",
	"9,2":  "7 Plus",
	"9,4":  "7 Plus",
	"10,1": "8",
	"10,4": "8",
	"10,2": "8 Plus",
	"10,5": "8 Plus",
	"10,3": "X",
	"10,6": "X",
	"11,8": "XR",
	"11,2": "XS",
	"11,4": "XS Max",
	"11,6": "XS Max",
	"12,1": "11",
	"12,3": "11 Pro",
	"12,5": "11 Pro Max",
	"12,8": "SE 2",
	"13,1": "12 mini",
	"13,2": "12",
	"13,3": "12 Pro",
	"13,4": "12 Pro Max",
	"14,4": "13 mini",
	"14,5": "13",
	"14,2": "13 Pro",
	"14,3": "13 Pro Max",
	"14,6": "SE 3rd Gen",
	"14,7": "14",
	"14,8": "14 Plus",
	"15,2": "14 Pro",
	"15,3": "14 Pro Max",
	"15,4": "15",
	"15,5": "15 Plus",
	"16,1": "15 Pro",
	"16,2": "15 Pro Max",
}

var iPadNames = map[string]string{
	"1,1":   "1",
	"2,1":   "2",
	"2,2":   "2",
	"2,3":   "2",
	"2,4":   "2",
	"2,5":   "mini 1",
	"2,6":   "mini 1",
	"2,7":   "mini 1",
	"3,1":   "3",
	"3,2":   "3",
	"3,3":   "3",
	"3,4":   "4",
	"3,5":   "4",
	"3,6":   "4",
	"4,1":   "Air 1",
	"4,2":   "Air 1",
	"4,3":   "Air 1",
	"4,4":   "mini 2",
	"4,5":   "mini 2",
	"4,6":   "mini 2",
	"4,7":   "mini 3",
	"4,8":   "mini 3",
	"4,9":   "mini 3",
	"5,1":   "mini 4",
	"5,2":   "mini 4",
	"5,3":   "Air 2",
	"5,4":   "Air 2",
	"6,7":   "Pro 12.9-inch",
	"6,8":   "Pro 12.9-inch",
	"6,3":   "Pro 9.7-inch",
	"6,4":   "Pro 9.7-inch",
	"6,11":  "2017",
	"6,12":  "2017",
	"7,1":   "Pro 12.9-inch 2017",
	"7,2":   "Pro 12.9-inch 2017",
	"7,3":   "Pro 10.5-inch 2017",
	"7,4":   "Pro 10.5-inch 2017",
	"7,5":   "2018",
	"7,6":   "2018",
	"7,11":  "2019",
	"7,12":  "2019",
	"8,1":   "Pro (11 inch, WiFi)",
	"8,3":   "Pro (11 inch, WiFi+Cellular)",
	"8,2":   "Pro (11 inch, 1TB, WiFi)",
	"8,4":   "Pro (11 inch, 1TB, WiFi+Cellular)",
	"8,5":   "Pro 3rd Gen (12.9 inch, WiFi)",
	"8,7":   "Pro 3rd Gen (12.9 inch, WiFi+Cellular)",
	"8,6":   "Pro 3rd Gen (12.9 inch, 1TB, WiFi)",
	"8,8":   "Pro 3rd Gen (12.9 inch, 1TB, WiFi+Cellular)",
	"8,9":   "Pro 2nd Gen (11 inch, WiFi)",
	"8,10":  "Pro 2nd Gen (11 inch, WiFi+Cellular)",
	"8,11":  "Pro 4th Gen (12.9 inch, WiFi)",
	"8,12":  "Pro 4th Gen (12.9 inch, WiFi+Cellular)",
	"11,1":  "mini 5",
	"11,2":  "mini 5",
	"11,3":  "Air 3",
	"11,4":  "Air 3",
	"11,6":  "2020",
	"11,7":  "2020",
	"12,1":  "2021",
	"12,2":  "2021",
	"13,1":  "Air 4",
	"13,2":  "Air 4",
	"13,4":  "Pro 3rd Gen (11 inch, WiFi)",
	"13,5":  "Pro 3rd Gen (11 inch, WiFi+Cellular)",
	"13,6":  "Pro 3rd Gen (11 inch, WiFi+Cellular)",
	"13,7":  "Pro 3rd Gen (11 inch, WiFi+Cellular)",
	"13,8":  "Pro 5th Gen (12.9 inch, WiFi)",
	"13,9":  "Pro 5th Gen (12.9 inch, WiFi+Cellular)",
	"13,10": "Pro 5th Gen (12.9 inch, WiFi+Cellular)",
	"13,11": "Pro 5th Gen (12.9 inch, WiFi+Cellular)",
	"13,16": "Air 5th Gen (WiFi)",
	"13,17": "Air 5th Gen (WiFi+Cellular)",
	"14,1":  "mini 6",
	"14,2":  "mini 6",
	"13,18": "2022",
	"13,19": "2022",
	"14,3":  "Pro 4th Gen (11 inch, WiFi)",
	"14,4":  "Pro 4th Gen (11 inch, WiFi+Cellular)",
	"14,5":  "Pro 6th Gen (12.9 inch, WiFi)",
	"14,6":  "Pro 6th Gen (12.9 inch, WiFi+Cellular)",
}

var iPodNames = map[string]string{
	"1,1": "Touch 1",
	"2,1": "Touch 2",
	"3,1": "Touch 3",
	"4,1": "Touch 4",
	"5,1": "Touch 5",
	"7,1": "Touch 6",
	"9,1": "Touch 7",
}
