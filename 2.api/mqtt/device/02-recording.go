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
	Longitude   string `json:"lng"`
	Latitude    string `json:"lat"`
	DeviceID    string `json:"deviceId"`
	SubscribeTo string `json:"subscribeTo"`
}

func Recording(payload, jwt, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "error/device/recording/" + payload

	var req request02
	err := jsoniter.UnmarshalFromString(payload, &req)
	if err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	time := requestTime.Format("2006-01-02 15:04:05")
	err = deviceService.Recording(req.Latitude, req.Longitude, req.DeviceID, time)
	if err != nil {
		logafa.Error("裝置回傳資料錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "裝置回傳資料錯誤")
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, "")
}
