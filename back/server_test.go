package main

import (
	"fmt"
	"testing"
)

func TestListAndroidDevices(t *testing.T) {
	devices := ListAndroidDevices()
	fmt.Println(devices)
	/*
		"id":"4UX02211060009xx"
		"type":"device"
		"usb":"338952192X"
		"product":"NOP-AN00P"
		"model":"NOP_AN00"
		"device":"HWNOP"
		"transport_id":"2"
		"name":"PORSCHE DESIGN HUAWEI Mate 40"
	*/
}
