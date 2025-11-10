package deviceMqtt

import (
	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request01 struct {
	DeviceType  string `json:"deviceType"`
	SubscribeTo string `json:"subscribeTo"`
}

func Create(payload, jwt, clientId, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "errReq/" + clientId

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
	var req request01
	if err := jsoniter.UnmarshalFromString(payload, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	deviceId, err := deviceService.Create(req.DeviceType, userData.MemberId)
	if err != nil {
		logafa.Error("裝置新增失敗, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "裝置新增失敗，請稍後嘗試")
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, deviceId)
}
