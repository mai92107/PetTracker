package request

import (
	"context"
	"time"
)

type RequestContext interface {

	// 提供context
	GetContext() context.Context
	SetContext(context.Context)
	Cancel()
	SetCancel(context.CancelFunc)

	// 取得用戶ID
	GetClientID() string

	// 取得用戶IP
	GetClientIP() string

	// 取得JWT
	GetJWT() string

	// 取得請求時間
	GetRequestTime() time.Time

	// 綁定Json到struct
	BindJSON(interface{}) error

	// 回覆成功
	Success(data interface{})

	// 回覆錯誤
	Error(code int, message string)
}
