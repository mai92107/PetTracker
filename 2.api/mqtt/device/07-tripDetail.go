package deviceMqtt

import (
	"net/http"
	"time"

	request "batchLog/0.core/commonResReq/req"
	response "batchLog/0.core/commonResReq/res"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"

	jsoniter "github.com/json-iterator/go"
)

type request07 struct {
	DeviceId string `json:"deviceId"`
	TripUuid string `json:"tripUuid"`
	request.MqttReq
}

func TripDetail(payload, jwt, clientId, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "errReq/" + clientId

	if payload == "" || payload == "{}" {
		logafa.Error("Payload 為空")
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Payload 為空")
		return
	}
	var req request07
	err := jsoniter.UnmarshalFromString(payload, &req)
	if err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}

	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "身份認證錯誤")
		return
	}
	info, err := deviceService.GetTripDetail(userInfo, req.DeviceId, req.TripUuid)
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "系統發生錯誤")
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, info)
}
