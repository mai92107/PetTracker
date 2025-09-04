package deviceService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	repo "batchLog/4.repo"
	"fmt"
)


func Create(identity, deviceType string, memberId int64)(string,error){

	if identity != model.ADMIN.ToString(){
		return "",fmt.Errorf("無權限新增裝置")
	}

	db := global.Repository.DB.Writing	
	// 取得用戶資料
	deviceId,err := repo.CreateDevice(db,deviceType,memberId)
	if err != nil {
		logafa.Error("新增使用者裝置發生錯誤，error: %+v",err)
        return "",err
    }
	return deviceId,nil
}
