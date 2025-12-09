package debug

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	tripService "batchLog/3.service/trip"
	"net/http"
)

type request01 struct {
	Duration int `json:"duration"`
}

func FlushToMaria(ctx request.RequestContext) {
	var req request01
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	userData, err := jwtUtil.GetUserDataFromJwt(ctx.GetJWT())
	if err != nil || userData.Identity != "ADMIN" {
		logafa.Error("身份認證錯誤", "error", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}
	tripService.FlushTripFmMongoToMaria(ctx.GetContext(), req.Duration, userData.GetExecutor())
	ctx.Success("")
}
