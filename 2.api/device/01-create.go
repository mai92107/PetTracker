package device

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
	"net/http"
)

type request01 struct {
	DeviceType string `json:"deviceType"`
}

func Create(ctx request.RequestContext) {
	var req request01
	if err := ctx.BindJSON(&req); err != nil {
		logafa.Error("Json 格式錯誤", "error", err)
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	userData, err := jwtUtil.GetUserDataFromJwt(ctx.GetJWT())
	if err != nil || userData.Identity != "ADMIN" {
		logafa.Error("身份認證錯誤", "error", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}
	deviceId, err := deviceService.Create(ctx.GetContext(), req.DeviceType, userData.MemberId)
	if err != nil {
		logafa.Error("裝置新增失敗", "error", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(deviceId)
}
