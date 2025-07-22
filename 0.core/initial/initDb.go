package initial

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMariaDB(setting jsonModal.MariaDbConfig) *global.DB {
	println(123456789)

	if !setting.InUse {
		return nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				setting.Reading.User, setting.Reading.Password,
				setting.Reading.Host, setting.Reading.Port,
				setting.Reading.Name)

	readingDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logafa.Error(" ❌ 無法連接讀取資料庫: %v", err)
		panic(err)
	}
	// 可以用相同 dsn，如果你未來區分寫入配置，再拆成不同 config
	writingDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logafa.Error(" ❌ 無法連接寫入資料庫: %v", err)
		panic(err)
	}
	logafa.Debug(" ✅ MariaDB 資料庫連接成功")
	return global.NewDBRepository(readingDb, writingDb)
}