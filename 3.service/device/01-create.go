package deviceService

import (
	"batchLog/0.core/global"
	repo "batchLog/4.repo"
	"context"
	"fmt"
)

func Create(ctx context.Context, deviceType string, memberId int64) (string, error) {

	if err := validateCreateRequest(deviceType); err != nil {
		return "", err
	}

	db := global.Repository.DB.MariaDb.Writing
	// 取得用戶資料
	deviceId, err := repo.CreateDevice(ctx, db, deviceType, memberId)
	if err != nil {
		return "", fmt.Errorf("新增使用者裝置發生錯誤, err: %w", err)
	}
	return deviceId, nil
}

func validateCreateRequest(deviceType string) error {
	if deviceType == "" {
		return fmt.Errorf("裝置名稱不可為空")
	}
	return nil
}
