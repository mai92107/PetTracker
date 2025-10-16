package deviceMqtt

import (
	response "batchLog/0.core/commonResponse"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
	"fmt"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type request01 struct {
	DeviceType string `json:"deviceType"`
}

func Create(request, jwt, ip string) {
	requestTime := time.Now().UTC()
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt("ERROR/"+jwt, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	topic := "create/" + fmt.Sprintf("%d", userData.MemberId)
	var req request01
	if err := jsoniter.UnmarshalFromString(request, &req); err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}

	deviceId, err := deviceService.Create(userData.Identity, req.DeviceType, userData.MemberId)
	if err != nil {
		logafa.Error("裝置新增失敗, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusInternalServerError, requestTime, "裝置新增失敗，請稍後嘗試")
		return
	}
	response.SuccessMqtt(topic, requestTime, deviceId)
}
