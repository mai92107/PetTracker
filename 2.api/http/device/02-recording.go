package deviceHttp

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

type request02 struct {
	Longitude string `json:"lng"`
	Latitude  string `json:"lat"`
	DeviceID  string `json:"deviceId"`
}

func Recording(c *gin.Context) {
	requestTime := time.Now().UTC()
	var req request02
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, requestTime, global.COMMON_REQUEST_ERROR)
		return
	}
	time := requestTime.Format("2006-01-02 15:04:05")
	jwt := c.GetHeader("jwt")
	_, err = jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.Error(c, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}

	deviceService.Recording(req.Latitude, req.Longitude, req.DeviceID, time)
	response.Success(c, requestTime, "")
}
