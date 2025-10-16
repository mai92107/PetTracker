package accountMqtt

import (
	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/logafa"
	mqttUtils "batchLog/2.api/mqtt"
	accountService "batchLog/3.service/account"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request01 struct {
	UserAccount string `json:"userAccount"`
	Password    string `json:"password"`
}

func Login(payload, ip string) {
	requestTime := time.Now().UTC()
	request, _ := mqttUtils.Decrypt(payload)

	topic := "account_login/" + payload

	var req request01
	if err := jsoniter.UnmarshalFromString(request, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	loginInfo, err := accountService.Login(ip, req.UserAccount, req.Password)
	if err != nil {
		logafa.Error("登入發生錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusInternalServerError, requestTime, "登入發生錯誤, "+err.Error())
		return
	}

	response.SuccessMqtt(topic, requestTime, loginInfo)
}
