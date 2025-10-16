package homeMqtt

import (
	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	"fmt"
	"net/http"
	"time"
)

func SayHello(jwt string) {
	requestedTime := time.Now().UTC()

	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt("ERROR/"+jwt, http.StatusForbidden, requestedTime, "身份認證錯誤")
		return
	}
	topic := "hello/" + fmt.Sprintf("%d", userInfo.MemberId)
	data := "hello"

	response.SuccessMqtt(topic, requestedTime, data)
}
