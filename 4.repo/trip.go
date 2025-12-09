package repo

import (
	gormTable "batchLog/0.core/gorm"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	"context"
	"fmt"
	"math"

	"gorm.io/gorm"
)

func GetTripList(ctx context.Context, tx *gorm.DB, deviceId string, pageable model.Pageable) ([]gormTable.TripSummary, int64, int64, error) {
	var deviceTrips []gormTable.TripSummary
	var totalCount int64
	var totalPage int64

	// 查總筆數
	if err := tx.WithContext(ctx).Model(&gormTable.TripSummary{}).
		Where("device_id = ?", deviceId).
		Count(&totalCount).Error; err != nil {
		logafa.Error("統計裝置行程數量失敗", "deviceId", deviceId, "error", err)
		return deviceTrips, totalCount, totalPage, fmt.Errorf("統計行程數量失敗")
	}

	// 如果總筆數為 0，直接回傳空陣列
	if totalCount == 0 {
		logafa.Info("裝置 %s 無任何行程紀錄", deviceId)
		return deviceTrips, totalCount, totalPage, nil
	}

	totalPage = int64(math.Ceil(float64(totalCount) / float64(pageable.Size)))

	// 2. 正式查詢資料（分頁 + 排序）
	err := tx.WithContext(ctx).Where("device_id = ?", deviceId).
		Offset(pageable.Offset()).    // 分頁
		Limit(pageable.Limit()).      // 每頁筆數
		Order(pageable.OrderBySQL()). // 排序
		Find(&deviceTrips).Error

	if err != nil {
		logafa.Error("查詢裝置行程失敗", "deviceId", deviceId, "error", err)
		return deviceTrips, totalCount, totalPage, fmt.Errorf("查詢行程失敗")
	}

	return deviceTrips, totalCount, totalPage, nil
}

func GetTripDetail(ctx context.Context, tx *gorm.DB, tripUuid string) (gormTable.TripSummary, error) {
	var trip gormTable.TripSummary

	err := tx.WithContext(ctx).Where("data_ref = ?", tripUuid).
		First(&trip).Error
	if err != nil {
		logafa.Error("查詢裝置行程失敗", "data_ref", tripUuid, "error", err)
		return trip, fmt.Errorf("查詢行程失敗")
	}

	return trip, nil
}
