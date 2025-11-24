package accountMqtt

import (
	response "batchLog/0.core/commonResReq/res"
	"batchLog/0.core/logafa"
	accountService "batchLog/3.service/account"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request02 struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	LastName    string `json:"lastName"`
	FirstName   string `json:"firstName"`
	NickName    string `json:"nickName"`
	SubscribeTo string `json:"subscribeTo"`
}

func Register(payload, jwt, clientId, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "errReq/" + clientId

	if payload == "" || payload == "{}" {
		logafa.Error("Payload 為空")
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Payload 為空")
		return
	}
	var req request02
	if err := jsoniter.UnmarshalFromString(payload, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %s", err.Error())
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	loginInfo, err := accountService.Register(ip, req.Username, req.Password, req.Email, req.LastName, req.FirstName, req.NickName)
	if err != nil {
		logafa.Error("註冊發生錯誤, error: %s", err.Error())
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "註冊發生錯誤")
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, loginInfo)
}
