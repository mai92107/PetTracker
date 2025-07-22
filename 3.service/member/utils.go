package memberService

import (
	"batchLog/0.core/global"
	gormTable "batchLog/0.core/gorm"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"fmt"
)

func FindByJwt(jwt string)(*gormTable.MemberInfo,error){

	tx := global.Repository.DB.Reading.Begin()
	defer func(){
		if r := recover();r != nil{
			logafa.Error("裝置追蹤失敗")
		}
	}()

	var userAccount *gormTable.Account
	var err error
	// 解讀 JWT
	userData,err := jwtUtil.GetUserDataFromJwt(jwt)
	if err != nil{
		logafa.Error("身份認證錯誤, error: %+v",err)
		return nil, fmt.Errorf("身份認證錯誤")
	}

	userAccount,err = repo.FindAccountByAccountName(tx,userData.AccountName)

    if err != nil {
		logafa.Error("查詢使用者資料發生錯誤，error: %+v",err)
        return nil,fmt.Errorf(global.COMMON_SYSTEM_ERROR)
    }
	memberInfo,err := repo.FindMemberByAccountUuid(tx,userAccount.Uuid.String())
	if err != nil {
		logafa.Error("查詢使用者資料發生錯誤，error: %+v",err)
        return nil,fmt.Errorf(global.COMMON_SYSTEM_ERROR)
    }
	if err := tx.Commit().Error ; err != nil{
		tx.Rollback()
		return nil,fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}
	return memberInfo,nil
}