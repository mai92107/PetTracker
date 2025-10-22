package systemMqtt

import (
	"net/http"
	"time"

	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/logafa"
	systemService "batchLog/3.service/system"

	jsoniter "github.com/json-iterator/go"
)

type request01 struct {
	SubscribeTo string `json:"subscribeTo"`
}

func SystemStatus(payload string) {
	requestTime := time.Now().UTC()

	errTopic := "error/system/status/" + payload
	var req request01
	if err := jsoniter.UnmarshalFromString(payload, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	data, err := systemService.SystemStatus()
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, data)
}
