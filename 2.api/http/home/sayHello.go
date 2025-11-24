package homeHttp

import (
	response "batchLog/0.core/commonResReq/res"
	"time"

	"github.com/gin-gonic/gin"
)

func SayHello(c *gin.Context) {
	requestedTime := time.Now().UTC()

	data := "hello"

	response.Success(c, requestedTime, data)
}
