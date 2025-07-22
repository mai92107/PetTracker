package deviceService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"fmt"
)


func Create(deviceName string, accountName string)(string,error){

	tx := global.Repository.DB.Writing.Begin()
	defer func(){
		if r := recover();r != nil{
			logafa.Error("裝置新增失敗")
		}
	}()

	// 取得使用者資料
	account,err := repo.FindAccountByAccountName(tx,accountName)
	if err != nil{
		return "",err
	}
	memberInfo,err := repo.FindMemberByAccountUuid(tx,account.Uuid.String())
	if err != nil{
		logafa.Error("查無此會員, error: %+v",err)
		return "",err
	}
	
	deviceId,err := repo.CreateDevice(tx,deviceName,memberInfo.Uuid)
	if err != nil {
		logafa.Error("新增使用者裝置發生錯誤，error: %+v",err)
        return "",err
    }

	if err := tx.Commit().Error ; err != nil{
		tx.Rollback()
		return "",fmt.Errorf(global.COMMON_SYSTEM_ERROR)

	}

	return deviceId,nil
}
