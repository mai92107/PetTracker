package deviceHttp

import (
	"net/http"
	"time"

	request "batchLog/0.core/commonResReq/req"
	response "batchLog/0.core/commonResReq/res"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	deviceService "batchLog/3.service/device"

	"github.com/gin-gonic/gin"
)

type request06 struct {
	DeviceId string `json:"deviceId"`
	request.PageInfo
}

func DeviceTrips(c *gin.Context) {
	requestTime := time.Now().UTC()

	var req request06
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, requestTime, global.COMMON_REQUEST_ERROR)
		return
	}
	jwt := c.GetHeader("jwt")
	userInfo, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.Error(c, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}
	datas, total, totalPages, err := deviceService.GetDeviceTrips(userInfo, req.DeviceId, model.NewPageable(&req.Page, &req.Size, req.Direction, req.OrderBy))
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.Error(c, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	pageInfo := response.GetPageResponse(req.PageInfo, total, totalPages)
	info := map[string]interface{}{
		"pageInfo": pageInfo,
		"trips":    datas,
	}

	response.Success(c, requestTime, info)
}
