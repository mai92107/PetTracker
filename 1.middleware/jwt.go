// middleware/jwt.go
package middleware

import (
	"net/http"
	"strings"
	"time"

	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"

	"github.com/gin-gonic/gin"
)

func JWTValidator(role string) gin.HandlerFunc {

	const ADMIN = "ADMIN"

	return func(c *gin.Context) {
		// 1. 取得 Authorization header
		authHeader := c.GetHeader("jwt")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		// 2. 檢查 Bearer
		if strings.HasPrefix(authHeader, "bearer") {
			authHeader = strings.Split(authHeader, " ")[1]
		}

		// 3. 解析 JWT
		claims, err := jwtUtil.GetUserDataFromJwt(authHeader)
		if err != nil {
			logafa.Warn("JWT 解析失敗: %v | token: %s", err, maskToken(authHeader))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":  "invalid token",
				"detail": err.Error(),
			})
			return
		}

		if !claims.IsAdmin() && role == ADMIN{
			logafa.Warn("用戶 %v 無權限執行此 %s 操作", claims.MemberId, c.Request.URL)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":  "無權限執行此操作",
				"detail": nil,
			})
			return
		}

		// 4. 驗證 exp
		now := time.Now().UTC()
		if !claims.ExpiresAt.After(now) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			return
		}

		// 6. 繼續
		c.Next()
	}
}

// 遮蔽 token（安全日誌）
func maskToken(token string) string {
	if len(token) < 10 {
		return "****"
	}
	return token[:6] + "..." + token[len(token)-4:]
}
