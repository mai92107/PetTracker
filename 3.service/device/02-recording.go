package deviceService

import (
	"batchLog/0.core/global"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	"batchLog/0.cron/persist"
	repo "batchLog/4.repo"
	"context"
	"fmt"
	"slices"
	"time"
)

func Recording(ctx context.Context, lat, lng float64, member jwtUtil.Claims, deviceId, recordTime, dataRef string, isEnd bool) (map[string]interface{}, error) {

	if isEnd {
		persist.SaveGpsFmRdsToMongo(ctx)
		persist.FlushTripFmMongoToMaria(ctx, 5)
		return GetTripDetail(ctx, member, deviceId, dataRef)
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		logafa.Error("載入當前地區失敗", "error", err)
		return nil, fmt.Errorf("載入當前地區失敗")
	}

	recordLocalTime, err := time.ParseInLocation(global.TIME_FORMAT, recordTime, loc)
	if err != nil {
		logafa.Error("時區轉換失敗", "error", err)
		return nil, fmt.Errorf("時區轉換失敗")
	}

	recordUtcTime := recordLocalTime.UTC()

	err = validateRecordingRequest(lat, lng, deviceId)
	if err != nil {
		logafa.Error("驗證失敗", "error", err)
		return nil, fmt.Errorf("驗證失敗")
	}
	tx := global.Repository.DB.MariaDb.Reading.Begin()
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("裝置追蹤失敗")
		}
	}()

	deviceIds, err := repo.GetDeviceIdsByMemberId(ctx, tx, member.MemberId)
	if err != nil {
		return nil, err
	}

	if len(deviceIds) == 0 || !slices.Contains(deviceIds, deviceId) {
		logafa.Error("查無該裝置", "user", member.MemberId, "device", deviceId)
		return nil, fmt.Errorf("使用者查無該裝置")
	}

	err = repo.SaveLocation(ctx, lat, lng, deviceId, recordUtcTime, dataRef)
	if err != nil {
		return nil, err
	}
	updateGlobalDeviceInfo(deviceId, recordTime)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logafa.Error("交易提交失敗", "error", err)
		return nil, fmt.Errorf("交易提交失敗")
	}
	return nil, nil
}

func updateGlobalDeviceInfo(deviceId string, now string) {
	global.ActiveDevicesLock.Lock()
	global.ActiveDevices[deviceId] = model.DeviceStatus{
		LastSeen: now,
	}
	global.ActiveDevicesLock.Unlock()
}

func validateRecordingRequest(lat, lng float64, deviceId string) error {
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
