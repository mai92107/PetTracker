package repo

import (
	"batchLog/0.core/global"
	gormTable "batchLog/0.core/gorm"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	"fmt"
)

func FindByJwt(jwt string) (*gormTable.Member, error) {

	db := global.Repository.DB.Reading
	// 解讀 JWT
	userData, err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil {
		logafa.Error("身份認證錯誤, error: %+v", err)
		return nil, fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}

	// 用 join 一次查出 Member
	var member gormTable.Member
	query := db.Table("member").
		Select("member.*").
		Joins("JOIN account ON account.uuid = member.account_uuid")

	switch userData.LoginType{
	case model.EMAIL.String():
		query = query.Where("account.email = ?", userData.AccountName)
	case model.USERNAME.String():
		query = query.Where("account.username = ?", userData.AccountName)
	}
	err = query.Take(&member).Error
	if err != nil {
		logafa.Error("查詢使用者資料發生錯誤，error: %+v", err)
		return nil, fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}
	return &member, nil
}