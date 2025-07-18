package device

import (
	response "batchLog/core/commonResponse"
	"batchLog/core/global"
	gormTable "batchLog/core/gorm"
	"batchLog/service/device"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


func Tracking(c *gin.Context){
	requestTime := time.Now().UTC()
	gps := gormTable.GPS{}
	err := c.ShouldBindJSON(&gps)
	if err != nil{
		response.Error(c,http.StatusBadRequest,requestTime,global.COMMON_REQUEST_ERROR)
		return
	}
	gps.RequestTime = requestTime.Format("2006-01-02 15:04:05")
	err = device.Tracking(gps, requestTime)
	if err != nil{
		response.Error(c,http.StatusInternalServerError,requestTime,global.COMMON_SYSTEM_ERROR)
		return
	}

	response.Success[string](c,requestTime,"")
}