// 1.middleware/worker.go
package middleware

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpWorkerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		<-global.NormalWorkerPool
		logafa.Debug("獲取 Worker (HTTP)")

		defer func() {
			global.NormalWorkerPool <- struct{}{}
			logafa.Debug("釋放 Worker (HTTP)")

			if r := recover(); r != nil {
				logafa.Error("HTTP Handler panic recovered: %v", r)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": global.COMMON_SYSTEM_ERROR,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func MqttWorkerMiddleware() func(request.RequestContext, func(ctx request.RequestContext)) {
	return func(ctx request.RequestContext, next func(ctx request.RequestContext)) {
		<-global.NormalWorkerPool
		logafa.Debug("獲取 Worker (MQTT)")

		defer func() {
			global.NormalWorkerPool <- struct{}{}
			logafa.Debug("釋放 Worker (MQTT)")

			if r := recover(); r != nil {
				logafa.Error("MQTT Handler panic recovered: %v", r)
				ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
			}
		}()
		next(ctx)
	}
}
