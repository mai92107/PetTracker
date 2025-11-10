package model

import "github.com/golang-jwt/jwt/v5"

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
	MemberId		int64		`json:"memberId"`
	AccountName 	string 		`json:"accountName"`
	LoginType		string		`json:"loginType"`
	Identity		string		`json:"identity"`
	LoginIp			string		`json:"loginIp"`
	jwt.RegisteredClaims
}

func (c Claims)IsAdmin()bool{
	return c.Identity == "ADMIN"
}