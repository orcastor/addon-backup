<p align="center">
  <a href="https://orcastor.github.io/doc/">
    <img src="https://orcastor.github.io/doc/logo.svg">
  </a>
</p>

<p align="center"><strong>OrcaS 手机备份插件</strong></p>

- 是否自动备份开关
- 设备管理 -> 备份管理

- 使用[libusb](https://github.com/gotmc/libusb)监听USB设备热插拔事件
- iOS 使用`idevicebackup2`
- Android 使用`adb+fuse` or `smb/nfs`