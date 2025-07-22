package redis

import (
	"batchLog/0.core/global"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 封裝方法：寫入一筆 資料（ZADD）
func ZAddData(key string, score float64, byteData []byte) error {
	return global.Repository.Cache.Writing.ZAdd(global.Repository.Cache.CTX, key, redis.Z{
		Score:  score,
		Member: byteData,
	}).Err()
}

// 封裝方法：依指定pattern 取得所有 key 值
func KeyScan(pattern string) ([]string, error) {
    var cursor uint64
    var keys []string

    for {
        var k []string
        var err error
        k, cursor, err = global.Repository.Cache.Reading.Scan(global.Repository.Cache.CTX, cursor, pattern, 100).Result()
        if err != nil {
            return nil, fmt.Errorf("scan keys failed: %w", err)
        }
        keys = append(keys, k...)
        if cursor == 0 {
            break
        }
    }
    return keys, nil
}


// 封裝方法：依 score 讀取區間資料（ZRANGE）
func ZRangeByScore(key string, startTs, endTs int64) ([]string, error) {
	raws, err := global.Repository.Cache.Reading.ZRangeByScore(global.Repository.Cache.CTX, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", startTs),
		Max: fmt.Sprintf("%d", endTs),
	}).Result()
	if err != nil {
		return nil, err
	}
	return raws, nil
}

// ✅ 移除指定 key 的資料指定時間區段資料
func ZRemRangeByScore(client *redis.Client,key string, startTs, endTs int64) error {
	_, err := client.ZRemRangeByScore(global.Repository.Cache.CTX, key, fmt.Sprintf("%v",startTs), fmt.Sprintf("%v",endTs)).Result()
	return err
}