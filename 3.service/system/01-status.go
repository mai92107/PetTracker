package systemService

import (
	"batchLog/0.core/global"
)

func SystemStatus() (map[string]interface{}, error) {
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
