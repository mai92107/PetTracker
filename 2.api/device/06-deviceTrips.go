package device

import (
	"net/http"

	request "batchLog/0.core/commonResReq/req"
	response "batchLog/0.core/commonResReq/res"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	deviceService "batchLog/3.service/device"
)

type request06 struct {
	DeviceId string `json:"deviceId"`
	request.PageInfo
}

func TripList(ctx request.RequestContext) {
	var req request06
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	jwt := ctx.GetJWT()
	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤", "error", err)
		ctx.Error(http.StatusForbidden, "身份認證錯誤")
		return
	}
	datas, total, totalPages, err := deviceService.GetTripList(ctx.GetContext(), userInfo, req.DeviceId, model.NewPageable(&req.Page, &req.Size, req.Direction, req.OrderBy))
	if err != nil {
		logafa.Error("系統發生錯誤", "error", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	pageInfo := response.GetPageResponse(req.PageInfo, total, totalPages)
	info := map[string]interface{}{
		"pageInfo": pageInfo,
		"trips":    datas,
	}

	ctx.Success(info)
}
