package model

import "github.com/golang-jwt/jwt/v5"



type Claims struct {
	Username 	string 	`json:"username"`
	DeviceID 	string 	`json:"deviceId"`
	LoginIp		string	`json:"loginIp"`
	jwt.RegisteredClaims
}