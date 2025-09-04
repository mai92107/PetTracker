package jwtUtil

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(accountName, identity string, memberId int64, ip string, currentTime time.Time, expireTime time.Duration)(string,error){
	var loginType model.LoginType
	if strings.Contains(accountName, "@"){
		loginType = model.EMAIL
	}else{
		loginType = model.USERNAME
	}
	// 產生 JWT
	claims := &model.Claims{
		LoginType: loginType.String(),
		AccountName: accountName,
		Identity: identity,
		MemberId: memberId,
		LoginIp: ip,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(expireTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(global.ConfigSetting.JsonSecretKey)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		logafa.Error("產生 JWT 發生錯誤, Error: %v",err)
		return "",err
	}
	return tokenStr,nil
}


func GetUserDataFromJwt(tokenStr string) (model.Claims, error) {
	claims := model.Claims{}
	// Parse token with claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.ConfigSetting.JsonSecretKey), nil
	})
	if err != nil {
		return model.Claims{}, fmt.Errorf("JWT 解析失敗: %w", err)
	}
	if !token.Valid {
		return model.Claims{}, fmt.Errorf("JWT 無效")
	}
	return claims, nil
}