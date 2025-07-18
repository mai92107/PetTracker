package router

import (
	"batchLog/api"
	"batchLog/api/account"
	"batchLog/api/device"
	"batchLog/api/home"
	"batchLog/core/global"
	"batchLog/factory"

	"github.com/gin-gonic/gin"
)


func RegisterRoutes(r *gin.Engine) {

	// 初始化 Factory -> Service -> Controller
	accountFactory := factory.NewAccountServiceFactory(&global.Repository.DB, &global.Repository.Cache)
	accountService := accountFactory.CreateService()
	accountController := account.NewAccountController(accountService)

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
		memberGroup.POST("/login",accountController.Login)
		memberGroup.POST("/register",accountController.Register)
	}

	trackGroup := r.Group("/device")
	{
		trackGroup.POST("/tracking", device.Tracking)
	}
}