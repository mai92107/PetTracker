package deviceMqtt

import (
	"fmt"
	"net/http"
	"time"

	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
)

func MqttOnlineDevice(request, jwt, ip string) {
	requestTime := time.Now().UTC()

	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil || userInfo.Identity != "ADMIN" {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt("ERROR/"+jwt, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	topic := "login/" + fmt.Sprintf("%d", userInfo.MemberId)

	deviceList, err := deviceService.MqttOnlineDevice()
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	response.SuccessMqtt(topic, requestTime, deviceList)
}
