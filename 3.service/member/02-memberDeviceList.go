package memberService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"context"
	"fmt"
)

func MemberDeviceList(ctx context.Context, memberId int64) ([]string, error) {
	tx := global.Repository.DB.MariaDb.Reading.Begin()
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("DB tx 啟動失敗")
		}
	}()
	deviceIds, err := repo.GetDeviceIdsByMemberId(ctx, tx, memberId)
	if err != nil {
		return []string{}, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return []string{}, fmt.Errorf("交易提交失敗, error: %+v", err)
	}
	return deviceIds, nil
}
