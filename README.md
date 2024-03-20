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

- iOS 使用[`idevicebackup2`](https://github.com/libimobiledevice/libimobiledevice)
- Android
  - 使用`adb` + [`better-adb-sync`](https://github.com/jb2170/better-adb-sync)
  - 使用`smb`协议
    - [自动备份华为手机系统及文件到NAS](https://www.oureiq.top:8812/2023/02/09/%E8%87%AA%E5%8A%A8%E5%A4%87%E4%BB%BD%E5%8D%8E%E4%B8%BA%E6%89%8B%E6%9C%BA%E7%B3%BB%E7%BB%9F%E5%8F%8A%E6%96%87%E4%BB%B6%E5%88%B0nas/)
    - [手把手教你把华为手机完整备份到NAS](https://www.cnblogs.com/djd66/p/16635579.html)

## 界面设计

- 是否自动备份开关
- 备份首页默认展示设备管理
  - 可以隐藏未连接设备
  - 按备份时间排序

### 注意事项 

- `yarn run build:pro`打包后的文件和webapp的放置到一起：
`ln -s $(addon-backup)/front/dist/ $(webapp)/dist/bak`
