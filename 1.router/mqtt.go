package router

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	mqttUtils "batchLog/2.api/mqtt"
	accountMqtt "batchLog/2.api/mqtt/account"
	deviceMqtt "batchLog/2.api/mqtt/device"
	homeMqtt "batchLog/2.api/mqtt/home"
	memberMqtt "batchLog/2.api/mqtt/member"
	systemMqtt "batchLog/2.api/mqtt/system_config"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func RouteFunction(topic string, payload string, qos byte) {

	if strings.HasPrefix(topic, "req") {
		requestType, jwt, ip := extractRequestFromTopic(topic)

		switch requestType {
		// no need verify jwt
		case "account_login":
			accountMqtt.Login(payload, ip)
		case "account_register":
			accountMqtt.Register(payload, ip)
		case "home_hello":
			homeMqtt.SayHello(payload)
		case "config_status":
			systemMqtt.SystemStatus(payload)

		// need verify jwt
		// admin
		case "device_create":
			deviceMqtt.Create(payload, jwt, ip)
		case "device_online":
			deviceMqtt.MqttOnlineDevice(payload, jwt, ip)
		case "device_status":
			deviceMqtt.DeviceStatus(payload, jwt, ip)

		// member
		case "device_recording":
			deviceMqtt.Recording(payload, jwt, ip)
		case "member_addDevice":
			memberMqtt.AddDevice(payload, jwt, ip)

		// debug utils
		case "encrypt":
			mqttUtils.Encrypt(payload, global.ConfigSetting.DefaultSecretKey)
		case "decrypt":
			mqttUtils.Decrypt(payload, global.ConfigSetting.DefaultSecretKey)

		default:
			logafa.Warn("⚠️ 未知的 request 類型: %s (topic: %s), payload: %+v", requestType, topic, payload)
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

func extractRequestFromTopic(topic string) (requestType, jwt, ip string) {

	// 取得 requestType 和 ip
	// hashedValue 的格式為 request/{requestType}/{jwt}/{ip}
	parts := strings.Split(topic, "/")
	if len(parts) < 4 {
		return "", "", ""
	}
	return parts[1], parts[2], parts[3]
}
