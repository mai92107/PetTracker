package model

import "github.com/golang-jwt/jwt/v5"



type Claims struct {
	AccountName 	string 	`json:"accountName"`
	DeviceID 		string 	`json:"deviceId"`
	LoginIp			string	`json:"loginIp"`
	jwt.RegisteredClaims
}