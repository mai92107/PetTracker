package mqttApi

// import (
// 	"batchLog/0.core/global"
// 	mqttService "batchLog/3.service/mqtt"
// 	"time"

// 	jsoniter "github.com/json-iterator/go"
// )

// func Recording(deviceId, locationJson,requestTime  string){
// 	global.ActiveDevicesLock.Lock()
// 	if status, ok := global.ActiveDevices[deviceId]; ok {
// 		status.LastSeen = time.Now().UTC()
// 		global.ActiveDevices[deviceId] = status
// 	}
// 	location := map[string]string{}
// 	jsoniter.UnmarshalFromString(locationJson, &location)

// 	mqttService.Recording(location["lat"],location["lng"],deviceId, requestTime)

// 	global.ActiveDevicesLock.Unlock()
// }
