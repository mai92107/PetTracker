// router/mqtt_router.go
package router

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	mqttUtils "batchLog/2.api/mqtt"
	accountMqtt "batchLog/2.api/mqtt/account"
	deviceMqtt "batchLog/2.api/mqtt/device"
	homeMqtt "batchLog/2.api/mqtt/home"
	memberMqtt "batchLog/2.api/mqtt/member"
	systemMqtt "batchLog/2.api/mqtt/system_config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttHandler func(payload, jwt, clientId, ip string)

type Permission int

const (
	PermGuest Permission = iota
	PermMember
	PermAdmin
)

type Route struct {
	Pattern    string
	Handler    MqttHandler
	Permission Permission
}

var mqttRoutes = []Route{
	// Guest (無需 JWT)
	{"account_login", accountMqtt.Login, PermGuest},
	{"account_register", accountMqtt.Register, PermGuest},
	{"home_hello", homeMqtt.SayHello, PermGuest},
	{"system_status", systemMqtt.SystemStatus, PermGuest},
	
	// Admin
	{"device_create", deviceMqtt.Create, PermAdmin},
	{"device_online", deviceMqtt.MqttOnlineDevice, PermAdmin},
	
	// Member
	{"device_recording", deviceMqtt.Recording, PermMember},
	{"member_addDevice", memberMqtt.AddDevice, PermMember},
	{"device_status", deviceMqtt.DeviceStatus, PermMember},
}

// topic sample : req/action/clientId/jwt/ip

func RouteFunction(action, payload, clientId, jwt, ip string) {
	// 查找路由
	for _, route := range mqttRoutes {
		if action != route.Pattern {
			continue
		}

		// 權限檢查
		if !checkPermission(route.Permission, jwt) {
			logafa.Warn("權限不足: %s (client: %s)", route.Pattern, clientId)
			sendBackErrMsg(clientId, "權限不足: %s (client: %s)", route.Pattern, clientId)
			return
		}

		route.Handler(payload, jwt, clientId, ip)
		return
	}

	routes := []string{}
	for _,route := range mqttRoutes{
		routes = append(routes, route.Pattern)
	}

	// === Debug 工具 ===
	switch action {
	case "encrypt":
		mqttUtils.Encrypt(payload, global.ConfigSetting.DefaultSecretKey)
	case "decrypt":
		mqttUtils.Decrypt(payload, global.ConfigSetting.DefaultSecretKey)
	default:
		logafa.Warn("未知 MQTT 請求: %s", action)
		sendBackErrMsg(clientId, "未知 MQTT 請求: %s, 核可請求為: %+v", action, routes)
	}
}

func OnMessageReceived(client mqtt.Client, msg mqtt.Message) {
	payload := string(msg.Payload())
	topic := msg.Topic()

	logafa.Debug("收到 MQTT 訊息！Topic: %s | Payload: %s", topic, payload)

	action, clientId, jwt, ip := extractInfoFromTopic(topic)
	if action == "" || ip == "" {
		logafa.Warn("無法解析 action 或 ip: %s", topic)
		sendBackErrMsg(clientId, "無法解析 action 或 ip: %s", topic)
		return
	}

	// 使用 worker pool 執行
	<-global.NormalWorkerPool
	go func() {
		defer func() {
			global.NormalWorkerPool <- struct{}{}
			if r := recover(); r != nil {
				logafa.Error("MQTT handler panic: %v\n on %s", r, topic)
			}
		}()
		RouteFunction(action, payload, clientId, jwt, ip)
	}()
}

func extractInfoFromTopic(topic string) (action, clientId, jwt, ip string) {
	parts := strings.Split(topic, "/")
	return parts[1], parts[2], parts[3], parts[4]
}

func checkPermission(perm Permission, jwt string) bool {
	switch perm {
	case PermGuest:
		return true
	case PermMember, PermAdmin:
		if jwt == "" {
			return false
		}
		claims, err := jwtUtil.GetUserDataFromJwt(jwt)
		if err != nil {
			logafa.Warn("JWT 解析失敗: %v", err)
			return false
		}
		if perm == PermAdmin && !claims.IsAdmin() {
			return false
		}
		return true
	default:
		return false
	}
}

func sendBackErrMsg(clientId, reason string, args ...interface{}) {
	requestTime := time.Now().UTC()
	errTopic := "errReq/" + clientId
	fullReason := fmt.Sprintf(reason, args...)
	response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, fullReason)
}
