package factory

import (
	"batchLog/core/global"
	"batchLog/repo"
	"batchLog/service/account"
)

type AccountRepositoryImpl struct {
	DB    *global.DB
	Cache *global.Cache
}

func NewAccountServiceFactory(db *global.DB, cache *global.Cache) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{
		DB:    db,
		Cache: cache,
	}
}

func (f *AccountRepositoryImpl) CreateService() account.AccountService {
	accountRepo := repo.NewAccountRepository(f.DB, f.Cache)
	deviceRepo := repo.NewDeviceRepository(f.DB, f.Cache)
	memberRepo := repo.NewMemberInfoRepository(f.DB, f.Cache)
	passwordRepo := repo.NewPasswordRepository(f.DB, f.Cache)

	return account.NewAccountService(accountRepo, deviceRepo, memberRepo, passwordRepo)
}
