package account

import (
	response "batchLog/core/commonResponse"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type request02 struct {
	Username 	string 	`json:"username"`
	Password 	string 	`json:"password"`
	Email		string	`json:"email"`
	LastName	string	`json:"lastName"`
	FirstName	string	`json:"firstName"`
	NickName	string	`json:"nickName"`
	Device		string`json:"device"`
}

func (ac *AccountController) Register(c *gin.Context) {
	requestTime := time.Now().UTC()
	var req request02
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest,requestTime,"Json 格式錯誤")
		return
	}
	ip := c.ClientIP()
	loginInfo,err := ac.accountService.Register(ip, req.Username,req.Password,req.Email,req.LastName,req.FirstName,req.NickName, req.Device)
	if err != nil{
		response.Error(c,http.StatusInternalServerError,requestTime,"註冊發生錯誤 : " + err.Error() )
		return 
	}
	response.Success[map[string]interface{}](c,requestTime,loginInfo)
}