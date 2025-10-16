package mqttApi

import (
	"batchLog/0.core/global"
	"batchLog/0.core/model"
	"time"
)

func Login(deviceId string, qos byte) {
	global.ActiveDevicesLock.Lock()
	global.ActiveDevices[deviceId] = model.DeviceStatus{
		QoS:      qos,
		LastSeen: time.Now().UTC(),
		Online:   true,
	}
	global.ActiveDevicesLock.Unlock()
}
