package deviceService

import (
	common "batchLog/0.core/commonFunction"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	repo "batchLog/4.repo"
	"fmt"
	"slices"
	"time"
)

func MqttDeviceStatus(deviceId string, member model.Claims) (map[string]any, error) {

	err := validateDeviceOwner(deviceId, member)
	if err != nil {
		return nil, err
	}
	isOnline, err := getDeviceOnline(deviceId)
	if err != nil {
		return nil, err
	}

	// 從全癒取得最新資料
	lastSeenFromGlobal, err := getDeviceInfo(deviceId)
	if err != nil {
		return nil, err
	}
	// 若 重啟 全域變數會消失 則從DB找最新紀錄
	lastSeenFromMongo, err := getRecordInfo(deviceId)
	if err != nil {
		return nil, err
	}
	lastSeenDefault := "----/--/-- --:--:--"

	lastSeen := common.Coalesce(lastSeenFromGlobal, lastSeenFromMongo, lastSeenDefault)
	
	return map[string]any{
		"lastSeen": lastSeen,
		"online":   isOnline,
	}, nil
}

func getDeviceOnline(deviceId string) (bool, error) {
	devices, err := repo.GetOnlineDevices()
	if err != nil {
		return false, err
	}
	exist := slices.Contains(devices, deviceId)
	return exist, nil
}

func getDeviceInfo(deviceId string) (string, error) {
	const timeout = 2 * time.Second
	start := time.Now()
	for {
		if global.ActiveDevicesLock.TryLock() {
			break
		}
		if time.Since(start) > timeout {
			// 鎖超時
			return "", fmt.Errorf("警告：MqttOnlineDevice() 嘗試加鎖超過 2 秒，放棄。")
		}
		time.Sleep(10 * time.Millisecond)
	}
	defer global.ActiveDevicesLock.Unlock()

	info, exist := global.ActiveDevices[deviceId]
	if !exist {
		return "", nil
	}
	return info.LastSeen, nil
}

func getRecordInfo(deviceId string) (string, error) {
	info, err := repo.GetLatestDeviceRecordByDeviceId(deviceId)
	if err != nil {
		return "", fmt.Errorf("無法取得裝置定位資訊, error: %+v", err)
	}
	if info == nil {
		return "", nil
	}
	return info.RecordedAt.Format(global.TIME_FORMAT), nil
}

func validateDeviceOwner(deviceId string, member model.Claims) error {
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
