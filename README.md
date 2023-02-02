<p align="center">
  <a href="https://orcastor.github.io/doc/">
    <img src="https://orcastor.github.io/doc/logo.svg">
  </a>
</p>

<p align="center"><strong>OrcaS 手机备份插件</strong></p>

## 监听USB设备热插拔事件

> 需要先安装[libusb](https://github.com/gotmc/libusb)的C库

### OS X

```bash
$ brew install libusb
```

### Windows

从[libusb.info](https://libusb.info)下载最新的二进制文件

### Linux

```bash
$ sudo apt-get install -y libusb-dev libusb-1.0-0-dev
```

## 实现细节

- iOS 使用`idevicebackup2`
- Android 使用`adb+fuse` or `smb/nfs`

## 界面设计

- 是否自动备份开关
- 设备管理 -> 备份管理