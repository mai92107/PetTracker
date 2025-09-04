package mqttApi

import (
	"time"

	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/global"

	"github.com/gin-gonic/gin"
)

func MqttStatus(c *gin.Context) {
	requestTime := time.Now().UTC()
	global.ActiveDevicesLock.Lock()
	status := "連線正常"
	if !global.IsConnected.Load() {
		status = "連線斷開"
	}
	global.ActiveDevicesLock.Unlock()

	data := map[string]interface{}{
		"message": "寵物追蹤系統運行中",
		"mqtt_status": status,
	}
	response.Success[map[string]interface{}](c,requestTime,data)
}