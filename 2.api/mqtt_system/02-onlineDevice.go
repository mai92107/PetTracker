package mqttApi

import (
	"time"

	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/global"

	"github.com/gin-gonic/gin"
)

func MqttOnlineDevice(c *gin.Context) {
	requestTime := time.Now().UTC()

	global.ActiveDevicesLock.Lock()
	deviceList := make([]string, 0, len(global.ActiveDevices))
	for deviceId, info := range global.ActiveDevices {
		if !info.Online{
			continue
		}
		deviceList = append(deviceList, deviceId)
	}
	global.ActiveDevicesLock.Unlock()

	response.Success[[]string](c,requestTime,deviceList)
}