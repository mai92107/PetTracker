package homeMqtt

import (
	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/logafa"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request01 struct {
	SubscribeTo string `json:"subscribeTo"`
}

func SayHello(payload string) {
	requestedTime := time.Now().UTC()
	errTopic := "error/home/sayHello/" + payload
	var req request01
	if err := jsoniter.UnmarshalFromString(payload, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestedTime, "Json 格式錯誤")
		return
	}
	data := "hello my world"
	response.SuccessMqtt(req.SubscribeTo, requestedTime, data)
}
