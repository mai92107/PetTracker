package memberService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"context"
	"fmt"
)

func AddDevice(ctx context.Context, memberId int64, deviceId, deviceName string) error {
	tx := global.Repository.DB.MariaDb.Reading.Begin()
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("DB tx 啟動失敗")
		}
	}()
	device, err := repo.FindDeviceByDeviceId(ctx, tx, deviceId)
	if err != nil {
		return err
	}
	writingDB := global.Repository.DB.MariaDb.Writing
	err = repo.AddDevice(ctx, writingDB, memberId, device.DeviceId, deviceName)
	if err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("交易提交失敗, error: %+v", err)
	}
	return nil
}
