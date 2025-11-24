package memberHttp

import (
	response "batchLog/0.core/commonResReq/res"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	memberService "batchLog/3.service/member"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MemberDevice(c *gin.Context) {
	requestTime := time.Now().UTC()

	jwt := c.GetHeader("jwt")
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		response.Error(c, http.StatusForbidden, requestTime, "身份認證錯誤")
		return
	}

	deviceIds, err := memberService.MemberDevices(userData.MemberId)
	if err != nil {
		logafa.Error("會員裝置取得失敗, error: %+v", err)
		response.Error(c, http.StatusInternalServerError, requestTime, "會員裝置取得失敗")
		return
	}

	response.Success(c, requestTime, deviceIds)
}
