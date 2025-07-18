package account

import (
	"batchLog/repo"
)

// ✅ 定義 Service 介面
type AccountService interface {
	Login(ip, accountName, password, deviceId string) (map[string]interface{}, error)
	Register(ip, username, password, email, lastName, firstName, nickName string, deviceName string)(map[string]interface{},error)
}

// ✅ 實作 Service
type AccountServiceImpl struct {
    accountRepo 			repo.AccountRepository
	deviceRepo				repo.DeviceRepository
	memberRepo				repo.MemberInfoRepository
	passwordRepo		repo.PasswordRepository
}

func NewAccountService(accountRepo repo.AccountRepository, deviceRepo repo.DeviceRepository, memberRepo	repo.MemberInfoRepository, passwordRepo repo.PasswordRepository) AccountService {
    return &AccountServiceImpl{
		accountRepo: accountRepo,
		deviceRepo: deviceRepo,
		memberRepo: memberRepo,
		passwordRepo: passwordRepo,
	}
}
