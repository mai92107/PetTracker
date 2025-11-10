package middleware

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"

	"github.com/gin-gonic/gin"
)

func WorkerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 搶 worker（會阻塞直到有空位）
		<-global.NormalWorkerPool

		// 確保最後釋放
		defer func() {
			global.NormalWorkerPool <- struct{}{}
			if r := recover(); r != nil {
				logafa.Error("Handler panic recovered: %v", r)
				c.JSON(500, gin.H{"error": "internal server error"})
			}
		}()

		// 執行真正的 handler
		logafa.Debug("%s 請完工人, 開始工作",c.Request.RequestURI)
		c.Next()
	}
}
