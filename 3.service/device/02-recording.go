package deviceService

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	repo "batchLog/4.repo"
	"fmt"
	"slices"
	"time"
)

func Recording(lat, lng float64, memberId int64, deviceId, recordTime, dataRef string) error {

	loc,err := time.LoadLocation("Local")
	if err != nil {
		logafa.Error("載入當前地區失敗, error: %+v", err)
		return fmt.Errorf("載入當前地區失敗")
	}

	recordLocalTime,err := time.ParseInLocation(global.TIME_FORMAT, recordTime, loc)
	if err != nil {
		logafa.Error("時區轉換失敗, error: %+v", err)
		return fmt.Errorf("時區轉換失敗")
	}

	recordUtcTime := recordLocalTime.UTC()

	err = validateRecordingRequest(lat, lng, deviceId)
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
		logafa.Error("使用者: %v, 查無該裝置:%v", memberId, deviceId)
		return fmt.Errorf("使用者查無該裝置")
	}

	err = repo.SaveLocation(lat, lng, deviceId, recordUtcTime, dataRef)
	if err != nil {
		return err
	}
	updateGlobalDeviceInfo(deviceId, recordTime)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("交易提交失敗, error: %+v", err)
	}
	return nil
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
