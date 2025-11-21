package repo

import (
	"batchLog/0.core/global"
	gormTable "batchLog/0.core/gorm"
	"batchLog/0.core/logafa"
	"batchLog/0.core/redis"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

func GetAllDeviceIds(tx *gorm.DB) ([]string, error) {
	var deviceIds []string
	err := tx.Model(&gormTable.Device{}).
		Pluck("device_id", &deviceIds).Error
	if err != nil {
		logafa.Error("查詢所有 deviceIds 失敗, error: %+v", err)
		return nil, fmt.Errorf("裝置ID查詢失敗")
	}
	return deviceIds, nil
}

func GetDeviceIdsByMemberId(tx *gorm.DB, memberId int64) ([]string, error) {
	var deviceIds []string

	err := tx.Model(&gormTable.MemberDevice{}).
		Where("member_id = ?", memberId).
		Pluck("device_id", &deviceIds).Error

	if err != nil {
		logafa.Error("查詢會員 deviceId 失敗, memberId: %v, error: %+v", memberId, err)
		return nil, fmt.Errorf("查無此會員或查詢失敗")
	}

	if len(deviceIds) == 0 {
		logafa.Warn("會員存在但無綁定 device, memberId: %v", memberId)
		return []string{}, nil
	}

	return deviceIds, nil
}

func FindDeviceByDeviceId(tx *gorm.DB, deviceId string) (*gormTable.Device, error) {
	var device gormTable.Device
	err := tx.First(&device, "device_id = ?", deviceId).Error
	if err != nil {
		logafa.Error("查無此裝置, error: %+v", err)
		return nil, fmt.Errorf("查無此裝置")
	}
	return &device, nil
}

func CreateDevice(tx *gorm.DB, deviceType string, memberId int64) (string, error) {
	device := gormTable.Device{
		Uuid:           uuid.New(),
		DeviceId:       generateDeviceId(),
		DeviceType:     deviceType,
		CreateByMember: memberId,
	}
	err := tx.Table("device").Create(&device).Error
	if err != nil {
		logafa.Error("建立裝置資料失敗, error: %+v", err)
		return "", fmt.Errorf("建立裝置資料失敗")
	}
	return device.DeviceId, nil
}

func generateDeviceId() string {

	prefix := redis.HGetData("device_setting", "device_prefix")
	sequence, err := global.Repository.Cache.Writing.HIncrBy(global.Repository.Cache.CTX, "device_setting", "device_sequence", 1).Result()
	if err != nil {
		logafa.Error("failed to increment sequence in Redis: %v", err)
		return ""
	}
	return fmt.Sprintf("%s-%06d", prefix, sequence)
}

func SaveLocation(lat, lng float64, deviceId string, recordTime time.Time, dataRef string) error {
	now := time.Now().UTC()
	// 存入 redis 臨時保存
	key := fmt.Sprintf("device:%s", deviceId)
	score := float64(now.UnixMilli())
	gps := gormTable.GPS{
		DeviceId:   deviceId,
		Latitude:   lat,
		Longitude:  lng,
		RecordTime: recordTime,
		DataRef:    dataRef,
	}
	byteData, err := jsoniter.Marshal(gps)
	if err != nil {
		logafa.Error("Json Marshal 失敗, error: %+v", err)
		return fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}
	// 存入 redis
	err = redis.ZAddData(key, score, byteData)
	if err != nil {
		logafa.Error("redis 儲存失敗, error: %+v", err)
		return fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}
	return nil
}

func GetOnlineDevices() ([]string, error) {
	keys, err := redis.KeyScan("device:*")
	if err != nil {
		logafa.Error("redis 掃描 device:* 失敗: %v", err)
		return nil, fmt.Errorf("%s: redis scan error", global.COMMON_SYSTEM_ERROR)
	}

	deviceIds := make([]string, 0, len(keys))
	for _, key := range keys {
		if !strings.HasPrefix(key, "device:") {
			continue // 防呆
		}
		parts := strings.SplitN(key, ":", 2) // 只切一次
		if len(parts) == 2 {
			deviceIds = append(deviceIds, parts[1])
		}
	}

	logafa.Info("目前在線裝置數量: %d", len(deviceIds))
	return deviceIds, nil
}
