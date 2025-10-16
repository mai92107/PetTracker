package router

import (
	"batchLog/0.core/logafa"
	mqttUtils "batchLog/2.api/mqtt"
	accountMqtt "batchLog/2.api/mqtt/account"
	deviceMqtt "batchLog/2.api/mqtt/device"
	homeMqtt "batchLog/2.api/mqtt/home"
	memberMqtt "batchLog/2.api/mqtt/member"
	systemMqtt "batchLog/2.api/mqtt/mqtt_system"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RouteFunction(topic string, payload string, qos byte) {

	if strings.HasPrefix(topic, "request") {
		requestType, jwt, ip := extractRequestFromTopic(topic)

		switch requestType {
		// no need verify jwt
		case "account_login":
			accountMqtt.Login(payload, ip)
		case "account_register":
			accountMqtt.Register(payload, jwt, ip)
		case "home_hello":
			homeMqtt.SayHello(jwt)
		case "config_status":
			systemMqtt.SystemStatus(payload, jwt, ip)

		// need verify jwt
		case "device_create":
			deviceMqtt.Create(payload, jwt, ip)
		case "device_recording":
			deviceMqtt.Recording(payload, jwt, ip)
		case "device_online":
			deviceMqtt.MqttOnlineDevice(payload, jwt, ip)
		case "device_status":
			deviceMqtt.DeviceStatus(payload, jwt, ip)
		case "member_addDevice":
			memberMqtt.AddDevice(payload, jwt, ip)

		default:
			logafa.Warn("⚠️ 未知的 request 類型: %s (topic: %s), payload: %+v", requestType, topic, payload)
		}

		// debug utils
		switch topic {
		case "request/encrypt":
			mqttUtils.Encrypt(payload)
		case "request/decrypt":
			mqttUtils.Decrypt(payload)
		}
	}
}

// 處理接收到的訊息
func OnMessageReceived(client mqtt.Client, msg mqtt.Message) {
	logafa.Debug("📥 收到 MQTT 訊息！")
	logafa.Debug("主題: %s", msg.Topic())
	logafa.Debug("內容: %s", string(msg.Payload()))

	RouteFunction(msg.Topic(), string(msg.Payload()), msg.Qos())
}

func extractRequestFromTopic(topic string) (string, string, string) {
	// topic 格式為：request/{hashedValue}
	parts := strings.Split(topic, "/")
	if len(parts) != 2 {
		return "", "", ""
	}
	hashedValue, _ := mqttUtils.Decrypt(parts[1])
	// 解碼 hashedValue 以取得 requestType 和 ip
	// hashedValue 的格式為 {requestType}-{jwt}-{ip}
	parts = strings.Split(hashedValue, "-")
	if len(parts) != 3 {
		return "", "", ""
	}
	return parts[0], parts[1], parts[2]
}
