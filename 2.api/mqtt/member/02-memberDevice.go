package memberMqtt

import (
	request "batchLog/0.core/commonResReq/req"
	response "batchLog/0.core/commonResReq/res"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	memberService "batchLog/3.service/member"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request02 struct {
	request.MqttReq
}

func MemberDevices(payload, jwt, clientId, ip string) {
	requestTime := time.Now().UTC()
	errTopic := "errReq/" + clientId

	if payload == "" || payload == "{}" {
		logafa.Error("Payload 為空")
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Payload 為空")
		return
	}
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	var req request02
	if err := jsoniter.UnmarshalFromString(payload, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	deviceIds, err := memberService.MemberDevices(userData.MemberId)
	if err != nil {
		logafa.Error("裝置新增錯誤, error: %+v", err)
		response.ErrorMqtt(errTopic, http.StatusInternalServerError, requestTime, "裝置新增錯誤")
		return
	}
	response.SuccessMqtt(req.SubscribeTo, requestTime, deviceIds)
}
