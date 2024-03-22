package main

import (
	"github.com/gin-gonic/gin"
	"github.com/orcastor/addon-backup/back"
	"github.com/orcastor/orcas/rpc/util"
)

type DeviceInfo struct {
	DeviceID      string `json:"device_id,omitempty"`
	Authorized    bool   `json:"authorized"`
	SerialNo      string `json:"serial_no,omitempty"`
	Connection    string `json:"connection,omitempty"` // USB / <IP>
	ProductName   string `json:"product_name,omitempty"`
	Brand         string `json:"brand,omitempty"`
	OS            string `json:"os,omitempty"` // ANDROID / IOS
	Total         string `json:"total,omitempty"`
	DataAvailable string `json:"data_available,omitempty"`
	SysAvailable  string `json:"sys_available,omitempty"`
}

func list(ctx *gin.Context) {
	var req struct {
	}
	ctx.BindJSON(&req)

	var devs []*DeviceInfo
	androidDevs := back.ListAndroidDevices()
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
	iosDevs := back.ListIOSDevices()
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
