package account

import (
	common "batchLog/core/commonFunction"
	"batchLog/core/global"
	gormTable "batchLog/core/gorm"
	jwtUtil "batchLog/core/jwt"
	"batchLog/core/logafa"
	"batchLog/core/redis"
	"fmt"
	"strings"
	"time"
)

func (account *AccountServiceImpl) Login(ip, accountName, password, deviceId string) (map[string]interface{}, error) {
	var userAccount *gormTable.Account
	var err error
	data := map[string]interface{}{}
	// 驗證帳號
	if strings.Contains(accountName, "@"){
		userAccount,err = account.accountRepo.FindByEmail(accountName)
	}else{
		userAccount,err = account.accountRepo.FindByUsername(accountName)
	}
    if err != nil {
        return data,fmt.Errorf("查無此用戶，請重新嘗試")
    }

	// 驗證密碼
    if !common.BcryptCompare(userAccount.Password, password){
		return data, fmt.Errorf("帳號或密碼錯誤")    }

	// 驗證裝置身份
	device, err := account.deviceRepo.FindByDeviceId(deviceId)
	if err != nil {
		return data, fmt.Errorf("裝置查詢失敗: %v", err)
	}
	member, err := account.memberRepo.FindByUuid(userAccount.Uuid.String())
	if err != nil {
		return data, fmt.Errorf("用戶資料查詢失敗: %v", err)
	}
	if device == nil || member.Uuid != device.MemberInfoUuid {
		return data, fmt.Errorf("非本人裝置，請勿嘗試登入")
	}

    now := time.Now().UTC()
	expireTime := 24 * time.Hour
	token, err := jwtUtil.GenerateJwt(accountName,deviceId,ip, now, expireTime)
	if err != nil{
		logafa.Error("jwt 產生失敗，error: %+v",err)
		return data,fmt.Errorf("系統錯誤")
	}

	// 儲存登入狀態至 Redis
	key := fmt.Sprintf("login:%s:%s", member.NickName, deviceId)
	data = map[string]interface{}{
		"token":     token,
		"loginTime": now,
		"expireAt": now.Add(expireTime),
	}
	err = redis.HSetData(key, data)
	if err != nil {
		logafa.Error("redis 設置失敗，error: %+v",err)
		return data,fmt.Errorf("系統錯誤")
	}
	// 設定過期時間
	global.Repository.Cache.Writing.Expire(global.Repository.Cache.CTX, key, 24*time.Hour)

	return data,nil
}