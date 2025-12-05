package redis

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"context"
)

func HSetData(ctx context.Context, key string, mapData map[string]interface{}) error {
	err := global.Repository.Cache.Writing.HSet(ctx, key, mapData).Err()
	if err != nil {
		logafa.Error("Redis HSet 寫入失敗, key: %s, data: %+v", key, mapData)
	}
	return err
}

func HSetFieldData(ctx context.Context, key, field, value string) error {
	err := global.Repository.Cache.Writing.HSet(ctx, key, field, value).Err()
	if err != nil {
		logafa.Error("Redis HSetFieldData 寫入失敗, key: %s, field: %s, value: %s", key, field, value)
	}
	return err
}

func HGetData(ctx context.Context, key, field string) string {
	value, err := global.Repository.Cache.Reading.HGet(ctx, key, field).Result()
	if err != nil {
		logafa.Error("Redis HGet 讀取失敗, key: %s, field: %s, error: %+v", key, field, err)
	}
	return value
}

func HGetAllData(ctx context.Context, key string) map[string]string {
	value, err := global.Repository.Cache.Reading.HGetAll(ctx, key).Result()
	if err != nil {
		logafa.Error("Redis HGetAll 讀取失敗, key: %s, error: %+v", key, err)
	}
	return value
}
