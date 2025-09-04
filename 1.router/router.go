package router

import (
	api "batchLog/2.api"
	accountApi "batchLog/2.api/account"
	deviceApi "batchLog/2.api/device"
	"batchLog/2.api/home"
	memberApi "batchLog/2.api/member"
	mqttApi "batchLog/2.api/mqtt_system"

	"github.com/gin-gonic/gin"
)


func RegisterRoutes(r *gin.Engine) {

	// 註冊路由
	// TODO: 未來需要檢查ip body header 在路徑後加上middleware檢查
	// 依類別分組
	r.GET("/health-check", api.CheckHealth)
	homeGroup := r.Group("/home")
	{
		homeGroup.GET("/say_hello", home.SayHello)
	}

	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login",accountApi.Login)
		accountGroup.POST("/register",accountApi.Register)
	}

	trackGroup := r.Group("/device")
	{
		trackGroup.POST("/create", deviceApi.Create)
		// trackGroup.POST("/tracking", deviceApi.Tracking)
	}

	memberGroup := r.Group("/member")
	{
		memberGroup.POST("/addDevice", memberApi.AddDevice)
	}

	mqttGroup := r.Group("/mqtt")
	{
		mqttGroup.GET("/status", mqttApi.MqttStatus)
		mqttGroup.GET("/onlineDevice", mqttApi.MqttOnlineDevice)
		mqttGroup.GET(":deviceId/status", mqttApi.DeviceStatus)
	}
}