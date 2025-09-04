package repo

import (
	common "batchLog/0.core/commonFunction"
	gormTable "batchLog/0.core/gorm"
	"batchLog/0.core/logafa"
	"fmt"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

func CreatePasswordHistory(tx *gorm.DB, accountUuid uuid.UUID, password string)error{
	hashedPassword,_ := common.BcryptHash(password)
	pastPassword := gormTable.PasswordHistory{
		AccountUuid: accountUuid,
		Password: hashedPassword,
	}
	err := tx.Create(&pastPassword).Error
	if err != nil {
		logafa.Error("建立使用者歷史密碼失敗, error: %+v",err)
		return fmt.Errorf("建立使用者歷史密碼失敗")
	}
	return nil
}
