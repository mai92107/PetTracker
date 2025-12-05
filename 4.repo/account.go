package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	common "batchLog/0.core/commonFunction"
	"batchLog/0.core/global"
	gormTable "batchLog/0.core/gorm"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model/role"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

func FindAccountByAccountName(ctx context.Context, tx *gorm.DB, userAccount string) (*gormTable.Account, error) {

	if strings.Contains(userAccount, "@") {
		return FindAccountByEmail(ctx, tx, userAccount)
	}
	return FindAccountByUsername(ctx, tx, userAccount)
}

func FindAccountByUsername(ctx context.Context, tx *gorm.DB, username string) (*gormTable.Account, error) {
	var account gormTable.Account
	err := tx.WithContext(ctx).First(&account, "username = ?", username).Error
	if err != nil {
		logafa.Error("查詢帳戶發生錯誤, error: %+v", err)
		return nil, fmt.Errorf("查詢帳戶發生錯誤")
	}
	return &account, nil
}

func FindAccountByEmail(ctx context.Context, tx *gorm.DB, email string) (*gormTable.Account, error) {
	var account gormTable.Account
	err := tx.WithContext(ctx).First(&account, "email = ?", email).Error
	if err != nil {
		logafa.Error("查詢帳戶發生錯誤, error: %+v", err)
		return nil, fmt.Errorf("查詢帳戶發生錯誤")
	}
	return &account, nil
}

func CreateAccount(ctx context.Context, tx *gorm.DB, memberId int64, username, password, email string) (uuid.UUID, error) {
	now := time.Now().UTC()
	hashedPassword, _ := common.BcryptHash(password)
	account := gormTable.Account{
		Uuid:          uuid.New(),
		MemberId:      memberId,
		Username:      username,
		Password:      hashedPassword,
		Email:         email,
		Identity:      role.MEMBER.ToString(),
		LastLoginTime: now,
	}
	err := tx.WithContext(ctx).Create(&account).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				return uuid.Nil, fmt.Errorf("使用者帳號 %s 已存在", username)
			}
			if strings.Contains(err.Error(), "email") {
				return uuid.Nil, fmt.Errorf("電子郵件 %s 已存在", email)
			}
		}
		logafa.Error("建立帳戶失敗, error: %+v", err)
		return uuid.Nil, fmt.Errorf("建立帳戶失敗")
	}
	return account.Uuid, nil
}

func UpdateLoginTime(ctx context.Context, tx *gorm.DB, accountUUID uuid.UUID) error {
	now := time.Now().UTC()

	err := tx.WithContext(ctx).Model(&gormTable.Account{}).
		Where("uuid = ?", accountUUID).
		Update("last_login_time", now).Error

	if err != nil {
		logafa.Error("更新登入時間失敗, error: %+v", err)
		return fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}
	return nil
}
