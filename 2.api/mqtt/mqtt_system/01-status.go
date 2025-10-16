package systemMqtt

import (
	"fmt"
	"net/http"
	"time"

	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	systemService "batchLog/3.service/system"
)

func SystemStatus(request, jwt, ip string) {
	requestTime := time.Now().UTC()
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt("ERROR/"+jwt, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	topic := "systemStatus/" + fmt.Sprintf("%d", userData.MemberId)
	data, err := systemService.SystemStatus()
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	response.SuccessMqtt(topic, requestTime, data)
}
