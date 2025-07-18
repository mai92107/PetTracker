package factory

import (
	"batchLog/core/global"
)

type MemberServiceFactory struct {
	DBReading    *global.DB
	DBWriting    *global.DB
	CacheReading *global.Cache
	CacheWriting *global.Cache
}

func NewMemberServiceFactory(dbReading, dbWriting *global.DB, cacheReading, cacheWriting *global.Cache) *MemberServiceFactory {
	return &MemberServiceFactory{
		DBReading:    dbReading,
		DBWriting:    dbWriting,
		CacheReading: cacheReading,
		CacheWriting: cacheWriting,
	}
}

// func (f *AccountServiceFactory) CreateService() service.AccountService {
// 	accountRepo := factory.NewAccountRepository(f.DBReading, f.DBWriting, f.CacheReading, f.CacheWriting)
// 	deviceRepo := factory.NewDeviceRepository(f.DBReading, f.DBWriting, f.CacheReading, f.CacheWriting)
// 	memberRepo := factory.NewMemberRepository(f.DBReading, f.DBWriting, f.CacheReading, f.CacheWriting)

// 	return service.NewAccountService(accountRepo, deviceRepo, memberRepo)
// }