package deviceService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/model"
	"fmt"
	"time"
)

func MqttDeviceStatus(deviceId string) (*model.DeviceStatus, error) {
	const timeout = 2 * time.Second
	start := time.Now()
	for {
		if global.ActiveDevicesLock.TryLock() {
			break
		}
		if time.Since(start) > timeout {
			// 鎖超時
			return nil, fmt.Errorf("警告：MqttOnlineDevice() 嘗試加鎖超過 2 秒，放棄。")
		}
		time.Sleep(10 * time.Millisecond)
	}
	defer global.ActiveDevicesLock.Unlock()

	info := global.ActiveDevices[deviceId]

	return &info, nil
}
