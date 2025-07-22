package deviceApi

import (
	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type request02 struct{
	longitude	string	`json:"lng"`
	latitude	string	`json:"lat"`
}
func Tracking(c *gin.Context){
	requestTime := time.Now().UTC()
	var req request02
	err := c.ShouldBindJSON(&req)
	if err != nil{
		response.Error(c,http.StatusBadRequest,requestTime,global.COMMON_REQUEST_ERROR)
		return
	}
	time := requestTime.Format("2006-01-02 15:04:05")
	jwt := c.GetHeader("jwt")
	userData,err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil{
		logafa.Error("身份認證錯誤, error: %+v",err)
		response.Error(c,http.StatusForbidden,requestTime,"身份認證錯誤")
		return
	}

	err = deviceService.Tracking(req.latitude,req.longitude, userData.DeviceID, userData.AccountName, time)
	if err != nil{
		response.Error(c,http.StatusInternalServerError,requestTime,global.COMMON_SYSTEM_ERROR)
		return
	}

	response.Success[string](c,requestTime,"")
}