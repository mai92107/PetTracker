package device

import (
	"net/http"

	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"
)

func DeviceList(ctx request.RequestContext) {
	deviceIds, err := deviceService.DeviceList(ctx.GetContext())
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(deviceIds)
}
