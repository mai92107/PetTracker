package deviceService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"fmt"
)

func Recording(lat, lng float64, deviceId, recordTime string) error {
	err := validateRecording(lat, lng, deviceId)
	if err != nil {
		return fmt.Errorf("驗證失敗, error: %+v", err)
	}
	tx := global.Repository.DB.MariaDb.Reading.Begin()
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

func validateRecording(lat, lng float64, deviceId string) error {
	if lat == 0 {
		return fmt.Errorf("lat 參數錯誤")
	}
	if lng == 0 {
		return fmt.Errorf("lng 參數錯誤")
	}
	if deviceId == "" {
		return fmt.Errorf("deviceId 參數錯誤")
	}
	return nil
}
