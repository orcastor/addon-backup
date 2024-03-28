package main

import (
	"github.com/gin-gonic/gin"
	"github.com/orcastor/addon-backup/back"
	"github.com/orcastor/orcas/rpc/util"
	"github.com/orcastor/phone_images/sdk"
)

type DeviceInfo struct {
	ID          string `json:"id,omitempty"`
	Authorized  bool   `json:"authorized"`
	SerialNo    string `json:"serial_no,omitempty"`
	Name        string `json:"name,omitempty"`
	Connection  string `json:"connection,omitempty"` // USB / <IP>
	ProductName string `json:"product_name,omitempty"`
	Brand       string `json:"brand,omitempty"`
	OS          string `json:"os,omitempty"` // ANDROID / IOS
	Version     string `json:"version,omitempty"`
	ImgURL      string `json:"img_url,omitempty"`

	Total         string `json:"total,omitempty"`
	DataAvailable string `json:"data_available,omitempty"`
	SysAvailable  string `json:"sys_available,omitempty"`

	IsCharging   bool   `json:"is_charing,omitempty"`
	BatteryLevel string `json:"battery_level,omitempty"`
	WillEncrypt  bool   `json:"will_encrypt,omitempty"`

	Progress  string `json:"progress,omitempty"`
	LastTime  string `json:"last_time,omitempty"`
	LastError string `json:"last_error,omitempty"`
}

/*
	{
	    "code": 0,
	    "data": {
	        "count": 1,
	        "devs": [
	            {
	                "id": "d2e6c832cd7d0a4cff535b59e4f567bbfef65dc2",
	                "authorized": true,
	                "serial_no": "F2LSQ3H8HFY5",
	                "name": "“PP”的 iPhone",
	                "connection": "USB",
	                "product_name": "iPhone 7 Plus",
	                "brand": "APPLE",
	                "os": "IOS",
	                "version": "15.8.2",
	                "total": "256000000000",
	                "data_available": "200961683456",
	                "sys_available": "0",
	                "is_charing": true,
	                "battery_level": "43"
	            }
	        ]
	    }
	}
*/
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
				ImgURL:        sdk.GetAndroidURL(sdk.GetAndroidProductName(dev["brand"], dev["name"])),
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
				Version:       dev["ProductVersion"],
				Total:         dev["TotalDiskCapacity"],
				DataAvailable: dev["TotalDataAvailable"],
				SysAvailable:  dev["TotalSystemAvailable"],
				IsCharging:    dev["BatteryIsCharging"] == "true",
				BatteryLevel:  dev["BatteryCurrentCapacity"],
				WillEncrypt:   dev["WillEncrypt"] == "1",
				ImgURL:        sdk.GetIOSURL(dev["product_name"], dev["ModelNumber"]),
			})
		}
	}

	util.Response(ctx, gin.H{
		"devs":  devs,
		"count": len(devs),
	})
}
