package service

import (
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"fmt"
	"slices"
)

// 驗證會員使否為管理者 或是 裝置擁有者
func ValidateDeviceOwner(deviceId string, member jwtUtil.Claims) error {
	if member.Identity == "ADMIN" {
		return nil
	}

	tx := global.Repository.DB.MariaDb.Reading.Begin()
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("裝置追蹤失敗")
		}
	}()
	deviceIds, err := repo.GetDeviceIdsByMemberId(tx, member.MemberId)
	if err != nil {
		return err
	}
	if !slices.Contains(deviceIds, deviceId) {
		logafa.Debug("用戶 %v 嘗試讀取裝置 %s 資訊", member.MemberId, deviceId)
		return fmt.Errorf("無權限執行此操作")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("交易提交失敗, error: %+v", err)
	}
	return nil
}
