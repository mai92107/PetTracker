package deviceMqtt

import (
	"fmt"
	"net/http"
	"time"

	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"

	jsoniter "github.com/json-iterator/go"
)

type request04 struct {
	DeviceID string `json:"deviceId"`
}

func DeviceStatus(request, jwt, ip string) {
	requestTime := time.Now().UTC()

	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil || userInfo.Identity != "ADMIN" {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt("ERROR/"+jwt, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	topic := "deviceStatus/" + fmt.Sprintf("%d", userInfo.MemberId)
	var req request04
	err = jsoniter.UnmarshalFromString(request, &req)
	if err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	deviceId := req.DeviceID
	info, err := deviceService.MqttDeviceStatus(deviceId)
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	response.SuccessMqtt(topic, requestTime, *info)
}
