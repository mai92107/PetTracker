package deviceService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"fmt"
	"slices"
)

func Recording(lat, lng float64, memberId int64, deviceId, recordTime string) error {
	err := validateRecording(lat, lng, deviceId)
	if err != nil {
		logafa.Error("驗證失敗, error: %+v", err)
		return fmt.Errorf("驗證失敗")
	}
	tx := global.Repository.DB.MariaDb.Reading.Begin()
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("裝置追蹤失敗")
		}
	}()

	deviceIds, err := repo.GetDeviceIdsByMemberId(tx, memberId)
	if err != nil {
		return err
	}

	if len(deviceIds) == 0 || !slices.Contains(deviceIds, deviceId) {
		logafa.Error("使用者: %v, 查無該裝置:%v",memberId, deviceId)
		return fmt.Errorf("使用者查無該裝置")
	}

	err = repo.SaveLocation(lat, lng, deviceId, recordTime)
	if err != nil {
		return err
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
