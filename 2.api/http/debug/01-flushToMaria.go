package debugHttp

import (
	response "batchLog/0.core/commonResReq/res"
	"batchLog/0.core/global"
	"batchLog/0.cron/persist"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type request01 struct {
	Duration int `json:"duration"`
}

func FlushToMaria(c *gin.Context) {
	requestTime := time.Now().UTC()
	var req request01
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, requestTime, global.COMMON_REQUEST_ERROR)
		return
	}
	persist.FlushTripFmMongoToMaria(req.Duration)
	response.Success(c, requestTime, "")
}
