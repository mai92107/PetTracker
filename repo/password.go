package repo

import (
	common "batchLog/core/commonFunction"
	"batchLog/core/global"
	gormTable "batchLog/core/gorm"
	"batchLog/core/logafa"
	"fmt"

	"github.com/google/uuid"
)

// ✅ 定義 Repository 介面
type PasswordRepository interface {
	CreateHistory(accountUuid uuid.UUID, password string)error
}

// ✅ 實作 Repository
type PasswordRepositoryImpl struct {
	DB    *global.DB
	Cache *global.Cache
}

func NewPasswordRepository(db *global.DB, cache *global.Cache) *PasswordRepositoryImpl {
	return &PasswordRepositoryImpl{
		DB: 		db,
		Cache: 		cache,
	}
}

func (r *PasswordRepositoryImpl)CreateHistory(accountUuid uuid.UUID, password string)error{
	hashedPassword,_ := common.BcryptHash(password)
	pastPassword := gormTable.PastPassword{
		AccountUuid: accountUuid,
		Password: hashedPassword,
	}
	err := r.DB.Writing.Table("past_passwords").Create(&pastPassword).Error
	if err != nil {
		logafa.Error("建立使用者歷史密碼失敗, error: %+v",err)
		return fmt.Errorf("建立使用者歷史密碼失敗")
	}
	return nil
}
