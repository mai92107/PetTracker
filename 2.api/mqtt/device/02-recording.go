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

type request02 struct {
	Longitude string `json:"lng"`
	Latitude  string `json:"lat"`
	DeviceID  string `json:"deviceId"`
}

func Recording(request, jwt, ip string) {
	requestTime := time.Now().UTC()
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.ErrorMqtt("ERROR/"+jwt, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	topic := "recording/" + fmt.Sprintf("%d", userData.MemberId)
	var req request02
	err = jsoniter.UnmarshalFromString(request, &req)
	if err != nil {
		logafa.Error("Json 格式錯誤, error: %+v", err)
		response.ErrorMqtt(topic, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	time := requestTime.Format("2006-01-02 15:04:05")
	deviceService.Recording(req.Latitude, req.Longitude, req.DeviceID, time)
	response.SuccessMqtt(topic, requestTime, "")
}
