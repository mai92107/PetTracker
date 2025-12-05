package device

import (
	"net/http"

	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
)

type request04 struct {
	DeviceId string `json:"deviceId"`
}

func DeviceStatus(ctx request.RequestContext) {
	var req request04
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	userInfo, err := jwtUtil.GetUserDataFromJwt(ctx.GetJWT())
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}
	info, err := deviceService.MqttDeviceStatus(ctx.GetContext(), req.DeviceId, userInfo)
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(info)
}
