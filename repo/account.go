package repo

import (
	common "batchLog/core/commonFunction"
	"batchLog/core/global"
	gormTable "batchLog/core/gorm"
	"batchLog/core/logafa"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ✅ 定義 Repository 介面
type AccountRepository interface {
    FindByUsername(username string) (*gormTable.Account, error)
    FindByEmail(email string) (*gormTable.Account, error)
	CreateAccount(username, password, email string)(uuid.UUID,error)
}

// ✅ 實作 Repository
type AccountRepositoryImpl struct {
	DB    *global.DB
	Cache *global.Cache
}

func NewAccountRepository(db *global.DB, cache *global.Cache) AccountRepository {
	return &AccountRepositoryImpl{
		DB:    db,
		Cache: cache,
	}}


func (r *AccountRepositoryImpl) FindByUsername(username string) (*gormTable.Account, error) {
    var account gormTable.Account
    err := r.DB.Reading.First(&account, "username = ?", username).Error
    return &account, err
}

func (r *AccountRepositoryImpl) FindByEmail(email string) (*gormTable.Account, error) {
    var account gormTable.Account
    err := r.DB.Reading.First(&account, "emil = ?", email).Error
    return &account, err
}

func (r *AccountRepositoryImpl)	CreateAccount(username, password, email string)(uuid.UUID,error) {
	now := time.Now().UTC()
	hashedPassword,_ := common.BcryptHash(password)
	account := gormTable.Account{
		Uuid:          uuid.New(),
		Username:      username,
		Password:      hashedPassword,
		Email:         email,
		LastLoginTime: &now,
	}
	err := r.DB.Writing.Table("accounts").Create(&account).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
            if strings.Contains(err.Error(), "username") {
                return uuid.Nil, fmt.Errorf("用戶名 %s 已存在", username)
            }
            if strings.Contains(err.Error(), "email") {
                return uuid.Nil, fmt.Errorf("電子郵件 %s 已存在", email)
            }
        }
		logafa.Error("建立帳戶失敗, error: %+v",err)
		return uuid.Nil,fmt.Errorf("建立帳戶失敗")
	}
	return account.Uuid,nil
}
