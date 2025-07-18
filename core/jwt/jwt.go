package jwtUtil

import (
	"batchLog/core/global"
	"batchLog/core/logafa"
	"batchLog/core/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(username, deviceId, ip string, currentTime time.Time, expireTime time.Duration)(string,error){
	// 產生 JWT
	claims := &model.Claims{
		Username: username,
		DeviceID: deviceId,
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

func ValidateJwt(jwt string)bool{
	return true
}

func GetUserDataFromJwt(jwt string)model.Claims{
	return model.Claims{}
}