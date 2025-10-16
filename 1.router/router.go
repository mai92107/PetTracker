package router

import (
	accountHttp "batchLog/2.api/http/account"
	deviceHttp "batchLog/2.api/http/device"
	homeHttp "batchLog/2.api/http/home"
	memberHttp "batchLog/2.api/http/member"
	systemHttp "batchLog/2.api/http/system_config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	// 註冊路由
	// TODO: 未來需要檢查ip body header 在路徑後加上middleware檢查
	// 依類別分組
	homeGroup := r.Group("/home")
	{
		homeGroup.GET("/say_hello", homeHttp.SayHello)
	}

	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login", accountHttp.Login)
		accountGroup.POST("/register", accountHttp.Register)
	}

	trackGroup := r.Group("/device")
	{
		trackGroup.POST("/create", deviceHttp.Create)
		trackGroup.POST("/recording", deviceHttp.Recording)
		trackGroup.GET("/onlineDevice", deviceHttp.MqttOnlineDevice)
		trackGroup.GET(":deviceId/status", deviceHttp.DeviceStatus)
	}

	memberGroup := r.Group("/member")
	{
		memberGroup.POST("/addDevice", memberHttp.AddDevice)
	}

	systemGroup := r.Group("/system")
	{
		systemGroup.GET("/status", systemHttp.SystemStatus)
	}
}
