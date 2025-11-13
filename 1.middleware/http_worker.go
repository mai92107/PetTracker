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
		defer func() {
			global.NormalWorkerPool <- struct{}{}
			logafa.Debug("工作完畢")
			if r := recover(); r != nil {
				logafa.Error("Handler panic recovered: %v", r)
				c.JSON(500, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
