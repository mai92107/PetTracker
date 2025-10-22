package deviceHttp

import (
	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type request01 struct {
	DeviceType string `json:"deviceType"`
}

func Create(c *gin.Context) {
	requestTime := time.Now().UTC()
	var req request01
	if err := c.ShouldBindJSON(&req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.Error(c, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	jwt := c.GetHeader("jwt")
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil || userData.Identity != "ADMIN" {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.Error(c, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	deviceId, err := deviceService.Create(req.DeviceType, userData.MemberId)
	if err != nil {
		logafa.Error("裝置新增失敗，請稍後嘗試, error: %+v", err)
		response.Error(c, http.StatusInternalServerError, requestTime, "裝置新增失敗，請稍後嘗試")
		return
	}
	response.Success(c, requestTime, deviceId)
}
