package systemService

import (
	"batchLog/0.core/global"
	"fmt"
	"time"
)

func SystemStatus() (map[string]interface{}, error) {
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
	status := "連線正常"
	if !global.IsConnected.Load() {
		status = "連線斷開"
	}

	data := map[string]interface{}{
		"message":     "寵物追蹤系統運行中",
		"mqtt_status": status,
	}
	return data, nil
}
