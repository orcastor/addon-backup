package back

import (
	"os/exec"
	"regexp"
	"strings"
)

/*
adb devices -l
* daemon not running; starting now at tcp:5037
* daemon started successfully
List of devices attached

List of devices attached
4UX02211060009xx       unauthorized usb:338952192X transport_id:1

List of devices attached
192.168.50.234:5555    device product:NOP-AN00P model:NOP_AN00 device:HWNOP transport_id:6

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
		for _, info := range strings.Split(row[1], " ") {
			kv := strings.Split(info, ":")
			switch len(kv) {
			case 1:
				m["type"] = kv[0]
			case 2:
				m[kv[0]] = kv[1]
			}
		}
		if _, ok := m["usb"]; ok {
			m["id"] = row[0] // id
			m["con"] = "USB"
		} else {
			m["id"] = GetAndroidDeviceInfo(row[0], "ro.boot.serialno")
			m["con"] = row[0] // ip:port
		}
		m["brand"] = GetAndroidDeviceInfo(row[0], "ro.product.brand")
		m["product_name"] = GetAndroidDeviceInfo(row[0], "ro.config.marketing_name")
		if m["product_name"] == "" {
			m["product_name"] = m["brand"] + " " + GetAndroidDeviceInfo(row[0], "ro.product.model")
		}
		m["name"] = GetAndroidDeviceInfo(row[0], "ro.product.name")
		GetAndroidDeviceDiskSpace(row[0], m)
		GetAndroidDeviceBatteryInfo(row[0], m)
		// TODO：获取WIFI信息
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

func GetAndroidDeviceDiskSpace(id string, m map[string]string) {
	output, _ := exec.Command("adb", "-s", id, "shell", "df", "/storage/emulated").CombinedOutput()
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
}

func GetAndroidDeviceBatteryInfo(id string, m map[string]string) {
	output, _ := exec.Command("adb", "-s", id, "shell", "dumpsys", "battery").CombinedOutput()
	info := strings.Split(string(output), "\n")
	if len(info) <= 0 {
		return
	}
	for _, line := range info {
		row := strings.Split(strings.Trim(line, "\r\n "), ": ")
		if len(row) < 2 {
			continue
		}
		key := strings.Trim(row[0], " \t\r\n")
		switch key {
		case "status", "level":
			m[key] = strings.Trim(row[1], " \t\r\n")
		}
	}
}
