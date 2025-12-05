package account

import (
	request "batchLog/0.core/commonResReq/req"
	"batchLog/0.core/global"
	accountService "batchLog/3.service/account"
	"net/http"
)

type request02 struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	NickName  string `json:"nickName"`
}

func Register(ctx request.RequestContext) {
	var req request02
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(http.StatusBadRequest, global.COMMON_REQUEST_ERROR)
		return
	}
	ip := ctx.GetClientIP()
	loginInfo, err := accountService.Register(ctx.GetContext(), ip, req.Username, req.Password, req.Email, req.LastName, req.FirstName, req.NickName)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, global.COMMON_SYSTEM_ERROR)
		return
	}
	ctx.Success(loginInfo)
}
