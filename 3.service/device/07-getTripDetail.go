package deviceService

import (
	common "batchLog/0.core/commonFunction"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	service "batchLog/3.service"
	repo "batchLog/4.repo"
	"fmt"
)

func GetTripDetail(member model.Claims, deviceId string, tripUuid string) (map[string]interface{}, error) {
	trip := map[string]interface{}{}

	err := validateTripDetailRequest(deviceId, tripUuid)
	if err != nil {
		return trip, err
	}
	err = service.ValidateDeviceOwner(deviceId, member)
	if err != nil {
		return trip, err
	}

	tx := global.Repository.DB.MariaDb.Reading.Begin()
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("DB tx 啟動失敗")
		}
	}()

	tripData, err := repo.GetTripDetail(tx, tripUuid)
	if err != nil {
		return trip, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logafa.Error("tx 提交失敗")
		return trip, err
	}

	trip = map[string]interface{}{
		"tripUuid":  tripData.DataRef,
		"deviceId":  tripData.DeviceID,
		"distance":  tripData.DistanceKM,
		"duration":  tripData.DurationMinutes,
		"startAt":   common.ToLocalTimeStr(tripData.StartTime),
		"endAt":     common.ToLocalTimeStr(tripData.EndTime),
		"point":     tripData.PointCount,
		"createdAt": common.ToLocalTime(tripData.CreatedAt),
		"updatedAt": common.ToLocalTime(tripData.UpdatedAt),
	}

	return trip, nil
}

func validateTripDetailRequest(deviceId string, tripUuid string) error {
	if deviceId == "" {
		return fmt.Errorf("deviceId 參數錯誤")
	}
	if tripUuid == "" {
		return fmt.Errorf("tripUuid 參數錯誤")
	}
	return nil
}
