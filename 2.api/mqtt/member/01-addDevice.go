package memberMqtt

import (
	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	memberService "batchLog/3.service/member"
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request01 struct {
	DeviceId   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
}

func AddDevice(request, jwt, ip string) {
	requestTime := time.Now().UTC()
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt("ERROR/"+jwt, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	topic := "addDevice/" + fmt.Sprintf("%d", userData.MemberId)
	var req request01
	if err := jsoniter.UnmarshalFromString(request, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	err = memberService.AddDevice(userData.MemberId, req.DeviceId, req.DeviceName)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusInternalServerError, requestTime, "裝置新增錯誤")
		return
	}
	response.SuccessMqtt(topic, requestTime, "")
}
