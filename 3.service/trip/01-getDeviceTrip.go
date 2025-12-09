package tripService

import (
	common "batchLog/0.core/commonFunction"
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	service "batchLog/3.service"
	repo "batchLog/4.repo"
	"context"
	"fmt"
)

func GetTripList(ctx context.Context, member jwtUtil.Claims, deviceId string, pageable model.Pageable) ([]map[string]interface{}, int64, int64, error) {
	trips := []map[string]interface{}{}
	var total int64
	var totalPages int64

	err := validateTripsRequest(deviceId)
	if err != nil {
		return trips, total, totalPages, err
	}

	err = service.ValidateDeviceOwner(ctx, deviceId, member)
	if err != nil {
		return trips, total, totalPages, err
	}

	tx := global.Repository.DB.MariaDb.Reading.Begin()
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("DB tx 啟動失敗")
		}
	}()

	tripsData, total, totalPages, err := repo.GetTripList(ctx, tx, deviceId, pageable)
	if err != nil {
		return trips, total, totalPages, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logafa.Error("tx 提交失敗")
		return trips, total, totalPages, err
	}

	for _, v := range tripsData {
		trips = append(trips, map[string]interface{}{
			"uuid":     v.DataRef,
			"time":     common.ToLocalTimeShortStr(v.EndTime),
			"duration": common.FormatDigits(v.DurationMinutes, 4),
			"distance": common.FormatDigits(v.DistanceKM, 4),
		})
	}

	return trips, total, totalPages, nil
}

func validateTripsRequest(deviceId string) error {
	if deviceId == "" {
		return fmt.Errorf("deviceId 參數錯誤")
	}
	return nil
}
