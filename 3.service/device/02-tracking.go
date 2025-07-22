package deviceService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"fmt"
)

func Tracking(lat, lng, deviceId, accountName, recordTime string)error{

	tx := global.Repository.DB.Reading.Begin()
	defer func(){
		if r := recover();r != nil{
			logafa.Error("裝置追蹤失敗")
		}
	}()

	account,err := repo.FindAccountByAccountName(tx,accountName)
	if err != nil{
		logafa.Error("查無此帳戶, error: %+v",err)
		return fmt.Errorf("查無此帳戶")
	}
	memberInfo,err := repo.FindMemberByAccountUuid(tx,account.Uuid.String())
	if err != nil{
		logafa.Error("查無此會員, error: %+v",err)
		return fmt.Errorf("查無此會員")
	}
	err = repo.SaveLocation(lat,lng,deviceId,memberInfo.NickName,recordTime)
	if err != nil{
		logafa.Error("裝置定位儲存失敗, error: %+v",err)
		return fmt.Errorf("裝置定位儲存失敗")
	}
	if err := tx.Commit().Error ; err != nil{
		tx.Rollback()
		return fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}
	return nil
}