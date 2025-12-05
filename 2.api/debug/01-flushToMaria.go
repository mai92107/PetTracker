package debug

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	"batchLog/0.cron/persist"
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
	persist.FlushTripFmMongoToMaria(ctx.GetContext(), req.Duration)
	ctx.Success("")
}
