package repo

import (
	"fmt"
	"strings"
	"time"

	common "batchLog/0.core/commonFunction"
	"batchLog/0.core/global"
	gormTable "batchLog/0.core/gorm"
	"batchLog/0.core/logafa"
	"batchLog/0.core/redis"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

func FindAccountByAccountName(tx *gorm.DB,userAccount string) (*gormTable.Account, error) {

	if strings.Contains(userAccount, "@"){
		return FindAccountByEmail(tx,userAccount)
	}
	return FindAccountByUsername(tx,userAccount)
}

func FindAccountByUsername(tx *gorm.DB,username string) (*gormTable.Account, error) {
	var account gormTable.Account
	err := tx.First(&account, "username = ?", username).Error
	if err != nil{
		logafa.Error("查詢帳戶發生錯誤, error: %+v",err)
		return nil, fmt.Errorf("查詢帳戶發生錯誤")
	}
	return &account, nil
}

func FindAccountByEmail(tx *gorm.DB,email string) (*gormTable.Account, error) {
	var account gormTable.Account
	err := tx.First(&account, "email = ?", email).Error
	if err != nil{
		logafa.Error("查詢帳戶發生錯誤, error: %+v",err)
		return nil, fmt.Errorf("查詢帳戶發生錯誤")
	}
	return &account, nil
}

func CreateAccount(tx *gorm.DB,username, password, email string) (uuid.UUID, error) {
	now := time.Now().UTC()
	hashedPassword, _ := common.BcryptHash(password)
	account := gormTable.Account{
		Uuid:          uuid.New(),
		Username:      username,
		Password:      hashedPassword,
		Email:         email,
		LastLoginTime: &now,
	}
	err := tx.Create(&account).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				return uuid.Nil, fmt.Errorf("用戶名 %s 已存在", username)
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

func SaveLoginStatus(nickname, deviceId, token string, now time.Time, expireTime time.Duration)(map[string]interface{},error){
	// 儲存登入狀態至 Redis
	key := fmt.Sprintf("login:%s:%s", nickname, deviceId)
	data := map[string]interface{}{
		"token":     token,
		"loginTime": now,
		"expireAt": now.Add(expireTime),
	}
	err := redis.HSetData(key, data)
	if err != nil {
		logafa.Error("redis 設置失敗，error: %+v",err)
		return data,fmt.Errorf("系統錯誤")
	}
	// 設定過期時間
	global.Repository.Cache.Writing.Expire(global.Repository.Cache.CTX, key, 24*time.Hour)
	return data, nil
}

