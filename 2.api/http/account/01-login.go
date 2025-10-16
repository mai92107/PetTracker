package accountHttp

import (
	response "batchLog/0.core/commonResponse"
	accountService "batchLog/3.service/account"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type request01 struct {
	UserAccount string `json:"userAccount"`
	Password    string `json:"password"`
}

func Login(c *gin.Context) {
	requestTime := time.Now().UTC()
	var req request01
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, requestTime, "Json 格式錯誤")
		return
	}
	// 擷取 IP 與 UA
	ip := c.ClientIP()

	loginInfo, err := accountService.Login(ip, req.UserAccount, req.Password)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, requestTime, "登入發生錯誤, "+err.Error())
		return
	}
	response.Success(c, requestTime, loginInfo)
}
