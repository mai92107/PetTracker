// middleware/jwt.go
package middleware

import (
	"net/http"
	"strings"
	"time"

	request "batchLog/0.core/commonResReq/req"
	response "batchLog/0.core/commonResReq/res"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model/role"

	"github.com/gin-gonic/gin"
)

func HttpJWTMiddleware(identity role.MemberIdentity) gin.HandlerFunc {
	return func(c *gin.Context) {
		if identity != role.GUEST {
			now := global.GetNow()
			authHeader := c.GetHeader("jwt")
			validateJwt(identity, authHeader, func(code int, msg string) {
				response.Error(c, code, now, msg)
			})
		}
		c.Next()
	}
}

func MqttJWTMiddleware(identity role.MemberIdentity) func(request.RequestContext, func(ctx request.RequestContext)) {
	return func(ctx request.RequestContext, next func(ctx request.RequestContext)) {
		if identity != role.GUEST {
			validateJwt(identity, ctx.GetJWT(), ctx.Error)
		}
		next(ctx)
	}
}

func validateJwt(identity role.MemberIdentity, authHeader string, sendError func(int, string)) {
	if authHeader == "" {
		sendError(http.StatusForbidden, "無法取得JWT")
		return
	}

	if strings.HasPrefix(authHeader, "bearer") {
		authHeader = strings.Split(authHeader, " ")[1]
	}

	claims, err := jwtUtil.GetUserDataFromJwt(authHeader)
	if err != nil {
		logafa.Warn("JWT 解析失敗", "error", err, "jwt", maskToken(authHeader))
		sendError(http.StatusUnauthorized, "invalid token")
		return
	}

	if !claims.IsAdmin() && identity == role.ADMIN {
		logafa.Warn("無權限執行此操作", "user", claims.MemberId)
		sendError(http.StatusUnauthorized, "forbidden")
		return
	}

	now := time.Now().UTC()
	if !claims.ExpiresAt.After(now) {
		sendError(http.StatusUnauthorized, "token expired")
		return
	}
}

// 遮蔽 token（安全日誌）
func maskToken(token string) string {
	if len(token) < 10 {
		return "****"
	}
	return token[:6] + "..." + token[len(token)-4:]
}
