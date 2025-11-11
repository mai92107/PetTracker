package accountService

import (
	common "batchLog/0.core/commonFunction"
	"batchLog/0.core/global"
	gormTable "batchLog/0.core/gorm"
	jwtUtil "batchLog/0.core/jwt"
	"batchLog/0.core/logafa"
	repo "batchLog/4.repo"
	"fmt"
	"time"
)

func Login(ip, accountName, password string) (map[string]interface{}, error) {

	// 參數驗證
	if err := validateLogin(accountName, password); err != nil {
		return nil, err
	}

	tx := global.Repository.DB.MariaDb.Reading.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logafa.Error("登入失敗, r: %+v", r)
		}
	}()

	var userAccount *gormTable.Account
	var err error
	data := map[string]interface{}{}
	// 驗證帳號
	userAccount, err = repo.FindAccountByAccountName(tx, accountName)
	if err != nil {
		tx.Rollback()
		return data, err
	}

	// 驗證密碼
	if !common.BcryptCompare(userAccount.Password, password) {
		tx.Rollback()

		return data, fmt.Errorf("密碼錯誤")
	}

	err = repo.UpdateLoginTime(tx, userAccount.Uuid)
	if err != nil {
		tx.Rollback()
		return data, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logafa.Error("執行交易發生錯誤, error: %+v", err)
		return data, fmt.Errorf(global.COMMON_SYSTEM_ERROR)
	}

	now := time.Now().UTC()
	expireTime := 24 * time.Hour
	token, err := jwtUtil.GenerateJwt(accountName, userAccount.Identity, userAccount.MemberId, ip, now, expireTime)
	if err != nil {
		return data, err
	}
	data = map[string]interface{}{
		"token":     token,
		"identity":  userAccount.Identity,
		"loginTime": now,
	}
	return data, nil
}

func validateLogin(accountName, password string) error {
	if accountName == "" || password == "" {
		return fmt.Errorf("帳號或密碼不可為空")
	}
	return nil
}
