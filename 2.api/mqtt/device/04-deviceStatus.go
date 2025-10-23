package deviceMqtt

import (
	"net/http"
	"time"

	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"

	jsoniter "github.com/json-iterator/go"
)

type request04 struct {
	DeviceID    string `json:"deviceId"`
	SubscribeTo string `json:"subscribeTo"`
}

func DeviceStatus(payload, jwt, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "errReq/device/deviceStatus/" + payload

	if jwt == "" {
		logafa.Error("JWT 參數錯誤, JWT: %s", jwt)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "JWT 參數錯誤")
		return
	}
	if payload == "" || payload == "{}" {
		logafa.Error("Payload 為空")
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Payload 為空")
		return
	}
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil || userData.Identity != "ADMIN" {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	var req request04
	err = jsoniter.UnmarshalFromString(payload, &req)
	if err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	deviceId := req.DeviceID
	info, err := deviceService.MqttDeviceStatus(deviceId)
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, *info)
}
