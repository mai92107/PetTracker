package deviceService

import (
	"batchLog/0.core/global"
	repo "batchLog/4.repo"
	"fmt"
)

func Create(deviceType string, memberId int64) (string, error) {

	if err := validateRequest(deviceType); err != nil {
		return "", err
	}

	db := global.Repository.DB.MariaDb.Writing
	// 取得用戶資料
	deviceId, err := repo.CreateDevice(db, deviceType, memberId)
	if err != nil {
		return "", fmt.Errorf("新增使用者裝置發生錯誤，error: %+v", err)
	}
	return deviceId, nil
}

func validateRequest(deviceType string) error {
	if deviceType == "" {
		return fmt.Errorf("裝置名稱不可為空")
	}
	return nil
}
