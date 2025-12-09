package jwtUtil

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginType int

const (
	EMAIL LoginType = iota
	USERNAME
)

func (lt LoginType) String() string {
	switch lt {
	case EMAIL:
		return "EMAIL"
	case USERNAME:
		return "USERNAME"
	default:
		return "UNKNOWN"
	}
}

type Claims struct {
	MemberId    int64  `json:"memberId"`
	AccountName string `json:"accountName"`
	LoginType   string `json:"loginType"`
	Identity    string `json:"identity"`
	LoginIp     string `json:"loginIp"`
	jwt.RegisteredClaims
}

func (c Claims) GetExecutor() string {
	return strconv.Itoa(int(c.MemberId))
}

func (c Claims) IsAdmin() bool {
	return c.Identity == "ADMIN"
}

func GenerateJwt(accountName, identity string, memberId int64, ip string, currentTime time.Time, expireTime time.Duration) (string, error) {
	var loginType LoginType
	if strings.Contains(accountName, "@") {
		loginType = EMAIL
	} else {
		loginType = USERNAME
	}
	// 產生 JWT
	claims := &Claims{
		LoginType:   loginType.String(),
		AccountName: accountName,
		Identity:    identity,
		MemberId:    memberId,
		LoginIp:     ip,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(expireTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(global.ConfigSetting.JwtSecretKey)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		logafa.Error("產生 JWT 發生錯誤, Error: %v", err)
		return "", err
	}
	return tokenStr, nil
}

func GetUserDataFromJwt(tokenStr string) (Claims, error) {
	claims := Claims{}
	// Parse token with claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.ConfigSetting.JwtSecretKey), nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("JWT 解析失敗: %w", err)
	}
	if !token.Valid {
		return Claims{}, fmt.Errorf("JWT 無效")
	}
	return claims, nil
}
