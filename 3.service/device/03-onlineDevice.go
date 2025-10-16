package deviceService

import (
	"batchLog/0.core/global"
	"fmt"
	"time"
)

func MqttOnlineDevice() ([]string, error) {
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

	deviceList := make([]string, 0, len(global.ActiveDevices))
	for deviceId, info := range global.ActiveDevices {
		if !info.Online {
			continue
		}
		deviceList = append(deviceList, deviceId)
	}
	return deviceList, nil
}
