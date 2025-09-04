package mqttApi

import (
	"time"

	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/global"
	"batchLog/0.core/model"

	"github.com/gin-gonic/gin"
)

func DeviceStatus(c *gin.Context) {
	requestTime := time.Now().UTC()
	deviceId := c.Param("deviceId")
	global.ActiveDevicesLock.Lock()
	info := global.ActiveDevices[deviceId]
	global.ActiveDevicesLock.Unlock()

	response.Success[model.DeviceStatus](c,requestTime,info)
}