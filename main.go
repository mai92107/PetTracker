package main

import (
	"batchLog/0.core/global"
	"batchLog/0.core/initial"
	"batchLog/0.core/logafa"
	router "batchLog/1.router"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	initial.InitAll()

	port := global.ConfigSetting.Port
	logafa.Info(" http server 已啟動, 使用 PORT: %s", port)

	r := gin.Default()
	router.RegisterRoutes(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// ✅ 使用 goroutine 執行伺服器（非阻塞）
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logafa.Error("伺服器啟動失敗, error: %v", err)
		}
	}()
	// ✅ 啟動優雅關閉邏輯
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // 等待訊號
	logafa.Info("收到終止訊號，開始優雅關閉...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logafa.Error("伺服器優雅關閉失敗, %+v", err)
	} else {
		logafa.Info("伺服器成功關閉")
	}
}
