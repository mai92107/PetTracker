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
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
	DeviceID  string  `json:"deviceId"`
	RecordAt  string  `json:"recordAt"`
	DataRef   string  `json:"dataRef"`
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
		logafa.Error("身份認證錯誤, error: %+v", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}

	deviceService.Recording(ctx.GetContext(), req.Latitude, req.Longitude, claim.MemberId, req.DeviceID, req.RecordAt, req.DataRef)
	ctx.Success("")
}
