package repo

import (
	"batchLog/core/global"
	gormTable "batchLog/core/gorm"
	"batchLog/core/logafa"
	"batchLog/core/redis"
	"fmt"

	"github.com/google/uuid"
)

// ✅ 定義 Repository 介面
type DeviceRepository interface {
    FindByDeviceId(deviceId string) (*gormTable.Device, error)
	generateDeviceId()string
	CreateDevice(deviceName string, memberInfo uuid.UUID)(string,error)
}

// ✅ 實作 Repository
type DeviceRepositoryImpl struct {
	DB    *global.DB
	Cache *global.Cache
}

func NewDeviceRepository(db *global.DB, cache *global.Cache) *DeviceRepositoryImpl {
	return &DeviceRepositoryImpl{
		DB: 		db,
		Cache: 		cache,
	}
}

func (r *DeviceRepositoryImpl) FindByDeviceId(deviceId string) (*gormTable.Device, error) {
    var device gormTable.Device
    err := r.DB.Reading.First(&device, "device_id = ?", deviceId).Error
    return &device, err
}

func (r *DeviceRepositoryImpl) generateDeviceId()string{

	prefix := redis.HGetData("device_setting","device_prefix")
	sequence, err := r.Cache.Writing.HIncrBy(r.Cache.CTX, "device_setting", "device_sequence",1).Result()
	if err != nil {
		logafa.Error("failed to increment sequence in Redis: %v", err)
		return ""
	}
	return fmt.Sprintf("%s-%06d", prefix, sequence)
}

func (r *DeviceRepositoryImpl) CreateDevice(deviceName string, memberInfoUuid uuid.UUID)(string,error){
	device := gormTable.Device{
		Uuid: uuid.New(),
		MemberInfoUuid: memberInfoUuid,
		DeviceId: r.generateDeviceId(),
		DeviceName: deviceName,
	}
	err := r.DB.Writing.Table("device").Create(&device).Error
	if err != nil {
		return "",fmt.Errorf("建立使用者裝置資料失敗, error: %+v",err)
	}
	return device.DeviceId,nil
}