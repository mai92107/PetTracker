// router/mqtt_router.go
package router

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model/role"
	middleware "batchLog/1.middleware"
	"batchLog/1.router/adapter"
	"batchLog/2.api/account"
	"batchLog/2.api/device"
	"batchLog/2.api/member"
	system "batchLog/2.api/system_config"
	"batchLog/2.api/test"
	"fmt"
	"net/http"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTT route
var mqttRoutes = map[string]Route{

	"hello": {Handler: executeMqtt(test.Hello), Permission: role.GUEST},

	// account
	"account_login":    {Handler: executeMqtt(account.Login), Permission: role.GUEST},
	"account_register": {Handler: executeMqtt(account.Register), Permission: role.GUEST},

	// device
	"device_create":    {Handler: executeMqtt(device.Create), Permission: role.ADMIN},
	"device_recording": {Handler: executeMqtt(device.Recording), Permission: role.MEMBER},
	"device_online":    {Handler: executeMqtt(device.OnlineDeviceList), Permission: role.ADMIN},
	"device_status":    {Handler: executeMqtt(device.DeviceStatus), Permission: role.MEMBER},
	"device_all":       {Handler: executeMqtt(device.DeviceList), Permission: role.ADMIN},

	// trip
	"trip_list": {Handler: executeMqtt(device.TripList), Permission: role.MEMBER},
	"trip_detail":  {Handler: executeMqtt(device.TripDetail), Permission: role.MEMBER},

	// Member
	"member_addDevice": {Handler: executeMqtt(member.AddDevice), Permission: role.MEMBER},
	"member_devices":   {Handler: executeMqtt(member.MemberDeviceList), Permission: role.MEMBER},

	// system
	"system_status": {Handler: executeMqtt(system.SystemStatus), Permission: role.GUEST},
}

type MqttHandler func(request.RequestContext)
type Route struct {
	Handler    MqttHandler
	Permission role.MemberIdentity
}

// topic sample : req/action/clientId/jwt/ip
func RouteFunction(ctx request.RequestContext, action string) {
	// route !!!!!!
	routeInfo, exist := mqttRoutes[action]
	if !exist || routeInfo.Handler == nil {
		logafa.Warn("查無此路徑", "action", action)
		sendBackErrMsg(ctx, "此功能暫未開放")
		return
	}

	// handler !!!!!!!
	handlerWithMiddlewares(
		ctx,
		routeInfo.Handler,

		// middleware
		middleware.MqttJWTMiddleware(routeInfo.Permission),
		middleware.MqttWorkerMiddleware(),
		middleware.MqttTimeoutMiddleware(3*time.Second),
	)
}

func OnMessageReceived(client mqtt.Client, msg mqtt.Message) {
	now := global.GetNow()
	payload := string(msg.Payload())
	topic := msg.Topic()

	logafa.Debug("收到 MQTT 訊息", "topic", topic, "payload", payload)
	action, clientId, jwt, ip := extractInfoFromTopic(topic)
	if action == "" || ip == "" {
		logafa.Warn("無法解析 action 或 ip: %s", topic)
		return
	}
	ctx := adapter.NewMQTTContext(payload, jwt, clientId, ip, now)

	RouteFunction(ctx, action)
}

func extractInfoFromTopic(topic string) (action, clientId, jwt, ip string) {
	parts := strings.Split(topic, "/")
	if len(parts) < 5 {
		return "", "", "", ""
	}
	return parts[1], parts[2], parts[3], parts[4]
}

func sendBackErrMsg(ctx request.RequestContext, reason string, args ...interface{}) {
	fullReason := fmt.Sprintf(reason, args...)
	ctx.Error(http.StatusBadRequest, fullReason)
}

func executeMqtt(handler func(request.RequestContext)) MqttHandler {
	return func(ctx request.RequestContext) {
		handler(ctx)
	}
}

func handlerWithMiddlewares(ctx request.RequestContext,
	handler func(request.RequestContext),
	middlewares ...func(request.RequestContext, func(request.RequestContext)),
) {
	// 用遞迴實作 middleware chain
	var chain func(index int, ctx request.RequestContext)
	chain = func(index int, ctx request.RequestContext) {
		if index < len(middlewares) {
			middlewares[index](ctx, func(ctx request.RequestContext) {
				chain(index+1, ctx)
			})
		} else {
			// 最後執行 handler
			handler(ctx)
		}
	}
	chain(0, ctx)
}
