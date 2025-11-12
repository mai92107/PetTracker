package deviceHttp

import (
	"net/http"
	"time"

	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"

	"github.com/gin-gonic/gin"
)

func DeviceStatus(c *gin.Context) {
	requestTime := time.Now().UTC()
	deviceId := c.Param("deviceId")

	jwt := c.GetHeader("jwt")
	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.Error(c, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	info, err := deviceService.MqttDeviceStatus(deviceId, userInfo)
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.Error(c, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	response.Success(c, requestTime, info)
}
