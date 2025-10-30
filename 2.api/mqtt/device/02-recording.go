package deviceMqtt

import (
	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request02 struct {
	Longitude   float64 `json:"lng"`
	Latitude    float64 `json:"lat"`
	DeviceID    string  `json:"deviceId"`
	SubscribeTo string  `json:"subscribeTo"`
	RecordAt    string  `json:"recordAt"`
}

func Recording(payload, jwt, clientId, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "errReq/" + clientId

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
	var req request02
	err := jsoniter.UnmarshalFromString(payload, &req)
	if err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	err = deviceService.Recording(req.Latitude, req.Longitude, req.DeviceID, req.RecordAt)
	if err != nil {
		logafa.Error("裝置回傳資料錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "裝置回傳資料錯誤")
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, "")
}
