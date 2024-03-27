package main

import (
	"github.com/gin-gonic/gin"
	"github.com/orcastor/addon-backup/back"
	"github.com/orcastor/orcas/rpc/util"
)

type DeviceInfo struct {
	ID            string `json:"id,omitempty"`
	Authorized    bool   `json:"authorized"`
	SerialNo      string `json:"serial_no,omitempty"`
	Name          string `json:"name,omitempty"`
	Connection    string `json:"connection,omitempty"` // USB / <IP>
	ProductName   string `json:"product_name,omitempty"`
	Brand         string `json:"brand,omitempty"`
	OS            string `json:"os,omitempty"` // ANDROID / IOS
	Total         string `json:"total,omitempty"`
	DataAvailable string `json:"data_available,omitempty"`
	SysAvailable  string `json:"sys_available,omitempty"`
	IsCharging    bool   `json:"is_charing,omitempty"`
	BatteryLevel  string `json:"battery_level,omitempty"`
	Encrypt       bool   `json:"encrypt,omitempty"`
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
				ID:         dev["id"],
				Authorized: false,
			})
		} else {
			devs = append(devs, &DeviceInfo{
				ID:            dev["id"],
				Authorized:    true,
				SerialNo:      dev["id"],
				Name:          dev["name"],
				Connection:    dev["con"],
				ProductName:   dev["product_name"],
				Brand:         dev["brand"],
				OS:            "ANDROID",
				Total:         dev["total"],
				DataAvailable: dev["available"],
				IsCharging:    dev["status"] == "2",
				BatteryLevel:  dev["level"],
			})
		}
	}
	iosDevs := back.ListIOSDevices()
	for _, dev := range iosDevs {
		if len(dev) <= 2 {
			devs = append(devs, &DeviceInfo{
				ID:         dev["id"],
				Authorized: false,
			})
		} else {
			devs = append(devs, &DeviceInfo{
				ID:            dev["id"],
				Authorized:    true,
				SerialNo:      dev["SerialNumber"],
				Name:          dev["DeviceName"],
				Connection:    dev["con"],
				ProductName:   dev["product_name"],
				Brand:         dev["brand"],
				OS:            "IOS",
				Total:         dev["TotalDiskCapacity"],
				DataAvailable: dev["TotalDataAvailable"],
				SysAvailable:  dev["TotalSystemAvailable"],
				IsCharging:    dev["BatteryIsCharging"] == "true",
				BatteryLevel:  dev["BatteryCurrentCapacity"],
				Encrypt:       dev["WillEncrypt"] == "1",
			})
		}
	}

	util.Response(ctx, gin.H{
		"devs":  devs,
		"count": len(devs),
	})
}
