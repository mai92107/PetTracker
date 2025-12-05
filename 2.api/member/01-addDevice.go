package member

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	memberService "batchLog/3.service/member"
	"net/http"
)

type request01 struct {
	DeviceId   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
}

func AddDevice(ctx request.RequestContext) {
	var req request01
	if err := ctx.BindJSON(&req); err != nil {
		// logafa.Error("Json 格式錯誤, error: %+v", err)
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	jwt := ctx.GetJWT()
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		// logafa.Error("身份認證錯誤, error: %+v", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}

	err = memberService.AddDevice(ctx.GetContext(), userData.MemberId, req.DeviceId, req.DeviceName)
	if err != nil {
		// logafa.Error("裝置新增錯誤, error: %+v", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success("")
}
