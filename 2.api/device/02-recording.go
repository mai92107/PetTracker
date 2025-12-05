package device

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
	"net/http"
)

type request02 struct {
	Longitude   float64 `json:"lng"`
	Latitude    float64 `json:"lat"`
	DeviceID    string  `json:"deviceId"`
	RecordAt    string  `json:"recordAt"`
	DataRef     string  `json:"dataRef"`
	SubscribeTo string  `json:"subscribeTo"`
}

func Recording(ctx request.RequestContext) {
	var req request02
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	claim, err := jwtUtil.GetUserDataFromJwt(ctx.GetJWT())
	if err != nil {
		logafa.Error("身份認證錯誤", "error", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}
	final := false
	if req.SubscribeTo != "" {
		final = true
	}

	info, err := deviceService.Recording(ctx.GetContext(), req.Latitude, req.Longitude, claim, req.DeviceID, req.RecordAt, req.DataRef, final)
	if err != nil {
		logafa.Error("系統發生錯誤", "error", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(info)
}
