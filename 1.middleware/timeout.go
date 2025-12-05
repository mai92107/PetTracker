package middleware

import (
	request "batchLog/0.core/commonResReq/req"
	response "batchLog/0.core/commonResReq/res"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 為每個 request 套用 timeout
func HttpTimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := global.GetNow()
		// 建立可取消 context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		c.Request = c.Request.WithContext(ctx)

		// 建立 channel 用來接結果
		finished := make(chan struct{})
		panicChan := make(chan interface{})
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			close(finished)
			cancel()
		}()
		// 用 select 監控三種結果
		select {
		case <-finished:
			return
		case p := <-panicChan:
			logafa.Error("發生Panic", "panic", p)

		case <-ctx.Done():
			// timeout，取消 handler 實作
			response.Error(c, http.StatusGatewayTimeout, now, "系統處理超時")
			c.Abort()
			return
		}
	}
}

func MqttTimeoutMiddleware(timeout time.Duration) func(request.RequestContext, func(ctx request.RequestContext)) {
	return func(ctx request.RequestContext, next func(ctx request.RequestContext)) {
		// 建立可取消 context
		newCtx, cancel := context.WithTimeout(ctx.GetContext(), timeout)
		ctx.SetContext(newCtx)
		ctx.SetCancel(cancel)

		// 建立 channel 用來接結果
		finished := make(chan struct{})
		panicChan := make(chan interface{})
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			next(ctx)
			close(finished)
		}()
		// 用 select 監控三種結果
		select {
		case <-finished:
			cancel()
			return
		case p := <-panicChan:
			cancel()
			logafa.Error("發生Panic", "panic", p)
			ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		case <-ctx.GetContext().Done():
			cancel()
			ctx.Error(http.StatusGatewayTimeout, "系統處理超時")
			return
		}
	}
}
