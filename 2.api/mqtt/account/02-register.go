package accountMqtt

import (
	response "batchLog/0.core/commonResponse"
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

func Register(payload, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "error/account/register/" + payload

	var req request02
	if err := jsoniter.UnmarshalFromString(payload, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	loginInfo, err := accountService.Register(ip, req.Username, req.Password, req.Email, req.LastName, req.FirstName, req.NickName)
	if err != nil {
		logafa.Error("註冊發生錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "註冊發生錯誤 : "+err.Error())
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, loginInfo)
}
