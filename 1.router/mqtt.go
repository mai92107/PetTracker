package router

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	mqttApi "batchLog/2.api/mqtt_device"
	"strings"
)

func RouteFunction(topic string, payload model.MqttPayload, qos byte){
	deviceID := extractDeviceIDFromTopic(topic)
	if deviceID == "" {
		logafa.Debug(" ⚠️ 無法解析 deviceId")
		return
	}

	switch payload.Subject{
	case model.LOGIN.ToString():
		mqttApi.Login(deviceID, qos)

	case model.LOCATION.ToString():
		mqttApi.Recording(deviceID,payload.Payload,payload.CurrentTime.Format(global.TIME_FORMAT))

	case model.OFFLINE.ToString():
		mqttApi.Offline(deviceID)
	}
}


func extractDeviceIDFromTopic(topic string) string {
	// topic 格式為：pet/{deviceId}/location
	parts := strings.Split(topic, "/")
	if len(parts) != 3 {
		return ""
	}
	return parts[1]
}