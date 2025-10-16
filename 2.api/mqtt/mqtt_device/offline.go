package mqttApi

import (
	"batchLog/0.core/global"
	"time"
)

func Offline(deviceId string){
	global.ActiveDevicesLock.Lock()
	if status, ok := global.ActiveDevices[deviceId]; ok {
		status.Online = false
		status.LastSeen = time.Now().UTC()
		global.ActiveDevices[deviceId] = status
	}
	global.ActiveDevicesLock.Unlock()	
}