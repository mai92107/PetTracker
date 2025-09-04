package accountApi

import (
	response "batchLog/0.core/commonResponse"
	accountService "batchLog/3.service/account"
	"fmt"
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
}

func Register(c *gin.Context) {
	requestTime := time.Now().UTC()
	var req request02
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest,requestTime,"Json 格式錯誤")
		return
	}
	err := validateRequest(req)
	if err != nil{
		response.Error(c,http.StatusBadRequest,requestTime,"傳入值驗證錯誤: " + err.Error())
		return
	}

	ip := c.ClientIP()
	loginInfo,err := accountService.Register(ip, req.Username,req.Password,req.Email,req.LastName,req.FirstName,req.NickName)
	if err != nil{
		response.Error(c,http.StatusInternalServerError,requestTime,"註冊發生錯誤 : " + err.Error() )
		return 
	}
	response.Success[map[string]interface{}](c,requestTime,loginInfo)
}

func validateRequest(req request02)error{
	if req.Email == ""{
		return fmt.Errorf("電子信箱不可為空")
	}
	if req.NickName == ""{
		return fmt.Errorf("使用者名稱不可為空")
	}
	if req.Username == ""{
		return fmt.Errorf("使用者帳號不可為空")
	}
	if req.Password == ""{
		return fmt.Errorf("使用者密碼不可為空")
	}
	return nil
}