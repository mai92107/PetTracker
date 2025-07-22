package main

import (
	"batchLog/0.core/global"
	"batchLog/0.core/initial"
	"batchLog/0.core/logafa"
	"batchLog/1.router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	initial.InitAll()

	port := global.ConfigSetting.Port
	logafa.Info(" http server 已啟動, 使用 PORT: %s", port)

	r := gin.Default()
	router.RegisterRoutes(r)

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		logafa.Error(" 伺服器無法啟動, PORT: %s, error: %v", port, err)
	}
}