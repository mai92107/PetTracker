package accountMqtt

import (
	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	accountService "batchLog/3.service/account"
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request02 struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	NickName  string `json:"nickName"`
}

func Register(request, jwt, ip string) {
	requestTime := time.Now().UTC()
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt("ERROR/"+jwt, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	topic := "register/" + fmt.Sprintf("%d", userData.MemberId)

	var req request02
	if err := jsoniter.UnmarshalFromString(request, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	loginInfo, err := accountService.Register(ip, req.Username, req.Password, req.Email, req.LastName, req.FirstName, req.NickName)
	if err != nil {
		logafa.Error("註冊發生錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusInternalServerError, requestTime, "註冊發生錯誤 : "+err.Error())
		return
	}

	response.SuccessMqtt(topic, requestTime, loginInfo)
}
