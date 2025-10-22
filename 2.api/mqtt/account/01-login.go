package accountMqtt

import (
	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/logafa"
	accountService "batchLog/3.service/account"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request01 struct {
	UserAccount string `json:"userAccount"`
	Password    string `json:"password"`
	SubscribeTo string `json:"subscribeTo"`
}

func Login(payload, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "error/account/login/" + payload

	var req request01
	if err := jsoniter.UnmarshalFromString(payload, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	loginInfo, err := accountService.Login(ip, req.UserAccount, req.Password)
	if err != nil {
		logafa.Error("登入發生錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "登入發生錯誤, "+err.Error())
		return
	}

	response.SuccessMqtt(req.SubscribeTo, requestTime, loginInfo)
}
