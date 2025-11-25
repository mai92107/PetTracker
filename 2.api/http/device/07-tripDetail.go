package deviceHttp

import (
	"net/http"
	"time"

	response "batchLog/0.core/commonResReq/res"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"

	"github.com/gin-gonic/gin"
)

type request07 struct {
	DeviceId string `json:"deviceId"`
	TripUuid string `json:"tripUuid"`
}

func TripDetail(c *gin.Context) {
	requestTime := time.Now().UTC()

	var req request07
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, requestTime, global.COMMON_REQUEST_ERROR)
		return
	}
	jwt := c.GetHeader("jwt")
	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.Error(c, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	info, err := deviceService.GetTripDetail(userInfo, req.DeviceId, req.TripUuid)
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.Error(c, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	response.Success(c, requestTime, info)
}
