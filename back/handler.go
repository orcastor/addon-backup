package back

import (
	"github.com/gin-gonic/gin"
	"github.com/orcastor/orcas/rpc/util"
)

type DeviceInfo struct {
	DeviceID      string
	Authorized    bool
	SerialNo      string
	Connection    string // USB / IP
	ProductName   string
	Brand         string
	OS            string // ANDROID / IOS
	Total         string
	DataAvailable string
	SysAvailable  string
}

func list(ctx *gin.Context) {
	var req struct {
	}
	ctx.BindJSON(&req)

	var devs []*DeviceInfo
	androidDevs := ListAndroidDevices()
	for _, dev := range androidDevs {
		if dev["type"] == "unauthorized" {
			devs = append(devs, &DeviceInfo{
				DeviceID:   dev["id"],
				Authorized: false,
			})
		} else {
			devs = append(devs, &DeviceInfo{
				DeviceID:      dev["id"],
				Authorized:    true,
				SerialNo:      dev["id"],
				Connection:    dev["con"],
				ProductName:   dev["name"],
				Brand:         dev["brand"],
				OS:            "ANDROID",
				Total:         dev["total"],
				DataAvailable: dev["available"],
			})
		}
	}
	iosDevs := ListIOSDevices()
	for _, dev := range iosDevs {
		if len(dev) <= 2 {
			devs = append(devs, &DeviceInfo{
				DeviceID:   dev["id"],
				Authorized: false,
			})
		} else {
			devs = append(devs, &DeviceInfo{
				DeviceID:      dev["id"],
				Authorized:    true,
				SerialNo:      dev["SerialNumber"],
				Connection:    dev["con"],
				ProductName:   dev["name"],
				Brand:         dev["brand"],
				OS:            "IOS",
				Total:         dev["TotalDiskCapacity"],
				DataAvailable: dev["TotalDataAvailable"],
				SysAvailable:  dev["TotalSystemAvailable"],
			})
		}
	}

	util.Response(ctx, gin.H{
		"devs":  devs,
		"count": len(devs),
	})
}
