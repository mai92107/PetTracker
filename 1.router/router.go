package router

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/model/role"
	middleware "batchLog/1.middleware"
	"batchLog/1.router/adapter"
	"batchLog/2.api/account"
	"batchLog/2.api/debug"
	"batchLog/2.api/device"
	"batchLog/2.api/member"
	system "batchLog/2.api/system_config"
	"batchLog/2.api/test"
	"batchLog/2.api/trip"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	// middleware
	r.Use(
		middleware.HttpTimeoutMiddleware(5*time.Second),
		middleware.HttpWorkerMiddleware(),
	)

	// 註冊路由
	// TODO: 未來需要檢查ip body header 在路徑後加上middleware檢查
	// 依類別分組
	homeGroup := r.Group("/home")
	{
		homeGroup.GET("/say_hello", required(role.GUEST), executeHttp(test.Hello))
	}

	accountGroup := r.Group("/account")
	{
		accountGroup.POST("/login", required(role.GUEST), executeHttp(account.Login))
		accountGroup.POST("/register", required(role.GUEST), executeHttp(account.Register))
	}

	deviceGroup := r.Group("/device")
	{
		deviceGroup.POST("/create", required(role.ADMIN), executeHttp(device.Create))
		deviceGroup.GET("/onlineDevice", required(role.ADMIN), executeHttp(device.OnlineDeviceList))
		deviceGroup.GET("/all", required(role.ADMIN), executeHttp(device.DeviceList))

		deviceGroup.POST("/recording", required(role.MEMBER), executeHttp(device.Recording))
		deviceGroup.GET("/:deviceId/status", required(role.MEMBER), executeHttp(device.DeviceStatus))
	}

	tripGroup := r.Group("/trip")
	{
		tripGroup.GET("/list", required(role.MEMBER), executeHttp(trip.TripList))
		tripGroup.GET("/detail", required(role.MEMBER), executeHttp(trip.TripDetail))

	}

	memberGroup := r.Group("/member")
	{
		memberGroup.POST("/addDevice", required(role.MEMBER), executeHttp(member.AddDevice))
		memberGroup.GET("/allDevice", required(role.MEMBER), executeHttp(member.MemberDeviceList))
	}

	systemGroup := r.Group("/system")
	{
		systemGroup.GET("/status", required(role.MEMBER), executeHttp(system.SystemStatus))
	}

	debugGroup := r.Group("/debug")
	{
		debugGroup.POST("/flush_to_maria", required(role.ADMIN), executeHttp(debug.FlushToMaria))
	}
}

func executeHttp(handler func(request.RequestContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := adapter.NewHttpContext(c)
		handler(ctx)
	}
}

func required(identity role.MemberIdentity) gin.HandlerFunc {
	return middleware.HttpJWTMiddleware(identity)
}
