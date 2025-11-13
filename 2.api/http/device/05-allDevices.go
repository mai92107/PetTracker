package deviceHttp

import (
	"net/http"
	"time"

	response "batchLog/0.core/commonResponse"
	"batchLog/0.core/logafa"
	deviceService "batchLog/3.service/device"

	"github.com/gin-gonic/gin"
)

func AllDevice(c *gin.Context) {
	requestTime := time.Now().UTC()

	deviceIds, err := deviceService.AllDevice()
	if err != nil {
		logafa.Error("系統發生錯誤, error: %+v", err)
		response.Error(c, http.StatusInternalServerError, requestTime, "系統發生錯誤, 請稍後嘗試")
		return
	}
	response.Success(c, requestTime, deviceIds)
}
