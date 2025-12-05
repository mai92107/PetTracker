package system

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
)

func SystemStatus(ctx request.RequestContext) {
	global.ActiveDevicesLock.Lock()
	status := "連線正常"
	if !global.IsConnected.Load() {
		status = "連線斷開"
	}
	global.ActiveDevicesLock.Unlock()

	data := map[string]interface{}{
		"message":     "寵物追蹤系統運行中",
		"mqtt_status": status,
	}
	ctx.Success(data)
}
