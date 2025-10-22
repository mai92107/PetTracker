package deviceService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"fmt"
)

func Recording(lat, lng, deviceId, recordTime string) error {
	tx := global.Repository.DB.Reading.Begin()
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("裝置追蹤失敗")
		}
	}()

	memberInfo, err := repo.FindMemberByDeviceId(tx, deviceId)
	if err != nil {
		return fmt.Errorf("查無此會員, error: %+v", err)
	}
	err = repo.SaveLocation(lat, lng, deviceId, memberInfo.NickName, recordTime)
	if err != nil {
		return fmt.Errorf("裝置定位儲存失敗, error: %+v", err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("交易提交失敗, error: %+v", err)
	}
	return nil
}
