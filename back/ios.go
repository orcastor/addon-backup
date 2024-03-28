package back

import (
	"os/exec"
	"strings"

	"github.com/orcastor/phone_images/sdk"
)

/*
# > idevice_id
9dd16339e3fb8357f5954fdeb83383e0e97aabxx (USB)
00008030-00140C56143040xx (USB)
*/
func ListIOSDevices() (devices []map[string]string) {
	output, _ := exec.Command("idevice_id").CombinedOutput()
	if strings.Index(string(output), "ERROR") >= 0 {
		return
	}
	devs := strings.Split(string(output), "\n")
	for _, dev := range devs {
		row := strings.Split(dev, " ")
		if len(row) <= 1 {
			continue
		}
		m := make(map[string]string, 0)
		m["id"] = row[0]
		m["con"] = strings.Trim(row[1], "()")
		// 基础信息
		GetIOSDeviceInfo(m, row[0])
		m["product_name"] = sdk.GetIOSProductName(m["ProductType"])
		m["brand"] = "APPLE"
		// 磁盘空间信息
		GetIOSDeviceInfo(m, row[0], "-q", "com.apple.disk_usage")
		GetIOSDeviceInfo(m, row[0], "-q", "com.apple.mobile.battery")
		GetIOSDeviceInfo(m, row[0], "-q", "com.apple.mobile.backup")
		devices = append(devices, m)
	}
	return
}

/*
# > ideviceinfo -u 9dd16339e3fb8357f5954fdeb83383e0e97aabxx

https://support.apple.com/zh-cn/HT202778
如果未信任：ERROR: Could not connect to lockdownd: Pairing dialog response pending (-19)

DeviceClass: iPhone
DeviceColor: #e4e7e8
DeviceName: iPhone  -- 自定义的名称
SerialNumber: F2LQT4JFGRX* F2LDC0HMN70*
PhoneNumber: +86 188 xxxx 8888
ProductType: iPhone8,2
ProductVersion: 15.7.1
ModelNumber: ML6J2
RegulatoryModelNumber: A1687
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
			fallthrough
		case "ProductType", "ProductVersion", "ModelNumber", "RegionInfo", "RegulatoryModelNumber":
			fallthrough
		case "TotalDiskCapacity", "TotalDataAvailable", "TotalSystemAvailable":
			fallthrough
		case "BatteryIsCharging", "BatteryCurrentCapacity":
			fallthrough
		case "WillEncrypt":
			m[kv[0]] = kv[1]
		}
	}
}
