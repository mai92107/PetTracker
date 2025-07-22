package router

import (
	api "batchLog/2.api"
	accountApi "batchLog/2.api/account"
	deviceApi "batchLog/2.api/device"
	"batchLog/2.api/home"

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

	memberGroup := r.Group("/account")
	{
		memberGroup.POST("/login",accountApi.Login)
		memberGroup.POST("/register",accountApi.Register)
	}

	trackGroup := r.Group("/device")
	{
		trackGroup.POST("/create", deviceApi.Create)
		trackGroup.POST("/tracking", deviceApi.Tracking)
	}
}