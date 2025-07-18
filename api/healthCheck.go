package api

import (
	response "batchLog/core/commonResponse"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckHealth(c *gin.Context){
	requestTime := time.Now().UTC()
	response.Success[string](c,requestTime,"")
}