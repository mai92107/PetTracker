package deviceMqtt

import (
	"net/http"
	"time"

	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"

	jsoniter "github.com/json-iterator/go"
)

type request03 struct {
	SubscribeTo string `json:"subscribeTo"`
}

func MqttOnlineDevice(payload, jwt, clientId, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "errReq/" + clientId

	if payload == "" || payload == "{}" {
		logafa.Error("Payload 為空")
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Payload 為空")
		return
	}

	var req request03
	if err := jsoniter.UnmarshalFromString(payload, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	deviceList, err := deviceService.MqttOnlineDevice()
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, deviceList)
}
