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
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
	DeviceID  string `json:"deviceId"`
	RecordAt  string  `json:"recordAt"`
}

func Recording(c *gin.Context) {
	requestTime := time.Now().UTC()
	var req request02
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, requestTime, global.COMMON_REQUEST_ERROR)
		return
	}
	jwt := c.GetHeader("jwt")
	claim, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.Error(c, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}

	deviceService.Recording(req.Latitude, req.Longitude, claim.MemberId, req.DeviceID, req.RecordAt)
	response.Success(c, requestTime, "")
}
