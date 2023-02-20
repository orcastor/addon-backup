package main

import (
	"fmt"
	"testing"
)

func TestListAndroidDevices(t *testing.T) {
	devices := ListAndroidDevices()
	fmt.Println(devices)
}
