package initial

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedis(setting jsonModal.RedisDbConfig) *global.Cache {
	if !setting.InUse {
		return nil
	}

	readClient := redis.NewClient(&redis.Options{
		Addr:     setting.Reading.Host + ":" + setting.Reading.Port,
		Password: setting.Reading.Password,
		DB:       0,
	})
	if err := readClient.Ping(context.Background()).Err(); err != nil {
		logafa.Error(" ❌ Redis Read Client 連線失敗: %v", err)
		panic(err)
	}

	writeClient := redis.NewClient(&redis.Options{
		Addr:     setting.Writing.Host + ":" + setting.Writing.Port,
		Password: setting.Writing.Password,
		DB:       0,
	})
	if err := writeClient.Ping(context.Background()).Err(); err != nil {
		logafa.Error(" ❌ Redis Write Client 連線失敗: %v", err)
		panic(err)
	}

	logafa.Debug(" ✅ Redis 資料庫連接成功")
	return &global.Cache{
		Reading: readClient,
		Writing: writeClient,
		CTX: context.Background(),
	}
}