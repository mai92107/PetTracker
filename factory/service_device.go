package factory

import (
	"batchLog/core/global"
)

type DeviceServiceFactory struct {
	DBReading    *global.DB
	DBWriting    *global.DB
	CacheReading *global.Cache
	CacheWriting *global.Cache
}

func NewDeviceServiceFactory(dbReading, dbWriting *global.DB, cacheReading, cacheWriting *global.Cache) *DeviceServiceFactory {
	return &DeviceServiceFactory{
		DBReading:    dbReading,
		DBWriting:    dbWriting,
		CacheReading: cacheReading,
		CacheWriting: cacheWriting,
	}
}

// func (f *DeviceServiceFactory) CreateService() device. {
// 	accountRepo := factory.NewAccountRepository(f.DBReading, f.DBWriting, f.CacheReading, f.CacheWriting)
// 	deviceRepo := factory.NewDeviceRepository(f.DBReading, f.DBWriting, f.CacheReading, f.CacheWriting)
// 	memberRepo := factory.NewMemberRepository(f.DBReading, f.DBWriting, f.CacheReading, f.CacheWriting)

// 	return service.NewAccountService(accountRepo, deviceRepo, memberRepo)
// }