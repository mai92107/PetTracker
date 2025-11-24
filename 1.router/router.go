package router

import (
	"batchLog/0.core/model/role"
	middleware "batchLog/1.middleware"
	accountHttp "batchLog/2.api/http/account"
	debugHttp "batchLog/2.api/http/debug"
	deviceHttp "batchLog/2.api/http/device"
	homeHttp "batchLog/2.api/http/home"
	memberHttp "batchLog/2.api/http/member"
	systemHttp "batchLog/2.api/http/system_config"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	const ADMIN = "ADMIN"
	const MEMBER = "MEMBER"

	r.Use(middleware.WorkerMiddleware())

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

	deviceGroup := r.Group("/device")
	{
		deviceGroup.POST("/create", middleware.JWTValidator(role.ADMIN), deviceHttp.Create)
		deviceGroup.GET("/onlineDevice", middleware.JWTValidator(role.ADMIN), deviceHttp.MqttOnlineDevice)
		deviceGroup.GET("/all", middleware.JWTValidator(role.ADMIN), deviceHttp.AllDevice)

		deviceGroup.POST("/recording", middleware.JWTValidator(role.MEMBER), deviceHttp.Recording)
		deviceGroup.GET("/:deviceId/status", middleware.JWTValidator(role.MEMBER), deviceHttp.DeviceStatus)
		deviceGroup.GET("/trips", middleware.JWTValidator(role.MEMBER), deviceHttp.DeviceTrips)
		deviceGroup.GET("/trip", middleware.JWTValidator(role.MEMBER), deviceHttp.TripDetail)

	}

	memberGroup := r.Group("/member")
	{
		memberGroup.POST("/addDevice", middleware.JWTValidator(role.MEMBER), memberHttp.AddDevice)
		memberGroup.GET("/allDevice", middleware.JWTValidator(role.MEMBER), memberHttp.MemberDevice)
	}

	systemGroup := r.Group("/system")
	{
		systemGroup.GET("/status", middleware.JWTValidator(role.MEMBER), systemHttp.SystemStatus)
	}

	debugGroup := r.Group("/debug")
	{
		debugGroup.POST("/flush_to_maria", middleware.JWTValidator(role.ADMIN), debugHttp.FlushToMaria)
	}
}
