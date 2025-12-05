package redis

import (
	"batchLog/0.core/global"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 封裝方法：寫入一筆 資料（ZADD）
// 並設定過期時間為 24 小時
func ZAddData(ctx context.Context, key string, score float64, byteData []byte) error {
	pipe := global.Repository.Cache.Writing.Pipeline()
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: byteData,
	})
	pipe.Expire(ctx, key, 24*time.Hour)

	_, err := pipe.Exec(ctx)
	return err
}

// 封裝方法：依指定pattern 取得所有 key 值
func KeyScan(ctx context.Context, pattern string) ([]string, error) {
	var cursor uint64
	var keys []string

	for {
		var k []string
		var err error
		k, cursor, err = global.Repository.Cache.Reading.Scan(ctx, cursor, pattern, 100).Result()
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
func ZRangeByScore(ctx context.Context, key string, startTs, endTs int64) ([]string, error) {
	raws, err := global.Repository.Cache.Reading.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", startTs),
		Max: fmt.Sprintf("%d", endTs),
	}).Result()
	if err != nil {
		return nil, err
	}
	return raws, nil
}

// ✅ 移除指定 key 的資料指定時間區段資料
func ZRemRangeByScore(ctx context.Context, client *redis.Client, key string, startTs, endTs int64) error {
	// 移除指定區間資料
	_, err := client.ZRemRangeByScore(ctx, key, fmt.Sprintf("%v", startTs), fmt.Sprintf("%v", endTs)).Result()
	if err != nil {
		return err
	}
	return nil
}
